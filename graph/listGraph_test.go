package graph

import (
	"github.com/Callidon/joseki/core"
	"testing"
    "math/rand"
    "fmt"
)

func TestAddListGraph(t *testing.T) {
	graph := NewListGraph()
	subj := core.NewURI("dblp", "Thomas")
	pred := core.NewURI("foaf", "age")
	obj := core.NewLiteral("22")
	triple := core.NewTriple(subj, pred, obj)
	graph.Add(triple)
}

func TestSimpleFilterListGraph(t *testing.T) {
	graph := NewListGraph()
	subj := core.NewURI("dblp", "Thomas")
	pred := core.NewURI("foaf", "age")
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
    subj := core.NewURI("dblp", "foo")
    datas := make([]core.Triple, 0)

    // create triples to be inserted
    for i := 0; i < nbDatas; i++ {
        triple := core.NewTriple(subj, core.NewURI("", string(rand.Intn(nbDatas))), core.NewLiteral(string(rand.Intn(nbDatas))))
        datas = append(datas, triple)
    }

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
        triple := core.NewTriple(core.NewURI("", string(rand.Intn(nbDatas))), core.NewURI("", string(rand.Intn(nbDatas))), core.NewLiteral(string(rand.Intn(nbDatas))))
        datas = append(datas, triple)
    }

    for i := 0; i < b.N; i++ {
        for _, triple := range datas {
            graph.Add(triple)
        }
    }
}

func BenchmarkFilterListGraph(b *testing.B) {
    graph := NewListGraph()
    nbDatas := 1000
    cpt := 0
    subj := core.NewURI("dblp", "foo")
    datas := make([]core.Triple, 0)

    // create triples to be inserted
    for i := 0; i < nbDatas; i++ {
        triple := core.NewTriple(subj, core.NewURI("", string(rand.Intn(nbDatas))), core.NewLiteral(string(rand.Intn(nbDatas))))
        datas = append(datas, triple)
    }

    for i := 0; i < b.N; i++ {
        for _ = range graph.Filter(subj, core.NewBlankNode("v"), core.NewBlankNode("w")) {
            cpt += 1
        }
    }
}
