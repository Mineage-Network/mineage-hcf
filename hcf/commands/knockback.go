package commands

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/mineage-network/mineage-hcf/hcf/knockback"
	"github.com/mineage-network/mineage-hcf/hcf/rank"
)

// KnockBack ...
type KnockBack struct {
	Type  kbType  `cmd:"type"`
	Value float64 `cmd:"value"`
}

// Run ...
func (kb KnockBack) Run(s cmd.Source, o *cmd.Output) {
	switch kb.Type {
	case "force":
		knockback.Knockback.SetForce(kb.Value)
	case "height":
		knockback.Knockback.SetHeight(kb.Value)
	case "hit-delay":
		knockback.Knockback.SetHitDelay(int(kb.Value))
	default:
		o.Error("Invalid data received")
	}
}

// Allow ...
func (KnockBack) Allow(s cmd.Source) bool {
	return allow(s, true, rank.Manager{})
}

type kbType string

// Type ...
func (kbType) Type() string {
	return "type"
}

// Options ...
func (kbType) Options(cmd.Source) []string {
	return []string{
		"force", "height", "hit-delay",
	}
}
