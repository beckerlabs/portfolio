package main

import "net/http"

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("GET /{$}", app.home)
	mux.HandleFunc("GET /about", app.about)
	mux.HandleFunc("GET /posts/{id}", app.postView)
	mux.HandleFunc("GET /posts", app.getAllPosts)

	return mux
}
