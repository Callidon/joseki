package graph

import (
	"github.com/Callidon/joseki/core"
	"os"
)

// Dummy implemntation of a RDF Graph, using a simple slice to store RDF Triples
// Very poorly optimized, should only be used for demonstration purpose
type ListGraph struct {
	triples []core.Triple
}

func NewListGraph() ListGraph {
	return ListGraph{make([]core.Triple, 0)}
}

func (g *ListGraph) LoadFromFile(file *os.File) {
	//TODO
}

// Add a new Triple pattern to the graph
func (g *ListGraph) Add(triple core.Triple) {
	g.triples = append(g.triples, triple)
}

// Fetch triples form the graph that match a BGP given in parameters
func (g *ListGraph) Filter(subject, predicate, object core.Node) chan core.Triple {
	results := make(chan core.Triple)
	ref_triple := core.NewTriple(subject, predicate, object)
	// search for matching triple pattern in graph
	go func() {
		for _, triple := range g.triples {
			test, err := ref_triple.Equivalent(triple)
			if (err == nil) && test {
				results <- triple
			}
		}
		close(results)
	}()
	return results
}

func (g *ListGraph) Serialize(format string) string {
	// TODO
	return ""
}
