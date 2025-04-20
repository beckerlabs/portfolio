package main

import (
	"net/http"

	"web.beckerlabs.dev/ui"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("GET /static/", http.FileServerFS(ui.Files))

	mux.HandleFunc("GET /{$}", app.home)
	mux.HandleFunc("GET /about", app.about)
	mux.HandleFunc("GET /posts", app.getPosts)

	posts, err := app.posts.LoadMarkdownPosts("./markdown")
	if err != nil {
		app.logger.Error(err.Error())
	}

	slugs := app.posts.GetSlugs(posts)
	for _, slug := range slugs {
		mux.HandleFunc("GET /posts/"+slug, app.postView)
	}
	mux.HandleFunc("GET /posts/{slug}", app.postView)

	return mux
}
