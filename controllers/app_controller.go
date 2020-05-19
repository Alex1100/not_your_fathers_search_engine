package controllers

import (
	// "database/sql"
	"log"
	"net/http"


	links "not_your_fathers_search_engine/controllers/links"
	cockroach_db "not_your_fathers_search_engine/services/linkgraph/store/cockroach_db"
)

// WORK IN PROGRESS
type CockRoachDataBase struct {
	DB *cockroach_db.CockroachDBGraph
}

func ExposeDB() *CockRoachDataBase {
	dsn := "postgresql://root@localhost:26257/not_your_fathers_search_engine?ssl=true&sslmode=require&sslrootcert=certs/ca.crt&sslkey=certs/client.maxroach.key&sslcert=certs/client.maxroach.crt"
	cdb, err := cockroach_db.NewCockroachDbGraph(dsn)
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

func (cdb *CockRoachDataBase) UpsertLink(w http.ResponseWriter, r *http.Request) {
	links.ExtendDBLinks(cdb.DB).UpsertLink(w, r)
}