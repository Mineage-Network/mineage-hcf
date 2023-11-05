package user

import (
	"fmt"
	"github.com/df-mc/dragonfly/server/player/bossbar"
	"github.com/df-mc/dragonfly/server/player/scoreboard"
	"github.com/mineage-network/mineage-hcf/hcf/backend/lang"
	"github.com/mineage-network/mineage-hcf/hcf/koth"
	"github.com/mineage-network/mineage-hcf/hcf/sotw"
	"github.com/mineage-network/mineage-hcf/hcf/util/class"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"math"
	"strings"
	"time"

	_ "unsafe"
)

// startClassTicker starts the user's tickers.
func (u *User) startTicker() {
	t := time.NewTicker(50 * time.Millisecond)
	l := u.Locale()
	for {
		select {
		case <-t.C:
			switch u.Class().(type) {
			case class.Bard:
				if e := u.BardEnergy(); e < 100-0.05 {
					u.bardEnergy.Store(e + 0.05)
				}

				i, _ := u.p.HeldItems()
				if e, ok := class.BardHoldEffectFromItem(i.Item()); ok {
					mates := u.NearbyTeammates(25)
					for _, m := range mates {
						m.AddEffect(e)
					}
				}
			}

			{
				bb := bossbar.New(compass(u.p.Rotation().Yaw()))
				u.p.SendBossBar(bb.WithHealthPercentage(1.0))
			}

			{
				sb := scoreboard.New(lang.Translatef(l, "scoreboard.title"))
				_, _ = sb.WriteString("§r\uE000")
				sb.RemovePadding()
				if fa, ok := u.Faction(); ok {
					if ft, ok := fa.FocusedFaction(); ok {
						_, _ = sb.WriteString(lang.Translatef(l, "scoreboard.focus.name", ft.DisplayName()))
						_, _ = sb.WriteString(lang.Translatef(l, "scoreboard.focus.dtr", ft.DTRString()))
						_, _ = sb.WriteString(lang.Translatef(l, "scoreboard.focus.online", len((Faction{ft}).Users()), ft.PlayerCount()))
						if h, ok := ft.Home(); ok {
							_, _ = sb.WriteString(lang.Translatef(l, "scoreboard.focus.home", h.X(), h.Z()))
						}
						_, _ = sb.WriteString("§3")
					}
				}

				if d, ok := sotw.Running(); ok && u.SOTW() {
					_, _ = sb.WriteString(lang.Translatef(l, "scoreboard.timer.sotw", parseDuration(time.Until(d))))
				}
				if d, ok := u.TimerEnabled(); ok {
					_, _ = sb.WriteString(lang.Translatef(l, "scoreboard.timer.pvp", parseDuration(d)))
				}
				if lo := u.Teleportations().Logout(); !lo.Expired() && lo.Teleporting() {
					_, _ = sb.WriteString(lang.Translatef(l, "scoreboard.logout", time.Until(lo.Expiration()).Seconds()))
				}
				if lo := u.Teleportations().Stuck(); !lo.Expired() && lo.Teleporting() {
					_, _ = sb.WriteString(lang.Translatef(l, "scoreboard.teleportation.stuck", time.Until(lo.Expiration()).Seconds()))
				}
				if tg := u.Tags().Combat(); tg.Active() {
					_, _ = sb.WriteString(lang.Translatef(l, "scoreboard.tag.spawn", tg.Remaining().Seconds()))
				}
				if h := u.Teleportations().Home(); !h.Expired() && h.Teleporting() {
					_, _ = sb.WriteString(lang.Translatef(l, "scoreboard.teleportation.home", time.Until(h.Expiration()).Seconds()))
				}
				if tg := u.Tags().Archer(); tg.Active() {
					_, _ = sb.WriteString(lang.Translatef(l, "scoreboard.tag.archer", tg.Remaining().Seconds()))
				}
				if cd := u.Cooldowns().Pearl(); cd.Active() {
					_, _ = sb.WriteString(lang.Translatef(l, "scoreboard.cooldown.pearl", cd.Remaining().Seconds()))
				}
				if cd := u.Cooldowns().SpecialAbilities(); cd.Active() {
					_, _ = sb.WriteString(lang.Translatef(l, "scoreboard.cooldown.abilities", cd.Remaining().Seconds()))
				}
				if cd := u.Cooldowns().GoldenApple(); cd.Active() {
					_, _ = sb.WriteString(lang.Translatef(l, "scoreboard.cooldown.golden.apple", cd.Remaining().Seconds()))
				}
				if class.Compare(u.Class(), class.Bard{}) {
					_, _ = sb.WriteString(lang.Translatef(l, "scoreboard.bard.energy", u.BardEnergy()))
				}

				if k, ok := koth.Running(); ok {
					t := time.Until(k.Time())
					if _, ok := k.Capturing(); !ok {
						t = time.Minute * 5
					}
					_, _ = sb.WriteString(lang.Translatef(l, "scoreboard.koth.running", k.Name(), parseDuration(t)))
				}

				//if len(sb.Lines()) == 5 {
				//	sb.Remove(4)
				//}

				_, _ = sb.WriteString("  ")
				_, _ = sb.WriteString(lang.Translatef(l, "scoreboard.footer"))
				_, _ = sb.WriteString("\uE000")

				for i, li := range sb.Lines() {
					if !strings.Contains(li, "\uE000") {
						sb.Set(i, " "+li)
					}
				}

				if len(sb.Lines()) > 3 {
					if !compareLines(sb.Lines(), u.scoreboard.Load().Lines()) {
						u.scoreboard.Store(sb)
						u.p.RemoveScoreboard()
						u.p.SendScoreboard(sb)
					}
				} else {
					u.p.RemoveScoreboard()
				}
			}
		case <-u.close:
			t.Stop()
			return
		}
	}
}

// compareLines ...
func compareLines(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, l := range a {
		if l != b[i] {
			return false
		}
	}
	return true
}

// parseDuration ...
func parseDuration(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60
	if hours > 0 {
		return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
	}
	return fmt.Sprintf("%02d:%02d", minutes, seconds)
}

// compass ...
func compass(yaw float64) string {
	wrap := func(angle float64) float64 {
		return angle + math.Ceil(-angle/360)*360
	}

	yaw = wrap(yaw)
	yaw = yaw * 2 / 10
	yaw += 72

	width := 25

	compass := strings.Repeat("|", 72)
	compassChars := []rune(compass)

	compassChars[0] = 'S'
	compassChars[9] = 'S'
	compassChars[9+1] = 'W'
	compassChars[18] = 'W'
	compassChars[18+9] = 'N'
	compassChars[18+9+1] = 'W'
	compassChars[36] = 'N'
	compassChars[36+9] = 'N'
	compassChars[36+9+1] = 'E'
	compassChars[54] = 'E'
	compassChars[54+9] = 'S'
	compassChars[54+9+1] = 'E'

	compass = strings.Repeat(string(compassChars), 3)

	directionInt := int(yaw)
	trimmedCompass := compass[directionInt-int(math.Floor(float64(width/2))) : directionInt-int(math.Floor(float64(width/2)))+width]

	trimmedCompass = strings.ReplaceAll(trimmedCompass, "|", "| ")
	trimmedCompass = strings.ReplaceAll(trimmedCompass, "| N|", "| N |")
	trimmedCompass = strings.ReplaceAll(trimmedCompass, "| NE|", "| NE |")
	trimmedCompass = strings.ReplaceAll(trimmedCompass, "| E|", "| E |")
	trimmedCompass = strings.ReplaceAll(trimmedCompass, "| SE|", "| SE |")
	trimmedCompass = strings.ReplaceAll(trimmedCompass, "| S|", "| S |")
	trimmedCompass = strings.ReplaceAll(trimmedCompass, "| SW|", "| SW |")
	trimmedCompass = strings.ReplaceAll(trimmedCompass, "| W|", "| W |")
	trimmedCompass = strings.ReplaceAll(trimmedCompass, "| NW|", "| NW |")

	return strings.ReplaceAll(trimmedCompass, "|", text.Colourf("<grey>%s</grey>", "|"))
}
