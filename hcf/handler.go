package hcf

import (
	"fmt"
	"github.com/bedrock-gophers/packethandler"
	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/entity"
	"github.com/df-mc/dragonfly/server/entity/effect"
	"github.com/df-mc/dragonfly/server/event"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/chat"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/particle"
	"github.com/df-mc/dragonfly/server/world/sound"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/mineage-network/mineage-hcf/hcf/area"
	"github.com/mineage-network/mineage-hcf/hcf/data"
	"github.com/mineage-network/mineage-hcf/hcf/knockback"
	"github.com/mineage-network/mineage-hcf/hcf/logger"
	"github.com/mineage-network/mineage-hcf/hcf/rank"
	"github.com/mineage-network/mineage-hcf/hcf/sotw"
	"github.com/mineage-network/mineage-hcf/hcf/user"
	"github.com/mineage-network/mineage-hcf/hcf/user/module"
	"github.com/mineage-network/mineage-hcf/hcf/util"
	"github.com/mineage-network/mineage-hcf/hcf/util/class"
	"github.com/pzurek/durafmt"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"log"
	"net/netip"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"
)

// packetHandler ...
type packetHandler struct {
	packethandler.NopHandler
	c *packethandler.Conn
}

// removeFlag removes a flag from the entity data.
func removeFlag(key uint32, index uint8, m protocol.EntityMetadata) {
	v := m[key]
	switch key {
	case protocol.EntityDataKeyPlayerFlags:
		m[key] = v.(byte) &^ (1 << index)
	default:
		m[key] = v.(int64) &^ (1 << int64(index))
	}
}

// HandleServerPacket ...
func (h packetHandler) HandleServerPacket(_ *event.Context, pk packet.Packet) {
	if pkt, ok := pk.(*packet.SetActorData); ok {
		u, ok := user.LookupName(h.c.IdentityData().DisplayName)
		if !ok {
			return
		}
		meta := protocol.EntityMetadata(pkt.EntityMetadata)

		for _, usr := range user.All() {
			if u.EntityRuntimeID(usr.Player()) == pkt.EntityRuntimeID {
				if meta.Flag(protocol.EntityDataKeyFlags, protocol.EntityDataFlagInvisible) {
					if usr.Tags().Archer().Active() {
						removeFlag(protocol.EntityDataKeyFlags, protocol.EntityDataFlagInvisible, meta)
						u.ViewPlayerSkin(usr.Player(), usr.Player().Skin())
						continue
					}
				}
				if _, ok := sotw.Running(); ok && usr.SOTW() {
					meta[protocol.EntityDataKeyName] = text.Colourf("<grey>%s</grey>", usr.Name())
				} else if _, ok := usr.TimerEnabled(); ok {
					meta[protocol.EntityDataKeyName] = text.Colourf("<grey>%s</grey>", usr.Name())
				}

				if fa, ok := usr.Faction(); ok {
					if uFaction, ok := u.Faction(); ok && uFaction.Compare(fa) {
						if meta.Flag(protocol.EntityDataKeyFlags, protocol.EntityDataFlagInvisible) {
							removeFlag(protocol.EntityDataKeyFlags, protocol.EntityDataFlagInvisible, meta)
							dimmedSkin := usr.Player().Skin()
							dimmedSkin.Pix = make([]byte, len(dimmedSkin.Pix))
							u.ViewPlayerSkin(usr.Player(), dimmedSkin)
						} else {
							u.ViewPlayerSkin(usr.Player(), usr.Player().Skin())
						}
						meta[protocol.EntityDataKeyName] = text.Colourf("<green>%s</green>", usr.Name())
					}
				}
			}
		}

		if users, ok := u.Focusing(); ok {
			for _, m := range users {
				if pkt.EntityRuntimeID == u.EntityRuntimeID(m.Player()) {
					meta[protocol.EntityDataKeyName] = text.Colourf("<purple>%s</purple>", m.Name())
				}
			}
		}
		pkt.EntityMetadata = meta
	}
}

// handler ...
type handler struct {
	player.NopHandler

	area   *module.Area
	class  *module.Class
	custom *module.Custom
	combat *module.Combat

	sign cube.Pos

	u   *user.User
	p   *player.Player
	hcf *HCF
}

// HandleJoin ...
func (h *handler) HandleJoin() {
	defer h.u.UpdateState()

	p := h.p

	if h.u.DeathBanned() {
		h.u.HandleDeathBan()
		return
	}

	if l, ok := logger.LookupXUID(p.XUID()); ok {
		p.Teleport(l.Player().Position())
		l.Reconnect()
	}

	for _, u := range user.All() {
		if u.Vanished() {
			h.u.Player().HideEntity(u.Player())
		}
	}

	// TODO: remove
	h.u.AddWaypoint("home", mgl64.Vec3{0, 66, 0})
}

// formatRegex is a regex used to clean color formatting on a string.
var formatRegex = regexp.MustCompile(`§[\da-gk-or]`)

// HandleChat ...
func (h *handler) HandleChat(ctx *event.Context, message *string) {
	ctx.Cancel()
	r := h.u.Ranks().Highest()

	if !h.u.CanSendMessage() {
		h.u.Message("user.message.cooldown")
		return
	}

	if _, ok := h.u.Mute(); ok {
		h.u.Message("user.message.mute")
		return
	}

	if msg := strings.TrimSpace(*message); len(msg) > 0 {
		msg = formatRegex.ReplaceAllString(msg, "")

		global := func() {
			if fa, ok := h.u.Faction(); ok {
				formatFaction := text.Colourf("<grey>[<green>%s</green>]</grey> %s", fa.DisplayName(), r.Chat(h.p.Name(), msg))
				formatEnemy := text.Colourf("<grey>[<red>%s</red>]</grey> %s", fa.DisplayName(), r.Chat(h.p.Name(), msg))
				for _, u := range user.All() {
					if userFaction, ok := u.Faction(); ok && userFaction == fa {
						u.Player().Message(formatFaction)
					} else {
						u.Player().Message(formatEnemy)
					}
				}
				chat.StdoutSubscriber{}.Message(formatEnemy)
			} else {
				_, _ = chat.Global.WriteString(r.Chat(h.p.Name(), msg))
			}
		}

		switch h.u.ChatType() {
		case user.ChatTypeGlobal():
			global()
		case user.ChatTypeFaction():
			fa, ok := h.u.Faction()
			if !ok {
				h.u.UpdateChatType(user.ChatTypeGlobal())
				global()
				return
			}
			for _, u := range fa.Users() {
				u.Player().Message(text.Colourf("<dark-aqua>%s: %s</dark-aqua>", h.p.Name(), msg))
			}
		}
		h.u.RenewLastMessage()
	}
}

// HandleStartBreak ...
func (h *handler) HandleStartBreak(ctx *event.Context, pos cube.Pos) {
	h.custom.HandleStartBreak(ctx, pos)
}

// formatItemName ...
func formatItemName(s string) string {
	split := strings.Split(s, "_")
	for i, str := range split {
		upperCasesPrefix := unicode.ToUpper(rune(str[0]))
		split[i] = string(upperCasesPrefix) + str[1:]
	}
	return strings.Join(split, " ")
}

// HandleSignEdit ...
func (h *handler) HandleSignEdit(ctx *event.Context, frontSide bool, _, newText string) {
	ctx.Cancel()
	if !frontSide {
		return
	}

	lines := strings.Split(newText, "\n")
	if len(lines) <= 0 {
		return
	}

	switch strings.ToLower(lines[0]) {
	case "[elevator]":
		if len(lines) < 2 {
			return
		}
		var newLines []string

		newLines = append(newLines, text.Colourf("<dark-red>[Elevator]</dark-red>"))
		switch strings.ToLower(lines[1]) {
		case "up":
			newLines = append(newLines, text.Colourf("Up"))
		case "down":
			newLines = append(newLines, text.Colourf("Down"))
		default:
			return
		}
		time.AfterFunc(time.Millisecond, func() {
			b := h.p.World().Block(h.sign)
			if s, ok := b.(block.Sign); ok {
				h.p.World().SetBlock(h.sign, block.Sign{Wood: s.Wood, Attach: s.Attach, Waxed: s.Waxed, Front: block.SignText{
					Text:       strings.Join(newLines, "\n"),
					BaseColour: s.Front.BaseColour,
					Glowing:    false,
				}, Back: s.Back}, nil)
			}
		})
	case "[shop]":
		if len(lines) < 4 {
			return
		}

		if !h.u.Ranks().Contains(rank.Admin{}) {
			h.u.World().SetBlock(h.sign, block.Air{}, nil)
			return
		}

		var newLines []string
		spl := strings.Split(lines[1], " ")
		choice := strings.ToLower(spl[0])
		q, _ := strconv.Atoi(spl[1])
		price, _ := strconv.Atoi(lines[3])
		switch choice {
		case "buy":
			newLines = append(newLines, text.Colourf("<blue>- Buy -</blue>"))
		case "sell":
			newLines = append(newLines, text.Colourf("<red>- Sell -</red>"))
		}

		newLines = append(newLines, formatItemName(lines[2]))
		newLines = append(newLines, fmt.Sprint(q))
		newLines = append(newLines, fmt.Sprintf("$%d", price))

		time.AfterFunc(time.Millisecond, func() {
			b := h.p.World().Block(h.sign)
			if s, ok := b.(block.Sign); ok {
				h.p.World().SetBlock(h.sign, block.Sign{Wood: s.Wood, Attach: s.Attach, Waxed: s.Waxed, Front: block.SignText{
					Text:       strings.Join(newLines, "\n"),
					BaseColour: s.Front.BaseColour,
					Glowing:    false,
				}, Back: s.Back}, nil)
			}
		})
	}
}

// HandleItemDrop ...
func (h *handler) HandleItemDrop(ctx *event.Context, _ world.Entity) {
	if h.u.Logging() {
		ctx.Cancel()
		return
	}
}

// HandleItemPickup ...
func (h *handler) HandleItemPickup(ctx *event.Context, st *item.Stack) {
	if h.u.Logging() {
		ctx.Cancel()
		return
	}
	h.custom.HandleItemPickup(ctx, st)
}

// HandleHurt ...
func (h *handler) HandleHurt(ctx *event.Context, dmg *float64, attackImmunity *time.Duration, src world.DamageSource) {
	if h.u.Logging() {
		ctx.Cancel()
		return
	}
	if s, ok := src.(entity.AttackDamageSource); ok {
		if pl, ok := s.Attacker.(*player.Player); ok {
			h, _ := pl.HeldItems()
			if _, ok := h.Item().(item.Sword); ok {
				*dmg = *dmg / 1.25
			}
		}
	}
	h.area.HandleHurt(ctx, dmg, attackImmunity, src)
	h.combat.HandleHurt(ctx, dmg, attackImmunity, src)
}

// HandleAttackEntity ...
func (h *handler) HandleAttackEntity(ctx *event.Context, e world.Entity, force, height *float64, critical *bool) {
	u := h.u
	w := u.World()

	if h.u.Logging() {
		ctx.Cancel()
		return
	}

	*force, *height = knockback.Knockback.RealForce(), knockback.Knockback.RealHeight()

	t, ok := e.(*player.Player)
	if !ok || t.AttackImmune() {
		return
	}

	target, ok := user.LookupEntity(t)
	if ok {
		if !h.u.CanAttack(target) {
			ctx.Cancel()
			return
		}

		h.u.Tags().Combat().Set(time.Second * 20)
		target.Tags().Combat().Set(time.Second * 20)
	} else if l, ok := logger.LookupEntity(e); ok {
		h.u.Tags().Combat().Set(time.Second * 20)
		t = l.Player()
	}

	held, left := h.u.HeldItems()

	if s, ok := held.Item().(item.Sword); ok && s.Tier == item.ToolTierGold && class.Compare(u.Class(), class.Rogue{}) && t.Rotation().Direction() == u.Rotation().Direction() {
		cd := u.Cooldowns().RogueAbility()
		if cd.Active() {
			u.Message("user.ability.cooldown", cd.Remaining().Seconds())
		} else {
			for i := 1; i <= 3; i++ {
				w.AddParticle(t.Position().Add(mgl64.Vec3{0, float64(i), 0}), particle.Dust{
					Colour: item.ColourRed().RGBA(),
				})
			}
			w.PlaySound(u.Position(), sound.ItemBreak{})
			t.Hurt(8, module.NoArmourAttackEntitySource{
				Attacker: u.Player(),
			})
			t.KnockBack(u.Position(), *force, *height)
			u.SetHeldItems(item.Stack{}, left)
			cd.Set(time.Second * 10)
		}
	}

	h.custom.HandleAttackEntity(ctx, e, force, height, critical)
}

// HandleTeleport ...
func (h *handler) HandleTeleport(ctx *event.Context, pos mgl64.Vec3) {
	if h.u.Logging() {
		ctx.Cancel()
		return
	}
	h.area.HandleTeleport(ctx, pos)
}

// HandleItemConsume ...
func (h *handler) HandleItemConsume(ctx *event.Context, i item.Stack) {
	if _, ok := i.Item().(item.GoldenApple); ok {
		cd := h.u.Cooldowns().GoldenApple()
		if cd.Active() {
			ctx.Cancel()
			h.p.Message(text.Colourf("<red>You can consume this item in %.1f seconds</red>", cd.Remaining().Seconds()))
			return
		} else {
			cd.Set(time.Second * 30)
		}
	}
	if _, ok := i.Item().(item.EnchantedApple); ok {
		cd := h.u.Cooldowns().NotchApple()
		if cd.Active() {
			ctx.Cancel()
			h.p.Message(text.Colourf("<red>You can consume this item in %s</red>", durafmt.HMSWithSeparator(cd.Remaining(), ":")))
			return
		} else {
			cd.Set(time.Hour * 4)
		}
	}
	if _, ok := i.Item().(item.Bucket); ok {
		ctx.Cancel()
		h.p.ReleaseItem()
		_, leftHeld := h.p.HeldItems()
		h.p.SetHeldItems(item.NewStack(item.Bucket{}, 1), leftHeld)
		for _, e := range h.p.Effects() {
			switch e.Type() {
			case effect.Saturation{}:
				h.p.RemoveEffect(e.Type())
			case effect.Hunger{}:
				h.p.RemoveEffect(e.Type())
			case effect.Blindness{}:
				h.p.RemoveEffect(e.Type())
			case effect.Nausea{}:
				h.p.RemoveEffect(e.Type())
			case effect.Weakness{}:
				h.p.RemoveEffect(e.Type())
			case effect.Poison{}:
				h.p.RemoveEffect(e.Type())
			case effect.Wither{}:
				h.p.RemoveEffect(e.Type())
			case effect.Slowness{}:
				h.p.RemoveEffect(e.Type())
			}
		}
	}
}

// HandleFoodLoss ...
func (*handler) HandleFoodLoss(ctx *event.Context, _ int, _ *int) {
	ctx.Cancel()
}

// HandleItemDamage ...
func (h *handler) HandleItemDamage(_ *event.Context, i item.Stack, n int) {
	dur := i.Durability()
	if _, ok := i.Item().(item.Armour); ok && dur != -1 && dur-n <= 0 {
		h.u.SetClass(class.Resolve(h.p))
	}
}

// HandleMove ...
func (h *handler) HandleMove(ctx *event.Context, newPos mgl64.Vec3, yaw, pitch float64) {
	if h.u.Logging() {
		ctx.Cancel()
		return
	}
	h.area.HandleMove(ctx, newPos, yaw, pitch)

	for _, w := range h.u.Waypoints() {
		w.Update(h.p)
	}
}

// HandleBlockPlace ...
func (h *handler) HandleBlockPlace(ctx *event.Context, pos cube.Pos, b world.Block) {
	if h.u.Logging() {
		ctx.Cancel()
		return
	}
	h.area.HandleBlockPlace(ctx, pos, b)
	h.custom.HandleBlockPlace(ctx, pos, b)

	if !ctx.Cancelled() {
		if _, ok := b.(block.Sign); ok {
			h.sign = pos
		}
	}
}

// HandleBlockBreak ...
func (h *handler) HandleBlockBreak(ctx *event.Context, pos cube.Pos, drops *[]item.Stack, xp *int) {
	if h.u.Logging() {
		ctx.Cancel()
		return
	}
	h.custom.HandleBlockBreak(ctx, pos, drops, xp)
	h.area.HandleBlockBreak(ctx, pos, drops, xp)
}

// HandleItemUse ...
func (h *handler) HandleItemUse(ctx *event.Context) {
	h.class.HandleItemUse(ctx)
	h.custom.HandleItemUse(ctx)
}

// HandlePunchAir ...
func (h *handler) HandlePunchAir(ctx *event.Context) {
	h.custom.HandlePunchAir(ctx)
	h.area.HandlePunchAir(ctx)
}

// HandleItemUseOnBlock ...
func (h *handler) HandleItemUseOnBlock(ctx *event.Context, pos cube.Pos, face cube.Face, clickPos mgl64.Vec3) {
	if h.u.Logging() {
		ctx.Cancel()
		return
	}

	h.area.HandleItemUseOnBlock(ctx, pos, face, clickPos)
	h.custom.HandleItemUseOnBlock(ctx, pos, face, clickPos)

	if s, ok := h.u.World().Block(pos).(block.Sign); ok {
		ctx.Cancel()
		cd := h.u.Cooldowns().ItemUse()
		if cd.Active(block.Sign{}) {
			return
		} else {
			cd.Set(block.Sign{}, time.Second/4)
		}

		lines := strings.Split(s.Front.Text, "\n")
		if len(lines) < 2 {
			return
		}
		if strings.Contains(strings.ToLower(lines[0]), "[elevator]") {
			blockFound := false
			if strings.Contains(strings.ToLower(lines[1]), "up") {
				for y := pos.Y() + 1; y < 256; y++ {
					if _, ok := h.u.World().Block(cube.Pos{pos.X(), y, pos.Z()}).(block.Air); ok {
						if !blockFound {
							h.p.Message(text.Colourf("<red>There is no block above the sign</red>"))
							return
						}
						if _, ok := h.u.World().Block(cube.Pos{pos.X(), y + 1, pos.Z()}).(block.Air); !ok {
							h.p.Message(text.Colourf("<red>There is not enough space to teleport up</red>"))
							return
						}
						h.p.Teleport(pos.Vec3Middle().Add(mgl64.Vec3{0, float64(y - pos.Y()), 0}))
						break
					} else {
						blockFound = true
					}
				}
			} else if strings.Contains(strings.ToLower(lines[1]), "down") {
				for y := pos.Y() - 1; y > 0; y-- {
					if _, ok := h.u.World().Block(cube.Pos{pos.X(), y, pos.Z()}).(block.Air); ok {
						if !blockFound {
							h.p.Message(text.Colourf("<red>There is not enough space to teleport down</red>"))
							return
						}
						if _, ok := h.u.World().Block(cube.Pos{pos.X(), y - 1, pos.Z()}).(block.Air); !ok {
							h.p.Message(text.Colourf("<red>There is not enough space to teleport down</red>"))
							return
						}
						h.p.Teleport(pos.Vec3Middle().Add(mgl64.Vec3{0, float64(y - pos.Y() - 1), 0}))
						break
					} else {
						blockFound = true
					}
				}
			}
		}

		title := strings.ToLower(util.StripMinecraftColour(lines[0]))
		if strings.Contains(title, "buy") ||
			strings.Contains(title, "sell") &&
				(area.Spawn(h.u.World()).Area().Vec3WithinOrEqualFloorXZ(h.u.Position())) {
			it, ok := world.ItemByName("minecraft:"+strings.ReplaceAll(strings.ToLower(lines[1]), " ", "_"), 0)
			if !ok {
				h.hcf.log.Error("shop: invalid item")
				return
			}

			q, err := strconv.Atoi(lines[2])
			if err != nil {
				h.hcf.log.Error("shop: invalid quantity")
				return
			}

			price, err := strconv.ParseFloat(strings.Trim(lines[3], "$"), 64)
			if err != nil {
				h.hcf.log.Error("shop: invalid price")
				return
			}

			choice := strings.ReplaceAll(title, " ", "")
			choice = strings.ReplaceAll(choice, "-", "")

			switch choice {
			case "buy":
				if h.u.Balance() < price {
					h.u.Message("shop.balance.insufficient")
					return
				}
				if !ok {
					h.hcf.log.Error("shop: invalid block to item conversion")
					return
				}
				h.u.ReduceBalance(price)
				h.u.AddItemOrDrop(item.NewStack(it, q))
				h.u.Message("shop.buy.success", q, lines[1])
			case "sell":
				inv := h.u.Player().Inventory()
				count := 0
				var items []item.Stack
				for _, slotItem := range inv.Slots() {
					n1, _ := it.EncodeItem()
					if slotItem.Empty() {
						continue
					}
					n2, _ := slotItem.Item().EncodeItem()
					if n1 == n2 {
						count += slotItem.Count()
						items = append(items, slotItem)
					}
				}
				if count >= q {
					h.u.IncreaseBalance(float64(count/q) * price)
					h.u.Message("shop.sell.success", count, lines[1])
				} else {
					h.u.Message("shop.sell.fail")
					return
				}
				for i, v := range items {
					if i >= count {
						break
					}
					amt := count - (count % q)
					if amt > 64 {
						amt = 64
					}
					err := inv.RemoveItemFunc(amt, func(stack item.Stack) bool {
						return stack.Equal(v)
					})
					if err != nil {
						// log.Fatal(err)
					}
				}
			}
		}
	}
}

// HandleQuit ...
func (h *handler) HandleQuit() {
	addr, _ := netip.ParseAddrPort(h.u.Address().String())
	ip := addr.Addr()

	connectionsMu.Lock()
	if connections[ip] <= 1 {
		delete(connections, ip)
	} else {
		connections[ip]--
	}
	connectionsMu.Unlock()

	err := data.SaveUser(h.u)
	if err != nil {
		log.Println(err)
	}
	if !h.u.Logged() && !area.Spawn(h.u.World()).Area().Vec3WithinOrEqualFloorXZ(h.u.Position()) {
		logger.NewLogger(h.p)
	}
	h.u.Close()
}
