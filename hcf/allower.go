package hcf

import (
	"encoding/hex"
	"github.com/hako/durafmt"
	"github.com/mineage-network/mineage-hcf/hcf/backend/lang"
	"github.com/mineage-network/mineage-hcf/hcf/data"
	"github.com/sandertv/gophertunnel/minecraft/protocol/login"
	"github.com/unickorn/strutils"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/sha3"
	"golang.org/x/text/language"
	"net"
	"net/netip"
	"strings"
)

// allower ensures that all players who join are whitelisted if whitelisting is enabled.
type allower struct {
	hcf *HCF
}

// Allow ...
func (a *allower) Allow(ip net.Addr, identity login.IdentityData, client login.ClientData) (string, bool) {
	addr, _ := netip.ParseAddrPort(ip.String())

	s := sha3.New256()
	s.Write(addr.Addr().AsSlice())
	s.Write([]byte(data.Salt))
	address := hex.EncodeToString(s.Sum(nil))

	l, _ := language.Parse(strings.Replace(client.LanguageCode, "_", "-", 1))
	filter := bson.M{
		"$or": bson.A{
			bson.M{"did": client.DeviceID},
			bson.M{"ssid": client.SelfSignedID},
			bson.M{"xuid": identity.XUID},
			bson.M{"address": address},
		},
	}
	users, err := data.SearchOfflineUsers(filter)
	if err != nil {
		panic(err)
	}
	for _, u := range users {
		if !u.Ban.Expired() {
			reason := strings.TrimSpace(u.Ban.Reason)
			if u.Ban.Permanent {
				description := lang.Translatef(l, "user.blacklist.description", reason)
				if u.XUID() == identity.XUID {
					return strutils.CenterLine(lang.Translatef(l, "user.blacklist.header") + "\n" + description), false
				}
				return strutils.CenterLine(lang.Translatef(l, "user.blacklist.header.alt") + "\n" + description), false
			}
			description := lang.Translatef(l, "user.ban.description", reason, durafmt.ParseShort(u.Ban.Remaining()))
			if u.XUID() == identity.XUID {
				return strutils.CenterLine(lang.Translatef(l, "user.ban.header") + "\n" + description), false
			}
			return strutils.CenterLine(lang.Translatef(l, "user.ban.header.alt") + "\n" + description), false
		}
	}
	return "", true
}
