package main

import (
	"context"
	"fmt"
	"log"

	"git.onebytedata.com/OneByteDataPlatform/one-byte-module/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb://ioengine:aW9lbmdpbmU@10.10.10.10:27017")

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	// Do stuff
	database := client.Database("IOEngine")
	collection := database.Collection("TagMapper")

	NewRule := models.TagMapper{
		RuleName:     "Test",
		Conditions:   "PatientName|CONTAINS|xyz",
		Replacements: "AccessionNumber|\"CONTA\"",
	}

	insertResult, err := collection.InsertOne(context.TODO(), NewRule)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Test Rule had been inserted: ", insertResult.InsertedID)

	filter := bson.D{bson.E{Key: "RuleName", Value: "Test"}}
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
	}

	defer cursor.Close(context.TODO())
	for cursor.Next(context.TODO()) {
		var rule bson.M
		if err = cursor.Decode(&rule); err != nil {
			log.Fatal(err)
		}
		fmt.Println(rule)
	}

	/*
		deleteResult, err := collection.DeleteMany(context.TODO(), filter)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Deleted %v documents in the tagMapper collection\n", deleteResult.DeletedCount)
	*/
	err = client.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connection to MongoDB closed.")
}
