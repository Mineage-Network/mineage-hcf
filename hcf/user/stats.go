package user

import "github.com/df-mc/atomic"

// Stats ...
type Stats struct {
	kills  atomic.Value[int]
	deaths atomic.Value[int]
}

// NewStats returns a new stats struct.
func NewStats(k, d int) *Stats {
	return &Stats{
		kills:  *atomic.NewValue(k),
		deaths: *atomic.NewValue(d),
	}
}

// DefaultStats ...
func DefaultStats() *Stats {
	return &Stats{
		kills:  *atomic.NewValue(0),
		deaths: *atomic.NewValue(0),
	}
}

// Stats ...
func (u *User) Stats() *Stats {
	return u.stats
}

// Kills ...
func (s *Stats) Kills() int {
	return s.kills.Load()
}

// Deaths ...
func (s *Stats) Deaths() int {
	return s.deaths.Load()
}

// AddKill ...
func (s *Stats) AddKill() {
	s.kills.Store(s.Kills() + 1)
}

// AddDeath ...
func (s *Stats) AddDeath() {
	s.deaths.Store(s.Deaths() + 1)
}
