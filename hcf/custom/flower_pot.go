package custom

import (
	"github.com/df-mc/dragonfly/server/block/model"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/mineage-network/mineage-hcf/hcf/util/nbtconv"
)

// FlowerPot ...
type FlowerPot struct {
	solid

	UpdateBit uint8

	world.NBTer
}

// EncodeItem ...
func (f FlowerPot) EncodeItem() (name string, meta int16) {
	return "minecraft:flower_pot", 0
}

// EncodeBlock ...
func (f FlowerPot) EncodeBlock() (string, map[string]any) {
	return "minecraft:flower_pot", map[string]any{"update_bit": f.UpdateBit}
}

// EncodeNBT ...
func (f FlowerPot) EncodeNBT() map[string]any {
	return map[string]any{
		"update_bit": f.UpdateBit,
	}
}

// DecodeNBT ...
func (f FlowerPot) DecodeNBT(m map[string]any) any {
	return FlowerPot{
		UpdateBit: nbtconv.Uint8(m, "update_bit"),
	}
}

// Hash ...
func (f FlowerPot) Hash() uint64 {
	return nextHash()
}

// Model ...
func (f FlowerPot) Model() world.BlockModel {
	return model.Solid{}
}
