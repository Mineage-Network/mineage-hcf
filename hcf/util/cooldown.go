package util

import (
	"github.com/df-mc/atomic"
	"time"
)

// Cooldown represents a cooldown.
type Cooldown struct {
	expiration atomic.Value[time.Time]
}

// NewCooldown returns a new cooldown.
func NewCooldown() *Cooldown {
	return &Cooldown{}
}

// Active returns true if the cooldown is active.
func (c *Cooldown) Active() bool {
	return c.expiration.Load().After(time.Now())
}

// Set sets the cooldown.
func (c *Cooldown) Set(d time.Duration) {
	c.expiration.Store(time.Now().Add(d))
}

// Reset resets the cooldown.
func (c *Cooldown) Reset() {
	c.expiration.Store(time.Time{})
}

// Remaining returns the remaining time of the cooldown.
func (c *Cooldown) Remaining() time.Duration {
	return time.Until(c.expiration.Load())
}

// MappedCooldown represents a cooldown mapped to a key.
type MappedCooldown[T comparable] map[T]*Cooldown

// NewMappedCooldown returns a new mapped cooldown.
func NewMappedCooldown[T comparable]() MappedCooldown[T] {
	return make(map[T]*Cooldown)
}

// Active returns true if the cooldown is active.
func (m MappedCooldown[T]) Active(key T) bool {
	cooldown, ok := m[key]
	return ok && cooldown.Active()
}

// Set sets the cooldown.
func (m MappedCooldown[T]) Set(key T, d time.Duration) {
	cooldown := m.Key(key)
	cooldown.Set(d)
	m[key] = cooldown
}

// Key returns the cooldown for the key.
func (m MappedCooldown[T]) Key(key T) *Cooldown {
	cooldown, ok := m[key]
	if !ok {
		newCD := &Cooldown{}
		m[key] = newCD
		return newCD
	}
	return cooldown
}

// Reset resets the cooldown.
func (m MappedCooldown[T]) Reset(key T) {
	delete(m, key)
}

// Remaining returns the remaining time of the cooldown.
func (m MappedCooldown[T]) Remaining(key T) time.Duration {
	cooldown, ok := m[key]
	if !ok {
		return 0
	}
	return cooldown.Remaining()
}

// All returns all cooldowns.
func (m MappedCooldown[T]) All() []*Cooldown {
	var cooldowns []*Cooldown
	for _, cooldown := range m {
		cooldowns = append(cooldowns, cooldown)
	}
	return cooldowns
}
