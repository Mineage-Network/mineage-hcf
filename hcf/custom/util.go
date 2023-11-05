package custom

import (
	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/block/model"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/item/enchantment"
	"github.com/df-mc/dragonfly/server/world"
)

var (
	hash uint64 = 1000
)

// nextHash ...
func nextHash() uint64 {
	hash = hash + 1
	return hash
}

// solid represents a block that is fully solid. It always returns a model.Solid when Model is called.
type solid struct{}

// Model ...
func (solid) Model() world.BlockModel {
	return model.Solid{}
}

// empty represents a block that is fully empty/transparent, such as air or a plant. It always returns a
// model.Empty when Model is called.
type empty struct{}

// Model ...
func (empty) Model() world.BlockModel {
	return model.Empty{}
}

// chest represents a block that has a model of a chest.
type chest struct{}

// Model ...
func (chest) Model() world.BlockModel {
	return model.Chest{}
}

// carpet represents a block that has a model of a carpet.
type carpet struct{}

// Model ...
func (carpet) Model() world.BlockModel {
	return model.Carpet{}
}

// tilledGrass represents a block that has a model of farmland or dirt paths.
type tilledGrass struct{}

// Model ...
func (tilledGrass) Model() world.BlockModel {
	return model.TilledGrass{}
}

// leaves represents a block that has a model of leaves. A full block but with no solid faces.
type leaves struct{}

// Model ...
func (leaves) Model() world.BlockModel {
	return model.Leaves{}
}

// thin represents a thin, partial block such as a glass pane or an iron bar, that connects to nearby solid faces.
type thin struct{}

// Model ...
func (thin) Model() world.BlockModel {
	return model.Thin{}
}

// newBreakInfo creates a BreakInfo struct with the properties passed. The XPDrops field is 0 by default. The blast
// resistance is set to the block's hardness*5 by default.
func newBreakInfo(hardness float64, harvestable func(item.Tool) bool, effective func(item.Tool) bool, drops func(item.Tool, []item.Enchantment) []item.Stack) block.BreakInfo {
	return block.BreakInfo{
		Hardness:        hardness,
		BlastResistance: hardness * 5,
		Harvestable:     harvestable,
		Effective:       effective,
		Drops:           drops,
	}
}

// pickaxeEffective is a convenience function for custom that are effectively mined with a pickaxe.
var pickaxeEffective = func(t item.Tool) bool {
	return t.ToolType() == item.TypePickaxe
}

// axeEffective is a convenience function for custom that are effectively mined with an axe.
var axeEffective = func(t item.Tool) bool {
	return t.ToolType() == item.TypeAxe
}

// shearsEffective is a convenience function for custom that are effectively mined with shears.
var shearsEffective = func(t item.Tool) bool {
	return t.ToolType() == item.TypeShears
}

// shovelEffective is a convenience function for custom that are effectively mined with a shovel.
var shovelEffective = func(t item.Tool) bool {
	return t.ToolType() == item.TypeShovel
}

// hoeEffective is a convenience function for custom that are effectively mined with a hoe.
var hoeEffective = func(t item.Tool) bool {
	return t.ToolType() == item.TypeHoe
}

// nothingEffective is a convenience function for custom that cannot be mined efficiently with any tool.
var nothingEffective = func(item.Tool) bool {
	return false
}

// alwaysHarvestable is a convenience function for custom that are harvestable using any item.
var alwaysHarvestable = func(t item.Tool) bool {
	return true
}

// neverHarvestable is a convenience function for custom that are not harvestable by any item.
var neverHarvestable = func(t item.Tool) bool {
	return false
}

// pickaxeHarvestable is a convenience function for custom that are harvestable using any kind of pickaxe.
var pickaxeHarvestable = pickaxeEffective

// simpleDrops returns a drops function that returns the items passed.
func simpleDrops(s ...item.Stack) func(item.Tool, []item.Enchantment) []item.Stack {
	return func(item.Tool, []item.Enchantment) []item.Stack {
		return s
	}
}

// oneOf returns a drops function that returns one of each of the item types passed.
func oneOf(i ...world.Item) func(item.Tool, []item.Enchantment) []item.Stack {
	return func(item.Tool, []item.Enchantment) []item.Stack {
		var s []item.Stack
		for _, it := range i {
			s = append(s, item.NewStack(it, 1))
		}
		return s
	}
}

// hasSilkTouch checks if an item has the silk touch enchantment.
func hasSilkTouch(enchantments []item.Enchantment) bool {
	for _, enchant := range enchantments {
		if _, ok := enchant.Type().(enchantment.SilkTouch); ok {
			return true
		}
	}
	return false
}

// silkTouchOneOf returns a drop function that returns 1x of the silk touch drop when silk touch exists, or 1x of the
// normal drop when it does not.
func silkTouchOneOf(normal, silkTouch world.Item) func(item.Tool, []item.Enchantment) []item.Stack {
	return func(t item.Tool, enchantments []item.Enchantment) []item.Stack {
		if hasSilkTouch(enchantments) {
			return []item.Stack{item.NewStack(silkTouch, 1)}
		}
		return []item.Stack{item.NewStack(normal, 1)}
	}
}

// silkTouchDrop returns a drop function that returns the silk touch drop when silk touch exists, or the
// normal drop when it does not.
func silkTouchDrop(normal, silkTouch item.Stack) func(item.Tool, []item.Enchantment) []item.Stack {
	return func(t item.Tool, enchantments []item.Enchantment) []item.Stack {
		if hasSilkTouch(enchantments) {
			return []item.Stack{silkTouch}
		}
		return []item.Stack{normal}
	}
}

// silkTouchOnlyDrop returns a drop function that returns the drop when silk touch exists.
func silkTouchOnlyDrop(it world.Item) func(t item.Tool, enchantments []item.Enchantment) []item.Stack {
	return func(t item.Tool, enchantments []item.Enchantment) []item.Stack {
		if hasSilkTouch(enchantments) {
			return []item.Stack{item.NewStack(it, 1)}
		}
		return nil
	}
}
