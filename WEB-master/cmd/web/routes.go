package main

import (
	"github.com/gorilla/mux"
	"net/http"
)

func (app *application) routes() http.Handler {

	mux := mux.NewRouter()

	mux.HandleFunc("/horror", app.horror)
	mux.HandleFunc("/comedy", app.comedy)
	mux.HandleFunc("/drama", app.drama)
	mux.HandleFunc("/scifi", app.sciFi)

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/createMovie", app.createMovie)
	mux.HandleFunc("/updateMovie", app.updateMovie).Methods("POST")
	mux.HandleFunc("/deleteMovie", app.deleteMovie).Methods("DELETE")
	mux.HandleFunc("/contacts", app.contacts)
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.PathPrefix("/static/").Handler(http.StripPrefix("/static", fileServer))

	return mux
}
