// Copyright (c) 2016 Thomas Minier. All rights reserved.
// Use of this source code is governed by a MIT License
// license that can be found in the LICENSE file.

package graph

import (
	"github.com/Callidon/joseki/rdf"
	"sync"
)

// ListGraph is dummy implementation of a RDF Graph, using a simple slice to store RDF Triples.
//
// Very poorly optimized, should only be used for demonstration or benchmarking purposes.
type ListGraph struct {
	triples []rdf.Triple
	*sync.Mutex
	*rdfReader
}

// NewListGraph creates a new List Graph.
func NewListGraph() *ListGraph {
	reader := newRDFReader()
	g := &ListGraph{make([]rdf.Triple, 0), &sync.Mutex{}, reader}
	reader.graph = g
	return g
}

// Add a new Triple pattern to the graph.
func (g *ListGraph) Add(triple rdf.Triple) {
	g.Lock()
	defer g.Unlock()
	g.triples = append(g.triples, triple)
}

// Delete triples from the graph that match a BGP given in parameters.
func (g *ListGraph) Delete(subject, object, predicate rdf.Node) {
	var newTriples []rdf.Triple
	refTriple := rdf.NewTriple(subject, predicate, object)
	g.Lock()
	defer g.Unlock()
	// resinsert into the graph the elements we doesn't want to delete
	for _, triple := range g.triples {
		if test, _ := triple.Equivalent(refTriple); !test {
			newTriples = append(newTriples, triple)
		}
	}
	g.triples = newTriples
}

// Filter fetch triples form the graph that match a BGP given in parameters.
func (g *ListGraph) Filter(subject, predicate, object rdf.Node) chan rdf.Triple {
	results := make(chan rdf.Triple)
	refTriple := rdf.NewTriple(subject, predicate, object)
	g.Lock()
	defer g.Unlock()
	// search for matching triple pattern in graph
	go func() {
		for _, triple := range g.triples {
			test, err := refTriple.Equivalent(triple)
			if (err == nil) && test {
				results <- triple
			}
		}
		close(results)
	}()
	return results
}

// Serialize the graph into a given format and return it as a string.
func (g *ListGraph) Serialize(format string) string {
	// TODO
	return ""
}
