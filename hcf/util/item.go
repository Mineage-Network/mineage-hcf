package util

import (
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
	"reflect"
	"strings"
	"unicode"
)

// ItemName returns the name of the item.
func ItemName(i world.Item) string {
	var s strings.Builder

	if it, ok := i.(item.Sword); ok {
		switch it.Tier {
		case item.ToolTierDiamond:
			s.WriteString("Diamond ")
		case item.ToolTierGold:
			s.WriteString("Golden ")
		case item.ToolTierIron:
			s.WriteString("Iron ")
		case item.ToolTierStone:
			s.WriteString("Stone ")
		case item.ToolTierWood:
			s.WriteString("Wooden ")
		}
	}

	t := reflect.TypeOf(i)
	if t == nil {
		return ""
	}
	name := t.Name()

	for _, r := range name {
		if unicode.IsUpper(r) && !strings.HasPrefix(name, string(r)) {
			s.WriteRune(' ')
		}
		s.WriteRune(r)
	}
	return s.String()
}
