package graph

import (
    "testing"
	"github.com/Callidon/joseki/core"
    "math/rand"
)

func TestAddHDTGraph(t *testing.T) {
    var node core.Node
	graph := NewHDTGraph()

	subj := core.NewURI("dblp", "Thomas")
    predA := core.NewURI("foaf", "age")
    predB := core.NewURI("schema", "livesIn")
    objA := core.NewLiteral("22")
    objB := core.NewLiteral("Nantes")
	tripleA := core.NewTriple(subj, predA, objA)
    tripleB := core.NewTriple(subj, predB, objB)
	graph.Add(tripleA)
    graph.Add(tripleB)

    // check for the structure of the tree (repartition of nodes & number of levels)
    if len(graph.root.sons) != 1 {
        t.Error("doesn't found exactly one subject after inserting two triples with the same subject")
    }
    if len(graph.root.sons[0].sons) != 2 {
        t.Error("doesn't found exactly two predicates after inserting two triples with different predicates")
    }
    if len(graph.root.sons[0].sons[1].sons) != 1 {
        t.Error("doesn't found exactly one subject")
    }
    if len(graph.root.sons[0].sons[1].sons[2].sons) > 0 {
        t.Error("the tree has 4 levels instead of only three (excluding the root level)")
    }

    // check for the values in the nodes
    node, _ = graph.dictionnary.extract(graph.root.sons[0].id)
    if test, _ := node.Equals(subj); ! test {
        t.Error("expected <dbpl:Thomas> to be the only subject node but found", node)
    }
    node, _ = graph.dictionnary.extract(graph.root.sons[0].sons[1].id)
    if test, _ := node.Equals(predA); ! test {
        t.Error("expected <foaf:age> to be the first predicate of <dblp:Thomas> but found", node)
    }
    node, _ = graph.dictionnary.extract(graph.root.sons[0].sons[3].id)
    if test, _ := node.Equals(predB); ! test {
        t.Error("expected <schema:livesIn> to be the second predicate of <dblp:Thomas> but found", node)
    }
    node, _ = graph.dictionnary.extract(graph.root.sons[0].sons[1].sons[2].id)
    if test, _ := node.Equals(objA); ! test {
        t.Error("expected \"20\" to be the object of <dblp:Thomas> <foaf:age> but found", node)
    }
    node, _ = graph.dictionnary.extract(graph.root.sons[0].sons[3].sons[4].id)
    if test, _ := node.Equals(objB); ! test {
        t.Error("expected \"Nantes\" to be the object of <dblp:Thomas> <schema:livesIn> but found", node)
    }
}

func TestFilterHDTGraph(t *testing.T) {
	graph := NewHDTGraph()

	subj := core.NewURI("dblp", "Thomas")
	pred := core.NewURI("foaf", "age")
	obj := core.NewLiteral("22")
	triple := core.NewTriple(subj, pred, obj)
	graph.Add(triple)

	/*triples, _ := graph.Filter(subj, pred, obj)

	if len(triples) != 1 {
		t.Error("expected length == 1 but got length ==", len(triples))
	}

	if test, _ := triples[0].Equals(triple); !test {
		t.Error(triple, "not in results :", triples)
	}*/
}

func BenchmarkAddHDTGraph(b *testing.B) {
    graph := NewHDTGraph()
    nbDatas := 1000
    datas := make([]core.Triple, 0)

    // create triples to be inserted
    for i := 0; i < nbDatas; i++ {
        triple := core.NewTriple(core.NewURI("", string(rand.Int())), core.NewURI("", string(rand.Int())), core.NewLiteral(string(rand.Int())))
        datas = append(datas, triple)
    }

    for i := 0; i < b.N; i++ {
        for _, triple := range datas {
            graph.Add(triple)
        }
    }
}
