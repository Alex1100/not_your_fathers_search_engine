package controllers

import (
	"log"
	"net/http"
	"os"

	links "not_your_fathers_search_engine/api/controllers/links"
	cockroachdb "not_your_fathers_search_engine/pkg/services/linkgraph/store/cockroach_db"
)

// CockRoachDataBase connector to embed CockRoachDB
type CockRoachDataBase struct {
	DB *cockroachdb.CockroachGraph
}

// ExposeDB expose db via Dependency Injection
func ExposeDB() *CockRoachDataBase {
	dbConnectionString := os.Getenv("db_link")
	cdb, err := cockroachdb.NewCockroachGraph(dbConnectionString)
	if err != nil {
		log.Fatal("DB NOT CONNECTED: ", err)
	}
	return &CockRoachDataBase{
		DB: cdb,
	}
}

// SearchLink finds all links from a given srcLink
func (cdb *CockRoachDataBase) SearchLink(w http.ResponseWriter, r *http.Request) {
	links.ExtendDBLinks(cdb.DB).SearchLink(w, r)
}
