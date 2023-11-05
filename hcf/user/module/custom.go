package module

import (
	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/entity"
	"github.com/df-mc/dragonfly/server/event"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/mineage-network/mineage-hcf/hcf/crate"
	ench "github.com/mineage-network/mineage-hcf/hcf/enchantment"
	"github.com/mineage-network/mineage-hcf/hcf/user"
	"github.com/mineage-network/mineage-hcf/hcf/util"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"time"
)

// Custom ...
type Custom struct {
	u *user.User
	player.NopHandler
}

// NewCustom ...
func NewCustom(u *user.User) *Custom {
	return &Custom{u: u}
}

// HandleAttackEntity ...
func (c *Custom) HandleAttackEntity(ctx *event.Context, e world.Entity, _, _ *float64, _ *bool) {
	if ctx.Cancelled() {
		return
	}

	u := c.u
	p := u.Player()

	t, ok := user.LookupEntity(e)
	if !ok {
		return
	}

	if t.FullInvis() {
		t.ToggleFullInvis()
		for _, us := range user.All() {
			us.View(t.Player())
		}
	}

	if u.FullInvis() {
		u.ToggleFullInvis()
		for _, us := range user.All() {
			us.View(u.Player())
		}
	}

	arm := p.Armour()
	for _, a := range arm.Slots() {
		for _, e := range a.Enchantments() {
			if att, ok := e.Type().(ench.AttackEnchantment); ok {
				att.AttackEntity(p, t.Player())
			}
		}
	}

	/*held, left := p.HeldItems()
	typ, ok := item3.SpecialItem(held)
	if ok {
		switch kind := typ.(type) {
		case item3.AntiBuildBoneType:
			sp := u.Cooldowns().SpecialAbilities()
			if sp.Active() {
				p.Message(text.Colourf("<red>You are on Partner Items cooldown for %.1f seconds</red>", sp.Remaining().Seconds()))
				break
			}
			cd := u.Cooldowns().SpecialItems()
			if cd.Active(kind) {
				p.Message(text.Colourf("<red>You are on bone cooldown for %.1f seconds</red>", cd.Remaining(kind).Seconds()))
				break
			}
			t.AddBoneHit(p)
			if t.Boned() {
				t.Player().Message(text.Colourf("<red>You have been boned by %s</red>", p.Name()))
				p.Message(text.Colourf("<green>You have boned %s</green>", t.Name()))
				cd.Set(kind, time.Minute)
				sp.Set(time.Second * 10)

				p.SetHeldItems(u.SubtractItem(held, 1), left)
			} else {
				p.Message(text.Colourf("<green>You have hit %s with a bone %d times</green>", t.Name(), t.BoneHits(p)))
			}
		case item3.PearlDisablerType:
			sp := u.Cooldowns().SpecialAbilities()
			if sp.Active() {
				p.Message(text.Colourf("<red>You are on Partner Items cooldown for %.1f seconds</red>", sp.Remaining().Seconds()))
				break
			}
			cd := u.Cooldowns().SpecialItems()
			if cd.Active(kind) {
				p.Message(text.Colourf("<red>You are on pearl disabler cooldown for %.1f seconds</red>", cd.Remaining(kind).Seconds()))
				break
			}
			// TODO: IDK how i wanna implement this
		case item3.ScramblerType:
			sp := u.Cooldowns().SpecialAbilities()
			if sp.Active() {
				p.Message(text.Colourf("<red>You are on Partner Items cooldown for %.1f seconds</red>", sp.Remaining().Seconds()))
				break
			}
			cd := u.Cooldowns().SpecialItems()
			if cd.Active(kind) {
				p.Message(text.Colourf("<red>You are on Scrambler cooldown for %.1f seconds</red>", cd.Remaining(kind).Seconds()))
				break
			}
			t.AddScramblerHit(u.Player())
			if t.ScramblerHits(u.Player()) >= 3 {
				t.ResetScramblerHits(u.Player())
				inv := t.Player().Inventory()
				for i := 36; i <= 44; i++ {
					j := rand.Intn(i+1-36) + 36
					it1, _ := inv.Item(i)
					it2, _ := inv.Item(j)
					inv.SetItem(i, it1)
					inv.SetItem(j, it2)
				}
			}
		}
	}*/
}

// HandleItemPickup ...
func (c *Custom) HandleItemPickup(_ *event.Context, st *item.Stack) {
	/*for _, sp := range item3.SpecialItems() {
		if _, ok := st.Value(sp.Key()); ok {
			*st = item3.NewSpecialItem(sp, st.Count())
		}
	}*/
}

// HandleBlockPlace ...
func (c *Custom) HandleBlockPlace(ctx *event.Context, pos cube.Pos, b world.Block) {
	u := c.u
	p := u.Player()
	/*w := u.World()

	switch b.(type) {
	case block.EnderChest:
		held, left := p.HeldItems()
		if _, ok := held.Value("PARTNER_PACKAGE"); !ok {
			break
		}

		keys := item3.SpecialItems()
		n, err := rdm.Int(rdm.Reader, big.NewInt(int64(len(keys))))
		if err != nil {
			panic(err)
		}
		i := item3.NewSpecialItem(keys[n.Int64()], rand.Intn(3)+1)

		ctx.Cancel()

		p.SetHeldItems(u.SubtractItem(held, 1), left)

		u.AddItemOrDrop(i)

		w.AddEntity(entity.NewFirework(pos.Vec3(), cube.Rotation{90, 90}, item.Firework{
			Duration: 0,
			Explosions: []item.FireworkExplosion{
				{
					Shape:   item.FireworkShapeStar(),
					Trail:   true,
					Colour:  util.RandomColour(),
					Twinkle: true,
				},
			},
		}))
		return
	}*/

	if ctx.Cancelled() {
		return
	}

	if u.Boned() {
		p.Message(text.Colourf("<red>You cannot place custom while boned</red>"))
		ctx.Cancel()
	}
}

// HandleBlockBreak ...
func (c *Custom) HandleBlockBreak(ctx *event.Context, pos cube.Pos, drops *[]item.Stack, xp *int) {
	u := c.u
	p := u.Player()

	if ctx.Cancelled() {
		return
	}

	if u.Boned() {
		p.Message(text.Colourf("<red>You cannot break custom while boned</red>"))
		ctx.Cancel()
	}
}

// HandleStartBreak ...
func (c *Custom) HandleStartBreak(ctx *event.Context, pos cube.Pos) {
	u := c.u
	p := u.Player()

	w := p.World()
	b := w.Block(pos)

	/*held, _ := p.HeldItems()
	typ, ok := item3.SpecialItem(held)
	if ok {
		cd := u.Cooldowns().SpecialAbilities()
		if cd.Active() {
			u.Player().SendJukeboxPopup(lang.Translatef(u.Locale(), "popup.cooldown.partner_item", cd.Remaining().Seconds()))
			ctx.Cancel()
			return
		}
		if spi := u.Cooldowns().SpecialItems(); spi.Active(typ) {
			u.Player().SendJukeboxPopup(lang.Translatef(u.Locale(), "popup.cooldown.partner_item.item", typ.Name(), spi.Remaining(typ).Seconds()))
			ctx.Cancel()
		} else {
			u.Player().SendJukeboxPopup(lang.Translatef(u.Locale(), "popup.ready.partner_item.item", typ.Name()))
			ctx.Cancel()
		}
	}*/

	for _, c := range crate.All() {
		if _, ok := b.(block.Chest); ok && pos.Vec3() == c.Position() {
			p.OpenBlockContainer(pos)
			ctx.Cancel()
		}
	}
}

// HandlePunchAir ...
func (c *Custom) HandlePunchAir(ctx *event.Context) {
	/*u := c.u
	p := u.Player()

	held, _ := p.HeldItems()
	typ, ok := item3.SpecialItem(held)
	if ok {
		cd := u.Cooldowns().SpecialAbilities()
		if cd.Active() {
			u.Player().SendJukeboxPopup(lang.Translatef(u.Locale(), "popup.cooldown.partner_item", cd.Remaining().Seconds()))
			ctx.Cancel()
			return
		}
		if spi := u.Cooldowns().SpecialItems(); spi.Active(typ) {
			u.Player().SendJukeboxPopup(lang.Translatef(u.Locale(), "popup.cooldown.partner_item.item", typ.Name(), spi.Remaining(typ).Seconds()))
			ctx.Cancel()
		} else {
			u.Player().SendJukeboxPopup(lang.Translatef(u.Locale(), "popup.ready.partner_item.item", typ.Name()))
			ctx.Cancel()
		}
	}*/
}

// HandleItemUseOnBlock ...
func (c *Custom) HandleItemUseOnBlock(ctx *event.Context, pos cube.Pos, face cube.Face, clickPos mgl64.Vec3) {
	u := c.u
	p := u.Player()
	w := p.World()

	i, left := p.HeldItems()
	b := w.Block(pos)

	for _, c := range crate.All() {
		if _, ok := b.(block.Chest); ok && pos.Vec3() == c.Position() && p.Sneaking() {
			ctx.Cancel()
			if _, ok := i.Value(c.EncodeCrate()); !ok {
				u.Player().Message(text.Colourf("<red>You need a %s key to open this crate</red>", util.StripMinecraftColour(c.Name())))
				break
			}
			u.AddItemOrDrop(ench.AddEnchantmentLore(c.Reward()))

			p.SetHeldItems(u.SubtractItem(i, 1), left)

			w.AddEntity(entity.NewFirework(c.PositionMiddle().Add(mgl64.Vec3{0, 1, 0}), cube.Rotation{90, 90}, item.Firework{
				Duration: 0,
				Explosions: []item.FireworkExplosion{
					{
						Shape:   item.FireworkShapeStar(),
						Trail:   true,
						Colour:  util.RandomColour(),
						Twinkle: true,
					},
				},
			}))
		}
	}

	switch usable := i.Item().(type) {
	case item.EnderPearl:
		if f, ok := b.(block.WoodFenceGate); ok && f.Open {
			if cd := u.Cooldowns().Pearl(); !cd.Active() {
				cd.Set(15 * time.Second)
				usable.Use(w, p, &item.UseContext{})
				p.SetHeldItems(u.SubtractItem(i, 1), left)
				ctx.Cancel()
			}
		}
	}

	if ctx.Cancelled() {
		return
	}

	switch b.(type) {
	case block.WoodFenceGate, block.Chest:
		if u.Boned() {
			p.Message(text.Colourf("<red>You cannot interact with custom while boned</red>"))
			ctx.Cancel()
		}
	}
}

// TODO:

// HandleItemUse ...
func (c *Custom) HandleItemUse(ctx *event.Context) {
	u := c.u
	p := u.Player()

	held, left := p.HeldItems()
	if v, ok := held.Value("MONEY_NOTE"); ok {
		u.IncreaseBalance(v.(float64))
		p.SetHeldItems(u.SubtractItem(held, 1), left)
		p.Message(text.Colourf("<green>You have deposited $%.0f into your bank account</green>", v.(float64)))
		return
	}

	/*if v, ok := item3.SpecialItem(held); ok {
		cd := u.Cooldowns().SpecialAbilities()
		if cd.Active() {
			p.Message(text.Colourf("<red>You are on Partner Items cooldown for %.1f seconds</red>", cd.Remaining().Seconds()))
			ctx.Cancel()
			return
		}
		switch typ := v.(type) {
		case item3.SigilType:
			if spi := u.Cooldowns().SpecialItems(); spi.Active(typ) {
				p.Message(text.Colourf("<red>You are on Sigil cooldown for %.1f seconds</red>", spi.Remaining(typ).Seconds()))
				ctx.Cancel()
				return
			} else {
				spi.Set(typ, time.Minute*2)
			}

			oe := u.NearbyEnemies(10)
			p.Message(text.Colourf("<green>The All-Powerful Allah has blinded the enemy!</green>"))
			for _, e := range oe {
				e.World().AddEntity(entity.NewLightningWithDamage(e.Position(), 3, false, 0))
				e.AddEffect(effect.New(effect.Poison{}, 1, time.Second*3))
				e.AddEffect(effect.New(effect.Blindness{}, 2, time.Second*7))
				e.AddEffect(effect.New(effect.Nausea{}, 2, time.Second*7))
				e.Player().Message(text.Colourf("<red>Those are the ones whom Allah has cursed; so He has made them deaf, and made their eyes blind! (47:23)</red>"))
			}

			p.SetHeldItems(u.SubtractItem(held, 1), left)

			cd.Set(time.Second * 10)
		case item3.SwitcherBallType:
			cd.Set(time.Second * 10)
		case item3.FullInvisibilityType:
			if spi := u.Cooldowns().SpecialItems(); spi.Active(typ) {
				p.Message(text.Colourf("<red>You are on Full Invisiblity cooldown for %.1f seconds</red>", spi.Remaining(typ).Seconds()))
				ctx.Cancel()
				return
			} else {
				spi.Set(typ, time.Minute*2)
			}
			for _, us := range user.All() {
				if us == c.u {
					continue
				}
				us.Hide(c.u.Player())
			}
			cd.Set(time.Second * 10)
		case item3.NinjaStarType:
			t, ok := u.LastAttacker()
			if !ok {
				p.Message(text.Colourf("<red>You have not been attacked recently</red>"))
				ctx.Cancel()
				return
			}
			if spi := u.Cooldowns().SpecialItems(); spi.Active(typ) {
				p.Message(text.Colourf("<red>You are on Ninja Star cooldown for %.1f seconds</red>", spi.Remaining(typ).Seconds()))
				ctx.Cancel()
				return
			} else {
				spi.Set(typ, time.Minute*2)
			}

			p.SetHeldItems(u.SubtractItem(held, 1), left)

			cd.Set(time.Second * 10)
			go func() {
				deadline := time.Now().Add(time.Second * 3)
				p.Message(text.Colourf("<green>Teleporting to <red>%s</red> in 3 seconds</green>", t.Name()))
				for range time.Tick(time.Second) {
					if _, ok := user.LookupName(t.Name()); !ok {
						p.Message(text.Colourf("<red>%s disconnected, ninja ability canceled</red>", t.Name()))
						return
					}
					remaining := deadline.Sub(time.Now())
					if remaining <= 0 {
						p.Teleport(t.Player().Position())
						break
					}
					p.Message(text.Colourf("<green>Teleporting to <red>%s</red> in %d seconds</green>", t.Name(), int(remaining.Seconds())+1))
				}
			}()
		}
	}*/
	switch held.Item().(type) {
	case item.EnderPearl:
		cd := u.Cooldowns().Pearl()
		if u.PearlDisabled() {
			ctx.Cancel()
			p.Message(text.Colourf("<red>You're pearl disabled!</red>"))
			u.TogglePearlDisable()
			cd.Set(15 * time.Second)
		}
		if cd.Active() {
			ctx.Cancel()
			return
		} else {
			cd.Set(15 * time.Second)
		}
	}
}
