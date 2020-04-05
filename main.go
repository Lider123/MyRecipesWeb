package main

import (
	"cloud.google.com/go/firestore"
	"context"
	"firebase.google.com/go"
	"github.com/go-playground/form"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type Ingredient struct {
	Id string `firestore:"id"`
	Name string `firestore:"name"`
	Count string `firestore:"count"`
}

type Recipe struct {
	Author string `firestore:"author"`
	Ingredients []Ingredient `firestore:"ingredients"`
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

func getContext() context.Context {
	return context.Background()
}

func addRecipe(recipe Recipe, client *firestore.Client, ctx context.Context) error {
	_, _, err := client.Collection("recipes").Add(ctx, recipe)
	return err
}

func handleHome(w http.ResponseWriter, r *http.Request) {
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
	tmpl := template.Must(template.ParseFiles("templates/base.html", "templates/home.html"))
	tmpl.Execute(w, data)
}

func makeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func handleRecipes(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Failed to parse form", http.StatusInternalServerError)
			return
		}

		var recipe Recipe
		decoder := form.NewDecoder()
		err = decoder.Decode(&recipe, r.PostForm)
		if err != nil {
			log.Panic(err)
		}
		recipe.UpdatedAt = makeTimestamp()
		for _, v := range recipe.Ingredients {
			v.Id = string(rand.Intn(1e8))
		}

		ctx := getContext()
		firestoreClient, err := getFirestoreClient(ctx)
		if err != nil {
			log.Fatalf("Failed adding alovelace: %v", err)
		} else {
			err = addRecipe(recipe, firestoreClient, ctx)
		}
		if err != nil {
			log.Fatalf("Failed adding recipe: %v", err)
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func handleRecipesAdd(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/base.html", "templates/add-recipe.html"))
	tmpl.Execute(w, nil)
}

func main() {
	http.HandleFunc("/", handleHome)
	http.HandleFunc("/recipes", handleRecipes)
	http.HandleFunc("/recipes/add", handleRecipesAdd)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	log.Fatal(http.ListenAndServe(":8000", nil))
}
