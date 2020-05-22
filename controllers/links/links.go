package links

import (
	"fmt"
	"net/http"
	"github.com/google/uuid"
	"time"
	"log"

	cockroach "not_your_fathers_search_engine/services/linkgraph/store/cockroach_db"
	graph "not_your_fathers_search_engine/services/linkgraph/graph"
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
	fmt.Println("HTTP srcURL IS: ", r.URL.Query().Get("srcURL"))
	srcURL := r.URL.Query().Get("srcURL")
	if len(srcURL) == 0 {
		fmt.Println("NO srcURL PROVIDED")
		return
	}

	links := crawler.StartCrawlProcess(srcURL)
	// linksGraph := make([]*graph.Link, 0)
	for _, link := range links {
		id, err := uuid.NewRandom()
		if err != nil {
			log.Fatalf("uuid.NewV4() failed with %s\n", err)
		}
		linkGraph := &graph.Link{
			ID: id,
			URL: link,
			RetrievedAt: time.Now(),
		}
		err = db.db.UpsertLink(linkGraph)
		
		if err != nil {
			fmt.Println("Epic failure, you should probably look into it: ", err)
		}
	}
	// need to save links in cockroach_db here
	fmt.Println("LINKS STORED/REFRESHED")
	return
}

func (db *CockroachDBGraph) UpsertLink(w http.ResponseWriter, r *http.Request) {
	return
}