package data

import (
	"fmt"
	"github.com/mineage-network/mineage-hcf/hcf/rank"
	"github.com/mineage-network/mineage-hcf/hcf/user"
	"github.com/mineage-network/mineage-hcf/hcf/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"math/rand"
	"strings"
	"time"
)

// User ...
type User struct {
	// xuid is the xuid of the user.
	xuid string
	// displayName is the display name of the user.
	displayName string
	// name is the name of the user.
	name string
	// deviceID is the device ID of the user.
	deviceID string
	// selfSignedID is the self-signed ID of the user.
	selfSignedID string
	// address is the hashed IP address of the user.
	address string
	// Ranks is the ranks manager of the User.
	Ranks *user.Ranks
	// Stats contains the stats of the user.
	Stats *user.Stats
	// Lives contains the lives of the user.
	Lives *user.Lives
	// loggerDeath is whether the user died from logger.
	loggerDeath bool
	// Mute is the mute information of the User.
	Mute user.Punishment
	// Ban is the ban information of the User.
	Ban user.Punishment
	// Balance is the balance of the User.
	Balance float64
	// Timer is the pvp timer of the User.
	Timer time.Time
	// SOTW is whether the user has their SOTW enabled
	SOTW bool
}

// NewOfflineUser creates a new offline user with the provided data.
func NewOfflineUser(name string) User {
	b := make([]byte, 16)
	for i := range b {
		b[i] = byte(rand.Intn(10))
	}
	return User{
		displayName: strings.ToLower(name),
		name:        strings.ToLower(name),
		Ranks:       user.NewRanks([]util.Rank{rank.Player{}}, map[util.Rank]time.Time{}),
		Stats:       user.DefaultStats(),
		Lives:       user.DefaultLives(),

		Timer: time.Now(),
	}
}

// SearchOfflineUsers searches for offline users using the given conditions.
func SearchOfflineUsers(cond any) ([]User, error) {
	var data []userData
	collection := db.Collection("users")
	cursor, err := collection.Find(ctx(), cond)
	if err != nil {
		return nil, err
	}

	if err := cursor.All(ctx(), &data); err != nil {
		return nil, err
	}

	users := make([]User, 0, len(data))
	for _, d := range data {
		u, _ := parseData(d)
		users = append(users, u)
	}
	return users, nil
}

// LoadOfflineUser loads an offline User from the database by checking XUID and Name. If the user does not exist, an error will be
// returned to the second return value.
func LoadOfflineUser(id string) (User, error) {
	collection := db.Collection("users")
	var data userData
	err := collection.FindOne(ctx(), bson.M{"$or": []bson.M{{"name": strings.ToLower(id)}, {"xuid": id}}}).Decode(&data)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return User{}, fmt.Errorf("user not found")
		}
	}
	return parseData(data)
}

// SaveOfflineUser saves an offline User to the database. If an error occurs, it will be returned.
func SaveOfflineUser(u User) error {
	var ranks []rankData
	for _, r := range u.Ranks.All() {
		data := rankData{Name: r.Name()}
		if e, ok := u.Ranks.Expiration(r); ok {
			data.Expiration, data.Expires = e, true
		}
		ranks = append(ranks, data)
	}

	users := db.Collection("users")
	data := userData{
		Name:         u.Name(),
		DisplayName:  u.DisplayName(),
		DeviceID:     u.DeviceID(),
		SelfSignedID: u.SelfSignedID(),
		LoggerDeath:  u.loggerDeath,
		Ranks:        ranks,
		Lives:        u.Lives,
		Punishments:  punishmentData{Mute: u.Mute, Ban: u.Ban},
		Balance:      u.Balance,
		SOTW:         u.SOTW,
		Timer:        u.Timer,
	}

	filter := bson.M{"xuid": u.XUID()}
	if len(u.XUID()) == 0 {
		filter = bson.M{"name": strings.ToLower(u.Name())}
	}
	update := bson.M{"$set": data}

	res, err := users.UpdateOne(ctx(), filter, update)
	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		_, err = users.InsertOne(ctx(), data)
		return err
	}
	return err
}

// parseData parses userData into an offline User.
func parseData(data userData) (User, error) {
	ranks := make([]util.Rank, 0, len(data.Ranks))
	expirations := make(map[util.Rank]time.Time)
	for _, dat := range data.Ranks {
		r, ok := rank.ByName(dat.Name)
		if !ok {
			return User{}, fmt.Errorf("load user: rank %s does not exist", dat.Name)
		}
		ranks = append(ranks, r)
		if dat.Expires {
			expirations[r] = dat.Expiration
		}
	}
	return User{
		displayName:  data.DisplayName,
		name:         data.Name,
		address:      data.Address,
		deviceID:     data.DeviceID,
		selfSignedID: data.SelfSignedID,

		Ranks:       user.NewRanks(ranks, expirations),
		Lives:       data.Lives,
		loggerDeath: data.LoggerDeath,

		Mute:    data.Punishments.Mute,
		Ban:     data.Punishments.Ban,
		Balance: data.Balance,
		SOTW:    data.SOTW,
		Timer:   data.Timer,
	}, nil
}

// XUID returns the XUID of the offline user.
func (u User) XUID() string {
	return u.xuid
}

// DisplayName returns the display name of the offline user.
func (u User) DisplayName() string {
	return u.displayName
}

// Name returns the name of the offline user.
func (u User) Name() string {
	return u.name
}

// Address returns the hashed and salted ip address of the offline user.
func (u User) Address() string {
	return u.address
}

// DeviceID returns the device ID of the offline user.
func (u User) DeviceID() string {
	return u.deviceID
}

// SelfSignedID returns the self-signed id of the offline user.
func (u User) SelfSignedID() string {
	return u.selfSignedID
}

// WithLoggerDeath sets LoggerDeath of the offline user.
func (u User) WithLoggerDeath(b bool) User {
	u.loggerDeath = b
	return u
}
