package enchantment

import (
	"github.com/df-mc/dragonfly/server/entity/effect"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/mineage-network/mineage-hcf/hcf/util"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

// Level ...
var Level = 2

// init ...
func init() {
	item.RegisterEnchantment(0, Protection{})
	item.RegisterEnchantment(9, Sharpness{})
	item.RegisterEnchantment(99, NightVision{})
	item.RegisterEnchantment(100, Speed{})
	item.RegisterEnchantment(101, FireResistance{})
	item.RegisterEnchantment(102, Invisibility{})
}

// AddEnchantmentLore ...
func AddEnchantmentLore(i item.Stack) item.Stack {
	var lores []string

	for _, e := range i.Enchantments() {
		typ := e.Type()

		var lvl string
		if e.Level() > 1 {
			lvl, _ = util.Itor(e.Level())
		}
		switch typ.(type) {
		case EffectEnchantment:
			lores = append(lores, text.Colourf("<yellow>%s %s</yellow>", typ.Name(), lvl))
		}
	}
	i = i.WithLore(lores...)
	return i
}

// EffectEnchantment ...
type EffectEnchantment interface {
	Effect() effect.Effect
}

// AttackEnchantment ...
type AttackEnchantment interface {
	AttackEntity(wearer world.Entity, ent world.Entity)
}
