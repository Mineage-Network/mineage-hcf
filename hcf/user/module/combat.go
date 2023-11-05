package module

import (
	"github.com/df-mc/dragonfly/server/entity"
	"github.com/df-mc/dragonfly/server/entity/effect"
	"github.com/df-mc/dragonfly/server/event"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/chat"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/sound"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/mineage-network/mineage-hcf/hcf/backend/lang"
	ench "github.com/mineage-network/mineage-hcf/hcf/enchantment"
	"github.com/mineage-network/mineage-hcf/hcf/knockback"
	"github.com/mineage-network/mineage-hcf/hcf/sotw"
	"github.com/mineage-network/mineage-hcf/hcf/user"
	"github.com/mineage-network/mineage-hcf/hcf/util"
	"github.com/mineage-network/mineage-hcf/hcf/util/class"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"golang.org/x/text/language"
	"math"
	"time"
)

// Combat ...
type Combat struct {
	u *user.User
	player.NopHandler
}

// NewCombat ...
func NewCombat(u *user.User) *Combat {
	return &Combat{u: u}
}

// HandleHurt ...
func (c *Combat) HandleHurt(ctx *event.Context, dmg *float64, attackImmunity *time.Duration, src world.DamageSource) {
	if ctx.Cancelled() {
		return
	}

	u := c.u
	p := u.Player()
	w := p.World()
	pos := p.Position()

	if _, ok := sotw.Running(); ok && u.SOTW() {
		ctx.Cancel()
		return
	}

	if _, ok := u.TimerEnabled(); ok {
		ctx.Cancel()
		return
	}

	*attackImmunity = knockback.Knockback.RealHitDelay()

	switch s := src.(type) {
	case effect.PoisonDamageSource:
		*attackImmunity = 0
		if p.Health()-p.FinalDamageFrom(*dmg, src) <= 0 {
			ctx.Cancel()
		}
		if p.AttackImmune() && p.AttackImmunity() <= 250 {
			ctx.Cancel()
		}
		return
	case entity.AttackDamageSource:
		t, ok := user.LookupEntity(s.Attacker)
		if !ok || !u.CanAttack(t) {
			ctx.Cancel()
			return
		}

		if u.Tags().Archer().Active() {
			*dmg = *dmg * 1.15
		}
		u.SetLastAttacker(t)
	case NoArmourAttackEntitySource:
		t, ok := user.LookupEntity(s.Attacker)
		if !ok || !u.CanAttack(t) {
			ctx.Cancel()
			return
		}
		u.SetLastAttacker(t)
	case entity.ProjectileDamageSource:
		if _, ok := s.Projectile.Type().(entity.FireworkType); ok {
			ctx.Cancel()
			return
		}

		t, ok := user.LookupEntity(s.Owner)
		if !ok || !u.CanAttack(t) {
			ctx.Cancel()
			return
		}

		u.Tags().Combat().Set(time.Second * 20)
		t.Tags().Combat().Set(time.Second * 20)

		u.SetLastAttacker(t)
		switch s.Projectile.Type().(type) {
		/*case ent.SwitcherBallType:
		if k, ok := koth.Running(); ok {
			if pl, ok := k.Capturing(); ok && pl == u {
				t.Player().Message(text.Colourf("<red>You cannot switch places with someone capturing a koth</red>"))
				break
			}
		}
		dist := p.Position().Sub(t.Position()).Len()
		if dist > 7 {
			t.Player().Message(text.Colourf("<red>You are too far away from %s</red>", p.Name()))
			break
		}

		ctx.Cancel()
		targetPos := t.Position()
		pos := p.Position()

		t.Player().Teleport(pos)
		p.Teleport(targetPos)*/
		case entity.ArrowType:
			arm := t.Player().Armour()
			for _, a := range arm.Slots() {
				for _, e := range a.Enchantments() {
					if att, ok := e.Type().(ench.AttackEnchantment); ok {
						att.AttackEntity(t.Player(), p)
					}
				}
			}
			if class.Compare(t.Class(), class.Archer{}) {
				if !class.Compare(u.Class(), class.Archer{}) {
					u.Tags().Archer().Set(time.Second * 10)
					u.UpdateState()
				}
				dist := p.Position().Sub(t.Position()).Len()
				d := math.Round(dist)
				if d > 20 {
					d = 20
				}
				dmg := d / 10
				t.Message("archer.tagged.target", math.Round(dist), dmg)
				p.Hurt(dmg*2, NoArmourAttackEntitySource{
					Attacker: t.Player(),
				})
				p.KnockBack(p.Position(), 0.38, 0.38)
			}
		}
	}
	if (p.Health()-p.FinalDamageFrom(*dmg, src) <= 0 || (src == entity.VoidDamageSource{})) && !ctx.Cancelled() {
		ctx.Cancel()

		w.PlaySound(pos, sound.Thunder{})

		npc := player.New(p.Name(), p.Skin(), pos)
		npc.Handle(npcHandler{})
		npc.SetAttackImmunity(time.Millisecond * 1400)
		npc.SetNameTag(p.NameTag())
		npc.SetScale(p.Scale())
		w.AddEntity(npc)

		for _, viewer := range w.Viewers(npc.Position()) {
			viewer.ViewEntityAction(npc, entity.DeathAction{})
		}
		time.AfterFunc(time.Second*2, func() {
			_ = npc.Close()
		})

		if att, ok := attackerFromSource(src); ok {
			npc.KnockBack(att.Position(), 0.5, 0.2)
		}

		for _, e := range p.Effects() {
			p.RemoveEffect(e.Type())
		}
		for _, et := range u.World().Entities() {
			if be, ok := et.(entity.Behaviour); ok {
				if pro, ok := be.(*entity.ProjectileBehaviour); ok {
					if pro.Owner() == p {
						u.World().RemoveEntity(et)
					}
				}
			}
		}
		for _, cd := range u.Cooldowns().Resetable() {
			cd.Reset()
		}
		for _, tag := range u.Tags().All() {
			tag.Reset()
		}

		u.DropContents()
		p.SetHeldItems(item.Stack{}, item.Stack{})
		p.ResetFallDistance()

		if fa, ok := u.Faction(); ok {
			fa.SetDTR(fa.DTR() - 1)
			fa.RemovePoints(1)
			fa.SetRegenerationTime(time.Now().Add(time.Minute * 15))
		}

		killer, ok := u.LastAttacker()
		if ok {
			killer.Stats().AddKill()

			if fa, ok := killer.Faction(); ok {
				fa.AddPoints(1)
			}

			held, _ := killer.HeldItems()
			heldName := held.CustomName()

			if len(heldName) <= 0 {
				heldName = util.ItemName(held.Item())
			}

			if held.Empty() || len(heldName) <= 0 {
				heldName = "their fist"
			}

			_, _ = chat.Global.WriteString(lang.Translatef(language.English, "user.kill", p.Name(), u.Stats().Kills(), killer.Name(), killer.Stats().Kills(), text.Colourf("<red>%s</red>", heldName)))
			u.ResetLastAttacker()
		} else {
			_, _ = chat.Global.WriteString(lang.Translatef(language.English, "user.suicide", p.Name(), u.Stats().Kills()))
		}

		u.ReduceLife()
		if u.DeathBanned() {
			u.HandleDeathBan()
			return
		}

		p.Heal(20, effect.InstantHealingSource{})
		p.Teleport(mgl64.Vec3{0, 80, 0})
		p.Extinguish()
		p.SetFood(20)
		u.SetClass(class.Resolve(p))
		u.EnableTimer()
		u.UpdateState()
	}
}
