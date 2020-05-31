package controllers

import (
	"log"
	"net/http"
	"os"

	links "not_your_fathers_search_engine/api/controllers/links"
	cockroach_db "not_your_fathers_search_engine/pkg/services/linkgraph/store/cockroach_db"
)

type CockRoachDataBase struct {
	DB *cockroach_db.CockroachDBGraph
}

func ExposeDB() *CockRoachDataBase {
	dbConnectionString := os.Getenv("db_link")
	cdb, err := cockroach_db.NewCockroachDbGraph(dbConnectionString)
	if err != nil {
		log.Fatal("DB NOT CONNECTED: %s", err)
	}
	return &CockRoachDataBase{
		DB: cdb,
	}
}

func (cdb *CockRoachDataBase) SearchLink(w http.ResponseWriter, r *http.Request) {
	links.ExtendDBLinks(cdb.DB).SearchLink(w, r)
}
