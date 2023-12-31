package util

import (
	"github.com/df-mc/dragonfly/server/entity/effect"
	"github.com/df-mc/dragonfly/server/item"
)

// Class represents a class.
type Class interface {
	Armour() [4]item.ArmourTier
	Effects() []effect.Effect
}
