package knockback

import (
	"fmt"
	"github.com/restartfu/gophig"
	"time"
)

// KnockBack ...
type KnockBack struct {
	conf          *gophig.Gophig
	Force, Height float64
	HitDelay      int
	MaxHeight     float64
}

var (
	Knockback *KnockBack
)

func init() {
	conf := gophig.NewGophig("./knockback", "toml", 0777)
	Knockback = &KnockBack{
		conf:      conf,
		Force:     4,
		Height:    4,
		MaxHeight: 2.5,
		HitDelay:  465,
	}
}

func (kb *KnockBack) RealForce() float64  { return kb.Force / 10 }
func (kb *KnockBack) RealHeight() float64 { return kb.Height / 10 }
func (kb *KnockBack) RealHitDelay() time.Duration {
	return time.Duration(kb.HitDelay * int(time.Millisecond))
}

// SetHitDelay ...
func (kb *KnockBack) SetHitDelay(v int) {
	kb.HitDelay = v
	err := kb.conf.SetConf(kb)
	if err != nil {
		fmt.Errorf("knockback error: %v", err)
	}
}

// SetForce ...
func (kb *KnockBack) SetForce(v float64) {
	kb.Force = v
	err := kb.conf.SetConf(kb)
	if err != nil {
		fmt.Errorf("knockback error: %v", err)
	}
}

// SetHeight ...
func (kb *KnockBack) SetHeight(v float64) {
	kb.Height = v
	err := kb.conf.SetConf(kb)
	if err != nil {
		fmt.Errorf("knockback error: %v", err)
	}
}
