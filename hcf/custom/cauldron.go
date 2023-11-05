package custom

import (
	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/block/model"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/mineage-network/mineage-hcf/hcf/util/nbtconv"
)

// Cauldron ...
type Cauldron struct {
	solid

	Liquid    string
	FillLevel uint8

	world.NBTer
}

// EncodeItem ...
func (c Cauldron) EncodeItem() (name string, meta int16) {
	return "minecraft:cauldron", 0
}

// BreakInfo ...
func (c Cauldron) BreakInfo() block.BreakInfo {
	return newBreakInfo(2, pickaxeHarvestable, pickaxeHarvestable, oneOf(c))
}

// EncodeBlock ...
func (c Cauldron) EncodeBlock() (string, map[string]any) {
	return "minecraft:cauldron", map[string]any{"cauldron_liquid": c.Liquid, "fill_level": c.FillLevel}
}

// EncodeNBT ...
func (c Cauldron) EncodeNBT() map[string]any {
	return map[string]any{
		"cauldron_liquid": c.Liquid,
		"fill_level":      c.FillLevel,
	}
}

// DecodeNBT ...
func (c Cauldron) DecodeNBT(m map[string]any) any {
	return Cauldron{
		Liquid:    nbtconv.String(m, "cauldron_liquid"),
		FillLevel: nbtconv.Uint8(m, "fill_level"),
	}
}

// Hash ...
func (c Cauldron) Hash() uint64 {
	return nextHash()
}

// Model ...
func (c Cauldron) Model() world.BlockModel {
	return model.Solid{}
}
