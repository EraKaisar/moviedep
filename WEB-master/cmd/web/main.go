package main

import (
	"context"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"tleukanov.net/snippetbox/pkg/models"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	movies   *models.MovieModel

	templateCache map[string]*template.Template
}

func main() {
	dsn := flag.String("dsn", "mongodb+srv://Aphrodi:mugedekter@cluster0.f3qp7um.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0", "MongoDB data source name")
	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	client, err := mongo.NewClient(options.Client().ApplyURI(*dsn))
	if err != nil {
		errorLog.Fatal(err)
	}
	ctx := context.Background()
	err = client.Connect(ctx)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer client.Disconnect(ctx)

	// Check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		errorLog.Fatal(err)
	}
	infoLog.Println("Connected to MongoDB!")

	// Initialize template cache
	templateCache, err := templateCache("./ui/html/")
	if err != nil {
		errorLog.Fatal(err)
	}

	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		movies:        &models.MovieModel{Collection: client.Database("Cluster0").Collection("movies")},
		templateCache: templateCache,
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Starting server on %s", *addr)
	err = srv.ListenAndServe()
	if err != nil {
		errorLog.Fatal(err)
	}
}
