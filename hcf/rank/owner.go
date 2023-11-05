package rank

import (
	"github.com/mineage-network/mineage-hcf/hcf/util"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

// Owner ...
type Owner struct{}

// Name ...
func (Owner) Name() string {
	return "owner"
}

// Chat ...
func (Owner) Chat(name, message string) string {
	return text.Colourf("<grey>[<dark-red>Owner</dark-red>]</grey> <dark-red>%s</dark-red><dark-grey>:</dark-grey> <dark-red>%s</dark-red>", name, message)
}

// Tag ..
func (Owner) Tag(name string) string {
	return text.Colourf("<dark-red>%s</dark-red>", name)
}

// Inherits ...
func (Owner) Inherits() util.Rank {
	return Manager{}
}
