package commands

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/mineage-network/mineage-hcf/hcf/backend/lang"
	"github.com/mineage-network/mineage-hcf/hcf/rank"
	"github.com/mineage-network/mineage-hcf/hcf/user"
	"strings"
)

// GameMode is a command for a player to change their own game mode or another player's.
type GameMode struct {
	GameMode gameMode                   `cmd:"gamemode"`
	Targets  cmd.Optional[[]cmd.Target] `cmd:"target"`
}

// Run ...
func (g GameMode) Run(s cmd.Source, o *cmd.Output) {
	var name string
	var mode world.GameMode
	switch strings.ToLower(string(g.GameMode)) {
	case "survival", "0", "s":
		name, mode = "survival", world.GameModeSurvival
	case "creative", "1", "c":
		name, mode = "creative", world.GameModeCreative
	case "adventure", "2", "a":
		name, mode = "adventure", world.GameModeAdventure
	case "spectator", "3", "sp":
		name, mode = "spectator", world.GameModeSpectator
	}

	l := locale(s)
	sn := s.(cmd.NamedTarget)
	targets := g.Targets.LoadOr(nil)
	if len(targets) > 1 {
		o.Error(lang.Translatef(l, "command.targets.exceed"))
		return
	}
	if len(targets) == 1 {
		target, ok := targets[0].(*player.Player)
		if !ok {
			o.Error(lang.Translatef(l, "command.target.unknown"))
			return
		}

		user.Alert(sn, "staff.alert.gamemode.change.other", target.Name(), name)

		target.SetGameMode(mode)
		o.Printf(lang.Translatef(l, "command.gamemode.update.other", target.Name(), name))
		return
	}
	if p, ok := s.(*player.Player); ok {
		user.Alert(sn, "staff.alert.gamemode.change", name)

		p.SetGameMode(mode)
		o.Printf(lang.Translatef(l, "command.gamemode.update.self", name))
		return
	}
	o.Error(lang.Translatef(l, "command.gamemode.console"))
}

// Allow ...
func (GameMode) Allow(s cmd.Source) bool {
	return allow(s, true, rank.Manager{})
}

type gameMode string

// Type ...
func (gameMode) Type() string {
	return "GameMode"
}

// Options ...
func (gameMode) Options(cmd.Source) []string {
	return []string{
		"survival", "0", "s",
		"creative", "1", "c",
		"adventure", "2", "a",
		"spectator", "3", "sp",
	}
}
