package rank

import "github.com/sandertv/gophertunnel/minecraft/text"

// Mod ...
type Mod struct{}

// Name ...
func (Mod) Name() string {
	return "mod"
}

// Chat ...
func (Mod) Chat(name, message string) string {
	return text.Colourf("<grey>[<purple>Mod</purple>]</grey> <purple>%s</purple><dark-grey>:</dark-grey> <purple>%s</purple>", name, message)
}

// Tag ...
func (Mod) Tag(name string) string {
	return text.Colourf("<purple>%s</purple>", name)
}
