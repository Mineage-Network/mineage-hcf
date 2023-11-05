package user

import (
	"github.com/df-mc/atomic"
	"github.com/mineage-network/mineage-hcf/hcf/util"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

// Tags ...
func (u *User) Tags() *Tags {
	return u.tags
}

// Tags ...
type Tags struct {
	archer atomic.Value[*util.Tag]
	combat atomic.Value[*util.Tag]
}

// NewTags ...
func NewTags(u *User) *Tags {
	return &Tags{
		archer: *atomic.NewValue(util.NewTag(
			func(t *util.Tag) {
				u.p.SetNameTag(text.Colourf("<red>%s</red>", u.p.Name()))
			},
			func(t *util.Tag) {
				u.p.SetNameTag(text.Colourf("<green>%s</green>", u.p.Name()))
				u.UpdateState()
			},
		)),
		combat: *atomic.NewValue(util.NewTag(nil, nil)),
	}
}

// All returns all tags.
func (t *Tags) All() []*util.Tag {
	return []*util.Tag{
		t.archer.Load(),
		t.combat.Load(),
	}
}

func (t *Tags) Archer() *util.Tag {
	return t.archer.Load()
}

func (t *Tags) Combat() *util.Tag {
	return t.combat.Load()
}
