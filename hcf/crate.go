package hcf

import (
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/item/enchantment"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/mineage-network/mineage-hcf/hcf/crate"
	ench "github.com/mineage-network/mineage-hcf/hcf/enchantment"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

var (
	kothEnchantments = []item.Enchantment{
		item.NewEnchantment(ench.Protection{}, 3),
		item.NewEnchantment(enchantment.Unbreaking{}, 3),
	}

	kothRewards = []crate.Reward{
		11: crate.NewReward(item.NewStack(item.Helmet{Tier: item.ArmourTierDiamond{}}, 1).WithEnchantments(
			append(kothEnchantments, item.NewEnchantment(ench.NightVision{}, 1), item.NewEnchantment(ench.Invisibility{}, 1))...).WithCustomName(text.Colourf("<dark-red>KOTH Helmet</dark-red>")), 20),
		12: crate.NewReward(item.NewStack(item.Chestplate{Tier: item.ArmourTierDiamond{}}, 1).WithEnchantments(
			append(kothEnchantments, item.NewEnchantment(ench.FireResistance{}, 1))...).WithCustomName(text.Colourf("<dark-red>KOTH Chestplate</dark-red>")), 20),
		13: crate.NewReward(item.NewStack(item.Leggings{Tier: item.ArmourTierDiamond{}}, 1).WithCustomName(text.Colourf("<dark-red>KOTH Leggings</dark-red>")), 20),
		14: crate.NewReward(item.NewStack(item.Boots{Tier: item.ArmourTierDiamond{}}, 1).WithEnchantments(
			append(kothEnchantments, item.NewEnchantment(ench.Speed{}, 2))...).WithCustomName(text.Colourf("<dark-red>KOTH Boots</dark-red>")), 20),
		15: crate.NewReward(item.NewStack(item.Sword{Tier: item.ToolTierDiamond}, 1).WithEnchantments(
			item.NewEnchantment(ench.Sharpness{}, 3), item.NewEnchantment(enchantment.Unbreaking{}, 3), item.NewEnchantment(enchantment.FireAspect{}, 2)).WithCustomName(text.Colourf("<dark-red>KOTH Fire</dark-red>")), 20),
	}
)

func init() {
	crate.Register(crate.NewCrate(text.Colourf("<b><dark-red>KOTH</dark-red></b>"), mgl64.Vec3{23, 70, 37}, kothRewards))
}
