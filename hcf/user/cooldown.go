package user

import (
	"github.com/df-mc/dragonfly/server/world"
	"github.com/mineage-network/mineage-hcf/hcf/util"
)

// specialItemType ...
type specialItemType interface {
	Name() string
	Item() world.Item
	Lore() []string
	Key() string
}

// Cooldowns ...
type Cooldowns struct {
	pearl            *util.Cooldown
	rogueAbility     *util.Cooldown
	specialAbilities *util.Cooldown
	goldenApple      *util.Cooldown
	notchApple       *util.Cooldown
	bardItems        util.MappedCooldown[world.Item]
	kits             util.MappedCooldown[string]
	itemUse          util.MappedCooldown[world.Item]
	specialItems     util.MappedCooldown[specialItemType]
}

// NewCooldowns returns a new cooldowns struct.
func NewCooldowns() *Cooldowns {
	return &Cooldowns{
		pearl:            util.NewCooldown(),
		rogueAbility:     util.NewCooldown(),
		specialAbilities: util.NewCooldown(),
		goldenApple:      util.NewCooldown(),
		notchApple:       util.NewCooldown(),
		bardItems:        util.NewMappedCooldown[world.Item](),
		kits:             util.NewMappedCooldown[string](),
		itemUse:          util.NewMappedCooldown[world.Item](),
		specialItems:     util.NewMappedCooldown[specialItemType](),
	}
}

// Cooldowns ...
func (u *User) Cooldowns() *Cooldowns {
	return u.cooldowns
}

// Resetable returns all resetable cooldowns.
func (c *Cooldowns) Resetable() []*util.Cooldown {
	return append(append(c.bardItems.All(), []*util.Cooldown{
		c.pearl,
		c.rogueAbility,
		c.specialAbilities,
		c.goldenApple,
	}...), c.specialItems.All()...)
}

// Pearl returns the ender pearl cooldown.
func (c *Cooldowns) Pearl() *util.Cooldown {
	return c.pearl
}

// RogueAbility returns the rogue ability cooldown.
func (c *Cooldowns) RogueAbility() *util.Cooldown {
	return c.rogueAbility
}

// GoldenApple returns the golden apple cooldown.
func (c *Cooldowns) GoldenApple() *util.Cooldown {
	return c.goldenApple
}

// NotchApple returns the notch apple cooldown.
func (c *Cooldowns) NotchApple() *util.Cooldown {
	return c.notchApple
}

// BardItems returns all bard item cooldowns.
func (c *Cooldowns) BardItems() util.MappedCooldown[world.Item] {
	return c.bardItems
}

// Kits returns all kits.
func (c *Cooldowns) Kits() util.MappedCooldown[string] {
	return c.kits
}

// ItemUse returns all item use cooldowns.
func (c *Cooldowns) ItemUse() util.MappedCooldown[world.Item] {
	return c.itemUse
}

// SpecialItems returns all partner item cooldowns.
func (c *Cooldowns) SpecialItems() util.MappedCooldown[specialItemType] {
	return c.specialItems
}

// SpecialAbilities returns the pearl protection item cooldown.
func (c *Cooldowns) SpecialAbilities() *util.Cooldown {
	return c.specialAbilities
}
