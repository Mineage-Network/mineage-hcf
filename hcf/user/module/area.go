package module

import (
	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/event"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/mineage-network/mineage-hcf/hcf/area"
	"github.com/mineage-network/mineage-hcf/hcf/crate"
	"github.com/mineage-network/mineage-hcf/hcf/factions"
	"github.com/mineage-network/mineage-hcf/hcf/koth"
	"github.com/mineage-network/mineage-hcf/hcf/rank"
	"github.com/mineage-network/mineage-hcf/hcf/user"
	"github.com/mineage-network/mineage-hcf/hcf/util"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"strings"
	"time"
)

// Area ...
type Area struct {
	u *user.User
	player.NopHandler
}

// NewArea ...
func NewArea(u *user.User) *Area {
	return &Area{u: u}
}

// HandleHurt ...
func (a *Area) HandleHurt(ctx *event.Context, _ *float64, _ *time.Duration, _ world.DamageSource) {
	u := a.u
	w := u.World()
	pos := u.Position()
	if area.Spawn(w).Area().Vec3WithinOrEqualXZ(pos) {
		ctx.Cancel()
		return
	}
}

// HandleMove ...
func (a *Area) HandleMove(ctx *event.Context, newPos mgl64.Vec3, _, _ float64) {
	u := a.u
	w := u.World()

	u.ClearWall()

	cubePos := cube.PosFromVec3(newPos)

	if a.u.Tags().Combat().Active() {
		a := area.Spawn(a.u.World()).Area()
		mul := util.NewAreaVec2(mgl64.Vec2{a.Min().X() - 10, a.Min().Y() - 10}, mgl64.Vec2{a.Max().X() + 10, a.Max().Y() + 10})
		if mul.Vec2WithinOrEqualFloor(mgl64.Vec2{u.Position().X(), u.Position().Z()}) {
			u.SendWall(cubePos, area.Overworld.Spawn(), item.ColourRed())
		}
	}

	if _, ok := u.TimerEnabled(); ok {
		for f, a := range factions.Claims() {
			mul := util.NewAreaVec2(mgl64.Vec2{a.Min().X() - 10, a.Min().Y() - 10}, mgl64.Vec2{a.Max().X() + 10, a.Max().Y() + 10})
			if mul.Vec2WithinOrEqualFloor(mgl64.Vec2{u.Position().X(), u.Position().Z()}) {
				u.SendWall(cubePos, area.NewZone(f.Name(), a), item.ColourBlue())
			}
		}
	}

	if !newPos.ApproxEqual(u.Position()) {
		u.Teleportations().Home().Cancel()
		u.Teleportations().Logout().Cancel()
	}
	if area.Spawn(w).Area().Vec3WithinOrEqualFloorXZ(newPos) && u.Tags().Combat().Active() {
		ctx.Cancel()
		return
	}

	if _, ok := u.TimerEnabled(); ok {
		for _, a := range factions.Claims() {
			if a.Vec3WithinOrEqualXZ(newPos) {
				ctx.Cancel()
				return
			}
		}
	}

	k, ok := koth.Running()
	if ok {
		r := u.Ranks().Highest()
		if k.Area().Vec3WithinOrEqualFloorXZ(newPos) {
			if k.StartCapturing(u, r.Tag(u.Name())) {
				user.Broadcast("koth.capturing", k.Name(), r.Tag(u.Name()))
			}
		} else {
			if k.StopCapturing(u) {
				user.Broadcast("koth.not.capturing", k.Name())
			}
		}
	}

	var areas []area.Zone

	for f, a := range factions.Claims() {
		name := text.Colourf("<red>%s</red>", f.DisplayName())
		if fac, ok := u.Faction(); ok && f == fac.Faction {
			name = text.Colourf("<green>%s</green>", f.DisplayName())
		}
		areas = append(areas, area.NewZone(name, a))
	}

	z, okZone := u.Zone()
	for _, a := range append(area.Protected(w), areas...) {
		if a.Area().Vec3WithinOrEqualFloorXZ(newPos) {
			if z.Area() != a.Area() {
				if okZone {
					if z.Area() != (util.AreaVec2{}) && ok {
						u.Message("area.leave", z.Name())
					}
				}
				u.SetZone(a)
				u.Message("area.enter", a.Name())
				return
			} else {
				return
			}
		}
	}

	if z.Area() != area.Wilderness(w).Area() {
		if z.Area() != (util.AreaVec2{}) {
			if okZone {
				u.Message("area.leave", z.Name())
			}
		}
		u.SetZone(area.Wilderness(w))
		u.Message("area.enter", area.Wilderness(w).Name())
	}
}

// HandleTeleport ...
func (a *Area) HandleTeleport(ctx *event.Context, pos mgl64.Vec3) {
	u := a.u

	if u.Tags().Combat().Active() {
		if area.Spawn(u.World()).Area().Vec3WithinOrEqualXZ(pos) {
			ctx.Cancel()
		}
	}

	if _, ok := u.TimerEnabled(); ok {
		for _, a := range factions.Claims() {
			if a.Vec3WithinOrEqualXZ(pos) {
				ctx.Cancel()
			}
		}
	}
}

// HandleBlockPlace ...
func (a *Area) HandleBlockPlace(ctx *event.Context, pos cube.Pos, _ world.Block) {
	u := a.u
	w := u.World()

	for f, c := range factions.Claims() {
		fa, _ := u.Faction()
		if f == fa.Faction {
			continue
		}
		if !f.Raidable() && c.Vec3WithinOrEqualXZ(pos.Vec3()) {
			ctx.Cancel()
			return
		}
	}
	for _, a := range area.Protected(w) {
		if a.Area().Vec3WithinOrEqualXZ(pos.Vec3()) {
			if !u.Ranks().Contains(rank.Admin{}) || u.Player().GameMode() != world.GameModeCreative {
				ctx.Cancel()
				return
			}
		}
	}
}

// HandleBlockBreak ...
func (a *Area) HandleBlockBreak(ctx *event.Context, pos cube.Pos, _ *[]item.Stack, _ *int) {
	u := a.u
	w := u.World()

	for f, c := range factions.Claims() {
		fa, _ := u.Faction()
		if f == fa.Faction {
			continue
		}
		if !f.Raidable() && c.Vec3WithinOrEqualXZ(pos.Vec3()) {
			ctx.Cancel()
			return
		}
	}
	for _, a := range area.Protected(w) {
		if a.Area().Vec3WithinOrEqualXZ(pos.Vec3()) {
			if !u.Ranks().Contains(rank.Admin{}) || u.Player().GameMode() != world.GameModeCreative {
				ctx.Cancel()
				return
			}
		}
	}
}

// HandleItemUseOnBlock ...
func (a *Area) HandleItemUseOnBlock(ctx *event.Context, pos cube.Pos, _ cube.Face, _ mgl64.Vec3) {
	u := a.u
	w := u.World()

	i, _ := u.HeldItems()
	if _, ok := i.Item().(item.Bucket); ok {
		for _, a := range area.Protected(w) {
			if a.Area().Vec3WithinOrEqualXZ(pos.Vec3()) {
				ctx.Cancel()
				return
			}
		}
	}
	switch it := i.Item().(type) {
	case item.Hoe:
		ctx.Cancel()
		cd := u.Cooldowns().ItemUse().Key(it)
		if cd.Active() {
			return
		} else {
			cd.Set(1 * time.Second)
		}
		if it.Tier == item.ToolTierDiamond {
			_, ok := i.Value("CLAIM_WAND")
			if !ok {
				return
			}

			f, ok := u.Faction()
			if !ok {
				return
			}

			if !strings.EqualFold(f.Leader().Name(), u.Name()) {
				u.Message("faction.not-leader")
				return
			}

			if _, ok = f.Claim(); ok {
				u.Message("faction.has-claim")
				break
			}

			for _, a := range area.Protected(w) {
				if a.Area().Vec3WithinOrEqualXZ(pos.Vec3()) {
					u.Message("faction.area.already-claimed")
					return
				}
				if a.Area().Vec3WithinOrEqualXZ(pos.Vec3().Add(mgl64.Vec3{-1, 0, -1})) {
					u.Message("faction.area.too-close")
					return
				}
			}
			for _, fa := range factions.All() {
				c, ok := fa.Claim()
				if !ok {
					continue
				}
				if c.Vec3WithinOrEqualXZ(pos.Vec3()) {
					u.Message("faction.area.already-claimed")
					return
				}
				if c.Vec3WithinOrEqualXZ(pos.Vec3().Add(mgl64.Vec3{-1, 0, -1})) {
					u.Message("faction.area.too-close")
					return
				}
			}

			pn := 1
			if u.Sneaking() {
				pn = 2
				cpos, _ := u.ClaimPositions()
				ar := util.NewAreaVec2(cpos[0], mgl64.Vec2{float64(pos.X()), float64(pos.Z())})
				x := ar.Max().X() - ar.Min().X()
				y := ar.Max().Y() - ar.Min().Y()
				actualArea := x * y
				if actualArea > 75*75 {
					u.Message("faction.claim.too-big")
					return
				}
				cost := int(actualArea * 5)
				u.Message("faction.claim.cost", cost)
			}
			u.SetClaimPosition(pn-1, mgl64.Vec2{float64(pos.X()), float64(pos.Z())})
			u.Message("faction.claim.set-position", pn, mgl64.Vec2{float64(pos.X()), float64(pos.Z())})
		}
	}
	b := w.Block(pos)

	switch b.(type) {
	case block.WoodFenceGate, block.Chest:
		if _, ok := b.(block.Chest); ok {
			for _, c := range crate.All() {
				if c.Position() == pos.Vec3() {
					return
				}
			}
		}
		if u.Boned() {
			u.Message("user.interaction.boned")
			ctx.Cancel()
			return
		}
		for f, c := range factions.Claims() {
			fa, _ := u.Faction()
			if f == fa.Faction {
				continue
			}
			if !f.Raidable() && c.Vec3WithinOrEqualXZ(pos.Vec3()) {
				ctx.Cancel()
				return
			}
		}
		for _, a := range area.Protected(w) {
			if a.Area().Vec3WithinOrEqualXZ(pos.Vec3()) {
				ctx.Cancel()
				return
			}
		}
	}
}

// HandlePunchAir ...
func (a *Area) HandlePunchAir(_ *event.Context) {
	u := a.u
	w := u.World()

	if !u.Sneaking() {
		return
	}

	i, _ := u.HeldItems()

	it, ok := i.Item().(item.Hoe)
	if !ok || it.Tier != item.ToolTierDiamond {
		return
	}
	_, ok = i.Value("CLAIM_WAND")
	if !ok {
		return
	}
	f, ok := u.Faction()
	if !ok {
		return
	}
	if !strings.EqualFold(f.Leader().Name(), u.Name()) {
		u.Message("faction.not-leader")
		return
	}
	_, ok = f.Claim()
	if ok {
		u.Message("faction.has-claim")
		return
	}
	pos, ok := u.ClaimPositions()
	if !ok {
		u.Message("faction.area.too-close")
		return
	}
	claim := util.NewAreaVec2(pos[0], pos[1])
	var blocksPos []cube.Pos
	min := claim.Min()
	max := claim.Max()
	for x := min[0]; x <= max[0]; x++ {
		for y := min[1]; y <= max[1]; y++ {
			blocksPos = append(blocksPos, cube.PosFromVec3(mgl64.Vec3{x, 0, y}))
		}
	}
	for _, a := range area.Protected(w) {
		for _, b := range blocksPos {
			if a.Area().Vec3WithinOrEqualXZ(b.Vec3()) {
				u.Message("faction.area.already-claimed")
				return
			}
			if a.Area().Vec3WithinOrEqualXZ(b.Vec3().Add(mgl64.Vec3{-1, 0, -1})) {
				u.Message("faction.area.too-close")
				return
			}
		}
		if a.Area().Vec2WithinOrEqual(pos[0]) || a.Area().Vec2WithinOrEqual(pos[1]) {
			u.Message("faction.area.already-claimed")
			return
		}
		if a.Area().Vec2WithinOrEqual(pos[0].Add(mgl64.Vec2{-1, -1})) || a.Area().Vec2WithinOrEqual(pos[1].Add(mgl64.Vec2{-1, -1})) {
			u.Message("faction.area.too-close")
			return
		}
	}

	for _, fa := range factions.All() {
		c, ok := fa.Claim()
		if !ok {
			continue
		}
		for _, b := range blocksPos {
			if c.Vec3WithinOrEqualXZ(b.Vec3()) {
				u.Message("faction.area.already-claimed")
				return
			}
			if c.Vec3WithinOrEqualXZ(b.Vec3().Add(mgl64.Vec3{-1, 0, -1})) {
				u.Message("faction.area.too-close")
				return
			}
		}
		if c.Vec2WithinOrEqual(pos[0]) || c.Vec2WithinOrEqual(pos[1]) {
			u.Message("faction.area.already-claimed")
			return
		}
		if c.Vec2WithinOrEqual(pos[0].Add(mgl64.Vec2{-1, -1})) || c.Vec2WithinOrEqual(pos[1].Add(mgl64.Vec2{-1, -1})) {
			u.Message("faction.area.too-close")
			return
		}
	}

	x := claim.Max().X() - claim.Min().X()
	y := claim.Max().Y() - claim.Min().Y()
	actualArea := x * y
	if actualArea > 75*75 {
		u.Message("faction.claim.too-big")
		return
	}
	cost := actualArea * 5

	if fa, ok := u.Faction(); ok {
		if float64(fa.Balance()) < cost {
			u.Message("faction.claim.no-money")
			return
		}
	}

	f.ReduceBalance(cost)
	f.SetClaim(claim)
	u.Message("command.claim.success", pos[0], pos[1], cost)
}
