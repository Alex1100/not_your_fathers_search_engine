package links

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"

	config "not_your_fathers_search_engine/config"
	crawler "not_your_fathers_search_engine/pkg/services/crawler"
	graph "not_your_fathers_search_engine/pkg/services/linkgraph/graph"
	cockroach "not_your_fathers_search_engine/pkg/services/linkgraph/store/cockroach_db"
	publisher "not_your_fathers_search_engine/pkg/services/publisher"
)

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

	for _, link := range links {
		id, err := uuid.NewRandom()
		if err != nil {
			log.Fatalf("uuid.NewV4() failed with %s\n", err)
		}
		linkGraph := &graph.Link{
			ID:          id,
			URL:         link,
			RetrievedAt: time.Now(),
		}

		buf := new(bytes.Buffer)
		appConfig := config.ReadConfig()
		pubSubConfig := appConfig.PubSubConfig
		err = publisher.Publish(buf, pubSubConfig.ProjectID, pubSubConfig.Topics.UpsertLink, link)

		// This is temporary, eventually I will add a microservice to
		// subscribe to the Google Cloud Pub/Sub and pick off
		// tasks/messages from a queue to upsert into CockroachDB
		err = db.db.UpsertLink(linkGraph)

		if err != nil {
			fmt.Println("Epic failure, you should probably look into it: ", err)
		}
	}
	// need to save links in cockroach_db here
	fmt.Println("LINKS PUBLISHED IN GOOGLE CLOUD PLATFORM'S PUB/SUB `upsert_link` TASK")
	return
}

func (db *CockroachDBGraph) UpsertLink(w http.ResponseWriter, r *http.Request) {
	return
}
