package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Flavor struct {
	ID    primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name  string `json:"name,omitempty" bson:"name,omitempty"`
}

var client *mongo.Client
var flavorsCollection *mongo.Collection

func main() {
	// Establecer conexión con MongoDB
	clientOptions := options.Client().ApplyURI("mongodb://mongo:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Obtener la colección de sabores
	flavorsCollection = client.Database("neveria").Collection("flavors")

	// Inicializar el enrutador de la API
	router := mux.NewRouter()

	// Rutas de la API
	router.HandleFunc("/flavors", GetFlavors).Methods("GET")
	router.HandleFunc("/flavor", AddFlavor).Methods("POST")
	router.HandleFunc("/flavors/{id}", DeleteFlavor).Methods("DELETE")

	// Iniciar el servidor
	fmt.Println("Servidor en ejecución en :8081")
	http.ListenAndServe(":8081", router)
}

func GetFlavors(w http.ResponseWriter, r *http.Request) {
	var flavors []Flavor
	cur, err := flavorsCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	defer cur.Close(context.Background())
	for cur.Next(context.TODO()) {
		var flavor Flavor
		err := cur.Decode(&flavor)
		if err != nil {
			log.Fatal(err)
		}
		flavors = append(flavors, flavor)
	}
	json.NewEncoder(w).Encode(flavors)
}

func AddFlavor(w http.ResponseWriter, r *http.Request) {
	var flavor Flavor
	json.NewDecoder(r.Body).Decode(&flavor)

	objID := primitive.NewObjectID()
    flavor.ID = objID

	_, err := flavorsCollection.InsertOne(context.TODO(), flavor)
	if err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(flavor)
}

func DeleteFlavor(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)

    flavorID := params["id"]

    objID, err := primitive.ObjectIDFromHex(flavorID)
    if err != nil {
        http.Error(w, "ID no válido", http.StatusBadRequest)
        return
    }

    _, err = flavorsCollection.DeleteOne(context.TODO(), bson.M{"_id": objID})
    if err != nil {
        http.Error(w, "Error al eliminar el sabor", http.StatusInternalServerError)
        return
    }
    fmt.Fprintf(w, "El sabor con ID %s ha sido eliminado", flavorID)
}

