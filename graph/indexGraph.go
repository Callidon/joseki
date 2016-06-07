package graph

import (
	"github.com/Callidon/joseki/core"
	"os"
)

type IndexGraph struct {
	triples map[string][]core.Triple
}

func NewIndexGraph() IndexGraph {
	return IndexGraph{make(map[string][]core.Triple)}
}

func (g *IndexGraph) LoadFromFile(file *os.File) {
	//TODO
}

func (g *IndexGraph) Add(triple core.Triple) {
	key := triple.Subject.String()
	_, isIndexed := g.triples[key]
	if isIndexed {
		g.triples[key] = make([]core.Triple, 0)
		g.triples[key] = append(g.triples[key], triple)
	} else {
		g.triples[key] = append(g.triples[key], triple)
	}
}

func (g *IndexGraph) Filter(subject, predicate, object core.Node) ([]core.Triple, error) {
	results := make([]core.Triple, 0)
	ref_triple := core.NewTriple(subject, predicate, object)
	_, ok := subject.(core.BlankNode)
	// search for every subject
	if ok {
		// search for matching triple pattern in graph
		for _, triples := range g.triples {
			for _, triple := range triples {
				test, err := ref_triple.Compare(triple)
				if err != nil {
					return nil, err
				} else if test {
					results = append(results, triple)
				}
			}
		}
	} else {
		// search with a specific subject
		triples, isIndexed := g.triples[subject.String()]
		if isIndexed {
            for _, triple := range triples {
				test, err := ref_triple.Compare(triple)
				if err != nil {
					return nil, err
				} else if test {
					results = append(results, triple)
				}
			}
		}
	}
	return results, nil
}

func (g *IndexGraph) Serialize(format string) string {
	// TODO
	return ""
}
