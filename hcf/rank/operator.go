package rank

import (
	"github.com/sandertv/gophertunnel/minecraft/text"
)

// Operator ...
type Operator struct{}

// Name ...
func (Operator) Name() string {
	return "operator"
}

// Chat ...
func (Operator) Chat(name, message string) string {
	return text.Colourf("<grey>%s</grey><white>: %s</white>", name, message)
}

// Tag ...
func (Operator) Tag(name string) string {
	return text.Colourf("<grey>%s</grey>", name)
}
