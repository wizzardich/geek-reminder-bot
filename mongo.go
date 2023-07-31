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
	ChannelID int64              `bson:"channel_id"`
}

func process(unit func(*mongo.Collection, *context.Context) error) {
	url := "mongodb://" + mongoRouterHost + ":27017"

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(url))

	if err != nil {
		log.Printf("Could not establish a connection to %s.\n", url)
		log.Fatal(err)
	}

	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Printf("Could not properly disconnect from %s.\n", url)
			log.Fatal(err)
		}
	}()

	databaseHandler := client.Database(database)
	collectionHandler := databaseHandler.Collection(collection)

	err = unit(collectionHandler, &ctx)

	if err != nil {
		log.Printf("Could not process the mongodb operation.\n")
		log.Fatal(err)
	}
}

func listChannels() *[]ChannelRecord {
	var channels []ChannelRecord
	log.Println("querying the database for known channels")

	querier := func(collection *mongo.Collection, ctx *context.Context) error {
		cursor, err := collection.Find(*ctx, bson.M{})

		if err != nil {
			return err
		}

		err = cursor.All(*ctx, &channels)

		return err
	}

	process(querier)

	log.Printf("known channels: %v\n", channels)

	return &channels
}

func registerChannel(id int64) {
	log.Printf("registering new channel %d", id)

	inserter := func(collection *mongo.Collection, ctx *context.Context) error {
		channelRecord := ChannelRecord{ChannelID: id}

		_, err := collection.InsertOne(*ctx, channelRecord)

		if err != nil {
			log.Printf("registering returned an error: %s\n", err.Error())
			log.Println("this probably means that channel is already known.")
		}

		return err
	}

	process(inserter)
}

func deregisterChannel(id int64) {
	log.Printf("deregistering channel %d", id)

	deleter := func(collection *mongo.Collection, ctx *context.Context) error {
		_, err := collection.DeleteOne(*ctx, bson.M{"channel_id": id})

		return err
	}

	process(deleter)
}
