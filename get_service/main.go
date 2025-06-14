package main

import (
	"context"
	"net/http"
	"os"
	"time"
	"fmt"
	
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type BookStore struct {
	ID          string `bson:"ID" json:"id"`
	BookName    string `bson:"BookName" json:"title"`
	BookAuthor  string `bson:"BookAuthor" json:"author"`
	BookEdition string `bson:"BookEdition,omitempty" json:"edition,omitempty"`
	BookPages   string `bson:"BookPages,omitempty" json:"pages,omitempty"`
	BookYear    string `bson:"BookYear,omitempty" json:"year,omitempty"`
}

func findAllBooks(coll *mongo.Collection) []map[string]interface{} {
	cursor, err := coll.Find(context.TODO(), bson.D{{}})
	if err != nil {
		panic(err)
	}
	var results []BookStore
	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}

	var ret []map[string]interface{}
	for _, res := range results {
		ret = append(ret, map[string]interface{}{
			"id":      res.ID,
			"title":   res.BookName,
			"author":  res.BookAuthor,
			"pages":   res.BookPages,
			"edition": res.BookEdition,
			"year":    res.BookYear,
		})
	}
	return ret
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	uri := os.Getenv("DATABASE_URI")
	if uri == "" {
		fmt.Println("DATABASE_URI not set")
		os.Exit(1)
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		fmt.Println("Failed to connect to MongoDB")
		os.Exit(1)
	}
	defer client.Disconnect(ctx)

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		fmt.Println("MongoDB not reachable")
		os.Exit(1)
	}

	coll := client.Database("exercise-3").Collection("information")

	e := echo.New()
	e.GET("/api/books", func(c echo.Context) error {
		books := findAllBooks(coll)
		return c.JSON(http.StatusOK, books)
	})

	e.Logger.Fatal(e.Start(":3030"))
}
