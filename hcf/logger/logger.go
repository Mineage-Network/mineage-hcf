package logger

import (
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/mineage-network/mineage-hcf/hcf/data"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"sync"
	"time"
)

var (
	loggerMu sync.Mutex
	loggers  = map[world.Entity]*Logger{}
)

// LookupXUID looks up a logger by XUID.
func LookupXUID(xuid string) (*Logger, bool) {
	loggerMu.Lock()
	logrs := loggers
	loggerMu.Unlock()

	for _, l := range logrs {
		if l.XUID() == xuid {
			return l, true
		}
	}
	return nil, false
}

// LookupEntity looks up a logger by entity.
func LookupEntity(e world.Entity) (*Logger, bool) {
	loggerMu.Lock()
	l, ok := loggers[e]
	loggerMu.Unlock()
	return l, ok
}

// Logger represents a logger.
type Logger struct {
	xuid  string
	p     *player.Player
	close chan struct{}
}

// NewLogger creates a new logger for the player passed.
func NewLogger(p *player.Player) {
	l := &Logger{
		xuid:  p.XUID(),
		close: make(chan struct{}),
	}
	c := player.New(p.Name(), p.Skin(), p.Position())
	c.SetNameTag(text.Colourf("<dark-red>%s <grey>(Combat Logger)</grey></dark-red>", p.Name()))
	c.Move(mgl64.Vec3{}, p.Rotation().Yaw(), p.Rotation().Pitch())
	for s, i := range p.Inventory().Slots() {
		_ = c.Inventory().SetItem(s, i)
	}
	for s, i := range p.Armour().Slots() {
		switch s {
		case 0:
			c.Armour().SetHelmet(i)
		case 1:
			c.Armour().SetChestplate(i)
		case 2:
			c.Armour().SetLeggings(i)
		case 3:
			c.Armour().SetBoots(i)
		}
	}
	l.p = c
	c.Handle(&handler{l: l})
	p.World().AddEntity(c)

	loggerMu.Lock()
	loggers[c] = l
	loggerMu.Unlock()

	go func() {
		select {
		case <-time.After(30 * time.Second):
			_ = l.Close()
		case <-l.close:
			return
		}
	}()
}

// XUID returns the logger's XUID.
func (l *Logger) XUID() string {
	return l.xuid
}

// Kill kills the logger.
func (l *Logger) Kill() {
	u, _ := data.LoadOfflineUser(l.XUID())
	_ = data.SaveOfflineUser(u.WithLoggerDeath(true))
	_ = l.Close()
}

// Reconnect reconnects the logger.
func (l *Logger) Reconnect() {
	_ = l.Close()
}

// Player returns the logger's player.
func (l *Logger) Player() *player.Player {
	return l.p
}

// Close closes the logger.
func (l *Logger) Close() error {
	loggerMu.Lock()
	delete(loggers, l.p)
	loggerMu.Unlock()

	close(l.close)
	return l.p.Close()
}
