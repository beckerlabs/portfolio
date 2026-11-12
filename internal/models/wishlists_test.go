package models

import (
	"strings"
	"testing"
)

func TestParseWishlists(t *testing.T) {
	content := []byte(`
people:
  - name: Jeremy
    items:
      - name: Dumbbells
        url: https://www.example.com/dumbbells
        price: "$299.99"
        notes: The 5-52.5 lb set
  - name: Tyson
    items: []
  - name: Alyssa
`)

	lists, err := parseWishlists(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(lists) != 3 {
		t.Fatalf("got %d people, want 3", len(lists))
	}
	if lists[0].Name != "Jeremy" || lists[1].Name != "Tyson" || lists[2].Name != "Alyssa" {
		t.Errorf("people out of order: %+v", lists)
	}

	item := lists[0].Items[0]
	if item.Name != "Dumbbells" || item.URL != "https://www.example.com/dumbbells" ||
		item.Price != "$299.99" || item.Notes != "The 5-52.5 lb set" {
		t.Errorf("unexpected item: %+v", item)
	}
	if len(lists[2].Items) != 0 {
		t.Errorf("Alyssa should have no items, got %d", len(lists[2].Items))
	}
}

func TestParseWishlistsRejectsBadData(t *testing.T) {
	tests := map[string]struct {
		yaml    string
		wantErr string
	}{
		"missing person name": {
			yaml:    "people:\n  - items: []\n",
			wantErr: "missing a name",
		},
		"missing item name": {
			yaml:    "people:\n  - name: Tyson\n    items:\n      - url: https://example.com\n",
			wantErr: "item #1 is missing a name",
		},
		"javascript url": {
			yaml:    "people:\n  - name: Tyson\n    items:\n      - name: Bad\n        url: javascript:alert(1)\n",
			wantErr: "invalid url",
		},
		"missing url": {
			yaml:    "people:\n  - name: Tyson\n    items:\n      - name: No link\n",
			wantErr: "invalid url",
		},
		"malformed yaml": {
			yaml:    "people: [\n",
			wantErr: "parsing wishlists",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			_, err := parseWishlists([]byte(tc.yaml))
			if err == nil || !strings.Contains(err.Error(), tc.wantErr) {
				t.Errorf("got error %v, want one containing %q", err, tc.wantErr)
			}
		})
	}
}

func TestLoadWishlistsRepoFile(t *testing.T) {
	m := &WishlistsModel{}
	lists, err := m.LoadWishlists("../../wishlists/wishlists.yaml")
	if err != nil {
		t.Fatalf("the committed wishlists.yaml is invalid: %v", err)
	}
	if len(lists) == 0 {
		t.Error("expected at least one person in wishlists.yaml")
	}
}
