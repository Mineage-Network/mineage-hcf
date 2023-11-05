package commands

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/mineage-network/mineage-hcf/hcf/data"
	"github.com/mineage-network/mineage-hcf/hcf/rank"
)

type DataReset struct {
	Kind dataKind `name:"kind"`
}

func (d DataReset) Run(cmd.Source, *cmd.Output) {
	switch d.Kind {
	case "users":
		data.ResetUsers()
	case "factions":
		data.ResetFactions()
	}
}

// Allow ...
func (DataReset) Allow(src cmd.Source) bool {
	return allow(src, true, rank.Operator{})
}

type dataKind string

// Type ...
func (dataKind) Type() string {
	return "data_kind"
}

// Options ...
func (dataKind) Options(cmd.Source) []string {
	return []string{"users", "factions"}
}
