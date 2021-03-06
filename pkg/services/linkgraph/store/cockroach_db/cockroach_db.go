package cockroachdb

import (
	"database/sql"
	"fmt"
	"time"

	"not_your_fathers_search_engine/pkg/services/linkgraph/graph"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"golang.org/x/xerrors"
)

// CockroachGraph implements a graph that persists its links and edges to a
// cockroachdb instance.
type CockroachGraph struct {
	DB *sql.DB
}

var dbQueries *Queries = GetQueries()

// NewCockroachGraph returns a CockroachGraph instance that connects to the cockroachdb
// instance specified by dsn.
func NewCockroachGraph(dsn string) (*CockroachGraph, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	return &CockroachGraph{DB: db}, nil
}

// Close terminates the connection to the backing cockroachdb instance.
func (c *CockroachGraph) Close() {
	err := c.DB.Close()
	if err != nil {
		fmt.Println("ERR CLOSING DB: ", err)
	}
}

// UpsertLink creates a new link or updates an existing link.
func (c *CockroachGraph) UpsertLink(link *graph.Link) error {
	row := c.DB.QueryRow(dbQueries.upsertLinkQuery, link.URL, link.RetrievedAt.UTC())
	fmt.Println("ROW IS: ", row)
	if err := row.Scan(&link.ID, &link.RetrievedAt); err != nil {
		return xerrors.Errorf("upsert link: %w", err)
	}

	link.RetrievedAt = link.RetrievedAt.UTC()
	return nil
}

// FindLink looks up a link by its ID.
func (c *CockroachGraph) FindLink(id uuid.UUID) (*graph.Link, error) {
	row := c.DB.QueryRow(dbQueries.findLinkQuery, id)
	link := &graph.Link{ID: id}
	if err := row.Scan(&link.URL, &link.RetrievedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, xerrors.Errorf("find link: %w", graph.ErrNotFound)
		}

		return nil, xerrors.Errorf("find link: %w", err)
	}

	link.RetrievedAt = link.RetrievedAt.UTC()
	return link, nil
}

// Links returns an iterator for the set of links whose IDs belong to the
// [fromID, toID) range and were last accessed before the provided value.
func (c *CockroachGraph) Links(fromID, toID uuid.UUID, accessedBefore time.Time) (graph.LinkIterator, error) {
	rows, err := c.DB.Query(dbQueries.linksInPartitionQuery, fromID, toID, accessedBefore.UTC())
	if err != nil {
		return nil, xerrors.Errorf("links: %w", err)
	}

	return &linkIterator{rows: rows}, nil
}

// UpsertEdge creates a new edge or updates an existing edge.
func (c *CockroachGraph) UpsertEdge(edge *graph.Edge) error {
	row := c.DB.QueryRow(dbQueries.upsertEdgeQuery, edge.Src, edge.Dst)
	if err := row.Scan(&edge.ID, &edge.UpdatedAt); err != nil {
		if isForeignKeyViolationError(err) {
			err = graph.ErrUnknownEdgeLinks
		}
		return xerrors.Errorf("upsert edge: %w", err)
	}

	edge.UpdatedAt = edge.UpdatedAt.UTC()
	return nil
}

// Edges returns an iterator for the set of edges whose source vertex IDs
// belong to the [fromID, toID) range and were last updated before the provided
// value.
func (c *CockroachGraph) Edges(fromID, toID uuid.UUID, updatedBefore time.Time) (graph.EdgeIterator, error) {
	rows, err := c.DB.Query(dbQueries.edgesInPartitionQuery, fromID, toID, updatedBefore.UTC())
	if err != nil {
		return nil, xerrors.Errorf("edges: %w", err)
	}

	return &edgeIterator{rows: rows}, nil
}

// RemoveStaleEdges removes any edge that originates from the specified link ID
// and was updated before the specified timestamp.
func (c *CockroachGraph) RemoveStaleEdges(fromID uuid.UUID, updatedBefore time.Time) error {
	_, err := c.DB.Exec(dbQueries.removeStaleEdgesQuery, fromID, updatedBefore.UTC())
	if err != nil {
		return xerrors.Errorf("remove stale edges: %w", err)
	}

	return nil
}

// isForeignKeyViolationError returns true if err indicates a foreign key
// constraint violation.
func isForeignKeyViolationError(err error) bool {
	pqErr, valid := err.(*pq.Error)
	if !valid {
		return false
	}

	return pqErr.Code.Name() == "foreign_key_violation"
}
