// Copyright (c) 2016 Thomas Minier. All rights reserved.
// Use of this source code is governed by a MIT License
// license that can be found in the LICENSE file.

package sparql

import (
  "github.com/Callidon/joseki/graph"
  "github.com/Callidon/joseki/rdf"
)

const (
	// Max size for the buffer of this package
	bufferSize = 100
)

// queryType is the type of a SPARQL query
type queryType float64

const (
  _ = iota
  // selectQuery is a SPARQL SELECT query
  selectQuery queryType = 1 << (10 * iota)
  // askQuery is a SPARQL ASK query
  askQuery
  // describeQuery is a SPARQL DESCRIBE query
  describeQuery
  // constructQuery is a SPARQL CONSTRUCT query
  constructQuery
)

// queryDescriptor is a virtual description of a SPARQL query.
// It holds informations about the query during its construction process
// and can be converted into a real query.
type queryDescriptor struct {
  graph graph.Graph
  qType queryType
  //bgps map[string]sparqlNode
}

// newQueryDescriptor creates a new queryDescriptor.
func newQueryDescriptor(graph graph.Graph, qType queryType) *queryDescriptor {
  return &queryDescriptor{graph, qType}
}

// From set the source's graph for the query.
// Multiple calls of this method will override the previous source each time.
func (q *queryDescriptor) From(graph graph.Graph) {
  q.graph = graph
}

// Where add multiples triples pattern as a new BGP evualuted by the query
// Multiple calls of this method will each time add a new BGP to the query.
func (q *queryDescriptor) Where(triples ...rdf.Triple) {
  // TODO
}

// build analyse the query execution plan and return its first node
func (q *queryDescriptor) build() sparqlNode {
  // TODO
  return nil
}
