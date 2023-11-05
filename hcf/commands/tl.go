package commands

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/mineage-network/mineage-hcf/hcf/user"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

// TL is a command that allows players to see the coordinates of their faction members.
type TL struct{}

// Run ...
func (TL) Run(s cmd.Source, o *cmd.Output) {
	p, ok := s.(*player.Player)
	if !ok {
		return
	}
	u, ok := user.Lookup(p)
	if !ok {
		return
	}
	fa, ok := u.Faction()
	if !ok {
		p.Message(text.Colourf("<red>%s</red>", "You are not in a faction."))
		return
	}
	for _, m := range fa.Members() {
		if uTarget, ok := user.LookupName(m.Name()); ok {
			uTarget.Player().Message(text.Colourf("<green>%s</green><grey>:</grey> <yellow>%d<grey>,</grey> %d<grey>,</grey> %d</yellow>", p.Name(), int(p.Position().X()), int(p.Position().Y()), int(p.Position().Z())))
		}
	}
}
