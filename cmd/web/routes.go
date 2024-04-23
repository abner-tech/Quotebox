package main

import (
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {

	mux := pat.New()
	mux.Get("/", http.HandlerFunc(app.home))
	mux.Get("/quote/create", http.HandlerFunc(app.createQuoteForm))
	mux.Post("/quote/create", http.HandlerFunc(app.createQuote))
	mux.Get("/quote/:id", http.HandlerFunc(app.showQuote))
	//create fileserver to serve our static content
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Get("/static/", http.StripPrefix("/static", fileServer))

	//return app.recoverPanicMiddleware(app.logRequestMiddleware(securityHeaderMiddleware(mux)))

	standardMiddleware := alice.New(
		app.recoverPanicMiddleware,
		app.logRequestMiddleware,
		securityHeaderMiddleware,
	)
	return standardMiddleware.Then(mux)

}
