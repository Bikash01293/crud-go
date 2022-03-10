package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var collection = ConnectDb()

type Products struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	ProductName string             `json:"productname,omitempty" bson:"productname,omitempty"`
}

func ConnectDb() *mongo.Collection {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}
	fmt.Println("db connected !")
	collection := client.Database("dbname").Collection("cname")
	return collection
}

func Home(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		GetHome(w, r)
	default:
		json.NewEncoder(w).Encode("Bad Request !")
	}
}

func GetHome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	var products []Products
	cur, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {

		var product Products
		err := cur.Decode(&product)
		if err != nil {
			log.Fatal(err)
		}

		products = append(products, product)
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(products)
}

func Create(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		CreateProduct(w, r)
	default:
		json.NewEncoder(w).Encode("Bad request !")
	}
}

func CreateProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	var product Products
	_ = json.NewDecoder(r.Body).Decode(&product)

	result, err := collection.InsertOne(context.TODO(), product)
	if err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(result)
}

func Update(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		UpdateProduct(w, r)
	default:
		json.NewEncoder(w).Encode("Bad Request !")
	}
}

func UpdateProduct(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GET params were:", r.URL.Query())
	params := r.URL.Query()
	json.NewEncoder(w).Encode(params)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		DeleteProduct(w, r)
	default:
		json.NewEncoder(w).Encode("Bad Request !")
	}
}

func DeleteProduct(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GET params were:", r.URL.Query())
	params := r.URL.Query()
	id := r.URL.Query().Get("id")
	fmt.Println("id =>", id)
	json.NewEncoder(w).Encode(params)
}

func main() {
	ConnectDb()
	http.HandleFunc("/home", Home)
	http.HandleFunc("/create", Create)
	http.HandleFunc("/update/{id}", Update)
	http.HandleFunc("/delete/{id}", Delete)

	fmt.Println("starting server...")
	if err := http.ListenAndServe("localhost:8080", nil); err != nil {
		log.Fatal(err)
	}
}

