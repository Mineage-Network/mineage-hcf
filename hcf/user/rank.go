package user

import (
	"github.com/mineage-network/mineage-hcf/hcf/rank"
	"github.com/mineage-network/mineage-hcf/hcf/util"
	"golang.org/x/exp/slices"
	"sort"
	"sync"
	"time"
)

// Ranks is a user-based rank manager for both offline users and online users.
type Ranks struct {
	rankMu          sync.Mutex
	ranks           []util.Rank
	rankExpirations map[util.Rank]time.Time
}

// NewRanks creates a new Ranks instance.
func NewRanks(ranks []util.Rank, expirations map[util.Rank]time.Time) *Ranks {
	return &Ranks{
		ranks:           ranks,
		rankExpirations: expirations,
	}
}

// Add adds a rank to the manager's rank list.
func (r *Ranks) Add(ra util.Rank) {
	r.rankMu.Lock()
	r.ranks = append(r.ranks, ra)
	r.rankMu.Unlock()
	r.sortRanks()
}

// Remove removes a rank from the manager's rank list. Users are responsible for updating the highest rank usages if
// changed.
func (r *Ranks) Remove(ra util.Rank) bool {
	if _, ok := ra.(rank.Player); ok {
		return false
	}

	r.rankMu.Lock()
	i := slices.IndexFunc(r.ranks, func(other util.Rank) bool {
		return ra == other
	})
	r.ranks = slices.Delete(r.ranks, i, i+1)
	delete(r.rankExpirations, ra)
	r.rankMu.Unlock()
	r.sortRanks()
	return true
}

// Staff returns true if the ranks contains a staff rank.
func (r *Ranks) Staff() bool {
	return r.Contains(rank.Mod{}, rank.Operator{})
}

// Contains returns true if the manager has any of the given ranks. Users are responsible for updating the highest rank
// usages if changed.
func (r *Ranks) Contains(ranks ...util.Rank) bool {
	r.rankMu.Lock()
	defer r.rankMu.Unlock()

	var actualRanks []util.Rank
	for _, ra := range r.ranks {
		r.propagateRanks(&actualRanks, ra)
	}

	for _, r := range ranks {
		if i := slices.IndexFunc(actualRanks, func(other util.Rank) bool {
			return r == other
		}); i >= 0 {
			return true
		}
	}
	return false
}

// Expiration returns the expiration time for a rank. If the rank does not expire, the second return value will be false.
func (r *Ranks) Expiration(ra util.Rank) (time.Time, bool) {
	r.rankMu.Lock()
	defer r.rankMu.Unlock()
	e, ok := r.rankExpirations[ra]
	return e, ok
}

// Expire sets the expiration time for a rank.
func (r *Ranks) Expire(ra util.Rank, t time.Time) {
	r.rankMu.Lock()
	defer r.rankMu.Unlock()
	r.rankExpirations[ra] = t
}

// Highest returns the highest rank the manager has, in terms of hierarchy.
func (r *Ranks) Highest() util.Rank {
	r.rankMu.Lock()
	defer r.rankMu.Unlock()
	return r.ranks[len(r.ranks)-1]
}

// All returns the user's ranks.
func (r *Ranks) All() []util.Rank {
	r.rankMu.Lock()
	defer r.rankMu.Unlock()
	return append(make([]util.Rank, 0, len(r.ranks)), r.ranks...)
}

// propagateRanks propagates ranks to the user's rank list.
func (r *Ranks) propagateRanks(actualRanks *[]util.Rank, ra util.Rank) {
	*actualRanks = append(*actualRanks, ra)
	if h, ok := ra.(util.HeirRank); ok {
		r.propagateRanks(actualRanks, h.Inherits())
	}
}

// sortRanks sorts the ranks in the user's rank list.
func (r *Ranks) sortRanks() {
	sort.SliceStable(r.ranks, func(i, j int) bool {
		return rank.Tier(r.ranks[i]) < rank.Tier(r.ranks[j])
	})
}
