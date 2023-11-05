package commands

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/mineage-network/mineage-hcf/hcf/sotw"
	"github.com/mineage-network/mineage-hcf/hcf/user"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

// Spawn is a command that teleports the player to the spawn.
type Spawn struct{}

// Run ...
func (Spawn) Run(s cmd.Source, o *cmd.Output) {
	p := s.(*player.Player)
	u, ok := user.Lookup(p)
	if !ok {
		return
	}

	if _, ok := sotw.Running(); ok && u.SOTW() {
		u.Player().Teleport(mgl64.Vec3{0, 66.5, 0})
		return
	}

	o.Print(text.Colourf("<red>You must be on SOTW timer to use /spawn<red>"))
}

// Allow ...
func (Spawn) Allow(s cmd.Source) bool {
	_, ok := s.(*player.Player)
	return ok
}
