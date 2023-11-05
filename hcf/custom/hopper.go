package custom

import (
	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/block/model"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/mineage-network/mineage-hcf/hcf/util/nbtconv"
)

// Hopper ...
type Hopper struct {
	solid

	Facing  cube.Face
	Powered bool

	world.NBTer
}

// EncodeItem ...
func (h Hopper) EncodeItem() (name string, meta int16) {
	return "minecraft:hopper", 0
}

// BreakInfo ...
func (h Hopper) BreakInfo() block.BreakInfo {
	return newBreakInfo(3, pickaxeHarvestable, pickaxeEffective, oneOf(h))
}

// EncodeBlock ...
func (h Hopper) EncodeBlock() (string, map[string]any) {
	return "minecraft:hopper", map[string]any{
		"facing_direction": int32(h.Facing),
		"toggle_bit":       h.Powered,
	}
}

// EncodeNBT ...
func (h Hopper) EncodeNBT() map[string]any {
	return map[string]any{
		"facing_direction": int32(h.Facing),
		"toggle_bit":       h.Powered,
	}
}

// DecodeNBT ...
func (h Hopper) DecodeNBT(m map[string]any) any {
	return Hopper{
		Facing:  cube.Face(nbtconv.Int32(m, "facing_direction")),
		Powered: nbtconv.Bool(m, "toggle_bit"),
	}
}

// Hash ...
func (h Hopper) Hash() uint64 {
	return nextHash()
}

// Model ...
func (h Hopper) Model() world.BlockModel {
	return model.Solid{}
}
