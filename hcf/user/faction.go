package user

import (
	"github.com/mineage-network/mineage-hcf/hcf/factions"
	"golang.org/x/exp/slices"
	"math"
	"time"
)

// Faction ...
type Faction struct {
	*factions.Faction
}

// Broadcast ...
func (f Faction) Broadcast(key string, args ...interface{}) {
	for _, u := range f.Users() {
		u.Message(key, args...)
	}
}

// Users ...
func (f Faction) Users() []*User {
	var members []*User
	for _, m := range f.Members() {
		if u, ok := LookupName(m.Name()); ok {
			if fa, ok := u.Faction(); ok && fa.Compare(f) {
				members = append(members, u)
			}
		}
	}
	return members
}

// Compare ...
func (f Faction) Compare(other any) bool {
	switch other := other.(type) {
	case Faction:
		return f.Faction == other.Faction
	case *factions.Faction:
		return f.Faction == other
	}
	return false
}

// SetFactionCreateDelay sets the time until the user can create a faction again.
func (u *User) SetFactionCreateDelay() {
	u.factionCreateDelay.Store(time.Now().Add(5 * time.Minute))
}

// FactionCreateDelay returns the time until the user can create a faction again.
func (u *User) FactionCreateDelay() (time.Time, bool) {
	d := u.factionCreateDelay.Load()
	return d, time.Now().Before(d)
}

// SetFaction ...
func (u *User) SetFaction(fa *factions.Faction) {
	u.faction.Store(Faction{fa})
}

// Nearby returns the nearby users of a certain distance from the user
func (u *User) Nearby(dist float64) []*User {
	var us []*User
	for _, o := range u.World().Entities() {
		if o.Position().ApproxFuncEqual(u.p.Position(), func(f float64, f2 float64) bool {
			return math.Max(f, f2)-math.Min(f, f2) < dist
		}) {
			target, ok := LookupEntity(o)
			if ok {
				us = append(us, target)
			}
		}
	}
	return us
}

// NearbyTeammates returns the nearby teammates of a certain distance from the user
func (u *User) NearbyTeammates(dist float64) []*User {
	var us []*User
	userFaction, ok := u.Faction()
	if !ok {
		return []*User{u}
	}
	for _, target := range u.Nearby(dist) {
		targetFaction, _ := target.Faction()
		if userFaction.Compare(targetFaction) {
			us = append(us, target)
		}
	}
	return us
}

// NearbyEnemies returns the nearby enemies of a certain distance from the user
func (u *User) NearbyEnemies(dist float64) []*User {
	var us []*User
	userFaction, ok := u.Faction()
	for _, target := range u.Nearby(dist) {
		targetFaction, _ := target.Faction()
		if target != u && (!ok || !userFaction.Compare(targetFaction)) {
			us = append(us, target)
		}
	}
	return us
}

// UnInvite removes an invitation.
func (u *User) UnInvite(name string) {
	i := slices.Index(u.invitations, name)
	if i == -1 {
		return
	}
	u.invitations = append(u.invitations[:i], u.invitations[i+1:]...)
}

// Invite invites a player to the faction.
func (u *User) Invite(name string) {
	u.invitations = append(u.invitations, name)
}

// Invitations returns the list of invitations.
func (u *User) Invitations() []string {
	return u.invitations
}
