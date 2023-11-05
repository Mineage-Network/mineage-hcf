package commands

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/mineage-network/mineage-hcf/hcf/data"
	"github.com/mineage-network/mineage-hcf/hcf/user"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"strings"
)

type Balance struct{}

func (Balance) Run(src cmd.Source, out *cmd.Output) {
	p, ok := src.(*player.Player)
	if !ok {
		return
	}
	u, ok := user.Lookup(p)
	if !ok {
		return
	}

	p.Message(text.Colourf("<green>Your balance is $%2.f.</green>", u.Balance()))
}

type BalancePayOnline struct {
	Sub    cmd.SubCommand `cmd:"pay"`
	Target []cmd.Target   `cmd:"target"`
	Amount float64        `cmd:"amount"`
}

func (b BalancePayOnline) Run(src cmd.Source, out *cmd.Output) {
	p, ok := src.(*player.Player)
	if !ok {
		return
	}
	u, ok := user.Lookup(p)
	if !ok {
		return
	}
	t, ok := b.Target[0].(*player.Player)
	if !ok {
		return
	}
	target, ok := user.Lookup(t)
	if !ok {
		return
	}
	if u == target {
		out.Error("You cannot pay yourself.")
		return
	}
	if b.Amount < 0 {
		out.Error("You cannot pay a negative amount.")
		return
	}
	if u.Balance() < b.Amount {
		out.Error("You do not have enough money.")
		return
	}

	u.ReduceBalance(b.Amount)
	target.IncreaseBalance(b.Amount)

	t.Message(text.Colourf("<green>%s has paid you $%2.f.</green>", u.Ranks().Highest().Tag(p.Name()), 0))
	p.Message(text.Colourf("<green>You have paid %s $%2.f.</green>", target.Ranks().Highest().Tag(t.Name()), 0))
}

type BalancePayOffline struct {
	Sub    cmd.SubCommand `cmd:"pay"`
	Target string         `cmd:"target"`
	Amount float64        `cmd:"amount"`
}

func (b BalancePayOffline) Run(src cmd.Source, out *cmd.Output) {
	p, ok := src.(*player.Player)
	if !ok {
		return
	}
	u, ok := user.Lookup(p)
	if !ok {
		return
	}
	if b.Amount < 0 {
		out.Error("You cannot pay a negative amount.")
		return
	}
	if u.Balance() < b.Amount {
		out.Error("You do not have enough money.")
		return
	}

	if strings.EqualFold(b.Target, p.Name()) {
		out.Error("You cannot pay yourself.")
		return
	}

	t, err := data.LoadOfflineUser(b.Target)
	if err != nil {
		out.Error("user has never joined the server")
		return
	}

	u.ReduceBalance(b.Amount)
	t.Balance += b.Amount

	p.Message(text.Colourf("<green>You have paid</green> %s <green>$%2.f.</green>", t.Ranks.Highest().Tag(t.DisplayName()), b.Amount))
}
