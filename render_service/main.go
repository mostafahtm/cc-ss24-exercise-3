package main

import (
	"context"
	"html/template"
	"io"
	"log"
	"net/http"
	"time"
	"os"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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

type Template struct {
	tmpl *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.tmpl.ExecuteTemplate(w, name, data)
}

// Generic method to perform "SELECT * FROM BOOKS" (if this was SQL, which
// it is not :D ), and then we convert it into an array of map. In Golang, you
// define a map by writing map[<key type>]<value type>{<key>:<value>}.
// interface{} is a special type in Golang, basically a wildcard...
func findAllBooks(coll *mongo.Collection) []map[string]interface{} {
	cursor, err := coll.Find(context.TODO(), bson.D{{}})
	var results []BookStore
	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}

	var ret []map[string]interface{}
	for _, res := range results {
		ret = append(ret, map[string]interface{}{
			"id":          res.ID,
			"title":    res.BookName,
			"author":  res.BookAuthor,
			"pages":   res.BookPages,
			"edition": res.BookEdition,
			"year":   res.BookYear,
		})
	}

	return ret
}

// New structs for authors and years views
type AuthorView struct {
	Author string
	Books  []string
}

type YearView struct {
	Year  string
	Books []BookStore
}

// Function to get all unique authors and their books
func findAllAuthors(coll *mongo.Collection) []AuthorView {
	cursor, err := coll.Find(context.TODO(), bson.D{{}})
	var results []BookStore
	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}

	// Create a map to group books by author
	authorMap := make(map[string][]string)
	for _, book := range results {
		authorMap[book.BookAuthor] = append(authorMap[book.BookAuthor], book.BookName)
	}

	// Convert map to slice of AuthorView
	var authors []AuthorView
	for author, books := range authorMap {
		authors = append(authors, AuthorView{
			Author: author,
			Books:  books,
		})
	}

	return authors
}

// Function to get all books grouped by year
func findAllYears(coll *mongo.Collection) []YearView {
	cursor, err := coll.Find(context.TODO(), bson.D{{}})
	var results []BookStore
	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}

	// Create a map to group books by year
	yearMap := make(map[string][]BookStore)
	for _, book := range results {
		yearMap[book.BookYear] = append(yearMap[book.BookYear], book)
	}

	// Convert map to slice of YearView
	var years []YearView
	for year, books := range yearMap {
		years = append(years, YearView{
			Year:  year,
			Books: books,
		})
	}

	return years
}


func main() {
	e := echo.New()

	e.Renderer = &Template{
		tmpl: template.Must(template.ParseGlob("views/*.html")),
	}
	e.Use(middleware.Logger())
	e.Static("/css", "css")

	// Load environment variable
	uri := os.Getenv("DATABASE_URI")
	if uri == "" {
		log.Fatal("DATABASE_URI not set")
	}

	// MongoDB setup
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("Failed to create MongoDB client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := client.Connect(ctx); err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Fatalf("Error on MongoDB disconnect: %v", err)
		}
	}()

	coll := client.Database("exercise-2").Collection("information")

	// Routes
	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "index", nil)
	})
	e.GET("/books", func(c echo.Context) error {
		books := findAllBooks(coll)
		return c.Render(200, "book-table", books)
	})
	e.GET("/authors", func(c echo.Context) error {
		authors := findAllAuthors(coll)
		return c.Render(200, "authors", authors)
	})
	e.GET("/years", func(c echo.Context) error {
		years := findAllYears(coll)
		return c.Render(200, "years", years)
	})
	e.GET("/search", func(c echo.Context) error {
		return c.Render(200, "search-bar", nil)
	})
	e.GET("/create", func(c echo.Context) error {
		return c.NoContent(http.StatusNoContent)
	})

	// Start server
	e.Logger.Fatal(e.Start(":3030"))
}

