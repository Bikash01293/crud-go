package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var collection = ConnectDb()

type Products struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	ProductName string             `json:"productname,omitempty" bson:"productname,omitempty"`
}
//Db connection
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
//Getting all the products
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

//Creating a new product according to the client request.
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


//Updating the product according to the client request.
func UpdateProduct(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	id := r.URL.Query().Get("id")
	fmt.Println("id =>", id)

	ids, _ := primitive.ObjectIDFromHex(id)

	var product Products

	filter := bson.M{"_id": ids}

	_ = json.NewDecoder(r.Body).Decode(&product)

	update := bson.D{
		{"$set", bson.D{
			{"productname", product.ProductName},
		}},
	}

	err := collection.FindOneAndUpdate(context.TODO(), filter, update).Decode(&product)

	if err != nil {
		log.Fatal(err)
	}

	product.ID = ids

	json.NewEncoder(w).Encode(product)

}

func Delete(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		DeleteProduct(w, r)
	default:
		json.NewEncoder(w).Encode("Bad Request !")
	}
}

//Deleting the product according to the client request.
func DeleteProduct(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	id := r.URL.Query().Get("id")
	fmt.Println("id =>", id)

	ids, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Fatal(err, w)
	}
	fmt.Println(ids)

	filter := bson.M{"_id": ids}

	

	var products Products
	cur, err := collection.Find(context.TODO(), filter)
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

		products = product
	}


	deleteResult, err := collection.DeleteOne(context.TODO(), filter)

	if err != nil {
		log.Fatal(err, w)
	}

	fmt.Println(deleteResult)

	w.WriteHeader(http.StatusAccepted)

	json.NewEncoder(w).Encode(products)
}

func main() {
	ConnectDb()
	http.HandleFunc("/products", Home)
	http.HandleFunc("/create", Create)
	http.HandleFunc("/update", Update)
	http.HandleFunc("/delete", Delete)

	fmt.Println("starting server...")
	if err := http.ListenAndServe("localhost:8080", nil); err != nil {
		log.Fatal(err)
	}
}
