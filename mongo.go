package main

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const database = "reminderbot"
const collection = "channels"

// ChannelRecord is a persistency channel records stored in the database
type ChannelRecord struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	ChannelID int                `bson:"channel_id"`
}

func process(unit func(*mongo.Collection, *context.Context) error) {
	url := "mongodb://" + mongoRouterHost + ":27017"
	client, err := mongo.NewClient(options.Client().ApplyURI(url))

	if err != nil {
		log.Printf("Could not create a mongo client on %s.\n", url)
		log.Fatal(err)
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	err = client.Connect(ctx)

	if err != nil {
		log.Printf("Could not establish a connection to %s.\n", url)
		log.Fatal(err)
	}

	defer client.Disconnect(ctx)

	databaseHandler := client.Database(database)
	collectionHandler := databaseHandler.Collection(collection)

	unit(collectionHandler, &ctx)
}

func listChannels() *[]ChannelRecord {
	var channels []ChannelRecord

	querier := func(collection *mongo.Collection, ctx *context.Context) error {
		cursor, err := collection.Find(*ctx, bson.M{})

		if err != nil {
			return err
		}

		err = cursor.All(*ctx, &channels)

		return err
	}

	process(querier)

	return &channels
}

func registerChannel(id int) {
	inserter := func(collection *mongo.Collection, ctx *context.Context) error {
		channelRecord := ChannelRecord{ChannelID: id}

		_, err := collection.InsertOne(*ctx, channelRecord)

		return err
	}

	process(inserter)
}

func deregisterChannel(id int) {
	deleter := func(collection *mongo.Collection, ctx *context.Context) error {
		_, err := collection.DeleteOne(*ctx, bson.M{"channel_id": id})

		return err
	}

	process(deleter)
}
