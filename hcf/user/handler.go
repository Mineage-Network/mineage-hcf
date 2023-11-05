package user

import (
	"github.com/df-mc/dragonfly/server/entity/effect"
	"github.com/df-mc/dragonfly/server/event"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	ench "github.com/mineage-network/mineage-hcf/hcf/enchantment"
	"github.com/mineage-network/mineage-hcf/hcf/util/class"
)

// NopArmourTier ...
type NopArmourTier struct{}

func (NopArmourTier) BaseDurability() float64      { return 0 }
func (NopArmourTier) Toughness() float64           { return 0 }
func (NopArmourTier) KnockBackResistance() float64 { return 0 }
func (NopArmourTier) EnchantmentValue() int        { return 0 }
func (NopArmourTier) Name() string                 { return "" }

// armourHandler ...
type armourHandler struct {
	u *User
	p *player.Player
}

// HandleTake ...
func (h *armourHandler) HandleTake(_ *event.Context, _ int, it item.Stack) {
	h.u.SetClass(nil)

	var effects []effect.Effect

	for _, e := range it.Enchantments() {
		if enc, ok := e.Type().(ench.EffectEnchantment); ok {
			effects = append(effects, enc.Effect())
		}
	}

	for _, e := range effects {
		typ := e.Type()
		if h.u.HasEffectLevel(typ, e.Level()) {
			h.p.RemoveEffect(typ)
		}
	}
}

// HandlePlace ...
func (h *armourHandler) HandlePlace(_ *event.Context, _ int, it item.Stack) {
	var helmetTier item.ArmourTier = NopArmourTier{}
	var chestplateTier item.ArmourTier = NopArmourTier{}
	var leggingsTier item.ArmourTier = NopArmourTier{}
	var bootsTier item.ArmourTier = NopArmourTier{}

	a := h.p.Armour()
	helmet, ok := a.Helmet().Item().(item.Helmet)
	if ok {
		helmetTier = helmet.Tier
	}
	chestplate, ok := a.Chestplate().Item().(item.Chestplate)
	if ok {
		chestplateTier = chestplate.Tier
	}
	leggings, ok := a.Leggings().Item().(item.Leggings)
	if ok {
		leggingsTier = leggings.Tier
	}
	boots, ok := a.Boots().Item().(item.Boots)
	if ok {
		bootsTier = boots.Tier
	}
	switch it := it.Item().(type) {
	case item.Helmet:
		helmetTier = it.Tier
	case item.Chestplate:
		chestplateTier = it.Tier
	case item.Leggings:
		leggingsTier = it.Tier
	case item.Boots:
		bootsTier = it.Tier
	}
	newArmour := [4]item.ArmourTier{helmetTier, chestplateTier, leggingsTier, bootsTier}
	h.u.SetClass(class.ResolveFromArmour(newArmour))

	var effects []effect.Effect

	for _, e := range it.Enchantments() {
		if enc, ok := e.Type().(ench.EffectEnchantment); ok {
			effects = append(effects, enc.Effect())
		}
	}

	for _, e := range effects {
		h.p.AddEffect(e)
	}
}

// HandleDrop ...
func (h *armourHandler) HandleDrop(_ *event.Context, _ int, it item.Stack) {
	h.u.SetClass(nil)

	var effects []effect.Effect

	for _, e := range it.Enchantments() {
		if enc, ok := e.Type().(ench.EffectEnchantment); ok {
			effects = append(effects, enc.Effect())
		}
	}

	for _, e := range effects {
		typ := e.Type()
		if h.u.HasEffectLevel(typ, e.Level()) {
			h.p.RemoveEffect(typ)
		}
	}
}
