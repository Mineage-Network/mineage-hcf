package commands

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/mineage-network/mineage-hcf/hcf/backend/lang"
	"github.com/mineage-network/mineage-hcf/hcf/rank"
	"github.com/mineage-network/mineage-hcf/hcf/user"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"strings"
)

// Whisper is a command that allows a player to send a private message to another player.
type Whisper struct {
	Target  []cmd.Target `cmd:"target"`
	Message cmd.Varargs  `cmd:"message"`
}

// Run ...
func (w Whisper) Run(s cmd.Source, o *cmd.Output) {
	l := locale(s)
	u, ok := user.Lookup(s.(*player.Player))
	if !ok {
		// The user somehow left in the middle of this, so just stop in our tracks.
		return
	}
	msg := strings.TrimSpace(string(w.Message))
	if len(msg) <= 0 {
		o.Error(lang.Translatef(l, "message.empty"))
		return
	}
	if len(w.Target) > 1 {
		o.Error(lang.Translatef(l, "command.targets.exceed"))
		return
	}

	tP, ok := w.Target[0].(*player.Player)
	if !ok {
		o.Error(lang.Translatef(l, "command.target.unknown"))
		return
	}
	t, ok := user.Lookup(tP)
	if !ok {
		o.Error(lang.Translatef(l, "command.target.unknown"))
		return
	}

	uTag, uMsg := text.Colourf("<white>%s</white>", u.Name()), text.Colourf("<white>%s</white>", msg)
	tTag, tMsg := text.Colourf("<white>%s</white>", t.Name()), text.Colourf("<white>%s</white>", msg)
	if _, ok := u.Ranks().Highest().(rank.Player); !ok {
		uMsg = t.Ranks().Highest().Tag(msg)
		uTag = u.Ranks().Highest().Tag(u.Name())
	}
	if _, ok := t.Ranks().Highest().(rank.Player); !ok {
		tMsg = u.Ranks().Highest().Tag(msg)
		tTag = t.Ranks().Highest().Tag(t.Name())
	}

	t.SetLastMessageFrom(u.Player())
	t.SendCustomSound("random.orb", 1, 1, false)
	u.Message("command.whisper.to", tTag, tMsg)
	t.Message("command.whisper.from", uTag, uMsg)
}

// Allow ...
func (Whisper) Allow(s cmd.Source) bool {
	_, ok := s.(*player.Player)
	return ok
}
