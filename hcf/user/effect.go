package user

import (
	"github.com/df-mc/dragonfly/server/entity/effect"
	"time"
)

// EffectNoLoss ...
type EffectNoLoss struct {
	e effect.Effect
	c chan struct{}
}

// HasEffect returns whether the user has the effect or not.
func (u *User) HasEffect(e effect.Type) bool {
	for _, ef := range u.p.Effects() {
		if ef.Type() == e {
			return true
		}
	}
	return false
}

// HasEffectLevel returns whether the user has the effect or not.
func (u *User) HasEffectLevel(e effect.Type, l int) bool {
	for _, ef := range u.p.Effects() {
		if ef.Type() == e && ef.Level() == l {
			return true
		}
	}
	return false
}

// HasEffectDuration returns whether the user has the effect or not.
func (u *User) HasEffectDuration(e effect.Type, d time.Duration) bool {
	for _, ef := range u.p.Effects() {
		if ef.Type() == e && ef.Duration() >= d {
			return true
		}
	}
	return false
}

// AddEffects adds a list of effects to the user.
func (u *User) AddEffects(effects ...effect.Effect) {
	for _, e := range effects {
		u.AddEffect(e)
	}
}

// AddEffect adds an effect to the user.
func (u *User) AddEffect(e effect.Effect) {
	for _, ef := range u.p.Effects() {
		if ef.Type() == e.Type() {
			if e.Level() > ef.Level() {
				u.p.AddEffect(e)
			} else if ef.Duration() < e.Duration()-time.Second {
				u.p.AddEffect(e)
			}
			return
		}
	}
	u.p.AddEffect(e)
}

// RemoveEffects removes all the provided effects from the user.
func (u *User) RemoveEffects(effects ...effect.Effect) {
	for _, e := range effects {
		for _, ef := range u.p.Effects() {
			if e.Type() == ef.Type() && e.Level() == ef.Level() {
				u.p.RemoveEffect(e.Type())
			}
		}
	}
}

// AddEffectNoLossCond adds an effect to the user.
func (u *User) AddEffectNoLossCond(e effect.Effect, cond func() bool) {
	var oldEffect effect.Effect
	for _, ef := range u.p.Effects() {
		if ef.Type() == e.Type() {
			oldEffect = ef
			break
		}
	}

	u.effectsNoLossMu.Lock()
	eff, ok := u.effectsNoLoss[e.Type()]
	u.effectsNoLossMu.Unlock()

	if ok {
		oldEffect = eff.e
		close(eff.c)
	}

	c := make(chan struct{})
	u.p.AddEffect(e)

	if oldEffect != (effect.Effect{}) {
		u.effectsNoLossMu.Lock()
		u.effectsNoLoss[e.Type()] = struct {
			e effect.Effect
			c chan struct{}
		}{e: oldEffect, c: c}
		u.effectsNoLossMu.Unlock()
	}

	go func() {
		select {
		case <-time.After(e.Duration()):
			if cond() && oldEffect != (effect.Effect{}) {
				u.effectsNoLossMu.Lock()
				delete(u.effectsNoLoss, e.Type())
				u.effectsNoLossMu.Unlock()

				u.p.RemoveEffect(oldEffect.Type())
				u.p.AddEffect(oldEffect)
			}
		case <-c:
			return
		}
	}()
}
