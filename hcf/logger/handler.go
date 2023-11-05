package logger

import (
	"github.com/df-mc/dragonfly/server/entity"
	"github.com/df-mc/dragonfly/server/event"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/sound"
	"github.com/mineage-network/mineage-hcf/hcf/user/module"
	"time"
)

// handler ...
type handler struct {
	player.NopHandler
	l *Logger
}

// HandleHurt ...
func (*handler) HandleHurt(_ *event.Context, dmg *float64, attackImmunity *time.Duration, _ world.DamageSource) {
	*dmg = *dmg / 1.25
	*attackImmunity = time.Millisecond * 470
}

// npcHandler ...
type npcHandler struct {
	player.NopHandler
}

// HandleItemPickup ...
func (npcHandler) HandleItemPickup(ctx *event.Context, _ *item.Stack) {
	ctx.Cancel()
}

// HandleDeath ...
func (l *handler) HandleDeath(src world.DamageSource, _ *bool) {
	p := l.l.p
	w := p.World()
	pos := p.Position()

	p.World().PlaySound(pos, sound.Explosion{})

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

	l.l.Kill()
}

// attackerFromSource ...
func attackerFromSource(src world.DamageSource) (world.Entity, bool) {
	switch s := src.(type) {
	case entity.AttackDamageSource:
		return s.Attacker, true
	case module.NoArmourAttackEntitySource:
		return s.Attacker, true
	}
	return nil, false
}
