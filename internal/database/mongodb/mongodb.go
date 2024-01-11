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

type MongodbParmas struct {
	URI string
}

type Mongodb struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func (m *Mongodb) Connect(p database.Params) error {
	// ensure getting mongodb params
	var params MongodbParmas
	var ok bool
	if params, ok = p.(MongodbParmas); !ok {
		return errors.New("cannot connect to mongodb as did not receive mongodb params type")
	}

	// TODO - fix context
	var err error
	m.client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(params.URI))
	if err != nil {
		return err
	}

	m.coll = m.client.Database("test").Collection("heroes")

	// initialize database with sample data
	if err = m.populate(); err != nil {
		return err
	}

	return nil

}

func (m *Mongodb) Disconnect() error {
	if err := m.client.Disconnect(context.TODO()); err != nil {
		return err
	}
	return nil
}

func (m Mongodb) GetById(id int) (*hero.Hero, error) {
	var h *hero.Hero
	err := m.coll.FindOne(context.TODO(), bson.D{{Key: "id", Value: id}}).Decode(h)
	if err != nil {
		return nil, err
	}
	return h, nil
}

func (m Mongodb) GetAll() ([]hero.Hero, error) {
	var heroes []hero.Hero

	cursor, err := m.coll.Find(context.TODO(), bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
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

func (m Mongodb) GetByName(name string) ([]hero.Hero, error) {
	var heroes []hero.Hero

	filter := bson.D{{Key: "name", Value: name}}
	cursor, err := m.coll.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
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

func (m Mongodb) UpdateHero(h hero.Hero) error {
	filter := bson.D{{Key: "id", Value: h.Id}}
	update := bson.D{{Key: "$set", Value: h}}
	_, err := m.coll.UpdateOne(context.TODO(), filter, update)
	return err
}

func (m Mongodb) DeleteHero(id int) error {
	filter := bson.D{{Key: "id", Value: id}}
	_, err := m.coll.DeleteOne(context.TODO(), filter)
	return err
}

func (m Mongodb) AddHero(h hero.Hero) error {
	_, err := m.coll.InsertOne(context.TODO(), h)
	return err
}

// populate with sample heroes, check if each exists before populating
func (m Mongodb) populate() error {
	opts := options.Update().SetUpsert(true)
	for _, h := range database.TestHeroes {
		filter := bson.D{{Key: "id", Value: h.Id}}
		update := bson.D{{"$setOnInsert", h}}
		_, err := m.coll.UpdateOne(context.TODO(), filter, update, opts)
		if err != nil {
			return err
		}
	}
	return nil
}
