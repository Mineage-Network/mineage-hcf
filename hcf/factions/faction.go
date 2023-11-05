package factions

import (
	"fmt"
	"github.com/df-mc/dragonfly/server"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/mineage-network/mineage-hcf/hcf/util"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"math"
	"strings"
	"sync"
	"time"
)

var (
	factionMu sync.Mutex
	factions  = map[*Faction]struct{}{}
)

// init ...
func init() {
	go func() {
		for {
			factionMu.Lock()
			for f := range factions {
				if !f.Frozen() && f.regenerationTime != 0 {
					f.dtr = f.MaxDTR()
					f.SetRegenerationTime(time.Time{})
				}
			}
			factionMu.Unlock()
			time.Sleep(time.Millisecond * 100)
		}
	}()
}

// Claims returns a map of all claims.
func Claims() map[*Faction]util.AreaVec2 {
	factionMu.Lock()
	defer factionMu.Unlock()
	claims := map[*Faction]util.AreaVec2{}
	for f := range factions {
		if f.claim == (util.AreaVec2{}) {
			continue
		}
		claims[f] = f.claim
	}
	return claims
}

// All returns all factions.
func All() []*Faction {
	factionMu.Lock()
	defer factionMu.Unlock()
	var all []*Faction
	for f := range factions {
		all = append(all, f)
	}
	return all
}

// LookupName returns a faction by name.
func LookupName(name string) (*Faction, bool) {
	factionMu.Lock()
	defer factionMu.Unlock()
	for f := range factions {
		if strings.EqualFold(f.name, name) {
			return f, true
		}
	}
	return nil, false
}

// LookupMemberName returns a faction by member name.
func LookupMemberName(name string) (*Faction, bool) {
	factionMu.Lock()
	defer factionMu.Unlock()
	for f := range factions {
		for _, member := range f.members {
			if strings.EqualFold(member.name, name) {
				return f, true
			}
		}
	}
	return nil, false
}

// LookupMember returns a faction by a player.Player.
func LookupMember(p *player.Player) (*Faction, bool) {
	factionMu.Lock()
	defer factionMu.Unlock()
	for f := range factions {
		for _, member := range f.members {
			if strings.EqualFold(member.name, p.Name()) {
				return f, true
			}
		}
	}
	return nil, false
}

// Faction ...
type Faction struct {
	displayName      string
	name             string
	dtr              float64
	home             mgl64.Vec3
	balance          float64
	regenerationTime int
	points           int
	claim            util.AreaVec2
	members          []*Member
	focus            string
	focusType        FocusType
}

// NewFaction ...
func NewFaction(
	name string,
	members []*Member,
	dtr float64,
	home mgl64.Vec3,
	balance float64,
	regenerationTime int,
	points int,
	claim util.AreaVec2,
) *Faction {
	f := &Faction{
		displayName:      name,
		name:             strings.ToLower(name),
		dtr:              dtr,
		home:             home,
		balance:          balance,
		regenerationTime: regenerationTime,
		points:           points,
		claim:            claim,
		members:          members,
	}

	factionMu.Lock()
	factions[f] = struct{}{}
	factionMu.Unlock()

	return f
}

// Name returns the name of the faction.
func (f *Faction) Name() string {
	return f.name
}

// DisplayName returns the display name of the faction.
func (f *Faction) DisplayName() string {
	return f.displayName
}

// FocusedFaction returns the focused faction.
func (f *Faction) FocusedFaction() (*Faction, bool) {
	if f.focusType == FocusTypeFaction() {
		return LookupName(f.focus)
	}
	return nil, false
}

// FocusedPlayer returns the focused player.
func (f *Faction) FocusedPlayer() (string, bool) {
	if f.focusType == FocusTypePlayer() {
		return f.focus, f.focus != ""
	}
	return "", false
}

// FocusFaction sets the focused faction.
func (f *Faction) FocusFaction(fac *Faction) {
	f.focus = fac.name
	f.focusType = FocusTypeFaction()
}

// FocusPlayer sets the focused player.
func (f *Faction) FocusPlayer(name string) {
	f.focus = name
	f.focusType = FocusTypePlayer()
}

// UnFocus unfocuses the faction.
func (f *Faction) UnFocus() {
	f.focus = ""
}

// DTR returns the DTR of the faction.
func (f *Faction) DTR() float64 {
	return math.Round(f.dtr*100) / 100
}

// SetDTR sets the DTR of the faction.
func (f *Faction) SetDTR(dtr float64) {
	f.dtr = dtr
}

// MaxDTR returns the max DTR of the faction.
func (f *Faction) MaxDTR() float64 {
	dtr := 1.1 * float64(f.PlayerCount())
	return math.Round(dtr*100) / 100
}

// DTRString returns the DTR string of the faction
func (f *Faction) DTRString() string {
	if f.DTR() == f.MaxDTR() {
		return text.Colourf("<green>%.1f%s</green>", f.DTR(), f.DTRDot())
	}
	if f.DTR() < 0 {
		return text.Colourf("<red>%.1f%s</red>", f.DTR(), f.DTRDot())
	}
	return text.Colourf("<white>%.1f%s</white>", f.DTR(), f.DTRDot())
}

// DTRDot returns the DTR dot of the faction.
func (f *Faction) DTRDot() string {
	if f.DTR() == f.MaxDTR() {
		return "<b><green>*</green></b>"
	}
	if f.DTR() < 0 {
		return "<b><red>*</red></b>"
	}
	return "<b><yellow>*</yellow></b>"
}

// SetClaim sets the claim of the faction.
func (f *Faction) SetClaim(claim util.AreaVec2) {
	f.claim = claim
}

// UnClaim unclaims the claim of the faction.
func (f *Faction) UnClaim() {
	f.claim = util.AreaVec2{}
}

// Claim returns the claim of the faction.
func (f *Faction) Claim() (util.AreaVec2, bool) {
	return f.claim, f.claim != util.AreaVec2{}
}

// AddPoints adds points to the faction.
func (f *Faction) AddPoints(points int) {
	f.points += points
}

// RemovePoints removes points from the faction.
func (f *Faction) RemovePoints(points int) {
	f.points -= points
}

// Points returns the points of the faction.
func (f *Faction) Points() int {
	return f.points
}

// IncreaseBalance increases the faction's balance.
func (f *Faction) IncreaseBalance(balance float64) {
	f.balance += balance
}

// ReduceBalance reduces the faction's balance.
func (f *Faction) ReduceBalance(balance float64) {
	f.balance -= balance
}

// Balance returns the balance of the faction.
func (f *Faction) Balance() float64 {
	return f.balance
}

// Home returns the home location of the faction.
func (f *Faction) Home() (mgl64.Vec3, bool) {
	return f.home, f.home != mgl64.Vec3{}
}

// SetHome sets the home location of the faction.
func (f *Faction) SetHome(home mgl64.Vec3) {
	f.home = home
}

// Leader returns the leader of the faction.
func (f *Faction) Leader() *Member {
	for _, m := range f.members {
		if m.rank == (RankLeader{}) {
			return m
		}
	}
	return nil
}

// CoLeaders returns the co-leaders of the faction.
func (f *Faction) CoLeaders() []*Member {
	var coleaders []*Member
	for _, m := range f.members {
		if m.rank == (RankCoLeader{}) {
			coleaders = append(coleaders, m)
		}
	}
	return coleaders
}

// Captains returns the captains of the faction.
func (f *Faction) Captains() []*Member {
	var captains []*Member
	for _, m := range f.members {
		if m.rank == (RankCaptain{}) {
			captains = append(captains, m)
		}
	}
	return captains
}

// Members returns the members of the faction.
func (f *Faction) Members() []*Member {
	return f.members
}

// AddMember adds a member to the faction.
func (f *Faction) AddMember(p *player.Player) {
	f.dtr += 1.1
	f.members = append(f.members, &Member{
		name:        strings.ToLower(p.Name()),
		displayName: p.Name(),
	})
}

// RemoveMember removes a member from the faction.
func (f *Faction) RemoveMember(p *player.Player) {
	f.dtr -= 1.1
	for i, m := range f.members {
		if m.name == strings.ToLower(p.Name()) {
			f.members = append(f.members[:i], f.members[i+1:]...)
			_ = m.Close()
			return
		}
	}
}

// RemoveMemberName removes a member from the faction by name.
func (f *Faction) RemoveMemberName(name string) {
	f.dtr -= 1.1
	for i, m := range f.members {
		if strings.EqualFold(m.name, name) {
			f.members = append(f.members[:i], f.members[i+1:]...)
			return
		}
	}
}

// Member returns a member of the faction.
func (f *Faction) Member(name string) (*Member, bool) {
	for _, m := range f.members {
		if strings.EqualFold(m.name, name) {
			return m, true
		}
	}
	return nil, false
}

// PlayerCount returns the number of players in the faction.
func (f *Faction) PlayerCount() int {
	return len(f.members)
}

// Frozen returns whether the faction is frozen.
func (f *Faction) Frozen() bool {
	return time.Now().Before(f.RegenerationTime())
}

// RegenerationTime returns the time until the faction's dtr regenerates.
func (f *Faction) RegenerationTime() time.Time {
	return time.UnixMilli(int64(f.regenerationTime))
}

// SetRegenerationTime sets the time until the faction's dtr regenerates.
func (f *Faction) SetRegenerationTime(regenerationTime time.Time) {
	f.regenerationTime = int(regenerationTime.UnixMilli())
}

// IsCoLeader returns whether the given player is a co-leader of the faction.
func (f *Faction) IsCoLeader(name string) bool {
	for _, m := range f.CoLeaders() {
		if strings.EqualFold(m.name, name) {
			return true
		}
	}
	return false
}

// IsCaptain returns whether the given player is a captain of the faction.
func (f *Faction) IsCaptain(name string) bool {
	for _, m := range f.Captains() {
		if strings.EqualFold(m.name, name) {
			return true
		}
	}
	return false
}

// Promote promotes a member of the faction.
func (f *Faction) Promote(name string) {
	m, ok := f.Member(name)
	if !ok {
		return
	}
	switch m.rank.(type) {
	case RankMember:
		m.rank = RankCaptain{}
	case RankCaptain:
		m.rank = RankCoLeader{}
	case RankCoLeader:
		l := f.Leader()
		l.rank = RankCoLeader{}
		m.rank = RankLeader{}
	}
}

// Demote demotes a member of the faction.
func (f *Faction) Demote(name string) {
	m, ok := f.Member(name)
	if !ok {
		return
	}
	switch m.rank.(type) {
	case RankCaptain:
		m.rank = RankMember{}
	case RankCoLeader:
		m.rank = RankCaptain{}
	}
}

// Raidable returns whether the faction is raidable.
func (f *Faction) Raidable() bool {
	return f.dtr <= 0
}

// Information returns a formatted string containing the information of the faction.
func (f *Faction) Information(srv *server.Server) string {
	var formattedRegenerationTime string
	var formattedDtr string
	var formattedLeader string
	var formattedCoLeaders []string
	var formattedCaptains []string
	var formattedMembers []string
	if time.Now().Before(f.RegenerationTime()) {
		formattedRegenerationTime = text.Colourf("\n <yellow>Time Until Regen</yellow> <blue>%s</blue>", time.Until(f.RegenerationTime()).Round(time.Second))
	}
	formattedDtr = f.DTRString()
	var onlineCount int
	for _, p := range f.Members() {
		_, ok := srv.PlayerByName(p.DisplayName())
		if ok {
			if p.Rank() == (RankLeader{}) {
				formattedLeader = text.Colourf("<green>%s</green>", p.DisplayName())
			} else if p.Rank() == (RankCoLeader{}) {
				formattedCoLeaders = append(formattedCoLeaders, text.Colourf("<green>%s</green>", p.DisplayName()))
			} else if p.Rank() == (RankCaptain{}) {
				formattedCaptains = append(formattedCaptains, text.Colourf("<green>%s</green>", p.DisplayName()))
			} else {
				formattedMembers = append(formattedMembers, text.Colourf("<green>%s</green>", p.DisplayName()))
			}
			onlineCount++
		} else {
			if p.Rank() == (RankLeader{}) {
				formattedLeader = text.Colourf("<grey>%s</grey>", p.DisplayName())
			} else if p.Rank() == (RankCoLeader{}) {
				formattedCoLeaders = append(formattedCoLeaders, text.Colourf("<green>%s</green>", p.DisplayName()))
			} else if p.Rank() == (RankCaptain{}) {
				formattedCaptains = append(formattedCaptains, text.Colourf("<grey>%s</grey>", p.DisplayName()))
			} else {
				formattedMembers = append(formattedMembers, text.Colourf("<grey>%s</grey>", p.DisplayName()))
			}
		}
	}
	if len(formattedCoLeaders) == 0 {
		formattedCoLeaders = []string{"None"}
	}
	if len(formattedCaptains) == 0 {
		formattedCaptains = []string{"None"}
	}
	if len(formattedMembers) == 0 {
		formattedMembers = []string{"None"}
	}
	var home string
	h, ok := f.Home()
	if !ok {
		home = "not set"
	} else {
		home = fmt.Sprintf("%.0f, %.0f, %.0f", h.X(), h.Y(), h.Z())
	}
	return text.Colourf(
		"\uE000\n <blue>%s</blue> <grey>[%d/%d]</grey> <dark-aqua>-</dark-aqua> <yellow>HQ:</yellow> %s\n "+
			"<yellow>Leader: </yellow>%s\n "+
			"<yellow>Co-Leaders: </yellow>%s\n "+
			"<yellow>Captains: </yellow>%s\n "+
			"<yellow>Members: </yellow>%s\n "+
			"<yellow>Balance: </yellow><blue>$%2.f</blue>\n "+
			"<yellow>Points: </yellow><blue>%d</blue>\n "+
			"<yellow>Deaths until Raidable: </yellow>%s%s\n\uE000", f.DisplayName(), onlineCount, f.PlayerCount(), home, formattedLeader, strings.Join(formattedCoLeaders, ", "), strings.Join(formattedCaptains, ", "), strings.Join(formattedMembers, ", "), f.Balance(), f.Points(), formattedDtr, formattedRegenerationTime)
}

// Close closes the faction.
func (f *Faction) Close() {
	factionMu.Lock()
	defer factionMu.Unlock()
	delete(factions, f)
}
