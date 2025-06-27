package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	fileserver :=  http.FileServer(http.Dir("./ui/static"))
	mux.Handle("/static/",http.StripPrefix("/static",fileserver))

	router := httprouter.New()

	router.NotFound = http.HandlerFunc(func( w http.ResponseWriter, r *http.Request){
		app.notFound(w)
	})

	router.Handler(http.MethodGet,"/static/*filepath",http.StripPrefix("/static",fileserver))


	router.HandlerFunc(http.MethodGet,"/",app.home)
	router.HandlerFunc(http.MethodGet,"/snippet/view/:id",app.snippetView)
	router.HandlerFunc(http.MethodGet,"/snippet/create",app.snippetCreate)
	router.HandlerFunc(http.MethodPost,"/snippet/create",app.snippetCreatePost)
	router.HandlerFunc(http.MethodGet,"/snippet/delete/:id",app.snippetDelete)
	
	standard := alice.New(app.recoverPanic,app.logRequest,secureHeaders)

	return standard.Then(router) // Middleware chain using alice package

	return app.recoverPanic(app.logRequest(secureHeaders(mux))) // Vanilla middleware chain
}
