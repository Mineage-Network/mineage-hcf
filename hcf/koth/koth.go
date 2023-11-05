package koth

import (
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player/chat"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/mineage-network/mineage-hcf/hcf/custom"
	"github.com/mineage-network/mineage-hcf/hcf/factions"
	"github.com/mineage-network/mineage-hcf/hcf/util"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"math/rand"
	"strings"
	"time"
)

// init ...
func init() {
	t := time.NewTicker(time.Hour * 4)
	n := rand.Intn(len(All()))
	k := All()[n]
	_, _ = chat.Global.WriteString(text.Colourf("%s is now active!", k.name))
	k.Start()
	go func() {
		for range t.C {
			if _, ok := Running(); !ok {
				continue
			}
			n := rand.Intn(len(All()))
			k := All()[n]
			_, _ = chat.Global.WriteString(text.Colourf("%s is now active!", k.name))
			k.Start()
		}
	}()
}

// User ...
type User interface {
	Name() string
	AddItemOrDrop(it item.Stack)
}

var (
	// Broadcast ...
	Broadcast = func(format string, a ...interface{}) {}

	// TODO: Add koths, this is just a placeholder.

	// Example ...
	Example = &KOTH{
		name:        text.Colourf("<red>Example</red>"),
		area:        util.NewAreaVec2(mgl64.Vec2{0, 0}, mgl64.Vec2{0, 0}),
		cancel:      make(chan struct{}, 0),
		coordinates: mgl64.Vec2{0, 0},
	}
)

// All returns all KOTHs.
func All() []*KOTH {
	return []*KOTH{Example}
}

// Running returns true if the KOTH passed is currently running.
func Running() (*KOTH, bool) {
	for _, k := range All() {
		if k.running {
			return k, true
		}
	}
	return nil, false
}

// Lookup returns a KOTH by its name.
func Lookup(name string) (*KOTH, bool) {
	for _, k := range All() {
		if strings.EqualFold(util.StripMinecraftColour(k.Name()), name) {
			return k, true
		}
	}
	return nil, false
}

// KOTH represents a King of the Hill event.
type KOTH struct {
	name        string
	capturing   User
	running     bool
	time        time.Time
	cancel      chan struct{}
	area        util.AreaVec2
	coordinates mgl64.Vec2
}

// Name returns the name of the KOTH.
func (k *KOTH) Name() string {
	return k.name
}

// Start starts the KOTH.
func (k *KOTH) Start() {
	k.running = true
	k.capturing = nil
	k.cancel = make(chan struct{}, 0)
}

// Stop stops the KOTH.
func (k *KOTH) Stop() {
	k.running = false
	k.capturing = nil
	k.time = time.Time{}
	close(k.cancel)
}

// IsCapturing returns true if the player passed is currently capturing the KOTH.
func (k *KOTH) IsCapturing(u User) bool {
	return k.capturing == u
}

// Capturing returns the player that is currently capturing the KOTH, if any.
func (k *KOTH) Capturing() (User, bool) {
	return k.capturing, k.capturing != nil
}

// StartCapturing starts the capturing of the KOTH.
func (k *KOTH) StartCapturing(p User, name string) bool {
	if k.capturing != nil || !k.running {
		return false
	}
	k.time = time.Now().Add(300 * time.Second)
	go func() {
		select {
		case <-time.After(300 * time.Second):
			k.capturing = nil
			k.running = false
			if fac, ok := factions.LookupMemberName(p.Name()); ok {
				fac.AddPoints(10)
			}
			Broadcast("koth.captured", k.Name(), name)
			p.AddItemOrDrop(item.NewStack(custom.TripwireHook{}, 3).WithValue("crate-key_KOTH", true).WithCustomName(text.Colourf("<red>KOTH Key</red>")))
		case <-k.cancel:
			k.capturing = nil
			return
		}
	}()
	k.capturing = p
	return true
}

// StopCapturing stops the capturing of the KOTH.
func (k *KOTH) StopCapturing(u User) bool {
	if !k.running {
		return false
	}
	if k.capturing == u {
		k.capturing = nil
		k.cancel <- struct{}{}
		return true
	}
	return false
}

// Time returns the time at which the KOTH will be captured.
func (k *KOTH) Time() time.Time {
	return k.time
}

// Area returns the area of the KOTH.
func (k *KOTH) Area() util.AreaVec2 {
	return k.area
}

// Coordinates returns the coordinates of the KOTH.
func (k *KOTH) Coordinates() mgl64.Vec2 {
	return k.coordinates
}
