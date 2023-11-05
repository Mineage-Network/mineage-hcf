package module

import (
	"github.com/df-mc/dragonfly/server/event"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/mineage-network/mineage-hcf/hcf/user"
	"github.com/mineage-network/mineage-hcf/hcf/util"
	"github.com/mineage-network/mineage-hcf/hcf/util/class"
	"time"
)

// Class ...
type Class struct {
	u *user.User
	player.NopHandler
}

// NewClass returns a new Class module for the user passed.
func NewClass(u *user.User) *Class {
	return &Class{u: u}
}

// HandleItemUse ...
func (c *Class) HandleItemUse(_ *event.Context) {
	u := c.u
	p := u.Player()

	if !class.Compare(u.Class(), class.Bard{}) {
		return
	}

	i, _ := p.HeldItems()
	if e, ok := class.BardEffectFromItem(i.Item()); ok {
		_, ok := u.TimerEnabled()
		if ok {
			return
		}
		ok = u.SOTW()
		if ok {
			return
		}
		if cd := u.Cooldowns().BardItems().Key(i.Item()); cd.Active() {
			u.Message("bard.ability.cooldown", cd.Remaining().Seconds())
			return
		}
		if u.BardEnergy() < 30 {
			u.Message("bard.energy.insufficient")
			return
		}
		u.ReduceBardEnergy(30)
		teammates := u.NearbyTeammates(25)
		for _, m := range teammates {
			c := m.Class()
			m.AddEffectNoLossCond(e, func() bool {
				u, ok := user.Lookup(m.Player())
				return ok && c == u.Class()
			})
		}

		lvl, _ := util.Itor(e.Level())
		u.Message("bard.ability.use", util.EffectName(e), lvl, len(teammates))
		p.SetHeldItems(i.Grow(-1), item.Stack{})
		u.Cooldowns().BardItems().Key(i.Item()).Set(15 * time.Second)
	}
}
