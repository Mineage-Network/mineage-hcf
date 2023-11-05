package commands

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/mineage-network/mineage-hcf/hcf/koth"
	"github.com/mineage-network/mineage-hcf/hcf/rank"
	"github.com/mineage-network/mineage-hcf/hcf/user"
	"github.com/mineage-network/mineage-hcf/hcf/util"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"strings"
)

// KothList is a command that lists all KOTHs.
type KothList struct {
	Sub cmd.SubCommand `cmd:"list"`
}

// KothStart is a command that starts a KOTH.
type KothStart struct {
	Sub  cmd.SubCommand `cmd:"start"`
	KOTH kothList       `cmd:"koth"`
}

// KothStop is a command that stops a KOTH.
type KothStop struct {
	Sub cmd.SubCommand `cmd:"stop"`
}

// Run ...
func (KothList) Run(s cmd.Source, o *cmd.Output) {
	all := []string{
		text.Colourf("<yellow>KOTHs</yellow><grey>:</grey>"),
	}
	for _, k := range koth.All() {
		coords := k.Coordinates()
		all = append(all, text.Colourf("%s<grey>:</grey> <yellow>%0.f, %0.f</yellow>", k.Name(), coords.X(), coords.Y()))
	}
	o.Print(strings.Join(all, "\n"))
}

// Run ...
func (k KothStart) Run(s cmd.Source, o *cmd.Output) {
	name := text.Colourf("<grey>%s</grey>", s.(cmd.NamedTarget).Name())
	if p, ok := s.(*player.Player); ok {
		if u, ok := user.Lookup(p); ok {
			r := u.Ranks().Highest()
			name = r.Tag(p.Name())
		}
	}
	if _, ok := koth.Running(); ok {
		o.Print(text.Colourf("<red>A KOTH is already running.</red>"))
		return
	}
	ko, ok := koth.Lookup(string(k.KOTH))
	if !ok {
		o.Print(text.Colourf("<red>Invalid KOTH.</red>"))
		return
	}
	ko.Start()

	for _, u := range user.All() {
		if ko.Area().Vec3WithinOrEqualXZ(u.Player().Position()) {
			r := u.Ranks().Highest()
			ko.StartCapturing(u, r.Tag(u.Player().Name()))
		}
	}

	coords := ko.Coordinates()
	user.Broadcast("koth.start", name, ko.Name(), coords.X(), coords.Y())
}

// Run ...
func (KothStop) Run(s cmd.Source, o *cmd.Output) {
	name := text.Colourf("<grey>%s</grey>", s.(cmd.NamedTarget).Name())
	if p, ok := s.(*player.Player); ok {
		if u, ok := user.Lookup(p); ok {
			r := u.Ranks().Highest()
			name = r.Tag(p.Name())
		}
	}
	if k, ok := koth.Running(); !ok {
		o.Print(text.Colourf("<red>No KOTH is running.</red>"))
		return
	} else {
		k.Stop()
		user.Broadcast("koth.stop", name, k.Name())
	}
}

// Allow ...
func (KothList) Allow(s cmd.Source) bool {
	return allow(s, true, nil)
}

// Allow ...
func (KothStart) Allow(s cmd.Source) bool {
	return allow(s, true, rank.Mod{})
}

// Allow ...
func (KothStop) Allow(s cmd.Source) bool {
	return allow(s, true, rank.Mod{})
}

type (
	kothList string
)

// Type ...
func (kothList) Type() string {
	return "koth_list"
}

// Options ...
func (kothList) Options(cmd.Source) []string {
	var opts []string
	for _, k := range koth.All() {
		opts = append(opts, util.StripMinecraftColour(strings.ToLower(k.Name())))
	}
	return opts
}
