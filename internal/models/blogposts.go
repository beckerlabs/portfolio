package models

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

type PostsModelInterface interface {
	CreateSidebarLinks(headers []string) template.HTML
	LoadMarkdownPosts(dir string) ([]BlogPost, error)
	LoadPostsSidebarData(dir string) (PostsSidebarData, error)
	GetSlugs(posts []BlogPost) []string
	GetBlogPostBySlug(slug string, posts []BlogPost) (BlogPost, error)
}

type BlogPost struct {
	Title       string
	Slug        string
	Category    string
	Content     template.HTML
	Created     string
	Description string
	Headers     []string
	Order       int
}

type PostsSidebarData struct {
	Categories []Category
}

type Category struct {
	Name  string
	Pages []BlogPost
	Order int
}

type PostsModel struct{}

func (m *PostsModel) LoadMarkdownPosts(dir string) ([]BlogPost, error) {
	var posts []BlogPost
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".md") {
			path := dir + "/" + file.Name()
			content, err := os.ReadFile(path)
			if err != nil {
				return nil, err
			}

			post, err := parseMarkdownFile(content)
			if err != nil {
				return nil, err
			}

			posts = append(posts, post)
		}
	}

	return posts, nil
}

func parseMarkdownFile(content []byte) (BlogPost, error) {
	sections := strings.SplitN(string(content), "---", 2)
	if len(sections) < 2 {
		return BlogPost{}, errors.New("invalid markdown format")
	}

	metadata := sections[0]
	mdContent := sections[1]

	// deal with rogue \r's
	metadata = strings.ReplaceAll(metadata, "\r", "")
	mdContent = strings.ReplaceAll(mdContent, "\r", "")

	title, slug, category, order, created, description := parseMetadata(metadata)

	htmlContent := mdToHTML([]byte(mdContent))
	headers := extractHeaders([]byte(mdContent))

	return BlogPost{
		Title:       title,
		Slug:        slug,
		Category:    category,
		Order:       order,
		Created:     created,
		Description: description,
		Content:     template.HTML(htmlContent),
		Headers:     headers,
	}, nil
}

func extractHeaders(content []byte) []string {
	var headers []string
	//match only level 2 markdown headers
	re := regexp.MustCompile(`(?m)^##\s+(.*)`)
	matches := re.FindAllSubmatch(content, -1)

	for _, match := range matches {
		// match[1] contains header text without the '##'
		headers = append(headers, string(match[1]))
	}

	return headers
}

func mdToHTML(md []byte) []byte {
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs
	parser := parser.NewWithExtensions(extensions)

	opts := html.RendererOptions{
		Flags: html.CommonFlags | html.HrefTargetBlank,
	}
	renderer := html.NewRenderer(opts)

	doc := parser.Parse(md)

	output := markdown.Render(doc, renderer)

	return output
}

func parseMetadata(metadata string) (
	title string,
	slug string,
	category string,
	order int,
	created string,
	description string,
) {
	re := regexp.MustCompile(`(?m)^(\w+):\s*(.+)`)
	matches := re.FindAllStringSubmatch(metadata, -1)

	metaDataMap := make(map[string]string)
	for _, match := range matches {
		if len(match) == 3 {
			metaDataMap[match[1]] = match[2]
		}
	}

	title = metaDataMap["Title"]
	slug = metaDataMap["Slug"]
	category = metaDataMap["Category"]
	orderStr := metaDataMap["Order"]
	created = metaDataMap["Created"]
	description = metaDataMap["Description"]

	orderStr = strings.TrimSpace(orderStr)
	order, err := strconv.Atoi(orderStr)
	if err != nil {
		log.Printf("Error converting order from string: %v", err)
		order = 9999
	}

	return title, slug, category, order, created, description
}

// Gets all markdown posts and pulls the categories from them to create a SidebarData type.
func (m *PostsModel) LoadPostsSidebarData(dir string) (PostsSidebarData, error) {
	var sidebar PostsSidebarData
	categoriesMap := make(map[string]*Category)

	posts, err := m.LoadMarkdownPosts(dir)
	if err != nil {
		return sidebar, err
	}

	for _, post := range posts {
		if post.Category != "" {
			if _, exists := categoriesMap[post.Category]; !exists {
				categoriesMap[post.Category] = &Category{
					Name:  post.Category,
					Pages: []BlogPost{post},
					Order: post.Order,
				}
			} else {
				categoriesMap[post.Category].Pages = append(categoriesMap[post.Category].Pages, post)
			}
		}
	}

	// convert map to slice
	for _, cat := range categoriesMap {
		sidebar.Categories = append(sidebar.Categories, *cat)
	}

	// sort categories by order
	sort.Slice(sidebar.Categories, func(i, j int) bool {
		return sidebar.Categories[i].Order < sidebar.Categories[j].Order
	})

	return sidebar, nil
}

func (m *PostsModel) CreateSidebarLinks(headers []string) template.HTML {
	var linksHTML string
	for _, header := range headers {
		sanitizedHeader := sanitizeHeaderForID(header)
		link := fmt.Sprintf(`<li><a href="#%s">%s</a></li>`, sanitizedHeader, header)
		linksHTML += link
	}
	return template.HTML(linksHTML)
}

func sanitizeHeaderForID(header string) string {
	// lowercase
	header = strings.ToLower(header)

	// replace spaces with hyphens
	header = strings.ReplaceAll(header, " ", "-")

	// remove any characters that are not alphanumeric or hyphens
	header = regexp.MustCompile(`[^a-z0-9\-]`).ReplaceAllString(header, "")

	return header
}

func dict(values ...interface{}) (map[string]interface{}, error) {
	if len(values)%2 != 0 {
		return nil, errors.New("invalid dict call")
	}
	dict := make(map[string]interface{}, len(values)/2)
	for i := 0; i < len(values); i += 2 {
		key, ok := values[i].(string)
		if !ok {
			return nil, errors.New("dict keys must be strings")
		}
		dict[key] = values[i+1]
	}
	return dict, nil
}

func (m *PostsModel) GetSlugs(posts []BlogPost) []string {
	var slugs []string
	for _, post := range posts {
		slug := post.Slug
		if slug != "" {
			slugs = append(slugs, slug)
		}
	}
	return slugs
}

func (m *PostsModel) GetBlogPostBySlug(slug string, posts []BlogPost) (BlogPost, error) {
	for _, post := range posts {
		if post.Slug == slug {
			return post, nil
		}
	}
	return BlogPost{}, fmt.Errorf("post with slug %s not found", slug)
}
