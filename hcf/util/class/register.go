package class

import (
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/mineage-network/mineage-hcf/hcf/util"
)

var (
	// classes ...
	classes []util.Class
)

// ResolveFromArmour resolves a class from a set of armour.
func ResolveFromArmour(a [4]item.ArmourTier) util.Class {
	for _, c := range classes {
		if a[0] == c.Armour()[0] && a[1] == c.Armour()[1] && a[2] == c.Armour()[2] && a[3] == c.Armour()[3] {
			return c
		}
	}
	return nil
}

// Resolve resolves a class from a player.
func Resolve(p *player.Player) util.Class {
	a := p.Armour()
	helmet, ok := a.Helmet().Item().(item.Helmet)
	if !ok {
		return nil
	}
	chestplate, ok := a.Chestplate().Item().(item.Chestplate)
	if !ok {
		return nil
	}
	leggings, ok := a.Leggings().Item().(item.Leggings)
	if !ok {
		return nil
	}
	boots, ok := a.Boots().Item().(item.Boots)
	if !ok {
		return nil
	}
	return ResolveFromArmour([4]item.ArmourTier{helmet.Tier, chestplate.Tier, leggings.Tier, boots.Tier})
}

// Compare compares two classes.
func Compare(a any, b util.Class) bool {
	if a == b {
		return true
	}
	return false
}

// CompareAny compares a class to multiple classes.
func CompareAny(a any, b ...util.Class) bool {
	for _, c := range b {
		if Compare(a, c) {
			return true
		}
	}
	return false
}

// Register registers a class.
func Register(c util.Class) {
	classes = append(classes, c)
}

// init ...
func init() {
	Register(Archer{})
	Register(Bard{})
	Register(Miner{})
	Register(Rogue{})
}
