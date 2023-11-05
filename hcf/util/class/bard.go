package class

import (
	"github.com/df-mc/dragonfly/server/entity/effect"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
	"time"
)

// Bard ...
type Bard struct{}

// Armour ...
func (Bard) Armour() [4]item.ArmourTier {
	return [4]item.ArmourTier{
		item.ArmourTierGold{},
		item.ArmourTierGold{},
		item.ArmourTierGold{},
		item.ArmourTierGold{},
	}
}

// Effects ...
func (Bard) Effects() []effect.Effect {
	return []effect.Effect{
		effect.New(effect.Speed{}, 2, time.Hour*999),
		effect.New(effect.Regeneration{}, 1, time.Hour*999),
		effect.New(effect.Resistance{}, 2, time.Hour*999),
	}
}

// BardEffectDuration ...
var bardEffectDuration = time.Second * 6

// bardItemsUse ...
var bardItemsUse = map[world.Item]effect.Effect{
	item.BlazePowder{}: effect.New(effect.Strength{}, 2, bardEffectDuration),
	item.Feather{}:     effect.New(effect.JumpBoost{}, 4, bardEffectDuration),
	item.Sugar{}:       effect.New(effect.Speed{}, 3, bardEffectDuration),
	item.GhastTear{}:   effect.New(effect.Regeneration{}, 3, bardEffectDuration),
	item.IronIngot{}:   effect.New(effect.Resistance{}, 3, bardEffectDuration),
}

// bardItemsHold ...
var bardItemsHold = map[world.Item]effect.Effect{
	item.MagmaCream{}:  effect.New(effect.FireResistance{}, 1, bardEffectDuration),
	item.BlazePowder{}: effect.New(effect.Strength{}, 1, bardEffectDuration),
	item.Feather{}:     effect.New(effect.JumpBoost{}, 2, bardEffectDuration),
	item.Sugar{}:       effect.New(effect.Speed{}, 2, bardEffectDuration),
	item.GhastTear{}:   effect.New(effect.Regeneration{}, 1, bardEffectDuration),
	item.IronIngot{}:   effect.New(effect.Resistance{}, 1, bardEffectDuration),
}

// BardEffectFromItem ...
func BardEffectFromItem(i world.Item) (effect.Effect, bool) {
	e, ok := bardItemsUse[i]
	return e, ok
}

// BardHoldEffectFromItem ...
func BardHoldEffectFromItem(i world.Item) (effect.Effect, bool) {
	e, ok := bardItemsHold[i]
	return e, ok
}
