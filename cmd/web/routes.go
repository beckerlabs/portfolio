package main

import "net/http"

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("GET /{$}", app.home)
	mux.HandleFunc("GET /about", app.about)

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
