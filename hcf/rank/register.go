package rank

import (
	"github.com/mineage-network/mineage-hcf/hcf/util"
	"golang.org/x/exp/slices"
)

var (
	// ranks contains all registered util.Rank implementations.
	ranks []util.Rank
	// ranksByName contains all registered util.Rank implementations indexed by their name.
	ranksByName = map[string]util.Rank{}
)

// All returns all registered ranks.
func All() []util.Rank {
	return ranks
}

// Register registers a rank to the ranks list. The hierarchy of ranks is determined by the order of registration.
func Register(ra util.Rank) {
	ranks = append(ranks, ra)
	ranksByName[ra.Name()] = ra
}

// ByName returns the rank with the given name. If no rank with the given name is registered, the second return value
// is false.
func ByName(name string) (util.Rank, bool) {
	ra, ok := ranksByName[name]
	return ra, ok
}

// Staff returns true if the rank provided is a staff rank.
func Staff(ra util.Rank) bool {
	return Tier(ra) >= Tier(Mod{})
}

// Tier returns the tier of a rank based on its registration hierarchy.
func Tier(ra util.Rank) int {
	return slices.IndexFunc(ranks, func(other util.Rank) bool {
		return ra == other
	})
}

// init registers all implemented ranks.
func init() {
	Register(Operator{})
	Register(Player{})

	// TODO: Implement nitro booster rank, purchasable ranks, media rank, famous rank.

	Register(Mod{})
	Register(Admin{})
	Register(Manager{})
	Register(Owner{})
}
