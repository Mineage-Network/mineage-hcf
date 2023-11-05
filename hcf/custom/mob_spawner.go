package custom

import (
	"github.com/df-mc/dragonfly/server/block/model"
	"github.com/df-mc/dragonfly/server/world"
)

// MobSpawner ...
type MobSpawner struct {
	solid

	world.NBTer
}

// EncodeItem ...
func (m MobSpawner) EncodeItem() (name string, meta int16) {
	return "minecraft:mob_spawner", 0
}

// EncodeBlock ...
func (m MobSpawner) EncodeBlock() (string, map[string]any) {
	return "minecraft:mob_spawner", map[string]any{}
}

// EncodeNBT ...
func (m MobSpawner) EncodeNBT() map[string]any {
	return map[string]any{}
}

// DecodeNBT ...
func (m MobSpawner) DecodeNBT(map[string]any) any {
	return MobSpawner{}
}

// Hash ...
func (m MobSpawner) Hash() uint64 {
	return nextHash()
}

// Model ...
func (m MobSpawner) Model() world.BlockModel {
	return model.Solid{}
}
