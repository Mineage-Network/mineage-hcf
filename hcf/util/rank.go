package util

// Rank represents a rank in-game. These can vary, and are used to specify the permissions of the user. It also contains
// the name of the rank, prefix, and colour.
type Rank interface {
	// Name returns the name of the rank, for example "Admin".
	Name() string
	// Chat returns the formatted chat message using the name and message provided.
	Chat(name, message string) string
	// Tag returns the formatted name-tag using the name provided.
	Tag(name string) string
}

// HeirRank represents a rank that inherits from another rank.
type HeirRank interface {
	// Inherits returns the rank that this rank inherits from.
	Inherits() Rank
}
