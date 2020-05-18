package links

import (
	"fmt"
	"net/http"

	crawler "not_your_fathers_search_engine/services/crawler"
)

// WORK IN PROGRESS

type LinkResource struct {
	SearchLink func (http.ResponseWriter, *http.Request)
	UpsertLink func (http.ResponseWriter, *http.Request)
}

func NewLinkResource() *LinkResource {
	return &LinkResource{
		SearchLink: SearchLink,
		UpsertLink: UpsertLink,
	}
}

func SearchLink(w http.ResponseWriter, r *http.Request) {
	// make a call to the link graph in memory service here

	// might need to change it up later and instead only call
	// cockroach or elasticsearch
	crawler.StartCrawlProcess("https://heatchek.io")
	return
}

func UpsertLink(w http.ResponseWriter, r *http.Request) {
	return
}