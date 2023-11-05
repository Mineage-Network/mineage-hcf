package user

import (
	"fmt"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/google/uuid"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"math"
	_ "unsafe"
)

// TODO: world change

// Waypoint ...
type Waypoint struct {
	name     string
	uuid     uuid.UUID
	entityId int64

	pos      mgl64.Vec3
	distance float64
	hidden   bool
}

// NewWaypoint ...
func NewWaypoint(name string, pos mgl64.Vec3) Waypoint {
	return Waypoint{
		name:     name,
		uuid:     uuid.New(),
		entityId: 100,
		pos:      pos.Add(mgl64.Vec3{0.5, 0, 0.5}),
		distance: 30,
		hidden:   false,
	}
}

// ShowTo ...
func (w *Waypoint) ShowTo(p *player.Player) {
	session_writePacket(player_session(p), &packet.PlayerList{
		ActionType: packet.PlayerListActionAdd,
		Entries: []protocol.PlayerListEntry{
			{
				UUID:           w.uuid,
				Username:       w.name,
				EntityUniqueID: w.entityId,
				//Skin: protocol.Skin{
				//	SkinID: "Standard_Custom",
				//},
			},
		},
	})

	m := protocol.NewEntityMetadata()
	m.SetFlag(protocol.EntityDataKeyFlags, protocol.EntityDataFlagNoAI)
	m[protocol.EntityDataKeyScale] = float32(0.01)

	pk := &packet.AddPlayer{
		UUID:            w.uuid,
		Username:        w.name,
		EntityRuntimeID: uint64(w.entityId),
		GameType:        packet.GameTypeAdventure,
		AbilityData: protocol.AbilityData{
			EntityUniqueID:     w.entityId,
			PlayerPermissions:  0,
			CommandPermissions: 0,
			Layers:             []protocol.AbilityLayer{},
		},
		Position:       vec64To32(waypointPosition(p, w).Sub(mgl64.Vec3{0, 1.62, 0})),
		HeldItem:       protocol.ItemInstance{},
		EntityMetadata: m,
	}
	session_writePacket(player_session(p), pk)

	session_writePacket(player_session(p), &packet.PlayerList{
		ActionType: packet.PlayerListActionRemove,
		Entries: []protocol.PlayerListEntry{
			{
				UUID: w.uuid,
			},
		},
	})
}

// HideTo ...
func (w *Waypoint) HideTo(p *player.Player) {
	session_writePacket(player_session(p), &packet.RemoveActor{
		EntityUniqueID: w.entityId,
	})
}

// Update ...
func (w *Waypoint) Update(p *player.Player) {
	dist := math.Floor(p.Position().Sub(w.pos).Len())

	var pos mgl32.Vec3
	if dist <= 10 {
		pos = vec64To32(w.pos.Add(mgl64.Vec3{0, 3.5, 0}))
	} else {
		pos = vec64To32(waypointPosition(p, w)).Add(mgl32.Vec3{0, 2.5, 0})
	}

	session_writePacket(player_session(p), &packet.MovePlayer{
		EntityRuntimeID: uint64(w.entityId),
		Position:        pos,
		Mode:            packet.MoveModeNormal,
		HeadYaw:         0, Yaw: 0, Pitch: 0,
	})

	m := protocol.NewEntityMetadata()
	m[protocol.EntityDataKeyName] = fmt.Sprintf(text.Colourf("<grey>%s [%vm]</grey>", w.name, dist))
	m[protocol.EntityDataKeyAlwaysShowNameTag] = byte(1)

	session_writePacket(player_session(p), &packet.SetActorData{
		EntityRuntimeID:  uint64(w.entityId),
		EntityMetadata:   m,
		EntityProperties: protocol.EntityProperties{},
	})
}

// AddWaypoint ...
func (u *User) AddWaypoint(name string, pos mgl64.Vec3) {
	w := NewWaypoint(name, pos)
	w.ShowTo(u.p)
	u.waypoints[name] = w
}

// RemoveWaypoint ...
func (u *User) RemoveWaypoint(name string) {
	for _, w := range u.Waypoints() {
		if w.name == name {
			w.HideTo(u.p)
		}
	}
	delete(u.waypoints, name)
}

// Waypoints ...
func (u *User) Waypoints() []Waypoint {
	var waypoints []Waypoint
	for _, waypoint := range u.waypoints {
		waypoints = append(waypoints, waypoint)
	}
	return waypoints
}

// waypointPosition ...
func waypointPosition(p *player.Player, w *Waypoint) mgl64.Vec3 {
	pos := eyePosition(p)
	wPos := pos.Add(mgl64.Vec3{0, p.EyeHeight(), 0})
	if pos.Sub(wPos).Len() <= w.distance {
		return mgl64.Vec3{wPos.X(), wPos.Y(), wPos.Z()}
	}
	actPos := pos.Add(wPos.Sub(pos).Normalize().Mul(w.distance))
	return mgl64.Vec3{actPos.X(), actPos.Y(), actPos.Z()}
}

// eyePosition ...
func eyePosition(e world.Entity) mgl64.Vec3 {
	pos := e.Position()
	if eyed, ok := e.(interface{ EyeHeight() float64 }); ok {
		pos = pos.Add(mgl64.Vec3{0, eyed.EyeHeight()})
	}
	return pos
}
