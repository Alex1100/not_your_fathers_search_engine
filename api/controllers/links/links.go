package links

import (
	"bytes"
	"fmt"
	"net/http"

	config "not_your_fathers_search_engine/config"
	publisher "not_your_fathers_search_engine/pkg/services/publisher"
)

// SearchLink grabs links and publishes list of links
// to Google Cloud Pub/Sub
func SearchLink(w http.ResponseWriter, r *http.Request) {
	fmt.Println("HTTP srcURL IS: ", r.URL.Query().Get("srcURL"))
	srcURL := r.URL.Query().Get("srcURL")
	if len(srcURL) == 0 {
		fmt.Println("NO srcURL PROVIDED")
		return
	}

	buf := new(bytes.Buffer)
	appConfig := config.ReadConfig()
	pubSubConfig := appConfig.PubSubConfig
	err := publisher.Publish(buf, pubSubConfig.ProjectID, pubSubConfig.Topics.CrawlFromSource, srcURL)
	if err != nil {
		fmt.Println("Epic failure, you should probably look into it: ", err)
	}

	// need to save links in cockroach_db within another application
	fmt.Println("LINKS PUBLISHED IN GOOGLE CLOUD PLATFORM'S PUB/SUB `upsert_link` TASK")
	// should make a get request to the load balancer / api gateway / router
	// and get the link data from the indexed / stored data in the db instance
	// return info
}
