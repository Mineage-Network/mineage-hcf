package user

import (
	"fmt"
	"github.com/df-mc/atomic"
	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/entity"
	"github.com/df-mc/dragonfly/server/entity/effect"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/scoreboard"
	"github.com/df-mc/dragonfly/server/player/skin"
	"github.com/df-mc/dragonfly/server/session"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/particle"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/google/uuid"
	"github.com/mineage-network/mineage-hcf/hcf/area"
	"github.com/mineage-network/mineage-hcf/hcf/backend/lang"
	"github.com/mineage-network/mineage-hcf/hcf/factions"
	"github.com/mineage-network/mineage-hcf/hcf/koth"
	"github.com/mineage-network/mineage-hcf/hcf/rank"
	"github.com/mineage-network/mineage-hcf/hcf/util"
	"github.com/mineage-network/mineage-hcf/hcf/util/class"
	"github.com/mineage-network/mineage-hcf/hcf/util/sets"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"golang.org/x/exp/maps"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"math"
	"math/rand"
	"net"
	"strings"
	"sync"
	"time"

	_ "unsafe"
)

var (
	userMu    sync.Mutex
	users     = map[*player.Player]*User{}
	staff     = map[*player.Player]*User{}
	admins    = map[*player.Player]*User{}
	usersXUID = map[string]*User{}

	frozen      = sets.New[string]()
	deathBanned = sets.New[string]()
)

// All returns all users.
func All() []*User {
	userMu.Lock()
	defer userMu.Unlock()
	u := make([]*User, 0, len(users))
	for _, user := range users {
		u = append(u, user)
	}
	return u
}

// LookupName looks up a user by name.
func LookupName(name string) (*User, bool) {
	userMu.Lock()
	defer userMu.Unlock()
	for p, u := range users {
		if strings.EqualFold(p.Name(), name) {
			return u, true
		}
	}
	return nil, false
}

// LookupEntity looks up a user by entity.
func LookupEntity(e world.Entity) (*User, bool) {
	if p, ok := e.(*player.Player); ok {
		return Lookup(p)
	}
	return nil, false
}

// LookupXUID looks up a user by XUID.
func LookupXUID(xuid string) (*User, bool) {
	userMu.Lock()
	defer userMu.Unlock()
	u, ok := usersXUID[xuid]
	return u, ok
}

// Lookup returns a user by player.
func Lookup(p *player.Player) (*User, bool) {
	userMu.Lock()
	u, ok := users[p]
	userMu.Unlock()
	return u, ok
}

// Broadcast sends a message to all users.
func Broadcast(key string, args ...interface{}) {
	for _, u := range All() {
		u.Message(key, args...)
	}
}

// Staff returns a slice of all staff online.
func Staff() []*User {
	userMu.Lock()
	defer userMu.Unlock()
	return maps.Values(staff)
}

// Admins returns a slice of all admins online.
func Admins() []*User {
	userMu.Lock()
	defer userMu.Unlock()
	return maps.Values(admins)
}

// Alert alerts all staff users with an action performed by a cmd.Source.
func Alert(s cmd.NamedTarget, key string, args ...any) {
	for _, u := range Admins() {
		u.Message("staff.alert",
			s.Name(),
			fmt.Sprintf(lang.Translate(u.Player().Locale(), key), args...),
		)
	}
}

// User ...
type User struct {
	s *session.Session
	p *player.Player

	hashedAddress string
	address       net.Addr
	zone          atomic.Value[area.Zone]

	balance   atomic.Value[float64]
	reclaimed bool

	lastPlayerSpecificSkin map[string]skin.Skin

	lastMessageFrom atomic.Value[string]

	lastAttackerName atomic.Value[string]
	lastAttackTime   atomic.Value[time.Time]
	lastMessage      atomic.Value[time.Time]

	stats *Stats

	ranks     *Ranks
	lives     atomic.Value[*Lives]
	mute, ban atomic.Value[Punishment]

	reportSince atomic.Value[time.Time]

	effectsNoLossMu sync.Mutex
	effectsNoLoss   map[effect.Type]EffectNoLoss

	scoreboard atomic.Value[*scoreboard.Scoreboard]

	claimPositions [2]mgl64.Vec2

	tags           *Tags
	cooldowns      *Cooldowns
	teleportations *Teleportations

	vanished    atomic.Bool
	deathBanned atomic.Bool

	factionCreateDelay atomic.Value[time.Time]

	boneHits map[string]int
	bone     *util.Cooldown

	scramblerHits map[string]int

	pearlDisabled bool

	bardEnergy atomic.Value[float64]
	class      atomic.Value[util.Class]

	invitations []string

	faction atomic.Value[Faction]

	logged  bool
	logging atomic.Bool

	frozen   atomic.Bool
	chatType atomic.Value[ChatType]

	wallBlocks   map[cube.Pos]float64
	wallBlocksMu sync.Mutex

	staffMode atomic.Bool

	waypoints map[string]Waypoint

	timer time.Time
	sotw  bool

	fullInvis           bool
	invisibilityUpdated atomic.Bool

	close chan struct{}
}

// NewUser ...
func NewUser(
	p *player.Player,
	r *Ranks,
	kills, deaths int,
	l *Lives,
	factionCreateDelay int64,
	kits map[string]time.Time,
	reclaimed bool,
	hashedAddress string,
	mute, ban Punishment,
	balance float64,
	timer time.Time,
	sotw bool,
) *User {
	u := &User{
		p:             p,
		s:             player_session(p),
		address:       p.Addr(),
		close:         make(chan struct{}),
		effectsNoLoss: map[effect.Type]EffectNoLoss{},

		cooldowns: NewCooldowns(),
		stats:     NewStats(kills, deaths),

		reclaimed:          reclaimed,
		scoreboard:         *atomic.NewValue(scoreboard.New()),
		factionCreateDelay: *atomic.NewValue(time.UnixMilli(factionCreateDelay)),
		ranks:              r,

		lastPlayerSpecificSkin: map[string]skin.Skin{},

		hashedAddress: hashedAddress,

		bone:     util.NewCooldown(),
		boneHits: map[string]int{},

		scramblerHits: map[string]int{},

		waypoints: map[string]Waypoint{},

		mute:     *atomic.NewValue(mute),
		ban:      *atomic.NewValue(ban),
		lives:    *atomic.NewValue(l),
		chatType: *atomic.NewValue(ChatTypeGlobal()),
		balance:  *atomic.NewValue(balance),
		logging:  *atomic.NewBool(true),

		wallBlocks: make(map[cube.Pos]float64),

		staffMode: *atomic.NewBool(false),

		timer: timer,
		sotw:  sotw,
	}

	u.p.Armour().Handle(&armourHandler{
		u: u,
		p: p,
	})

	u.ranks.sortRanks()

	u.tags = NewTags(u)
	u.teleportations = NewTeleportations(u)

	for k, d := range kits {
		u.cooldowns.kits.Set(k, time.Until(d))
	}

	userMu.Lock()
	users[p] = u
	if u.ranks.Staff() {
		staff[p] = u
		for _, s := range staff {
			l := s.Player().Locale()
			s.Player().Message(lang.Translatef(l,
				"staff.joined",
				cases.Title(l).String(u.Ranks().Highest().Name()),
				u.Player().Name(),
			))
		}
	}
	if u.ranks.Contains(rank.Admin{}, rank.Operator{}) {
		admins[p] = u
	}
	usersXUID[p.XUID()] = u
	if frozen.Contains(p.XUID()) {
		u.p.SetImmobile()
		u.frozen.Toggle()
	}
	if deathBanned.Contains(p.XUID()) {
		u.deathBanned.Toggle()
	}
	userMu.Unlock()

	p.SetNameTag(text.Colourf("<red>%s</red>", p.Name()))

	fa, ok := factions.LookupMember(p)
	if ok {
		u.SetFaction(fa)
		for _, m := range (Faction{fa}).Users() {
			m.UpdateState()
		}
	}
	u.SetClass(class.Resolve(p))

	go u.startTicker()
	u.logging.Store(false)
	return u
}

// Address returns the address of the user.
func (u *User) Address() net.Addr {
	return u.address
}

// Vanished returns whether the user is vanished or not.
func (u *User) Vanished() bool {
	return u.vanished.Load()
}

// ToggleVanish toggles the user's vanish state.
func (u *User) ToggleVanish() {
	u.vanished.Toggle()
}

// ToggleStaffMode ...
func (u *User) ToggleStaffMode() {
	u.ToggleVanish()
	// TODO: Send items, scoreboard (if needed)
	// and a message.
	u.staffMode.Toggle()
}

// Logged returns whether the user is logged.
func (u *User) Logged() bool {
	return u.logged
}

// Logging returns whether the user is logging.
func (u *User) Logging() bool {
	return u.logging.Load()
}

// ToggleLogging ...
func (u *User) ToggleLogging() {
	u.logging.Toggle()
}

// Balance returns the user's balance.
func (u *User) Balance() float64 {
	return u.balance.Load()
}

// ReduceBalance reduces the user's balance.
func (u *User) ReduceBalance(amount float64) {
	u.balance.Store(u.Balance() - amount)
}

// IncreaseBalance increases the user's balance.
func (u *User) IncreaseBalance(amount float64) {
	u.balance.Store(u.Balance() + amount)
}

// Focusing returns whether the user is focusing.
func (u *User) Focusing() ([]*User, bool) {
	fa, ok := u.Faction()
	if !ok {
		return nil, false
	}
	if p, ok := fa.FocusedPlayer(); ok {
		f, ok := LookupName(p)
		if !ok {
			return nil, false
		}
		return []*User{f}, true
	}
	if f, ok := fa.FocusedFaction(); ok {
		return (Faction{f}).Users(), true
	}
	return nil, false
}

// UpdateState updates the user's state to its viewers.
func (u *User) UpdateState() {
	for _, v := range u.viewers() {
		v.ViewEntityState(u.p)
	}
}

// AddItemOrDrop adds an item to the user's inventory or drops it if the inventory is full.
func (u *User) AddItemOrDrop(it item.Stack) {
	if _, err := u.p.Inventory().AddItem(it); err != nil {
		u.DropItem(it)
	}
}

// UpdateChatType updates the chat type for the user.
func (u *User) UpdateChatType(t ChatType) {
	u.chatType.Store(t)
}

// ChatType returns the chat type the user is currently using.
func (u *User) ChatType() ChatType {
	return u.chatType.Load()
}

// CanSendMessage returns true if the user can send a message.
func (u *User) CanSendMessage() bool {
	if u.Ranks().Contains(rank.Operator{}) {
		return true
	}
	return time.Since(u.lastMessage.Load()) > time.Second*2
}

// RenewLastMessage renews the last time a message was sent from a player.
func (u *User) RenewLastMessage() {
	u.lastMessage.Store(time.Now())
}

// SetLastMessageFrom sets the player passed as the last player who messaged the user.
func (u *User) SetLastMessageFrom(p *player.Player) {
	u.lastMessageFrom.Store(p.XUID())
}

// LastMessageFrom returns the last user that messaged the user.
func (u *User) LastMessageFrom() (*User, bool) {
	u, ok := LookupXUID(u.lastMessageFrom.Load())
	return u, ok
}

// CanAttack returns whether the user can attack another User.
func (u *User) CanAttack(t *User) bool {
	_, uTimer := u.TimerEnabled()
	_, tTimer := t.TimerEnabled()
	if uTimer || tTimer || t.SOTW() || u.SOTW() {
		return false
	}
	w := u.World()
	if w != t.World() {
		return false
	}
	if area.Spawn(w).Area().Vec3WithinOrEqualFloorXZ(u.Position()) || area.Spawn(w).Area().Vec3WithinOrEqualFloorXZ(t.Position()) {
		return false
	}
	uFaction, ok := u.Faction()
	if ok {
		if tFaction, ok := t.Faction(); ok && uFaction == tFaction {
			return false
		}
	}
	return true
}

// Reclaimed returns whether the user has been reclaimed.
func (u *User) Reclaimed() bool {
	return u.reclaimed
}

// Reclaim reclaims the user.
func (u *User) Reclaim() {
	u.reclaimed = true
}

// ResetReclaim resets the reclaim status of the user.
func (u *User) ResetReclaim() {
	u.reclaimed = false
}

// DropItem drops the item stack provided on the ground.
func (u *User) DropItem(it item.Stack) {
	p := u.p
	w, pos := p.World(), p.Position()
	ent := entity.NewItem(it, pos)
	ent.SetVelocity(mgl64.Vec3{rand.Float64()*0.2 - 0.1, 0.2, rand.Float64()*0.2 - 0.1})
	w.AddEntity(ent)
}

// Boned returns whether the user has been boned.
func (u *User) Boned() bool {
	return u.bone.Active()
}

// BoneHits returns the number of bone hits of the user.
func (u *User) BoneHits(p *player.Player) int {
	hits, ok := u.boneHits[p.XUID()]
	if !ok {
		return 0
	}
	return hits
}

// AddBoneHit adds a bone hit to the user.
func (u *User) AddBoneHit(p *player.Player) {
	u.boneHits[p.XUID()]++
	if u.boneHits[p.XUID()] >= 3 {
		u.ResetBoneHits(p)
		u.bone.Set(15 * time.Second)
	}
}

// ResetBoneHits resets the bone hits of the user.
func (u *User) ResetBoneHits(p *player.Player) {
	u.boneHits[p.XUID()] = 0
}

// ScramblerHits returns the number of scrambler hits of the user.
func (u *User) ScramblerHits(p *player.Player) int {
	hits, ok := u.scramblerHits[p.XUID()]
	if !ok {
		return 0
	}
	return hits
}

// AddScramblerHit adds a scrambler hit to the user.
func (u *User) AddScramblerHit(p *player.Player) {
	u.scramblerHits[p.XUID()]++
}

// ResetScramblerHits resets the scrambler hits of the user.
func (u *User) ResetScramblerHits(p *player.Player) {
	u.scramblerHits[p.XUID()] = 0
}

// PearlDisabled returns whether the user is pearl disabled.
func (u *User) PearlDisabled() bool {
	return u.pearlDisabled
}

// TogglePearlDisable toggles the pearl disabler
func (u *User) TogglePearlDisable() {
	u.pearlDisabled = !u.pearlDisabled
}

// FullInvis returns whether the user is full invis.
func (u *User) FullInvis() bool {
	return u.fullInvis
}

// ToggleFullInvis toggles full invis for the user
func (u *User) ToggleFullInvis() {
	u.fullInvis = !u.fullInvis
}

// Message sends a message to the user.
func (u *User) Message(key string, args ...interface{}) {
	u.p.Message(lang.Translatef(u.Locale(), key, args...))
}

// SetZone sets the zone of the user.
func (u *User) SetZone(z area.Zone) {
	u.zone.Store(z)
}

// Zone returns the zone of the user.
func (u *User) Zone() (area.Zone, bool) {
	z := u.zone.Load()
	return z, len(z.Name()) > 0
}

// ResetClaimPositions resets the claim positions of the user.
func (u *User) ResetClaimPositions() {
	u.claimPositions = [2]mgl64.Vec2{}
}

// ClaimPositions returns the claim positions of the user.
func (u *User) ClaimPositions() ([2]mgl64.Vec2, bool) {
	return u.claimPositions, u.claimPositions[0] != mgl64.Vec2{} && u.claimPositions[1] != mgl64.Vec2{}
}

// SetClaimPosition sets a claim position of the user.
func (u *User) SetClaimPosition(i int, v mgl64.Vec2) {
	u.claimPositions[i] = v
}

// Teleportations returns the teleportations of the user.
func (u *User) Teleportations() *Teleportations {
	return u.teleportations
}

// SubtractItem subtracts an item from the user's inventory.
func (u *User) SubtractItem(s item.Stack, d int) item.Stack {
	if !u.p.GameMode().CreativeInventory() && d != 0 {
		return s.Grow(-d)
	}
	return s
}

// World returns the world of the user.
func (u *User) World() *world.World {
	return u.p.World()
}

// Faction returns the faction of the user.
func (u *User) Faction() (Faction, bool) {
	f := u.faction.Load()
	return f, f.Faction != nil
}

// Ranks returns the ranks of the user.
func (u *User) Ranks() *Ranks {
	return u.ranks
}

// Lives ...
func (u *User) Lives() *Lives {
	return u.lives.Load()
}

// ReduceBardEnergy reduces the bard energy of the user.
func (u *User) ReduceBardEnergy(amount float64) {
	u.bardEnergy.Store(math.Max(0, u.BardEnergy()-amount))
}

// LastAttacker returns the last attacker of the user.
func (u *User) LastAttacker() (*User, bool) {
	if time.Since(u.lastAttackTime.Load()) > 15*time.Second {
		return nil, false
	}
	name := u.lastAttackerName.Load()
	if len(name) == 0 {
		return nil, false
	}
	return LookupName(name)
}

// SetLastAttacker sets the last attacker of the user.
func (u *User) SetLastAttacker(t *User) {
	u.lastAttackerName.Store(t.Name())
	u.lastAttackTime.Store(time.Now())
}

// ResetLastAttacker resets the last attacker of the user.
func (u *User) ResetLastAttacker() {
	u.lastAttackerName.Store("")
	u.lastAttackTime.Store(time.Time{})
}

// TimerEnabled returns whether the user has timer enabled.
func (u *User) TimerEnabled() (time.Duration, bool) {
	if time.Now().Before(u.timer) {
		return time.Until(u.timer), true
	}
	return time.Until(time.Now()), false
}

// DisableTimer disables the timer of the user.
func (u *User) DisableTimer() {
	u.timer = time.Now()
}

// EnableTimer enables the timer of the player.
func (u *User) EnableTimer() {
	u.timer = time.Now().Add(time.Hour)
}

// TimerExpiry returns the expiry time of the timer
func (u *User) TimerExpiry() time.Time {
	return u.timer
}

// ToggleSOTW toggles the SOTW of the user.
func (u *User) ToggleSOTW() {
	u.sotw = !u.sotw
}

// SOTW returns if the user has SOTW enabled.
func (u *User) SOTW() bool {
	return u.sotw
}

// SetClass sets the class of the user.
func (u *User) SetClass(c util.Class) {
	lastClass := u.Class()
	if lastClass != c {
		if class.CompareAny(c, class.Bard{}, class.Archer{}, class.Rogue{}, class.Miner{}) {
			u.AddEffects(c.Effects()...)
		} else if class.CompareAny(lastClass, class.Bard{}, class.Archer{}, class.Rogue{}, class.Miner{}) {
			u.bardEnergy.Store(0)
			u.RemoveEffects(lastClass.Effects()...)
		}
		u.class.Store(c)
	}
}

// SendWall ...
func (u *User) SendWall(newPos cube.Pos, z area.Zone, color item.Colour) {
	areaMin := cube.Pos{int(z.Area().Min().X()), 0, int(z.Area().Min().Y())}
	areaMax := cube.Pos{int(z.Area().Max().X()), 255, int(z.Area().Max().Y())}
	wallBlock := block.StainedGlass{Colour: color}
	const wallLength, wallHeight = 15, 10

	if newPos.X() >= areaMin.X() && newPos.X() <= areaMax.X() { // edges of the top and bottom walls (relative to South)
		zCoord := areaMin.Z()
		if newPos.Z() >= areaMax.Z() {
			zCoord = areaMax.Z()
		}
		for horizontal := newPos.X() - wallLength; horizontal < newPos.X()+wallLength; horizontal++ {
			if horizontal >= areaMin.X() && horizontal <= areaMax.X() {
				for vertical := newPos.Y(); vertical < newPos.Y()+wallHeight; vertical++ {
					blockPos := cube.Pos{horizontal, vertical, zCoord}
					if blockReplaceable(u.p.World().Block(blockPos)) {
						u.s.ViewBlockUpdate(blockPos, wallBlock, 0)
						u.wallBlocksMu.Lock()
						u.wallBlocks[blockPos] = rand.Float64() * float64(rand.Intn(1)+1)
						u.wallBlocksMu.Unlock()
					}
				}
			}
		}
	}
	if newPos.Z() >= areaMin.Z() && newPos.Z() <= areaMax.Z() { // edges of the left and right walls (relative to South)
		xCoord := areaMin.X()
		if newPos.X() >= areaMax.X() {
			xCoord = areaMax.X()
		}
		for horizontal := newPos.Z() - wallLength; horizontal < newPos.Z()+wallLength; horizontal++ {
			if horizontal >= areaMin.Z() && horizontal <= areaMax.Z() {
				for vertical := newPos.Y(); vertical < newPos.Y()+wallHeight; vertical++ {
					blockPos := cube.Pos{xCoord, vertical, horizontal}
					if blockReplaceable(u.p.World().Block(blockPos)) {
						u.s.ViewBlockUpdate(blockPos, wallBlock, 0)
						u.wallBlocksMu.Lock()
						u.wallBlocks[blockPos] = rand.Float64() * float64(rand.Intn(1)+1)
						u.wallBlocksMu.Unlock()
					}
				}
			}
		}
	}
}

// ClearWall clears the users walls or lowers the remaining duration.
func (u *User) ClearWall() {
	u.wallBlocksMu.Lock()
	for p, duration := range u.wallBlocks {
		if duration <= 0 {
			delete(u.wallBlocks, p)
			u.s.ViewBlockUpdate(p, block.Air{}, 0)
			u.p.ShowParticle(p.Vec3(), particle.BlockForceField{})
			continue
		}
		u.wallBlocks[p] = duration - 0.1
	}
	u.wallBlocksMu.Unlock()
}

// blockReplaceable checks if the combat wall should replace a block.
func blockReplaceable(b world.Block) bool {
	_, air := b.(block.Air)
	_, doubleFlower := b.(block.DoubleFlower)
	_, flower := b.(block.Flower)
	_, tallGrass := b.(block.TallGrass)
	_, doubleTallGrass := b.(block.DoubleTallGrass)
	_, deadBush := b.(block.DeadBush)
	//_, cobweb := b.(block.Cobweb)
	//_, sapling := b.(block.Sapling)
	_, torch := b.(block.Torch)
	_, fire := b.(block.Fire)
	return air || tallGrass || deadBush || torch || fire || flower || doubleFlower || doubleTallGrass
}

// skinToProtocol converts a skin to its protocol representation.
func skinToProtocol(s skin.Skin) protocol.Skin {
	var animations []protocol.SkinAnimation
	for _, animation := range s.Animations {
		protocolAnim := protocol.SkinAnimation{
			ImageWidth:  uint32(animation.Bounds().Max.X),
			ImageHeight: uint32(animation.Bounds().Max.Y),
			ImageData:   animation.Pix,
			FrameCount:  float32(animation.FrameCount),
		}
		switch animation.Type() {
		case skin.AnimationHead:
			protocolAnim.AnimationType = protocol.SkinAnimationHead
		case skin.AnimationBody32x32:
			protocolAnim.AnimationType = protocol.SkinAnimationBody32x32
		case skin.AnimationBody128x128:
			protocolAnim.AnimationType = protocol.SkinAnimationBody128x128
		}
		protocolAnim.ExpressionType = uint32(animation.AnimationExpression)
		animations = append(animations, protocolAnim)
	}

	return protocol.Skin{
		PlayFabID:          s.PlayFabID,
		SkinID:             uuid.New().String(),
		SkinResourcePatch:  s.ModelConfig.Encode(),
		SkinImageWidth:     uint32(s.Bounds().Max.X),
		SkinImageHeight:    uint32(s.Bounds().Max.Y),
		SkinData:           s.Pix,
		CapeImageWidth:     uint32(s.Cape.Bounds().Max.X),
		CapeImageHeight:    uint32(s.Cape.Bounds().Max.Y),
		CapeData:           s.Cape.Pix,
		SkinGeometry:       s.Model,
		PersonaSkin:        s.Persona,
		CapeID:             uuid.New().String(),
		FullID:             uuid.New().String(),
		Animations:         animations,
		Trusted:            true,
		OverrideAppearance: true,
	}
}

// compareSkin ...
func compareSkin(s1, s2 skin.Skin) bool {
	if len(s1.Pix) != len(s2.Pix) {
		return false
	}
	for i := range s1.Pix {
		if s1.Pix[i] != s2.Pix[i] {
			return false
		}
	}
	return true
}

// ViewPlayerSkin sends the skin of the player to the user.
func (u *User) ViewPlayerSkin(p *player.Player, s skin.Skin) {
	if sk, ok := u.lastPlayerSpecificSkin[p.XUID()]; ok && compareSkin(sk, s) {
		return
	}
	u.lastPlayerSpecificSkin[p.XUID()] = s
	session_writePacket(u.s, &packet.PlayerSkin{
		UUID: p.UUID(),
		Skin: skinToProtocol(s),
	})
}

// BardEnergy returns the bard energy of the user.
func (u *User) BardEnergy() float64 {
	return u.bardEnergy.Load()
}

// DropContents drops the contents of the user.
func (u *User) DropContents() {
	drop_contents(u.p)
}

// Locale returns the locale of the user.
func (u *User) Locale() language.Tag {
	return u.p.Locale()
}

// Class returns the user's class.
func (u *User) Class() util.Class {
	return u.class.Load()
}

// Frozen returns if the user is frozen.
func (u *User) Frozen() bool {
	return u.frozen.Load()
}

// ToggleFreeze toggles the frozen state of the user.
func (u *User) ToggleFreeze() {
	u.frozen.Toggle()
}

// SetMute sets the mute data of the user.
func (u *User) SetMute(p Punishment) {
	u.mute.Store(p)
}

// Mute returns the mute data of the user and true if the data is valid. Otherwise, it will return false.
func (u *User) Mute() (Punishment, bool) {
	p := u.mute.Load()
	if p.Expired() {
		u.SetMute(Punishment{})
		return Punishment{}, false
	}
	return p, true
}

// SetBan sets the ban data of the user.
func (u *User) SetBan(p Punishment) {
	u.ban.Store(p)
}

// Ban returns the ban data of the user, this should only be valid once, when the user gets banned.
func (u *User) Ban() Punishment {
	return u.ban.Load()
}

// RenewReportSince renews the last time the user has made a report.
func (u *User) RenewReportSince() {
	u.reportSince.Store(time.Now())
}

// ReportSince returns the last time the user has made a report.
func (u *User) ReportSince() time.Time {
	return u.reportSince.Load()
}

// DeviceID returns the device ID of the user.
func (u *User) DeviceID() string {
	return u.s.ClientData().DeviceID
}

// SelfSignedID returns the self-signed ID of the user.
func (u *User) SelfSignedID() string {
	return u.s.ClientData().SelfSignedID
}

// HashedAddress returns the hashed IP address of the user.
func (u *User) HashedAddress() string {
	return u.hashedAddress
}

// Rotate rotates the user with the specified yaw and pitch deltas.
func (u *User) Rotate(deltaYaw, deltaPitch float64) {
	rot := u.p.Rotation()
	session_writePacket(u.s, &packet.MovePlayer{
		EntityRuntimeID: 1, // Always 1 on Dragonfly.
		Position:        vec64To32(u.p.Position().Add(mgl64.Vec3{0, 1.62})),
		Pitch:           float32(rot.Pitch() + deltaPitch),
		Yaw:             float32(rot.Yaw() + deltaYaw),
		HeadYaw:         float32(rot.Yaw() + deltaYaw),
		Mode:            packet.MoveModeTeleport,
		OnGround:        u.p.OnGround(),
	})
	u.p.Move(mgl64.Vec3{}, deltaYaw, deltaPitch)
}

// Hide will hide a player from the user
func (u *User) Hide(p *player.Player) {
	u.s.HideEntity(p)
}

// View will un-hide a player from the user
func (u *User) View(p *player.Player) {
	u.s.ViewEntity(p)
}

// EntityRuntimeID returns the entity runtime ID of the user.
func (u *User) EntityRuntimeID(e world.Entity) uint64 {
	return session_entityRuntimeID(u.s, e)
}

// Close closes the user.
func (u *User) Close() {
	for _, k := range koth.All() {
		k.StopCapturing(u)
	}
	u.effectsNoLossMu.Lock()
	for _, e := range u.effectsNoLoss {
		close(e.c)
	}
	u.effectsNoLossMu.Unlock()

	close(u.close)
	userMu.Lock()
	delete(users, u.p)
	userMu.Unlock()
}

// viewers returns a list of all viewers of the Player.
func (u *User) viewers() []world.Viewer {
	viewers := u.p.World().Viewers(u.p.Position())
	for _, v := range viewers {
		if v == u.s {
			return viewers
		}
	}
	return append(viewers, u.s)
}

// vec64To32 converts a mgl64.Vec3 to a mgl32.Vec3.
func vec64To32(vec3 mgl64.Vec3) mgl32.Vec3 {
	return mgl32.Vec3{float32(vec3[0]), float32(vec3[1]), float32(vec3[2])}
}

// noinspection ALL
//
//go:linkname player_session github.com/df-mc/dragonfly/server/player.(*Player).session
func player_session(*player.Player) *session.Session

// noinspection ALL
//
//go:linkname session_writePacket github.com/df-mc/dragonfly/server/session.(*Session).writePacket
func session_writePacket(*session.Session, packet.Packet)

// noinspection ALL
//
//go:linkname drop_contents github.com/df-mc/dragonfly/server/player.(*Player).dropContents
func drop_contents(*player.Player)

// noinspection ALL
//
//go:linkname session_entityRuntimeID github.com/df-mc/dragonfly/server/session.(*Session).entityRuntimeID
func session_entityRuntimeID(*session.Session, world.Entity) uint64
