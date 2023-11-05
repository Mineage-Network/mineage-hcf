package area

import (
	"github.com/go-gl/mathgl/mgl64"
	"github.com/mineage-network/mineage-hcf/hcf/util"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

var (
	Overworld = Manager{
		spawn:      NewZone(text.Colourf("<green>Spawn</green>"), util.NewAreaVec2(mgl64.Vec2{107, -107}, mgl64.Vec2{-107, 107})),
		warZone:    NewZone(text.Colourf("<red>Warzone</red>"), util.NewAreaVec2(mgl64.Vec2{300, 300}, mgl64.Vec2{-300, -300})),
		wilderness: NewZone(text.Colourf("<grey>Wilderness</grey>"), util.NewAreaVec2(mgl64.Vec2{3000, 3000}, mgl64.Vec2{-3000, -3000})),
		roads: []Zone{
			NewZone(text.Colourf("<red>Road</red>"), util.NewAreaVec2(mgl64.Vec2{-107, -15}, mgl64.Vec2{-2540, 15})),
			NewZone(text.Colourf("<red>Road</red>"), util.NewAreaVec2(mgl64.Vec2{107, 15}, mgl64.Vec2{2540, -15})),
			NewZone(text.Colourf("<red>Road</red>"), util.NewAreaVec2(mgl64.Vec2{15, -107}, mgl64.Vec2{-15, -2540})),
			NewZone(text.Colourf("<red>Road</red>"), util.NewAreaVec2(mgl64.Vec2{-15, 107}, mgl64.Vec2{15, 2540})),
		},
		koths: []Zone{
			// TODO: Add koths, this is just a placeholder.
			NewZone(text.Colourf("<red>Example KOTH</red>"), util.NewAreaVec2(mgl64.Vec2{0, 0}, mgl64.Vec2{0, 0})),
		},
	}
	Nether = Manager{
		spawn:      NewZone(text.Colourf("<green>Spawn</green>"), util.NewAreaVec2(mgl64.Vec2{60, 65}, mgl64.Vec2{-65, -60})),
		warZone:    NewZone(text.Colourf("<red>Warzone</red>"), util.NewAreaVec2(mgl64.Vec2{300, 300}, mgl64.Vec2{-300, -300})),
		wilderness: NewZone(text.Colourf("<grey>Wilderness</grey>"), util.NewAreaVec2(mgl64.Vec2{3000, 3000}, mgl64.Vec2{-3000, -3000})),
		roads: []Zone{
			NewZone(text.Colourf("<red>Road</red>"), util.NewAreaVec2(mgl64.Vec2{-66, -7}, mgl64.Vec2{-2540, 7})),
			NewZone(text.Colourf("<red>Road</red>"), util.NewAreaVec2(mgl64.Vec2{61, 7}, mgl64.Vec2{2540, -7})),
			NewZone(text.Colourf("<red>Road</red>"), util.NewAreaVec2(mgl64.Vec2{6, -61}, mgl64.Vec2{-8, -2540})),
			NewZone(text.Colourf("<red>Road</red>"), util.NewAreaVec2(mgl64.Vec2{-8, 66}, mgl64.Vec2{7, 2540})),
		},
	}
)

// Manager ...
type Manager struct {
	spawn      Zone
	warZone    Zone
	wilderness Zone
	roads      []Zone
	koths      []Zone
}

// Spawn ...
func (m Manager) Spawn() Zone {
	return m.spawn
}

// WarZone ...
func (m Manager) WarZone() Zone {
	return m.warZone
}

// Wilderness ...
func (m Manager) Wilderness() Zone {
	return m.wilderness
}

// Roads ...
func (m Manager) Roads() []Zone {
	return m.roads
}

// KOTHs ...
func (m Manager) KOTHs() []Zone {
	return m.koths
}
