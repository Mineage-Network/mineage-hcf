package user

import "github.com/mineage-network/mineage-hcf/hcf/rank"

// Lives ...
type Lives struct {
	Current int8 `bson:"current"`
	Default int8 `bson:"default"`
}

// DefaultLives ...
func DefaultLives() *Lives {
	return &Lives{
		Current: 1,
		Default: 1,
	}
}

// HandleDeathBan ...
func (u *User) HandleDeathBan() {
	// TODO: Teleport to death ban arena here
	// and give the diamond kit, with the arenas own scoreboard.
}

// DeathBanned ...
func (u *User) DeathBanned() bool {
	return u.deathBanned.Load() || deathBanned.Contains(u.XUID())
}

// ReduceLife ...
func (u *User) ReduceLife() {
	l := u.lives.Load()
	if l.Current <= 0 {
		// The user is death-banned here.
		u.deathBanned.Toggle()
		deathBanned.Add(u.XUID())
		l.Current = -1
		return
	}
	l.Current -= 1
}

// sortLives ...
func (u *User) sortLives(new int8) {
	l := u.lives.Load()
	l.Current += new
	l.Default = new
}

// rankLives ...
func (u *User) rankLives() int8 {
	// TODO: Add other ranks here.
	switch (u.Ranks().Highest()).(type) {
	case rank.Player:
		return 1
	}
	return 1
}
