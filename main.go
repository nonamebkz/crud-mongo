package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Trainer struct {
	Name string `bson:"name"`
	Age  int    `bson:"age"`
	City string `bson:"city"`
	Id   string `json:"id", bson:"_id"`
}

func connectMongo() (*mongo.Client, error) {
	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	fmt.Println("Connected to MongoDB!")

	return client, err
}

func insertData(w http.ResponseWriter, r *http.Request) {
	// menerima parameter
	reqBody, _ := ioutil.ReadAll(r.Body)
	var dataBaru Trainer
	json.Unmarshal(reqBody, &dataBaru)

	client, err := connectMongo()
	if err != nil {
		log.Fatal(err)
	}

	collection := client.Database("test").Collection("trainers")

	// ash := Trainer{"Ash", 10, "Pallet Town"}

	insertResult, err := collection.InsertOne(context.TODO(), dataBaru)
	if err != nil {
		log.Fatal(err)
	}

	// disconnect database
	err = client.Disconnect(context.TODO())

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connection to MongoDB closed.")

	dataResponse := make(map[string]interface{})
	pesan := fmt.Sprintf("Inserted a single document: %s", insertResult.InsertedID)
	dataResponse["pesan"] = pesan
	dataResponse["data"] = dataBaru
	json.NewEncoder(w).Encode(dataResponse)
}

func insertListData(w http.ResponseWriter, r *http.Request) {
	client, err := connectMongo()
	if err != nil {
		log.Fatal(err)
	}

	collection := client.Database("test").Collection("trainers")

	// misty := Trainer{"Misty", 10, "Cerulean City"}
	// brock := Trainer{"Brock", 15, "Pewter City"}

	reqBody, _ := ioutil.ReadAll(r.Body)
	var dataBaru []Trainer
	data := reflect.ValueOf(dataBaru).Interface().([]interface{})
	json.Unmarshal(reqBody, &data)
	// trainers := []interface{}{misty, brock}

	insertManyResult, err := collection.InsertMany(context.TODO(), data)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Inserted multiple documents: ", insertManyResult.InsertedIDs)
	err = client.Disconnect(context.TODO())

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connection to MongoDB closed.")

}

func updateData(w http.ResponseWriter, r *http.Request) {
	client, err := connectMongo()
	if err != nil {
		log.Fatal(err)
	}

	collection := client.Database("test").Collection("trainers")
	reqBody, _ := ioutil.ReadAll(r.Body)
	var dataBaru Trainer
	json.Unmarshal(reqBody, &dataBaru)
	filter := bson.M{"name": dataBaru.Name} //filter bedasarkan nama yg di ambil dari struct Trainer

	// dataBaru := Trainer{Name: "andre", Age: 27, City: "bekasi"}
	update := bson.M{"$set": dataBaru}

	updateResult, err := collection.UpdateMany(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
	}
	// disconnect
	fmt.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)

	err = client.Disconnect(context.TODO())

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connection to MongoDB closed.")

}

func deleteData(w http.ResponseWriter, r *http.Request) {
	client, err := connectMongo()
	if err != nil {
		log.Fatal(err)
	}

	collection := client.Database("test").Collection("trainers")
	reqBody, _ := ioutil.ReadAll(r.Body)
	var dataBaru Trainer
	json.Unmarshal(reqBody, &dataBaru)
	docID, _ := primitive.ObjectIDFromHex(dataBaru.Id) // delete bedasarkan ID
	filter := bson.M{"_id": docID}
	collection.DeleteOne(context.Background(), filter)

	// disconnect
	fmt.Println("Inserted multiple documents: ", insertManyResult.InsertedIDs)
	err = client.Disconnect(context.TODO())

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connection to MongoDB closed.")

}
func foundData(w http.ResponseWriter, r *http.Request) {
	client, err := connectMongo()
	if err != nil {
		log.Fatal(err)
	}

	collection := client.Database("test").Collection("trainers")
	reqBody, _ := ioutil.ReadAll(r.Body)
	var dataBaru Trainer
	json.Unmarshal(reqBody, &dataBaru)
	// docID, _ := primitive.ObjectIDFromHex(dataBaru.Id) // delete bedasarkan ID
	filter := bson.M{"name": dataBaru.Name}

	collection.FindOne(context.Background(), filter).Decode(&dataBaru)
	json.NewEncoder(w).Encode(dataBaru)

	// disconnect
	fmt.Println("Inserted multiple documents: ", insertManyResult.InsertedIDs)
	err = client.Disconnect(context.TODO())

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connection to MongoDB closed.")

}
func main() {
	// router
	http.HandleFunc("/insert", insertData)
	http.HandleFunc("/find", foundData)
	http.HandleFunc("/delete", deleteData)
	http.HandleFunc("/update", updateData)

	log.Fatal(http.ListenAndServe(":8889", nil))
}
