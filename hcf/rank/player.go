package rank

import "github.com/sandertv/gophertunnel/minecraft/text"

// Player ...
type Player struct{}

// Name ...
func (Player) Name() string {
	return "player"
}

// Chat ...
func (Player) Chat(name, message string) string {
	return text.Colourf("<grey>%s</grey><white>: %s</white>", name, message)
}

// Tag ...
func (Player) Tag(name string) string {
	return text.Colourf("<grey>%s</grey>", name)
}
