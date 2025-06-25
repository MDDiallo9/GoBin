package main

import (
	"net/http"

	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	fileserver :=  http.FileServer(http.Dir("./ui/static"))
	mux.Handle("/static/",http.StripPrefix("/static",fileserver))

	mux.HandleFunc("/",app.home)
	mux.HandleFunc("/snippet/view",app.snippetView)
	mux.HandleFunc("/snippet/create",app.snippetCreate)
	
	standard := alice.New(app.recoverPanic,app.logRequest,secureHeaders)

	return standard.Then(mux) // Middleware chain using alice package

	return app.recoverPanic(app.logRequest(secureHeaders(mux))) // Vanilla middleware chain
}
