package util

import (
	"github.com/df-mc/dragonfly/server/item"
	"math/rand"
	"strings"
)

// colours is a list of all colours.
var colours = []item.Colour{
	item.ColourWhite(),
	item.ColourOrange(),
	item.ColourMagenta(),
	item.ColourLightBlue(),
	item.ColourYellow(),
	item.ColourLime(),
	item.ColourPink(),
	item.ColourGrey(),
	item.ColourLightGrey(),
	item.ColourCyan(),
	item.ColourPurple(),
	item.ColourBlue(),
	item.ColourBrown(),
	item.ColourGreen(),
	item.ColourRed(),
	item.ColourBlack(),
}

// RandomColour returns a random colour.
func RandomColour() item.Colour {
	return colours[rand.Intn(len(colours))]
}

// StripMinecraftColour strips the minecraft colour from a string.
func StripMinecraftColour(s string) string {
	var str string
	for i := 0; i < len(s); i++ {
		if s[i] == '§' {
			i++
			continue
		}
		if s[i] != 'Â' {
			str += string(s[i])
		}
	}
	return strings.TrimSpace(str)
}
