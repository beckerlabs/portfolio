Title: Integrating a Markdown Blog Into An Existing Website
Description: Why use MySQL when Markdown do trick?
Slug: md-blog-issues
Category: Golang
Order: 1
Created: 2024-04-03
---

I just rewrote the blogging functionality of this website.  As all great programmers do, I stole the code from the internet.  I tried to integrate that work into this existing website, and as you'd expect, ran into some issues.  This blog post covers a few of the things I ran into, how I fixed them, and why I did this in the first place.

## When all you have is a hammer...

I started off my portfolio website by essentially porting over some of the finished app in Alex Edward's [Let's Go](https://lets-go.alexedwards.net/).  This book is a great introduction into using Go for web projects and I highly recommend reading it if you're a nerd like me who reads programming books.

The book uses MySQL as a database backend to store "Snippets".  This works great, but I think it's overkill for writing blog posts on a small portfolio website.  I spend a lot of time writing things down with Markdown, whether through README's, release notes or in Obsidian.  My blog posts aren't terribly fancy, so I figured Markdown was a good medium for these as well.

I went searching for an implementation of this online and quickly found [this blog post](https://fluxsec.red/how-I-developed-a-markdown-blog-with-go-and-HTMX#).  The implementation mostly fit my needs, so I went about integrating that basic premise into my website.

## CSS and design is hard

I've been using this website as a way to teach myself [TailwindCSS](https://tailwindcss.com/).  My eye for visual design kind of ended in 8th grade when I stopped making signatures for forums, but I've enjoyed using Tailwind to style this page.  I like just having to change the look of things inside of the HTML.

With the change in how blog posts are written, the HTML templating also changed.  I'm now passing through the content as a `template.HTML` type and just bringing that whole template into my blogpost view, like so:

```html
{{ define "blogPost" }}
{{ with .BlogPost}}
<h1>{{ .Title }}</h1>
<p>{{ .Description }}</p>
<p>{{ .Created }}</p>
<hr />
<div>
    {{ .Content }}
</div>
{{ end }}
{{ end }}
```

If you've used TailwindCSS before, you're probably wondering how to apply your styling to the content from the blog post.  Since we just output the whole `template.HTML` type, we don't have access to the HTML markup to add classes to.

## Enter child selectors

After being mildly frustrated for a minute or two, I remembered that child selectors exist in CSS so there has to be an analog for it in TailwindCSS.  I went searching through the docs and found the [page that references child selectors](https://tailwindcss.com/docs/hover-focus-and-other-states#child-selectors).  As the documentation notes, this isn't how you would usually add utility classes to elements. You normally want to apply it straight to the element.  But in this case, our only choice is to use child selectors.

The docs spell this out nicely, but I'm essentially selecting children of the div to format and make the website look pretty.  I'm really only interested in things like headings, links and code blocks.  Here's what that looks like in my case.  Warning, you might not like it, but this is what peak Tailwind looks like.

```html
<div class="[&_a]:text-blue-400 [&_a]:hover:text-blue-600 [&_a]:underline [&_a]:decoration-1 [&_a]:decoration-blue-400 
              [&_h2]:text-xl [&_h2]:text-gray-400 
              [&_p]:py-4
              [&_pre]:bg-black [&_pre]:text-green-400">
    {{ .Content }}
</div>
```

## Not as bad as I thought
This actually seems to be the extent of the problems that I had for this task.  I'm still fighting the layout of the blog posts, but that is not necessarily from changing the way things are done.  Just because I'm bad at CSS...