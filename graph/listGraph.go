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
	lock    *sync.Mutex
}

// NewListGraph creates a new List Graph.
func NewListGraph() ListGraph {
	return ListGraph{make([]rdf.Triple, 0), &sync.Mutex{}}
}

// LoadFromFile load the content of a RDF graph stored in a file into the current graph.
func (g *ListGraph) LoadFromFile(filename, format string) {
	loadFromFile(g, filename, format)
}

// Add a new Triple pattern to the graph.
func (g *ListGraph) Add(triple rdf.Triple) {
	g.lock.Lock()
	defer g.lock.Unlock()
	g.triples = append(g.triples, triple)
}

// Delete triples from the graph that match a BGP given in parameters.
func (g *ListGraph) Delete(subject, object, predicate rdf.Node) {
	var newTriples []rdf.Triple
	refTriple := rdf.NewTriple(subject, predicate, object)
	g.lock.Lock()
	defer g.lock.Unlock()
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
	g.lock.Lock()
	defer g.lock.Unlock()
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
