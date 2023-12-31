package commands

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/mineage-network/mineage-hcf/hcf/backend/lang"
	"github.com/mineage-network/mineage-hcf/hcf/rank"
	"github.com/mineage-network/mineage-hcf/hcf/user"
)

// vanishGameMode is the game mode used by vanished players.
type vanishGameMode struct{}

func (vanishGameMode) AllowsEditing() bool      { return true }
func (vanishGameMode) AllowsTakingDamage() bool { return false }
func (vanishGameMode) CreativeInventory() bool  { return false }
func (vanishGameMode) HasCollision() bool       { return false }
func (vanishGameMode) AllowsFlying() bool       { return true }
func (vanishGameMode) AllowsInteraction() bool  { return true }
func (vanishGameMode) Visible() bool            { return true }

// Vanish is a command that hides a staff from everyone else.
type Vanish struct{}

// Run ...
func (Vanish) Run(s cmd.Source, o *cmd.Output) {
	u, ok := user.Lookup(s.(*player.Player))
	if !ok {
		// The user somehow left in the middle of this, so just stop in our tracks.
		return
	}
	sn := s.(cmd.NamedTarget)
	if u.Vanished() {
		user.Alert(sn, "staff.alert.vanish.off")
		u.Player().SetGameMode(world.GameModeSurvival)
		o.Print(lang.Translatef(u.Player().Locale(), "command.vanish.disabled"))
	} else {
		user.Alert(sn, "staff.alert.vanish.on")
		u.Player().SetGameMode(vanishGameMode{})
		o.Print(lang.Translatef(u.Player().Locale(), "command.vanish.enabled"))
	}
	for _, t := range user.All() {
		if !u.Vanished() {
			t.Player().HideEntity(u.Player())
			continue
		}
		t.Player().ShowEntity(u.Player())
	}
	u.ToggleVanish()
}

// Allow ...
func (Vanish) Allow(s cmd.Source) bool {
	return allow(s, false, rank.Mod{})
}
