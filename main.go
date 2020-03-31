package main

import (
	"cloud.google.com/go/firestore"
	"context"
	"firebase.google.com/go"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"html/template"
	"log"
	"net/http"
)

type Ingredient struct {
	Id string `firestore:"id"`
	Name string `firestore:"name"`
	Count string `firestore:"count"`
}

type Recipe struct {
	Author string `firestore:"author"`
	Ingredients []*Ingredient `firestore:"ingredients"`
	Photo string `firestore:"photo"`
	Text string `firestore:"text"`
	Title string `firestore:"title"`
	UpdatedAt int64 `firestore:"updatedAt"`
}

func getFirestoreClient(ctx context.Context) (*firestore.Client, error) {
	sa := option.WithCredentialsFile("service_account.json")
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		return nil, err
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func fetchRecipes(client *firestore.Client, ctx context.Context) ([]Recipe, error) {
	iter := client.Collection("recipes").Documents(ctx)
	var recipes []Recipe
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		var r Recipe
		doc.DataTo(&r)
		recipes = append(recipes, r)
	}
	return recipes, nil
}

func handleMain(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	client, err := getFirestoreClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create Firestore client: %v", err)
	}

	recipes, err := fetchRecipes(client, ctx)
	if err != nil {
		log.Fatalf("Failed to fetch recipes: %v", err)
	}
	defer client.Close()

	data := struct {
		Recipes []Recipe
	}{recipes}
	tmpl := template.Must(template.ParseFiles("templates/main.html"))
	tmpl.Execute(w, data)
}

func main() {
	http.HandleFunc("/", handleMain)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	log.Fatal(http.ListenAndServe(":8000", nil))
}
