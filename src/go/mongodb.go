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

var uri string = os.Getenv("MONGO_URI")

type Counter struct {
	Total int32  `bson:"total"`
	Page  string `bson:"page"`
}

func GetPageList() []Counter {
	var result []Counter

	// Use the SetServerAPIOptions() method to set the Stable API version to 1
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)

	// Create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}

	collection := client.Database("Counter").Collection("count")
	cursor, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		panic(err)
	}

	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var elem Counter
		if err = cursor.Decode(&elem); err != nil {
			panic(err)
		}
		result = append(result, elem)
	}

	if err = cursor.Err(); err != nil {
		panic(err)
	}

	if err = client.Disconnect(context.TODO()); err != nil {
		panic(err)
	}

	// log.Println(result)
	return result
}

func UpdateCounter(page string) int32 {
	var oldValue int32

	// Use the SetServerAPIOptions() method to set the Stable API version to 1
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)

	// Create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}

	collection := client.Database("Counter").Collection("count")
	filter := bson.D{{"page", page}}
	var result Counter
	row := collection.FindOne(context.TODO(), filter)
	if row.Err() != nil {
		if row.Err() == mongo.ErrNoDocuments {
			oldValue = 0
		} else {
			panic(row.Err())
		}
	} else {
		row.Decode(&result)
		oldValue = result.Total
	}

	log.Println(result)

	var newValue int32 = oldValue + 1

	updateFilter := bson.D{{"page", page}}
	update := bson.D{{"$set", bson.D{{"total", newValue}, {"page", page}}}}
	updateOpts := options.Update().SetUpsert(true)

	if _, err = collection.UpdateOne(context.TODO(), updateFilter, update, updateOpts); err != nil {
		panic(err)
	}

	if err = client.Disconnect(context.TODO()); err != nil {
		panic(err)
	}

	return newValue
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	var total int32 = UpdateCounter(r.URL.Path)
	var pages []Counter = GetPageList()
	var pageList string = ""
	for _, page := range pages {
		if (page.Page != "") {
			log.Printf("Page: %s, Total: %d\n", page.Page, page.Total)
			pageList += fmt.Sprintf("<li><a href=\"%s\">%s</a>: %d</li>", page.Page, page.Page, page.Total)
		}
	}
	var content string = "<html><body><h1>Counter</h1><p>Visits: %d</p><hr><ul>%s</ul></body></html>"
	log.Printf("Total: %d, Path: %s\n", total, r.URL.Path)
	fmt.Fprintf(w, content, total, pageList)
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
