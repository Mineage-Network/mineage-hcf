package util

import (
	"github.com/df-mc/dragonfly/server/player"
	"github.com/go-gl/mathgl/mgl64"
	"time"
)

// TeleportationFunc is a function called when a teleportation is performed.
type TeleportationFunc func(t *Teleportation)

// Teleportation represents a teleportation.
type Teleportation struct {
	expiration  time.Time
	pos         mgl64.Vec3
	teleporting bool

	f TeleportationFunc
	c chan struct{}
}

// NewTeleportation returns a new teleportation.
func NewTeleportation(f TeleportationFunc) *Teleportation {
	return &Teleportation{
		f: f,
		c: make(chan struct{}),
	}
}

// Teleport teleports the player to the position after the duration has passed.
func (t *Teleportation) Teleport(p *player.Player, dur time.Duration, pos mgl64.Vec3) {
	t.expiration = time.Now().Add(dur)
	t.c = make(chan struct{})
	t.pos = pos
	t.teleporting = true

	go func() {
		select {
		case <-time.After(dur):
			if t.f != nil {
				t.f(t)
			}
			p.Teleport(pos)
			t.teleporting = false
		case <-t.c:
			t.teleporting = false
			return
		}
	}()
}

// Teleporting returns true if the player is currently teleporting.
func (t *Teleportation) Teleporting() bool {
	return t.teleporting
}

// Expired returns true if the teleportation has expired.
func (t *Teleportation) Expired() bool {
	return time.Now().After(t.expiration)
}

// Expiration returns the expiration time of the teleportation.
func (t *Teleportation) Expiration() time.Time {
	return t.expiration
}

// Pos returns the position the player will be teleported to.
func (t *Teleportation) Pos() mgl64.Vec3 {
	return t.pos
}

// C returns the channel that is closed when the teleportation is cancelled.
func (t *Teleportation) C() <-chan struct{} {
	return t.c
}

// Cancel cancels the teleportation.
func (t *Teleportation) Cancel() {
	if t.teleporting {
		t.c <- struct{}{}
	}
}
