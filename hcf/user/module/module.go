package module

import (
	"github.com/df-mc/dragonfly/server/entity"
	"github.com/df-mc/dragonfly/server/event"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
)

// npcHandler ...
type npcHandler struct {
	player.NopHandler
}

// HandleItemPickup ...
func (npcHandler) HandleItemPickup(ctx *event.Context, _ *item.Stack) {
	ctx.Cancel()
}

// NoArmourAttackEntitySource ...
type NoArmourAttackEntitySource struct {
	Attacker world.Entity
}

// Fire ...
func (NoArmourAttackEntitySource) Fire() bool {
	return false
}

// ReducedByArmour ...
func (NoArmourAttackEntitySource) ReducedByArmour() bool {
	return false
}

// ReducedByResistance ...
func (NoArmourAttackEntitySource) ReducedByResistance() bool {
	return false
}

// attackerFromSource returns the Attacker from a DamageSource. If the source is not an entity false is
// returned.
func attackerFromSource(src world.DamageSource) (world.Entity, bool) {
	switch s := src.(type) {
	case entity.AttackDamageSource:
		return s.Attacker, true
	case NoArmourAttackEntitySource:
		return s.Attacker, true
	}
	return nil, false
}
