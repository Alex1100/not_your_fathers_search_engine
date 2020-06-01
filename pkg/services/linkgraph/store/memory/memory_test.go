package memory

import (
	"testing"

	gc "gopkg.in/check.v1"
	"not_your_fathers_search_engine/pkg/services/linkgraph/graph/graphtest"
)

var _ = gc.Suite(new(InMemoryGraphTestSuite))

func Test(t *testing.T) { gc.TestingT(t) }

type InMemoryGraphTestSuite struct {
	graphtest.SuiteBase
}

func (s *InMemoryGraphTestSuite) SetUpTest(c *gc.C) {
	s.SetGraph(NewInMemoryGraph())
}
