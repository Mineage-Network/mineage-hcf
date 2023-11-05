package factions

// FocusType represents different focus factions may have.
type FocusType struct {
	focus int
}

// FocusTypePlayer returns the player focus type.
func FocusTypePlayer() FocusType {
	return FocusType{focus: 1}
}

// FocusTypeFaction returns the faction focus type.
func FocusTypeFaction() FocusType {
	return FocusType{focus: 2}
}

// FocusTypes returns a list of all focus types.
func FocusTypes() []FocusType {
	return []FocusType{
		FocusTypePlayer(),
		FocusTypeFaction(),
	}
}

// String ...
func (d FocusType) String() string {
	switch d.focus {
	case 1:
		return "Player"
	case 2:
		return "Faction"
	}
	panic("should never happen")
}
