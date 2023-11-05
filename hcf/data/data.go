package data

import (
	"context"
	"github.com/mineage-network/mineage-hcf/hcf/factions"
	"github.com/mineage-network/mineage-hcf/hcf/user"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

// rankData ...
type rankData struct {
	Name       string    `bson:"name"`
	Expires    bool      `bson:"expires"`
	Expiration time.Time `bson:"expiration"`
}

// punishmentData ...
type punishmentData struct {
	Mute user.Punishment `bson:"mute"`
	Ban  user.Punishment `bson:"ban"`
}

// ResetFactions resets all factions.
func ResetFactions() {
	for _, f := range factions.All() {
		_ = DeleteFaction(f)
	}
}

// ResetUsers resets all users.
func ResetUsers() {
	users := db.Collection("users")
	_, _ = users.DeleteMany(context.Background(), bson.M{})
}

// Salt is the database salt.
// TODO: Change this in production.
const Salt = "FUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUUCK"

// db is the Upper database session.
var db *mongo.Database

// init creates the Upper database connection.
func init() {
	// TODO: Add database URI.
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(""))
	if err != nil {
		panic(err)
	}
	db = client.Database("hcf")
}
