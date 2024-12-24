package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"github.com/graphql-go/graphql"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
    "github.com/rs/cors"
)

// State struct
type State struct {
	Name string `bson:"name" json:"name"`
}

// MongoDB collection
var collection *mongo.Collection

// Define GraphQL schema
var stateType = graphql.NewObject(graphql.ObjectConfig{
	Name: "State",
	Fields: graphql.Fields{
		"name": &graphql.Field{Type: graphql.String},
	},
})

var schema, _ = graphql.NewSchema(graphql.SchemaConfig{
	Query: graphql.NewObject(graphql.ObjectConfig{
		Name: "RootQuery",
		Fields: graphql.Fields{
			"states": &graphql.Field{
				Type: graphql.NewList(stateType),
				Args: graphql.FieldConfigArgument{
					"filter": &graphql.ArgumentConfig{Type: graphql.String},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					filter, _ := p.Args["filter"].(string)
					ctx := context.Background()

					query := bson.M{}
					if filter != "" {
						// filter logic 
						query = bson.M{"name": bson.M{"$regex": "^" + strings.ToLower(filter), "$options": "i"}}
					}

					cursor, err := collection.Find(ctx, query)
					if err != nil {
						return nil, err
					}
					defer cursor.Close(ctx)

					var states []State
					for cursor.Next(ctx) {
						var state State
						cursor.Decode(&state)
						states = append(states, state)
					}
					return states, nil
				},
			},
		},
	}),
})

func main() {
	log.Println("Connecting to DB ")
	uri := "mongodb://admin:admin123@localhost:27017"
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	log.Println("Successfully connected to DB ")

	defer client.Disconnect(context.Background())
	collection = client.Database("us_data").Collection("states")

	addStatesToDb()
	// Create the CORS handler
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:4200"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}).Handler

http.HandleFunc("/graphql", func(w http.ResponseWriter, r *http.Request) {
	// Read the request body
	var request map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Println("Error decoding request body:", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Extract the query string
	query, ok := request["query"].(string)

	if !ok || query == "" {
		log.Println("Query missing or empty")
		http.Error(w, "Query missing or empty", http.StatusBadRequest)
		return
	}
	// Extract the variables
	variables, ok := request["variables"].(map[string]interface{})
	if !ok || variables["filter"] == nil {
	log.Println("Missing filter variable in request")
	http.Error(w, "Missing filter variable in request", http.StatusBadRequest)
	return
	}

	log.Printf("Received filter: %v", variables["filter"])


	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
		VariableValues: variables,
	})
	json.NewEncoder(w).Encode(result)
})

log.Println("Server running at http://localhost:8080/graphql")
http.ListenAndServe(":8080", corsHandler(http.DefaultServeMux))
}

func addStatesToDb() {
	states := []State{
		{Name: "Alabama"}, {Name: "Alaska"}, {Name: "Arizona"}, {Name: "Arkansas"},
		{Name: "California"}, {Name: "Colorado"}, {Name: "Connecticut"}, {Name: "Delaware"},
		{Name: "Florida"}, {Name: "Georgia"}, {Name: "Hawaii"}, {Name: "Idaho"},
		{Name: "Illinois"}, {Name: "Indiana"}, {Name: "Iowa"}, {Name: "Kansas"},
		{Name: "Kentucky"}, {Name: "Louisiana"}, {Name: "Maine"}, {Name: "Maryland"},
		{Name: "Massachusetts"}, {Name: "Michigan"}, {Name: "Minnesota"}, {Name: "Mississippi"},
		{Name: "Missouri"}, {Name: "Montana"}, {Name: "Nebraska"}, {Name: "Nevada"},
		{Name: "New Hampshire"}, {Name: "New Jersey"}, {Name: "New Mexico"}, {Name: "New York"},
		{Name: "North Carolina"}, {Name: "North Dakota"}, {Name: "Ohio"}, {Name: "Oklahoma"},
		{Name: "Oregon"}, {Name: "Pennsylvania"}, {Name: "Rhode Island"}, {Name: "South Carolina"},
		{Name: "South Dakota"}, {Name: "Tennessee"}, {Name: "Texas"}, {Name: "Utah"},
		{Name: "Vermont"}, {Name: "Virginia"}, {Name: "Washington"}, {Name: "West Virginia"},
		{Name: "Wisconsin"}, {Name: "Wyoming"}, {Name: "American Samoa"}, {Name: "Guam"},
		{Name: "Northern Mariana Islands"}, {Name: "Puerto Rico"}, {Name: "U.S. Virgin Islands"},
	}
	collection.DeleteMany(context.Background(), bson.M{})
	for _, state := range states {
		collection.InsertOne(context.Background(), state)
	}
}