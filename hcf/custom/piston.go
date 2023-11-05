package custom

import (
	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/block/model"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/mineage-network/mineage-hcf/hcf/util/nbtconv"
)

// Piston ...
type Piston struct {
	solid

	Facing cube.Face
	Sticky bool

	world.NBTer
}

// EncodeItem ...
func (p Piston) EncodeItem() (name string, meta int16) {
	if p.Sticky {
		return "minecraft:sticky_piston", 0
	}
	return "minecraft:piston", 0
}

// BreakInfo ...
func (p Piston) BreakInfo() block.BreakInfo {
	return newBreakInfo(1.5, alwaysHarvestable, pickaxeEffective, oneOf(Piston{Sticky: p.Sticky}))
}

// EncodeBlock ...
func (p Piston) EncodeBlock() (string, map[string]any) {
	if p.Sticky {
		return "minecraft:sticky_piston", map[string]any{"facing_direction": int32(p.Facing)}
	}
	return "minecraft:piston", map[string]any{"facing_direction": int32(p.Facing)}
}

// EncodeNBT ...
func (p Piston) EncodeNBT() map[string]any {
	return map[string]any{
		"facing_direction": int32(p.Facing),
		"sticky_bit":       p.Sticky,
	}
}

// DecodeNBT ...
func (p Piston) DecodeNBT(m map[string]any) any {
	return Piston{
		Facing: cube.Face(nbtconv.Int32(m, "facing_direction")),
		Sticky: nbtconv.Bool(m, "Sticky"),
	}
}

// Hash ...
func (p Piston) Hash() uint64 {
	return nextHash()
}

// Model ...
func (p Piston) Model() world.BlockModel {
	return model.Solid{}
}
