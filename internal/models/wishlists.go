package models

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type WishlistsModelInterface interface {
	LoadWishlists(path string) ([]Wishlist, error)
}

// Wishlist is one person's section on the wishlist page.
type Wishlist struct {
	Name  string         `yaml:"name"`
	Items []WishlistItem `yaml:"items"`
}

// WishlistItem is a single thing someone wants, with a link to buy it.
type WishlistItem struct {
	Name  string `yaml:"name"`
	URL   string `yaml:"url"`
	Price string `yaml:"price"`
	Notes string `yaml:"notes"`
}

type wishlistsFile struct {
	People []Wishlist `yaml:"people"`
}

type WishlistsModel struct{}

// LoadWishlists reads and validates the wishlists YAML file at path.
func (m *WishlistsModel) LoadWishlists(path string) ([]Wishlist, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return parseWishlists(content)
}

func parseWishlists(content []byte) ([]Wishlist, error) {
	var file wishlistsFile
	if err := yaml.Unmarshal(content, &file); err != nil {
		return nil, fmt.Errorf("parsing wishlists: %w", err)
	}

	for _, person := range file.People {
		if strings.TrimSpace(person.Name) == "" {
			return nil, fmt.Errorf("wishlist entry is missing a name")
		}
		for i, item := range person.Items {
			if strings.TrimSpace(item.Name) == "" {
				return nil, fmt.Errorf("%s's item #%d is missing a name", person.Name, i+1)
			}
			if !isHTTPURL(item.URL) {
				return nil, fmt.Errorf("%s's item %q has an invalid url %q (must start with http:// or https://)", person.Name, item.Name, item.URL)
			}
		}
	}

	return file.People, nil
}

func isHTTPURL(raw string) bool {
	u, err := url.Parse(raw)
	if err != nil {
		return false
	}
	return (u.Scheme == "http" || u.Scheme == "https") && u.Host != ""
}
