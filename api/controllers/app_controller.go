package controllers

import (
	"net/http"

	links "not_your_fathers_search_engine/api/controllers/links"
)

// SearchLink finds all links from a given srcLink
func SearchLink(w http.ResponseWriter, r *http.Request) {
	links.SearchLink(w, r)
}
