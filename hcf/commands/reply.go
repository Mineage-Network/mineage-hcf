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

// Reply is a command that allows a player to reply to their most recent private message.
type Reply struct {
	Message cmd.Varargs `cmd:"message"`
}

// Run ...
func (r Reply) Run(s cmd.Source, o *cmd.Output) {
	u, ok := user.Lookup(s.(*player.Player))
	if !ok {
		// The user somehow left in the middle of this, so just stop in our tracks.
		return
	}
	l := u.Player().Locale()
	msg := strings.TrimSpace(string(r.Message))
	if len(msg) <= 0 {
		o.Error(lang.Translatef(l, "message.empty"))
		return
	}

	t, ok := u.LastMessageFrom()
	if !ok {
		o.Error(lang.Translatef(l, "command.reply.none"))
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
func (Reply) Allow(s cmd.Source) bool {
	_, ok := s.(*player.Player)
	return ok
}
