package commands

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/mineage-network/mineage-hcf/hcf/data"
	"github.com/mineage-network/mineage-hcf/hcf/rank"
	"github.com/mineage-network/mineage-hcf/hcf/sotw"
	"github.com/mineage-network/mineage-hcf/hcf/user"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"go.mongodb.org/mongo-driver/bson"
)

// SOTWStart is a command to start SOTW.
type SOTWStart struct {
	Sub cmd.SubCommand `cmd:"start"`
}

// SOTWEnd is a command to end SOTW.
type SOTWEnd struct {
	Sub cmd.SubCommand `cmd:"end"`
}

// SOTWDisable is a command to disable SOTW.
type SOTWDisable struct {
	Sub cmd.SubCommand `cmd:"disable"`
}

// Run ...
func (c SOTWStart) Run(s cmd.Source, o *cmd.Output) {
	if _, ok := sotw.Running(); ok {
		o.Print(text.Colourf("<red>SOTW is already running!</red>"))
		return
	}
	sotw.Start()

	offline, err := data.SearchOfflineUsers(bson.M{})
	if err != nil {
		panic(err)
	}
	for _, u := range offline {
		u.SOTW = true
		_ = data.SaveOfflineUser(u)
	}
	for _, u := range user.All() {
		if !u.SOTW() {
			u.ToggleSOTW()
		}
	}
	user.Broadcast("sotw.commenced")
}

// Run ...
func (c SOTWEnd) Run(s cmd.Source, o *cmd.Output) {
	if _, ok := sotw.Running(); !ok {
		o.Print(text.Colourf("<red>SOTW is not running!</red>"))
		return
	}
	sotw.End()

	offline, err := data.SearchOfflineUsers(bson.M{})
	if err != nil {
		panic(err)
	}
	for _, u := range offline {
		u.SOTW = false
		err = data.SaveOfflineUser(u)
		if err != nil {
			panic(err)
		}
	}
	for _, u := range user.All() {
		if u.SOTW() {
			u.ToggleSOTW()
		}
	}
	user.Broadcast("sotw.ended")
}

// Run ...
func (c SOTWDisable) Run(s cmd.Source, o *cmd.Output) {
	u, ok := user.Lookup(s.(*player.Player))
	if !ok {
		return
	}
	if !u.SOTW() {
		u.Message("sotw.disabled.already")
		return
	}
	u.Message("sotw.disabled")
	u.ToggleSOTW()
}

// Allow ...
func (SOTWStart) Allow(s cmd.Source) bool {
	return allow(s, true, rank.Manager{})
}

// Allow ...
func (SOTWEnd) Allow(s cmd.Source) bool {
	return allow(s, true, rank.Manager{})
}

// Allow ...
func (SOTWDisable) Allow(s cmd.Source) bool {
	_, ok := s.(*player.Player)
	return ok
}
