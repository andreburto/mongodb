// Connects to MongoDB and sets a Stable API version
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func UpdateCounter() int32 {
	uri := os.Getenv("MONGO_URI")

	// Use the SetServerAPIOptions() method to set the Stable API version to 1
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)

	// Create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), opts)

	if err != nil {
		panic(err)
	}

	collection := client.Database("Counter").Collection("count")

	filter := bson.D{{}} // Empty filter to retrieve all documents
	row := collection.FindOne(context.TODO(), filter)
	if err := row.Err(); err != nil {
		panic(err)
	}

	var result bson.M
	if err := row.Decode(&result); err != nil {
		panic(err)
	}

	var oldValue int32 = result["total"].(int32)
	var newValue int32 = oldValue + 1

	updateFilter := bson.D{{"id", result["id"]}}
	update := bson.D{{"$set", bson.D{{"total", newValue}}}}

	_, err = collection.UpdateOne(context.TODO(), updateFilter, update)
	if err != nil {
		panic(err)
	}

	if err = client.Disconnect(context.TODO()); err != nil {
		panic(err)
	}

	//fmt.Printf("Old Total: %d\nNew Total: %d\n", oldValue, newValue)
	return newValue
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	var content string = "<html><body><h1>Counter</h1><p>Visits: %d</p></body></html>"
	var total int32 = UpdateCounter()
	log.Printf("Total: %d, Path: %s\n", total, r.URL.Path)
	fmt.Fprintf(w, content, total)
}

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
}

// Replace the placeholder with your Atlas connection string
func main() {
	http.HandleFunc("/", handleRoot)
	http.HandleFunc("/favicon.ico", faviconHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
