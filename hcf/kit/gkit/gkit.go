package gkit

import (
	"github.com/df-mc/dragonfly/server/entity"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/go-gl/mathgl/mgl64"
	ench "github.com/mineage-network/mineage-hcf/hcf/enchantment"
	"math/rand"
)

// GKit contains all the items, armour, and effects obtained by a gkit.
type GKit interface {
	// Name is the name of the gkit.
	Name() string
	// Texture is the texture of the gkit.
	Texture() string
	// Items returns the items provided by the gkit.
	Items(*player.Player) (items [36]item.Stack)
	// Armour contains the armour applied by using the gkit.
	// The item stacks are ordered helmet, chestplate, leggings, and then boots.
	Armour(*player.Player) [4]item.Stack
}

// All returns with all gkits.
func All() []GKit {
	return []GKit{}
}

// Apply ...
func Apply(gKit GKit, p *player.Player) {
	inv := p.Inventory()
	armour := gKit.Armour(p)
	for slot, it := range gKit.Items(p) {
		if it.Empty() {
			continue
		}
		it = ench.AddEnchantmentLore(it)
		if inv.Slots()[slot].Item() != nil {
			dropItem(p, it)
		} else {
			_ = inv.SetItem(slot, it)
		}
	}
	arm := p.Armour()
	for slot, it := range armour {
		if it.Empty() {
			continue
		}
		it = ench.AddEnchantmentLore(it)
		if arm.Slots()[slot].Item() != nil {
			dropItem(p, it)
		} else {
			switch slot {
			case 0:
				arm.SetHelmet(it)
				arm.Inventory().Handler().HandlePlace(nil, 0, it)
			case 1:
				arm.SetChestplate(it)
				arm.Inventory().Handler().HandlePlace(nil, 1, it)
			case 2:
				arm.SetLeggings(it)
				arm.Inventory().Handler().HandlePlace(nil, 2, it)
			case 3:
				arm.SetBoots(it)
				arm.Inventory().Handler().HandlePlace(nil, 3, it)
			}
		}
	}
}

// dropItem ...
func dropItem(p *player.Player, it item.Stack) {
	w, pos := p.World(), p.Position()
	ent := entity.NewItem(it, pos)
	ent.SetVelocity(mgl64.Vec3{rand.Float64()*0.2 - 0.1, 0.2, rand.Float64()*0.2 - 0.1})
	w.AddEntity(ent)
}
