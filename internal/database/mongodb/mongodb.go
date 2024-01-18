package mongodb

import (
	"context"
	"errors"
	"gorestapi/internal/database"
	"gorestapi/internal/hero"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongodbParams struct {
	URI string
}

type Mongodb struct {
	Client *mongo.Client
	coll   *mongo.Collection
}

func (m *Mongodb) Connect(ctx context.Context, p database.Params) error {
	// ensure getting mongodb params
	var params MongodbParams
	var ok bool
	if params, ok = p.(MongodbParams); !ok {
		return errors.New("cannot connect to mongodb as did not receive mongodb params type")
	}

	var err error
	m.Client, err = mongo.Connect(ctx, options.Client().ApplyURI(params.URI))
	if err != nil {
		return err
	}

	m.coll = m.Client.Database("test").Collection("heroes")

	// initialize database with sample data
	if err = m.populate(ctx); err != nil {
		return err
	}

	return nil
}

func (m *Mongodb) Disconnect(ctx context.Context) error {
	if err := m.Client.Disconnect(ctx); err != nil {
		return err
	}
	return nil
}

func (m Mongodb) GetById(ctx context.Context, id int) (*hero.Hero, error) {
	var h hero.Hero
	err := m.coll.FindOne(ctx, bson.D{{Key: "id", Value: id}}).Decode(&h)
	if err != nil {
		return nil, err
	}
	return &h, nil
}

func (m Mongodb) GetAll(ctx context.Context) ([]hero.Hero, error) {
	var heroes []hero.Hero

	cursor, err := m.coll.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var hero hero.Hero
		err := cursor.Decode(&hero)
		if err != nil {
			return nil, err
		}
		heroes = append(heroes, hero)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return heroes, nil
}

func (m Mongodb) GetByName(ctx context.Context, name string) ([]hero.Hero, error) {
	var heroes []hero.Hero

	regexPattern := "^" + name
	// i option denotes case-insensitive
	filter := bson.D{{Key: "name", Value: bson.D{{Key: "$regex", Value: regexPattern}, {Key: "$options", Value: "i"}}}}
	cursor, err := m.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var hero hero.Hero
		err := cursor.Decode(&hero)
		if err != nil {
			return nil, err
		}
		heroes = append(heroes, hero)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return heroes, nil
}

func (m Mongodb) UpdateHero(ctx context.Context, h hero.Hero) error {
	filter := bson.D{{Key: "id", Value: h.Id}}
	update := bson.D{{Key: "$set", Value: h}}
	_, err := m.coll.UpdateOne(ctx, filter, update)
	return err
}

func (m Mongodb) DeleteHero(ctx context.Context, id int) error {
	filter := bson.D{{Key: "id", Value: id}}
	_, err := m.coll.DeleteOne(ctx, filter)
	return err
}

func (m Mongodb) AddHero(ctx context.Context, h hero.Hero) error {
	_, err := m.coll.InsertOne(ctx, h)
	return err
}

// populate with sample heroes, check if each exists before populating
func (m Mongodb) populate(ctx context.Context) error {
	opts := options.Update().SetUpsert(true)
	for _, h := range database.TestHeroes {
		filter := bson.D{{Key: "id", Value: h.Id}}
		update := bson.D{{Key: "$setOnInsert", Value: h}}
		_, err := m.coll.UpdateOne(ctx, filter, update, opts)
		if err != nil {
			return err
		}
	}
	return nil
}
