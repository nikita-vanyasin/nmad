package main

import (
	"context"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoDBStorage struct {
	c   *mongo.Client
	col *mongo.Collection
}

func NewMongoDBStorage() (Storage, error) {
	opts := options.Client().ApplyURI(CONFIG.MongoDBConnectURL)
	c, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		return nil, errors.WithMessage(err, "could not connect to MongoDB")
	}

	return &mongoDBStorage{
		c:   c,
		col: c.Database("nmad").Collection("nomad_location"),
	}, nil
}

func (s *mongoDBStorage) Close(ctx context.Context) error {
	return s.c.Disconnect(ctx)
}

func (s *mongoDBStorage) Save(ctx context.Context, nl NomadLocation) error {
	_, err := s.col.InsertOne(ctx, nl)
	return err
}

func (s *mongoDBStorage) List(ctx context.Context, chatID string) ([]NomadLocation, error) {
	return s.listAggregate(ctx, chatID)
}

func (s *mongoDBStorage) listAggregate(ctx context.Context, chatID string, stages ...bson.D) ([]NomadLocation, error) {
	orderStage := bson.D{{"$sort", bson.D{{"at", -1}}}}
	matchStage := bson.D{{"$match", bson.D{{"chat_id", chatID}}}}
	groupStage := bson.D{{"$group", bson.D{
		{"_id", "$username"},
		{"lastentry", bson.D{{"$first", "$$ROOT"}}},
	}}}
	replaceRoot := bson.D{{"$replaceRoot", bson.D{{"newRoot", "$lastentry"}}}}
	pipeline := mongo.Pipeline{orderStage, matchStage, groupStage, replaceRoot}
	c, err := s.col.Aggregate(ctx, append(pipeline, stages...))
	if err != nil {
		return nil, errors.WithMessage(err, "Find")
	}
	result := make([]NomadLocation, 0)
	err = c.All(ctx, &result)
	if err != nil {
		return nil, errors.WithMessage(err, "All")
	}
	return result, nil
}

func (s *mongoDBStorage) ListByCountry(ctx context.Context, chatID, country string) ([]NomadLocation, error) {
	matchCountry := bson.D{{"$match", bson.D{{"country", country}}}}
	return s.listAggregate(ctx, chatID, matchCountry)
}

func (s *mongoDBStorage) ListByCity(ctx context.Context, chatID, city string) ([]NomadLocation, error) {
	matchCity := bson.D{{"$match", bson.D{{"city", city}}}}
	return s.listAggregate(ctx, chatID, matchCity)
}
