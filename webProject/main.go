package main

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"time"

	pd "github.com/PipedreamHQ/pipedream-go"
	"github.com/gorilla/mux"
	"github.com/mongodb/mongo-go-drive/bson/primitive"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Declarando a variavel que iremos utilizar para conversar com o Client do MongoDB
var client *mongo.Client

// Para visualizar melhor joguei em um html com Must, assim fazendo os testes
var temp = template.Must(template.ParseGlob("templates/*.html"))

// Declarando nossa Struct dos Planetas
type Planeta struct {
	ID      primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Nome    string             `json:"nome,omitempty" bson:"nome,omitempty"`
	Clima   string             `json:"clima,omitempty" bson:"clima,omitempty"`
	Terreno string             `json:"terreno,omitempty" bson:"terreno,omitempty"`
}

// função responsável pela conexão com o html
func index(w http.ResponseWriter, r *http.Request) {
	planetas := []Planeta{
		{Nome: "Planeta A", Clima: "Frio", Terreno: "montanhoso"},
	}
	temp.ExecuteTemplate(w, "Index", planetas)
}

// Função para criar o Planeta
func CreatePlaneta(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var planeta Planeta
	json.NewDecoder(request.Body).Decode(&planeta)
	collection := client.Database("PlanetasDB").Collection("Constelacao")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	resultado, _ := collection.InsertOne(ctx, planeta)
	json.NewEncoder(response).Encode(resultado)

}

// Função para buscar o Planeta no banco de dados
func GetPlaneta(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	var planeta Planeta
	collection := client.Database("PlanetasDB").Collection("Constelacao")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	err := collection.FindOne(ctx, Planeta{ID: id}).Decode(&planeta)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(planeta)
}

// Função para busca da collection constelação com todos os planetas
func GetConstelacao(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var constelacao []Planeta
	collection := client.Database("PlanetasDB").Collection("constelacao")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	cursor, err := collection.Find(ctx, bson.M{})

	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}

	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var planeta Planeta
		cursor.Decode(&planeta)
		constelacao = append(constelacao, planeta)
	}

	if err := cursor.Err(); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(constelacao)

}

func main() {
	fmt.Println("Star Wart API Starting")
	http.HandleFunc("/", index)
	http.ListenAndServe(":8000", nil)

	fmt.Println(pd.Steps)

	data := make(map[string]interface{})
	data["name"] = "Luke"
	pd.Export("data", data)

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	clientOptions := options.Client().ApplyURI()
	client, _ = mongo.Connect(ctx, "mongodb://localhost:27017")
	router := mux.NewRouter()
	router.HandleFunc("/planeta", CreatePlaneta).Methods("POST")
	router.HandleFunc("/constelação", GetConstelacao).Methods("GET")
	router.HandleFunc("/planeta/{id}", GetPlaneta).Methods("GET")
	http.ListenAndServe(":12345", router)

}
