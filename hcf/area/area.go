package area

import (
	"github.com/df-mc/dragonfly/server/world"
)

// Spawn ...
func Spawn(w *world.World) Zone {
	switch w.Dimension() {
	case world.Overworld:
		return Overworld.Spawn()
	case world.Nether:
		return Nether.Spawn()
	}
	return Zone{}
}

// WarZone ...
func WarZone(w *world.World) Zone {
	switch w.Dimension() {
	case world.Overworld:
		return Overworld.WarZone()
	case world.Nether:
		return Nether.WarZone()
	}
	return Zone{}
}

// Wilderness ...
func Wilderness(w *world.World) Zone {
	switch w.Dimension() {
	case world.Overworld:
		return Overworld.Wilderness()
	case world.Nether:
		return Nether.Wilderness()
	}
	return Zone{}
}

// Roads ...
func Roads(w *world.World) []Zone {
	switch w.Dimension() {
	case world.Overworld:
		return Overworld.Roads()
	case world.Nether:
		return Nether.Roads()
	}
	return []Zone{}
}

// KOTHs ...
func KOTHs(w *world.World) []Zone {
	switch w.Dimension() {
	case world.Overworld:
		return Overworld.KOTHs()
	case world.Nether:
		return Nether.KOTHs()
	}
	return []Zone{}
}

// Protected ...
func Protected(w *world.World) []Zone {
	return append(Roads(w), append(KOTHs(w), []Zone{
		Spawn(w),
		WarZone(w),
	}...)...)
}
