package factions

// Rank represents a rank in a faction.
type Rank interface {
	// Name returns the name of the rank.
	Name() string
}

// RankLeader ...
type RankLeader struct{}

// Name ...
func (r RankLeader) Name() string {
	return "Leader"
}

// RankCoLeader ...
type RankCoLeader struct{}

// Name ...
func (r RankCoLeader) Name() string {
	return "Co-Leader"
}

// RankCaptain ...
type RankCaptain struct{}

// Name ...
func (r RankCaptain) Name() string {
	return "Captain"
}

// RankMember ...
type RankMember struct{}

// Name ...
func (r RankMember) Name() string {
	return "Member"
}
