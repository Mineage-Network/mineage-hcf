package main

import (
	"github.com/df-mc/dragonfly/server"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player/chat"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/mineage-network/mineage-hcf/hcf"
	"github.com/mineage-network/mineage-hcf/hcf/backend/lang"
	"github.com/mineage-network/mineage-hcf/hcf/commands"
	"github.com/mineage-network/mineage-hcf/hcf/custom"
	"github.com/restartfu/gophig"
	"github.com/sirupsen/logrus"
	"golang.org/x/text/language"
	"log"
	"net/http"
	"os"

	_ "net/http/pprof"
)

func init() {
	// TODO: ...
	// world.RegisterBlock(custom.Cauldron{})

	for _, b := range []world.Block{
		custom.Hopper{},
		custom.Piston{},
		custom.FlowerPot{},
		custom.MobSpawner{},
	} {
		world.RegisterBlock(b)
	}

	for _, i := range []world.Item{
		custom.Hopper{},
		custom.Piston{},
		custom.FlowerPot{},
		custom.MobSpawner{},
		custom.TripwireHook{},
	} {
		world.RegisterItem(i)
	}
}

// main ...
func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	lang.Register(language.English)

	l := logrus.New()
	l.Formatter = &logrus.TextFormatter{ForceColors: true}
	l.Level = logrus.InfoLevel

	conf := hcf.DefaultConfig()
	conf.Network.Address = ":19133"

	g := gophig.NewGophig("config", "toml", 0644)
	if err := g.GetConf(&conf); os.IsNotExist(err) {
		if err := g.SetConf(conf); err != nil {
			panic(err)
		}
	} else if err != nil {
		panic(err)
	}

	h := hcf.New(conf, l)
	chat.Global.Subscribe(chat.StdoutSubscriber{})

	registerCommands(h.Server())

	if err := h.Start(); err != nil {
		panic(err)
	}
}

// registerCommands ...
func registerCommands(srv *server.Server) {
	for _, c := range []cmd.Command{
		cmd.New("logout", "", nil, commands.Logout{}),
		cmd.New("datareset", "", nil, commands.DataReset{}),
		cmd.New("pvp", "", nil, commands.PvpEnable{}),
		cmd.New("tl", "", nil, commands.TL{}),
		cmd.New("kick", "", nil, commands.Kick{}),
		cmd.New("whisper", "", []string{"w", "tell", "msg"}, commands.Whisper{}),
		cmd.New("reply", "", []string{"r"}, commands.Reply{}),
		cmd.New("gamemode", "", nil, commands.GameMode{}),
		cmd.New("sotw", "", nil,
			commands.SOTWStart{},
			commands.SOTWEnd{},
			commands.SOTWDisable{},
		),
		cmd.New("spawn", "", nil, commands.Spawn{}),
		cmd.New("balance", "", []string{"bal"},
			commands.Balance{},
			commands.BalancePayOnline{},
			commands.BalancePayOffline{}),
		cmd.New("koth", "", nil,
			commands.KothList{},
			commands.KothStart{},
			commands.KothStop{},
		),
		cmd.New("f", "", []string{"faction"},
			commands.FactionCreate{},
			commands.NewFactionInformation(srv),
			commands.FactionDisband{},
			commands.FactionInvite{},
			commands.FactionJoin{},
			commands.NewFactionWho(srv),
			commands.FactionLeave{},
			commands.FactionKick{},
			commands.FactionPromote{},
			commands.FactionDemote{},
			commands.FactionTop{},
			commands.FactionClaim{},
			commands.FactionUnClaim{},
			commands.FactionSetHome{},
			commands.FactionHome{},
			commands.FactionList{},
			commands.FactionUnFocus{},
			commands.FactionFocusPlayer{},
			commands.FactionFocusFaction{},
			commands.FactionChat{},
			commands.FactionWithdraw{},
			commands.FactionDeposit{},
			commands.FactionWithdrawAll{},
			commands.FactionDepositAll{},
			commands.FactionStuck{},
		),
		cmd.New("teleport", "", []string{"tp"},
			commands.TeleportToPos{},
			commands.TeleportTargetsToPos{},
			commands.TeleportTargetsToTarget{},
			commands.TeleportToTarget{},
		),
		cmd.New("knockback", "", []string{"kb"}, commands.KnockBack{}),
	} {
		cmd.Register(c)
	}
}
