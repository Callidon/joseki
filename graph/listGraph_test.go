package graph

import (
	"github.com/Callidon/joseki/core"
	"testing"
    "math/rand"
)

func TestAddListGraph(t *testing.T) {
	graph := NewListGraph()
	subj := core.NewURI("dblp:Thomas")
	pred := core.NewURI("foaf:age")
	obj := core.NewLiteral("22")
	triple := core.NewTriple(subj, pred, obj)
	graph.Add(triple)

    if test, err := graph.triples[0].Equals(triple); !test && (err != nil) {
        t.Error(triple, "hasn't been inserted into the graph")
    }
}

func TestSimpleFilterListGraph(t *testing.T) {
	graph := NewListGraph()
	subj := core.NewURI("dblp:Thomas")
	pred := core.NewURI("foaf:age")
	obj := core.NewLiteral("22")
	triple := core.NewTriple(subj, pred, obj)
	graph.Add(triple)

    for result := range graph.Filter(subj, pred, obj) {
        if test, _ := result.Equals(triple); !test {
    		t.Error(triple, "not in results :", result)
    	}
    }
}

func TestComplexFilterListGraph(t *testing.T) {
    graph := NewListGraph()
    nbDatas := 1000
    cpt := 0
    subj := core.NewURI("dblp:foo")

    // insert random triples in the graph
    for i := 0; i < nbDatas; i++ {
        triple := core.NewTriple(subj, core.NewURI(string(rand.Intn(nbDatas))), core.NewLiteral(string(rand.Intn(nbDatas))))
        graph.Add(triple)
    }

    // select all triple of the graph
    for _ = range graph.Filter(subj, core.NewBlankNode("v"), core.NewBlankNode("w")) {
        cpt += 1
    }

    if cpt != nbDatas {
        t.Error("expected ", nbDatas, "results but got ", cpt, "results")
    }
}

func BenchmarkAddListGraph(b *testing.B) {
    graph := NewListGraph()
    nbDatas := 1000
    datas := make([]core.Triple, 0)

    // create triples to be inserted
    for i := 0; i < nbDatas; i++ {
        triple := core.NewTriple(core.NewURI(string(rand.Intn(nbDatas))), core.NewURI(string(rand.Intn(nbDatas))), core.NewLiteral(string(rand.Intn(nbDatas))))
        datas = append(datas, triple)
    }

    for i := 0; i < b.N; i++ {
        for _, triple := range datas {
            graph.Add(triple)
        }
    }
}

func BenchmarkAllFilterListGraph(b *testing.B) {
    graph := NewListGraph()
    nbDatas := 1000
    cpt := 0
    subj := core.NewURI("dblp:foo")

    // insert random triples in the graph
    for i := 0; i < nbDatas; i++ {
        triple := core.NewTriple(subj, core.NewURI(string(rand.Intn(nbDatas))), core.NewLiteral(string(rand.Intn(nbDatas))))
        graph.Add(triple)
    }

    for i := 0; i < b.N; i++ {
        // select all triple of the graph
        for _ = range graph.Filter(subj, core.NewBlankNode("v"), core.NewBlankNode("w")) {
            cpt += 1
        }
    }
}

func BenchmarkSpecificFilterListGraph(b *testing.B) {
    graph := NewListGraph()
    nbDatas := 1000
    cpt := 0
    subj := core.NewURI("dblp:foo")
    pred := core.NewURI("foaf:age")
    obj := core.NewURI("22")

    // insert random triples in the graph
    for i := 0; i < nbDatas; i++ {
        triple := core.NewTriple(subj, core.NewURI(string(rand.Intn(nbDatas))), core.NewLiteral(string(rand.Intn(nbDatas))))
        graph.Add(triple)
    }
    // insert a specific triple at the end
    triple := core.NewTriple(subj, pred, obj)
    graph.Add(triple)

    for i := 0; i < b.N; i++ {
        // fetch the last inserted triple into the graph
        for _ = range graph.Filter(subj, pred, obj) {
            cpt += 1
        }
    }
}
