package class

import (
	"github.com/df-mc/dragonfly/server/entity/effect"
	"github.com/df-mc/dragonfly/server/item"
	"time"
)

// Miner ...
type Miner struct{}

// Armour ...
func (Miner) Armour() [4]item.ArmourTier {
	return [4]item.ArmourTier{
		item.ArmourTierIron{},
		item.ArmourTierIron{},
		item.ArmourTierIron{},
		item.ArmourTierIron{},
	}
}

// Effects ...
func (Miner) Effects() []effect.Effect {
	return []effect.Effect{
		effect.New(effect.Haste{}, 2, time.Hour*999),
		effect.New(effect.NightVision{}, 2, time.Hour*999),
		effect.New(effect.FireResistance{}, 2, time.Hour*999),
	}
}
