package commands

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/mineage-network/mineage-hcf/hcf/user"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

// PvpEnable is a command to enable PVP.
type PvpEnable struct {
	Sub cmd.SubCommand `cmd:"enable"`
}

// Run ...
func (c PvpEnable) Run(s cmd.Source, o *cmd.Output) {
	u, ok := user.Lookup(s.(*player.Player))
	if !ok {
		return
	}
	if _, ok := u.TimerEnabled(); ok {
		u.DisableTimer()
		u.Player().Message(text.Colourf("<green>You have enabled PvP.</green>"))
	} else {
		u.Player().Message(text.Colourf("<red>You have already have enabled PvP.</red>"))
	}
}

// Allow ...
func (PvpEnable) Allow(s cmd.Source) bool {
	_, ok := s.(*player.Player)
	return ok
}
