// Connects to MongoDB and sets a Stable API version
package main

import (
	"context"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Replace the placeholder with your Atlas connection string
func main() {
	uri := os.Getenv("MONGO_URI")

	// Use the SetServerAPIOptions() method to set the Stable API version to 1
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)

	// Create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), opts)

	if err != nil {
		panic(err)
	}

	// Send a ping to confirm a successful connection
	var resulta bson.M
	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Decode(&resulta); err != nil {
		panic(err)
	}
	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")

	// Extract rows from the Counter collection and count the documents
	collection := client.Database("Counter").Collection("count")

	filter := bson.D{{}} // Empty filter to retrieve all documents
	count, err := collection.CountDocuments(context.TODO(), filter, nil)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Total documents in the Counter collection: %d\n", count)

	row := collection.FindOne(context.TODO(), filter)
	if err := row.Err(); err != nil {
		panic(err)
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	var resultb bson.M
	if err := row.Decode(&resultb); err != nil {
		panic(err)
	}

	fmt.Println(resultb)

}
