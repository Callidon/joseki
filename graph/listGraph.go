package graph

import (
	"os"
	"github.com/Callidon/joseki/core"
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

func (g *ListGraph) Add(triple core.Triple) {
	g.triples = append(g.triples, triple)
}

func (g *ListGraph) Filter(subject, predicate, object core.Node) ([]core.Triple, error) {
	results := make([]core.Triple, 0)
	ref_triple := core.NewTriple(subject, predicate, object)
	// search for matching triple pattern in graph
	for _, triple := range g.triples {
		test, err := ref_triple.Compare(triple)
		if err != nil {
			return nil, err
		} else if test {
			results = append(results, triple)
		}
	}
	return results, nil
}

func (g *ListGraph) Serialize(format string) string {
    // TODO
    return ""
}
