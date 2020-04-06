package main

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
