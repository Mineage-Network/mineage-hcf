package data

import (
	"context"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/mineage-network/mineage-hcf/hcf/factions"
	"github.com/mineage-network/mineage-hcf/hcf/util"
	"go.mongodb.org/mongo-driver/bson"
)

// ctx ...
func ctx() context.Context {
	return context.Background()
}

// areaData ...
type areaData struct {
	Min mgl64.Vec2 `bson:"min"`
	Max mgl64.Vec2 `bson:"max"`
}

// factionData ...
type factionData struct {
	DisplayName      string              `bson:"displayName"`
	Name             string              `bson:"name"`
	DTR              float64             `bson:"dtr"`
	Home             mgl64.Vec3          `bson:"home"`
	Balance          float64             `bson:"balance"`
	RegenerationTime int                 `bson:"regenerationTime"`
	Points           int                 `bson:"points"`
	Claim            areaData            `bson:"claim"`
	Members          []factionMemberData `bson:"members"`
}

// factionMemberData ...
type factionMemberData struct {
	Name        string `bson:"name"`
	DisplayName string `bson:"displayName"`
	Rank        int    `bson:"rank"`
}

// init ...
func init() {
	var facs []factionData
	cursor, err := db.Collection("factions").Find(context.Background(), bson.M{})
	if err != nil {
		panic(err)
	}
	err = cursor.All(context.Background(), &facs)
	if err != nil {
		panic(err)
	}
	for _, t := range facs {
		members := make([]*factions.Member, 0)
		for _, p := range t.Members {
			var rank factions.Rank
			switch p.Rank {
			case 0:
				rank = factions.RankMember{}
			case 1:
				rank = factions.RankCaptain{}
			case 2:
				rank = factions.RankCoLeader{}
			case 3:
				rank = factions.RankLeader{}
			}
			members = append(members, factions.NewMember(p.Name, p.DisplayName, rank))
		}
		claim := util.AreaVec2{}
		if t.Claim.Min != (mgl64.Vec2{}) && t.Claim.Max != (mgl64.Vec2{}) {
			min := t.Claim.Min
			max := t.Claim.Max
			claim = util.NewAreaVec2(mgl64.Vec2{min.X(), max.X()}, mgl64.Vec2{min.Y(), max.Y()})
		}
		f := factions.NewFaction(t.DisplayName, members, t.DTR, t.Home, t.Balance, t.RegenerationTime, t.Points, claim)
		if len(members) <= 0 {
			_ = DeleteFaction(f)
		}
	}
}

// SaveFaction saves a Faction struct to the database.
func SaveFaction(f *factions.Faction) error {
	facs := db.Collection("factions")

	claim, _ := f.Claim()

	h, _ := f.Home()
	data := factionData{
		DisplayName:      f.DisplayName(),
		Name:             f.Name(),
		DTR:              f.DTR(),
		Balance:          f.Balance(),
		Home:             h,
		RegenerationTime: int(f.RegenerationTime().UnixMilli()),
		Points:           f.Points(),
		Claim: areaData{
			Min: mgl64.Vec2{claim.Min().X(), claim.Max().X()},
			Max: mgl64.Vec2{claim.Min().Y(), claim.Max().Y()},
		},
		Members: []factionMemberData{},
	}
	for _, p := range f.Members() {
		var rank int
		switch p.Rank().(type) {
		case factions.RankMember:
			rank = 0
		case factions.RankCaptain:
			rank = 1
		case factions.RankCoLeader:
			rank = 2
		case factions.RankLeader:
			rank = 3
		}
		d := factionMemberData{
			Name:        p.Name(),
			DisplayName: p.DisplayName(),
			Rank:        rank,
		}
		data.Members = append(data.Members, d)
	}

	filter := bson.M{"name": bson.M{"$eq": f.Name()}}
	update := bson.M{"$set": data}

	res, err := facs.UpdateOne(ctx(), filter, update)
	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		_, err = facs.InsertOne(ctx(), data)
		return err
	}
	return err
}

// DeleteFaction deletes a faction from the database.
func DeleteFaction(f *factions.Faction) error {
	facs := db.Collection("factions")
	defer f.Close()
	filter := bson.M{"name": bson.M{"$eq": f.Name()}}
	_, err := facs.DeleteOne(ctx(), filter)
	return err
}
