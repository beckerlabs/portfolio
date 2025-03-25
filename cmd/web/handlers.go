package main

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"web.beckerlabs.dev/internal/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	posts, err := app.posts.Latest()
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := app.newTemplateData(r)
	data.Posts = posts

	app.render(w, r, http.StatusOK, "home.tmpl.html", data)
}

func (app *application) about(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	app.render(w, r, http.StatusOK, "about.tmpl.html", data)
}

func (app *application) postView(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.NotFound(w, r)
		return
	}

	intId, err := strconv.Atoi(id)
	if err != nil || intId < 1 {
		http.NotFound(w, r)
		return
	}
	post, err := app.posts.Get(intId)
	post.StructuredContent = strings.Split(post.Content, "\n")

	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			log.Printf("Error retrieving post: %v", err)
			app.serverError(w, r, err)
		}
		return
	}

	data := app.newTemplateData(r)
	data.Post = post

	app.render(w, r, http.StatusOK, "view.tmpl.html", data)
}

func (app *application) getAllPosts(w http.ResponseWriter, r *http.Request) {
	posts, err := app.posts.GetAll()
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := app.newTemplateData(r)
	data.Posts = posts

	app.render(w, r, http.StatusOK, "posts.tmpl.html", data)
}
