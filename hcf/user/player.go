package user

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
)

// Player is a user that is a player.
func (u *User) Player() *player.Player {
	return u.p
}

// Name returns the name of the user.
func (u *User) Name() string {
	return u.p.Name()
}

// XUID returns the XUID of the user.
func (u *User) XUID() string {
	return u.p.XUID()
}

// Hurt applies damage to the user.
func (u *User) Hurt(dmg float64, src world.DamageSource) {
	u.p.Hurt(dmg, src)
}

// HeldItems returns the items that the user is holding.
func (u *User) HeldItems() (item.Stack, item.Stack) {
	return u.p.HeldItems()
}

// SetHeldItems sets the items that the user is holding.
func (u *User) SetHeldItems(mainHand, offHand item.Stack) {
	u.p.SetHeldItems(mainHand, offHand)
}

// KnockBack applies knockback to the user.
func (u *User) KnockBack(src mgl64.Vec3, force, height float64) {
	u.p.KnockBack(src, force, height)
}

// Sneaking returns whether the user is sneaking.
func (u *User) Sneaking() bool {
	return u.p.Sneaking()
}

// Position returns the position of the user.
func (u *User) Position() mgl64.Vec3 {
	return u.p.Position()
}

// Rotation returns the rotation of the user.
func (u *User) Rotation() cube.Rotation {
	return u.p.Rotation()
}
