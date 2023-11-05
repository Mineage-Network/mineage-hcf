package rank

import (
	"github.com/mineage-network/mineage-hcf/hcf/util"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

// Manager ...
type Manager struct{}

// Name ...
func (Manager) Name() string {
	return "manager"
}

// Chat ...
func (Manager) Chat(name, message string) string {
	return text.Colourf("<grey>[<dark-red>Manager</dark-red>]</grey> <dark-red>%s</dark-red><dark-grey>:</dark-grey> <dark-red>%s</dark-red>", name, message)
}

// Tag ...
func (Manager) Tag(name string) string {
	return text.Colourf("<dark-red>%s</dark-red>", name)
}

// Inherits ...
func (Manager) Inherits() util.Rank {
	return Admin{}
}
