package user

import (
	"github.com/df-mc/atomic"
	"github.com/mineage-network/mineage-hcf/hcf/util"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

// Teleportations ...
type Teleportations struct {
	home   atomic.Value[*util.Teleportation]
	logout atomic.Value[*util.Teleportation]
	stuck  atomic.Value[*util.Teleportation]
}

// NewTeleportations returns a new teleportations struct.
func NewTeleportations(u *User) *Teleportations {
	return &Teleportations{
		home: *atomic.NewValue(util.NewTeleportation(func(t *util.Teleportation) {
			u.p.Message(text.Colourf("<green>You have been teleported home.</green>"))
		})),
		logout: *atomic.NewValue(util.NewTeleportation(func(t *util.Teleportation) {
			u.logged = true
			u.p.Disconnect(text.Colourf("<red>You have been logged out.</red>"))
		})),
		stuck: *atomic.NewValue(util.NewTeleportation(func(t *util.Teleportation) {
			u.p.Message(text.Colourf("<red>You have been teleported to a safe place.</red>"))
		})),
	}
}

// Home ...
func (t *Teleportations) Home() *util.Teleportation {
	return t.home.Load()
}

// Logout ...
func (t *Teleportations) Logout() *util.Teleportation {
	return t.logout.Load()
}

// Stuck ...
func (t *Teleportations) Stuck() *util.Teleportation {
	return t.stuck.Load()
}
