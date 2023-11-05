package commands

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/mineage-network/mineage-hcf/hcf/backend/lang"
	"github.com/mineage-network/mineage-hcf/hcf/rank"
	"github.com/mineage-network/mineage-hcf/hcf/user"
)

// Kick is a command that disconnects another player from the server.
type Kick struct {
	Targets []cmd.Target `cmd:"target"`
}

// Run ...
func (k Kick) Run(s cmd.Source, o *cmd.Output) {
	l, single := locale(s), true
	if len(k.Targets) > 1 {
		if p, ok := s.(*player.Player); ok {
			if u, ok := user.Lookup(p); ok && !u.Ranks().Contains(rank.Operator{}) {
				o.Error(lang.Translatef(l, "command.targets.exceed"))
				return
			}
		}
		single = false
	}

	var kicked int
	for _, p := range k.Targets {
		if p, ok := p.(*player.Player); ok {
			u, ok := user.Lookup(p)
			if !ok || u.Ranks().Contains(rank.Operator{}) {
				o.Print(lang.Translatef(l, "command.kick.fail"))
				continue
			}
			p.Disconnect(lang.Translatef(p.Locale(), "command.kick.reason"))
			if single {
				o.Print(lang.Translatef(l, "command.kick.success", p.Name()))
				return
			}
			kicked++
		} else if single {
			o.Print(lang.Translatef(l, "command.target.unknown"))
			return
		}
	}
	if !single {
		return
	}
	o.Print(lang.Translatef(l, "command.kick.multiple", kicked))
}

// Allow ...
func (Kick) Allow(s cmd.Source) bool {
	return allow(s, true, rank.Mod{})
}
