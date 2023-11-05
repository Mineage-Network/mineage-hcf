package commands

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/mineage-network/mineage-hcf/hcf/area"
	"github.com/mineage-network/mineage-hcf/hcf/user"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"time"
)

type Logout struct{}

// Run ...
func (l Logout) Run(s cmd.Source, o *cmd.Output) {
	p, ok := s.(*player.Player)
	if !ok {
		return
	}
	u, ok := user.Lookup(p)
	if !ok {
		return
	}

	for _, a := range area.Protected(p.World()) {
		if a.Area().Vec3WithinOrEqualXZ(p.Position()) {
			u.ToggleLogging()
			p.Disconnect(text.Colourf("<red>You have been logged out.</red>"))
			return
		}
	}

	if u.Teleportations().Logout().Teleporting() {
		o.Error("You are already logging out.")
		return
	}
	u.Teleportations().Logout().Teleport(p, time.Second*30, p.Position())
}
