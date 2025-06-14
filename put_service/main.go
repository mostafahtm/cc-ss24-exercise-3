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
	e.PUT("/api/books/:id", func(c echo.Context) error {
		id := c.Param("id")
		var book BookStore
		if err := c.Bind(&book); err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request body"})
		}
		update := bson.M{"$set": book}
		res, err := coll.UpdateOne(context.TODO(), bson.M{"ID": id}, update)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to update book"})
		}
		if res.MatchedCount == 0 {
			return c.NoContent(http.StatusNoContent)
		}
		return c.JSON(http.StatusOK, echo.Map{"message": "Book updated successfully"})
	})
	e.Logger.Fatal(e.Start(":3030"))
}
