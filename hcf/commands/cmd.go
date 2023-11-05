package commands

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/mineage-network/mineage-hcf/hcf/rank"
	"github.com/mineage-network/mineage-hcf/hcf/user"
	"github.com/mineage-network/mineage-hcf/hcf/util"
	"golang.org/x/text/language"
)

// locale ...
func locale(s cmd.Source) language.Tag {
	if p, ok := s.(*player.Player); ok {
		return p.Locale()
	}
	return language.English
}

// allow ...
func allow(src cmd.Source, console bool, roles ...util.Rank) bool {
	p, ok := src.(*player.Player)
	if !ok {
		return console
	}
	if len(roles) == 0 {
		return true
	}
	u, ok := user.Lookup(p)
	return ok && u.Ranks().Contains(append(roles, rank.Operator{})...)
}
