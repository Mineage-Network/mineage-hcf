package hcf

import (
	"github.com/bedrock-gophers/packethandler"
	"github.com/df-mc/dragonfly/server"
	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/entity"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/item/inventory"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/chat"
	"github.com/df-mc/dragonfly/server/player/playerdb"
	"github.com/df-mc/dragonfly/server/world/mcdb"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/mineage-network/mineage-hcf/hcf/crate"
	"github.com/mineage-network/mineage-hcf/hcf/data"
	ench "github.com/mineage-network/mineage-hcf/hcf/enchantment"
	"github.com/mineage-network/mineage-hcf/hcf/factions"
	"github.com/mineage-network/mineage-hcf/hcf/koth"
	"github.com/mineage-network/mineage-hcf/hcf/sotw"
	"github.com/mineage-network/mineage-hcf/hcf/user"
	"github.com/mineage-network/mineage-hcf/hcf/user/module"
	"github.com/mineage-network/mineage-hcf/hcf/worlds"
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/resource"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"github.com/sirupsen/logrus"
	"net/netip"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var (
	connectionsMu sync.Mutex
	connections   = make(map[netip.Addr]int)
)

// HCF ...
type HCF struct {
	log    *logrus.Logger
	srv    *server.Server
	config server.UserConfig
}

var (
	playerProvider *playerdb.Provider
	worldProvider  *mcdb.DB
)

// New ...
func New(c Config, log *logrus.Logger) *HCF {
	c.Server.Name = text.Colourf("<aqua><b>HCF</b></aqua>")
	c.Server.ShutdownMessage = text.Colourf("<red><b>SERVER RESTART</b></red>")

	conf, _ := c.Config(log)

	h := &HCF{
		log:    log,
		config: c.UserConfig,
	}

	conf.Allower = &allower{hcf: h}
	conf.Entities = entity.DefaultRegistry

	p, err := resource.ReadPath(c.Resources.Folder)
	if err != nil {
		panic(err)
	}
	conf.Resources = append(conf.Resources, p)

	pkt := packethandler.NewPacketListener()
	pkt.Listen(&conf, ":19132", []minecraft.Protocol{minecraft.DefaultProtocol})

	go func() {
		for {
			c, err := pkt.Accept()
			if err != nil {
				return
			}
			c.Handle(&packetHandler{c: c})
		}
	}()

	koth.Broadcast = user.Broadcast
	h.srv = conf.New()
	return h
}

// Server ...
func (h *HCF) Server() *server.Server {
	return h.srv
}

// Start starts the server.
func (h *HCF) Start() error {
	t := time.NewTicker(time.Minute * 10)
	go func() {
		for range t.C {
			for _, u := range user.All() {
				//_ = playerProvider.Save(u.Player().UUID(), u.Player().Data())
				//_ = worldProvider.SavePlayerSpawnPosition(u.Player().UUID(), cube.PosFromVec3(u.Player().Position()))
				_ = data.SaveUser(u)
			}
			for _, f := range factions.All() {
				_ = data.SaveFaction(f)
			}

			// TODO: Spawn entities here with the provided coordinates, workaround for mob spawners.
			//for _, pos := range []mgl64.Vec3{} {
			//	h.srv.World().AddEntity()
			//}

			for i := 5; i > 0; i-- {
				_, _ = chat.Global.WriteString(text.Colourf("<b><grey>»</grey></b> <green>All ground entities will be cleared in %d seconds.</green>", i))
				time.Sleep(time.Second)
			}

			count := 0
			for _, e := range h.srv.World().Entities() {
				if e, ok := e.(*entity.Ent); ok {
					if e.Type() == (entity.ItemType{}) {
						count++
						_ = e.Close()
					}
				}
			}
			_, _ = chat.Global.WriteString(text.Colourf("<b><grey>»</grey></b> <green>Cleared %d entities from the ground.</green>", count))
		}
	}()
	c := make(chan os.Signal, 2)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-c
		sotw.Save()
		for _, f := range factions.All() {
			err := data.SaveFaction(f)
			if err != nil {
				h.log.Errorf("save faction: %v", err)
			}
		}
		for _, e := range h.srv.World().Entities() {
			if _, ok := e.Type().(*player.Type); !ok {
				_ = e.Close()
			}
		}
		if err := h.srv.Close(); err != nil {
			h.log.Errorf("close server: %v", err)
		}
	}()

	h.srv.Listen()

	w := h.srv.World()

	w.Handle(&worlds.Handler{})
	w.StopRaining()
	w.StopThundering()
	w.SetTime(6000)
	w.StopTime()
	w.StopWeatherCycle()
	w.SetSpawn(cube.Pos{0, 66, 0})

	for _, c := range crate.All() {
		b := block.NewChest()
		b.CustomName = text.Colourf("<b>%s <grey>Crate</grey></b>", c.Name())

		*b.Inventory() = *inventory.New(27, nil)

		var items [27]item.Stack
		for i, r := range c.Rewards() {
			st := ench.AddEnchantmentLore(r.Stack())
			st = st.WithLore(append(st.Lore(), text.Colourf("<gold>Chances: %d%%</gold>", r.Chance()))...)
			items[i] = st
		}

		for i, s := range items {
			if s.Empty() {
				items[i] = item.NewStack(block.StainedGlass{Colour: item.ColourLightBlue()}, 1)
			}
		}

		for s, i := range items {
			_ = b.Inventory().SetItem(s, i)
		}

		b.Inventory().Handle(crate.Handler{})

		w.SetBlock(cube.PosFromVec3(c.Position()), b, nil)

		t := entity.NewText(text.Colourf("<b>%s <grey>Crate</grey></b>\n<gold>Click to view crate</gold>\n<gold>Sneak click to use key</gold>", c.Name()), c.PositionMiddle().Add(mgl64.Vec3{0, 1.5, 0}))
		w.AddEntity(t)
	}

	for h.srv.Accept(h.accept) {
		// NOOP
	}
	return nil
}

// accept ...
func (h *HCF) accept(p *player.Player) {
	/*inv := p.Inventory()
	for slot, i := range inv.Slots() {
		for _, sp := range it.SpecialItems() {
			if _, ok := i.Value(sp.Key()); ok {
				_ = inv.SetItem(slot, it.NewSpecialItem(sp, i.Count()))
			}
		}
	}*/

	p.ShowCoordinates()

	u, err := data.LoadUser(p)
	if err != nil {
		h.log.Fatalf("new user: %v", err)
	}

	for _, k := range u.Cooldowns().Kits().All() {
		k.Reset()
	}

	ha := &handler{
		area:   module.NewArea(u),
		class:  module.NewClass(u),
		combat: module.NewCombat(u),
		custom: module.NewCustom(u),

		u:   u,
		p:   p,
		hcf: h,
	}

	ha.HandleJoin()
	p.Handle(ha)
}
