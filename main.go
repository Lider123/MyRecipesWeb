package main

import (
	"github.com/go-playground/form"
	"html/template"
	"log"
	"math/rand"
	"net/http"
)


func handleHome(w http.ResponseWriter, r *http.Request) {
	ctx := getContext()
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
