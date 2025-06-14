package main

import (
	"context"
	"net/http"
	"os"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type BookStore struct {
	ID          string `bson:"ID" json:"id"`
	BookName    string `bson:"BookName" json:"title"`
	BookAuthor  string `bson:"BookAuthor" json:"author"`
	BookEdition string `bson:"BookEdition,omitempty" json:"edition,omitempty"`
	BookPages   string `bson:"BookPages,omitempty" json:"pages,omitempty"`
	BookYear    string `bson:"BookYear,omitempty" json:"year,omitempty"`
}

func main() {
	e := echo.New()
	uri := os.Getenv("DATABASE_URI")
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	coll := client.Database("exercise-3").Collection("information")

	e.POST("/api/books", func(c echo.Context) error {
		var newBook BookStore
		if err := c.Bind(&newBook); err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request body"})
		}
		cursor, err := coll.Find(context.TODO(), bson.M{"ID": newBook.ID})
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Database error"})
		}
		var existing []BookStore
		if err = cursor.All(context.TODO(), &existing); err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Database error"})
		}
		if len(existing) > 0 {
			return c.JSON(http.StatusOK, echo.Map{"message": "Duplicate book entry, not inserted"})
		}
		_, err = coll.InsertOne(context.TODO(), newBook)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to insert book"})
		}
		return c.JSON(http.StatusCreated, newBook)
	})
	e.Logger.Fatal(e.Start(":3030"))
}