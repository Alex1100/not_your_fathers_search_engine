package links

import (
	"fmt"
	"net/http"

	cockroach "not_your_fathers_search_engine/services/linkgraph/store/cockroach_db"
	crawler "not_your_fathers_search_engine/services/crawler"
)

// WORK IN PROGRESS

type CockroachDBGraph struct {
	db *cockroach.CockroachDBGraph
}

func ExtendDBLinks(db *cockroach.CockroachDBGraph) *CockroachDBGraph {
	return &CockroachDBGraph{db: db}
}

func (db *CockroachDBGraph) SearchLink(w http.ResponseWriter, r *http.Request) {
	// make a call to the link graph in memory service here

	// might need to change it up later and instead only call
	// cockroach or elasticsearch
	links := crawler.StartCrawlProcess("https://heatchek.io")
	fmt.Println("LINKS ARE: ", links)
	return
}

func (db *CockroachDBGraph) UpsertLink(w http.ResponseWriter, r *http.Request) {
	return
}