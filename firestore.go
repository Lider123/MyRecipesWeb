package main

import (
	"cloud.google.com/go/firestore"
	"context"
	firebase "firebase.google.com/go"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

func getContext() context.Context {
	return context.Background()
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

func addRecipe(recipe Recipe, client *firestore.Client, ctx context.Context) error {
	_, _, err := client.Collection("recipes").Add(ctx, recipe)
	return err
}
