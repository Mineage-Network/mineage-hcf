package data

import (
	"encoding/hex"
	"fmt"
	"github.com/mineage-network/mineage-hcf/hcf/rank"
	"github.com/mineage-network/mineage-hcf/hcf/user"
	"github.com/mineage-network/mineage-hcf/hcf/util"
	"net/netip"
	"strings"
	"time"

	"github.com/df-mc/dragonfly/server/player"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/sha3"
)

// userData ...
type userData struct {
	XUID               string     `bson:"xuid"`
	Name               string     `bson:"name"`
	DisplayName        string     `bson:"displayName"`
	DeviceID           string     `bson:"did"`
	SelfSignedID       string     `bson:"ssid"`
	Address            string     `bson:"address"`
	Ranks              []rankData `bson:"ranks"`
	FactionCreateDelay int64      `bson:"factionCreateDelay"`
	Stats              struct {
		Kills  int `bson:"kills"`
		Deaths int `bson:"deaths"`
	} `bson:"stats"`
	Lives       *user.Lives          `bson:"lives"`
	Reclaimed   bool                 `bson:"reclaimed"`
	Kits        map[string]time.Time `bson:"kits"`
	LoggerDeath bool                 `bson:"loggerDeath"`
	Punishments punishmentData       `bson:"punishments"`
	Balance     float64              `bson:"balance"`
	Timer       time.Time            `bson:"timer"`
	SOTW        bool                 `bson:"sotw"`
}

// LoadUser ...
func LoadUser(
	p *player.Player,
) (*user.User, error) {
	users := db.Collection("users")
	result, err := users.Find(ctx(), bson.M{"$or": []bson.M{{"name": strings.ToLower(p.Name())}, {"xuid": p.XUID()}}})
	addr, _ := netip.ParseAddrPort(p.Addr().String())

	s := sha3.New256()
	s.Write(addr.Addr().AsSlice())
	s.Write([]byte(Salt))
	address := hex.EncodeToString(s.Sum(nil))

	if result.Err() != nil || err != nil {
		return user.NewUser(p, user.NewRanks(
			[]util.Rank{rank.Player{}},
			map[util.Rank]time.Time{}),
			0, 0, user.DefaultLives(), 0,
			map[string]time.Time{},
			false, address, user.Punishment{}, user.Punishment{}, 250, time.Now(), false), nil
	}

	var data userData
	err = users.FindOne(ctx(), bson.M{"name": strings.ToLower(p.Name())}).Decode(&data)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return user.NewUser(p, user.NewRanks(
				[]util.Rank{rank.Player{}},
				map[util.Rank]time.Time{}),
				0, 0, user.DefaultLives(), 0,
				map[string]time.Time{},
				false, address, user.Punishment{}, user.Punishment{}, 250, time.Now(), false), nil
		}
		return nil, err
	}
	var ranks []util.Rank
	expirations := make(map[util.Rank]time.Time)
	for _, dat := range data.Ranks {
		if dat.Expires && time.Now().After(dat.Expiration) {
			continue
		}
		r, ok := rank.ByName(dat.Name)
		if !ok {
			return nil, fmt.Errorf("load user: rank %s does not exist", dat.Name)
		}
		ranks = append(ranks, r)
		if dat.Expires {
			expirations[r] = dat.Expiration
		}
	}

	u := user.NewUser(p, user.NewRanks(ranks, expirations),
		data.Stats.Kills,
		data.Stats.Deaths,
		&user.Lives{
			Current: 1,
			Default: 1,
		},
		data.FactionCreateDelay,
		data.Kits,
		data.Reclaimed,
		data.Address,
		data.Punishments.Mute,
		data.Punishments.Ban, data.Balance, data.Timer, data.SOTW)
	if data.LoggerDeath {
		p.Teleport(p.World().Spawn().Vec3Middle())
		p.Inventory().Clear()
		p.Armour().Clear()
	}
	return u, nil
}

// SaveUser ...
func SaveUser(u *user.User) error {
	var ranks []rankData
	for _, r := range u.Ranks().All() {
		data := rankData{Name: r.Name()}
		if e, ok := u.Ranks().Expiration(r); ok {
			data.Expiration, data.Expires = e, true
		}
		ranks = append(ranks, data)
	}

	var kits = map[string]time.Time{}
	for k, c := range u.Cooldowns().Kits() {
		kits[k] = time.Now().Add(c.Remaining())
	}

	m, _ := u.Mute()
	data := userData{
		XUID:         u.XUID(),
		DisplayName:  u.Name(),
		Name:         strings.ToLower(u.Name()),
		DeviceID:     u.DeviceID(),
		SelfSignedID: u.SelfSignedID(),
		Address:      u.HashedAddress(),
		Reclaimed:    u.Reclaimed(),
		Ranks:        ranks,
		Lives:        u.Lives(),
		Kits:         kits,
		Punishments: punishmentData{
			Mute: m,
			Ban:  u.Ban(),
		},
		Balance: u.Balance(),
		SOTW:    u.SOTW(),
		Timer:   u.TimerExpiry(),
	}

	d, _ := u.FactionCreateDelay()
	data.FactionCreateDelay = d.UnixMilli()

	data.Stats.Kills = u.Stats().Kills()
	data.Stats.Deaths = u.Stats().Deaths()

	users := db.Collection("users")
	filter := bson.M{"name": bson.M{"$eq": data.Name}}
	update := bson.M{"$set": data}

	res, err := users.UpdateOne(ctx(), filter, update)
	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		_, err = users.InsertOne(ctx(), data)
		return err
	}
	return nil
}
