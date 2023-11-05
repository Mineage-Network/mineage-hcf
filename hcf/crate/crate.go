package crate

import (
	"fmt"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/mineage-network/mineage-hcf/hcf/util"
	"math/rand"
)

// crates is a list of all crates.
var crates []Crate

// Register registers a crate.
func Register(c Crate) {
	crates = append(crates, c)
}

// All returns all crates.
func All() []Crate {
	return crates
}

// Crate represents a crate.
type Crate struct {
	name    string
	pos     mgl64.Vec3
	rewards []Reward
	stacks  []item.Stack
	value   string
}

// NewCrate returns a new crate.
func NewCrate(name string, pos mgl64.Vec3, rewards []Reward) Crate {
	c := Crate{
		name:    name,
		pos:     pos,
		rewards: rewards,
		value:   fmt.Sprintf("crate-key_%s", util.StripMinecraftColour(name)),
	}

	for _, r := range rewards {
		for i := 0; i < r.chance; i++ {
			c.stacks = append(c.stacks, r.stack)
		}
	}
	return c
}

// Name returns the name of the crate.
func (c Crate) Name() string {
	return c.name
}

// Position returns the position of the crate.
func (c Crate) Position() mgl64.Vec3 {
	return c.pos
}

// PositionMiddle returns the middle position of the crate.
func (c Crate) PositionMiddle() mgl64.Vec3 {
	return cube.PosFromVec3(c.pos).Vec3Middle()
}

// Rewards returns the rewards of the crate.
func (c Crate) Rewards() []Reward {
	return c.rewards
}

// Reward returns a random reward from the crate.
func (c Crate) Reward() item.Stack {
	return c.stacks[rand.Intn(len(c.stacks))]
}

// EncodeCrate is used to encode the crate into a string.
func (c Crate) EncodeCrate() string {
	return c.value
}
