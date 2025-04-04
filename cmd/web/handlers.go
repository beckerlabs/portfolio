package main

import (
	"net/http"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// posts, err := app.posts.Latest()
	// if err != nil {
	// 	app.serverError(w, r, err)
	// 	return
	// }

	data := app.newTemplateData(r)
	// data.Posts = posts

	app.render(w, r, http.StatusOK, "home.tmpl.html", data)
}

func (app *application) about(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	app.render(w, r, http.StatusOK, "about.tmpl.html", data)
}

func (app *application) postView(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	// Get the slug from the path and make sure it's not empty.
	slug := r.URL.Path[len("/posts/"):]
	if slug == "" {
		http.NotFound(w, r)
		return
	}

	posts, err := app.posts.LoadMarkdownPosts("./markdown")
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// Get the post with the given slug.
	post, err := app.posts.GetBlogPostBySlug(slug, posts)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	data.BlogPost = post

	postsSidebarData, err := app.posts.LoadPostsSidebarData("./markdown")
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	data.PostsSidebar = postsSidebarData

	sidebarLinks := app.posts.CreateSidebarLinks(post.Headers)
	data.SidebarLinks = sidebarLinks

	// Pass through the content

	// data.Post = post

	app.render(w, r, http.StatusOK, "view.tmpl.html", data)
}
