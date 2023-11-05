package factions

// Member ...
type Member struct {
	name        string
	displayName string
	rank        Rank
}

// NewMember ...
func NewMember(name, displayName string, rank Rank) *Member {
	return &Member{
		name:        name,
		displayName: displayName,
		rank:        rank,
	}
}

// Name returns the name of the member.
func (m *Member) Name() string {
	return m.name
}

// DisplayName returns the display name of the member.
func (m *Member) DisplayName() string {
	return m.displayName
}

// Rank returns the rank of the member.
func (m *Member) Rank() Rank {
	return m.rank
}

// Close ...
func (m *Member) Close() error {
	return nil
}
