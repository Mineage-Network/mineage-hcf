package crate

import "github.com/df-mc/dragonfly/server/item"

// Reward represents a reward that can be given by a crate.
type Reward struct {
	stack  item.Stack
	chance int
}

// NewReward creates a new reward.
func NewReward(stack item.Stack, chance int) Reward {
	return Reward{stack: stack, chance: chance}
}

// Stack returns the item stack of the reward.
func (r Reward) Stack() item.Stack {
	return r.stack
}

// Chance returns the chance of the reward being given.
func (r Reward) Chance() int {
	return r.chance
}
