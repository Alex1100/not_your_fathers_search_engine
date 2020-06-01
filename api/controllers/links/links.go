package links

import (
	"bytes"
	"fmt"
	"net/http"

	config "not_your_fathers_search_engine/config"
	crawler "not_your_fathers_search_engine/pkg/services/crawler"
	cockroach "not_your_fathers_search_engine/pkg/services/linkgraph/store/cockroach_db"
	publisher "not_your_fathers_search_engine/pkg/services/publisher"
)

// CockroachDBGraph provides a way for the links package
// methods acess the DB
type CockroachDBGraph struct {
	db *cockroach.CockroachGraph
}

// ExtendDBLinks dependency injection for db access
func ExtendDBLinks(db *cockroach.CockroachGraph) *CockroachDBGraph {
	return &CockroachDBGraph{db: db}
}

// SearchLink grabs links and publishes list of links
// to Google Cloud Pub/Sub
func (db *CockroachDBGraph) SearchLink(w http.ResponseWriter, r *http.Request) {
	fmt.Println("HTTP srcURL IS: ", r.URL.Query().Get("srcURL"))
	srcURL := r.URL.Query().Get("srcURL")
	if len(srcURL) == 0 {
		fmt.Println("NO srcURL PROVIDED")
		return
	}

	links := crawler.StartCrawlProcess(srcURL)
	buf := new(bytes.Buffer)
	appConfig := config.ReadConfig()
	pubSubConfig := appConfig.PubSubConfig
	err := publisher.Publish(buf, pubSubConfig.ProjectID, pubSubConfig.Topics.UpsertLink, links)
	if err != nil {
		fmt.Println("Epic failure, you should probably look into it: ", err)
	}

	// for _, link := range links {
	// 	id, err := uuid.NewRandom()
	// 	if err != nil {
	// 		log.Fatalf("uuid.NewV4() failed with %s\n", err)
	// 	}
	// 	linkGraph := &graph.Link{
	// 		ID:          id,
	// 		URL:         string(link),
	// 		RetrievedAt: time.Now(),
	// 	}

	// 	// This is temporary, eventually I will add a microservice to
	// 	// subscribe to the Google Cloud Pub/Sub and pick off
	// 	// tasks/messages from a queue to upsert into CockroachDB
	// 	err = db.db.UpsertLink(linkGraph)

	// 	if err != nil {
	// 		fmt.Println("Epic failure, you should probably look into it: ", err)
	// 	}
	// }
	// need to save links in cockroach_db within another application
	fmt.Println("LINKS PUBLISHED IN GOOGLE CLOUD PLATFORM'S PUB/SUB `upsert_link` TASK")
}
