package rank

import (
	"github.com/mineage-network/mineage-hcf/hcf/util"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

// Admin ...
type Admin struct{}

// Name ...
func (Admin) Name() string {
	return "admin"
}

// Chat ...
func (Admin) Chat(name, message string) string {
	return text.Colourf("<grey>[<red>Admin</red>]</grey> <red>%s</red><dark-grey>:</dark-grey> <red>%s</red>", name, message)
}

// Tag ...
func (Admin) Tag(name string) string {
	return text.Colourf("<red>%s</red>", name)
}

// Inherits ...
func (Admin) Inherits() util.Rank {
	return Mod{}
}
