package cockroachdb

import (
	"not_your_fathers_search_engine/pkg/services/linkgraph/graph"
)

// Queries Includes all Queries
type Queries struct {
	upsertLinkQuery       string
	findLinkQuery         string
	linksInPartitionQuery string
	upsertEdgeQuery       string
	edgesInPartitionQuery string
	removeStaleEdgesQuery string
}

// GetQueries Access All Queries
func GetQueries() *Queries {
	return &Queries{
		upsertLinkQuery:       upsertLinkQuery,
		findLinkQuery:         findLinkQuery,
		linksInPartitionQuery: linksInPartitionQuery,
		upsertEdgeQuery:       upsertEdgeQuery,
		edgesInPartitionQuery: edgesInPartitionQuery,
		removeStaleEdgesQuery: removeStaleEdgesQuery,
	}
}

var (
	upsertLinkQuery = `
		INSERT INTO links (url, retrieved_at) VALUES ($1, $2) 
		ON CONFLICT (url) DO UPDATE SET retrieved_at=GREATEST(links.retrieved_at, $2)
		RETURNING id, retrieved_at
	`
	findLinkQuery = `
		SELECT url, retrieved_at FROM links WHERE id=$1
	`
	linksInPartitionQuery = `
		SELECT id, url, retrieved_at FROM links WHERE id >= $1 AND id < $2 AND retrieved_at < $3
	`

	upsertEdgeQuery = `
		INSERT INTO edges (src, dst, updated_at) VALUES ($1, $2, NOW())
		ON CONFLICT (src,dst) DO UPDATE SET updated_at=NOW()
		RETURNING id, updated_at
	`
	edgesInPartitionQuery = `
		SELECT id, src, dst, updated_at FROM edges WHERE src >= $1 AND src < $2 AND updated_at < $3
	`
	removeStaleEdgesQuery = `
		DELETE FROM edges WHERE src=$1 AND updated_at < $2
	`

	// Compile-time check for ensuring CockroachGraph implements Graph.
	_ graph.Graph = (*CockroachGraph)(nil)
)
