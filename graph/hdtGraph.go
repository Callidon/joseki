package graph

import (
	"github.com/Callidon/joseki/core"
	"os"
)

type bitmapNode struct {
    id int
    sons []*bitmapNode
}

type bitmapTriple struct {
    subject_id int
    predicate_id int
    object_id int
}

// Implementation of a RDF Graph based on the HDT-MR model proposed by Giménez-García et al
// For more details, see http://dataweb.infor.uva.es/projects/hdt-mr/
type HDTGraph struct {
    dictionnary bimap
	triples map[string][]core.Triple
}

// Return a new Bitmap Node without any son
func newBitmapNode(id int) bitmapNode {
    return bitmapNode{id, make([]*bitmapNode, 0)}
}

// Add a son to a Bitmap Node
func (n *bitmapNode) addSon(node *bitmapNode) {
    n.sons = append(n.sons, node)
}

// Return a new empty HDT Graph
func NewHDTGraph() HDTGraph {
	return HDTGraph{newBimap(), make(map[string][]core.Triple)}
}

func (g *HDTGraph) LoadFromFile(file *os.File) {
	//TODO
}

func (g *HDTGraph) Add(triple core.Triple) {
	key := triple.Subject.String()
	_, isIndexed := g.triples[key]
	if isIndexed {
		g.triples[key] = make([]core.Triple, 0)
		g.triples[key] = append(g.triples[key], triple)
	} else {
		g.triples[key] = append(g.triples[key], triple)
	}
}

func (g *HDTGraph) Filter(subject, predicate, object core.Node) ([]core.Triple, error) {
	results := make([]core.Triple, 0)
	ref_triple := core.NewTriple(subject, predicate, object)
	_, ok := subject.(core.BlankNode)
	// search for every subject
	if ok {
		// search for matching triple pattern in graph
		for _, triples := range g.triples {
			for _, triple := range triples {
				test, err := ref_triple.Equivalent(triple)
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
				test, err := ref_triple.Equivalent(triple)
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

func (g *HDTGraph) Serialize(format string) string {
	// TODO
	return ""
}
