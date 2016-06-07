package graph

import (
	"github.com/Callidon/joseki/core"
	"testing"
)

func TestAddIndexGraph(t *testing.T) {
	graph := NewIndexGraph()
	subj := core.NewURI("dblp", "Thomas")
	pred := core.NewURI("foaf", "age")
	obj := core.NewLiteral("22")
	triple := core.NewTriple(subj, pred, obj)
	graph.Add(triple)
}

func TestFilterIndexGraph(t *testing.T) {
	graph := NewIndexGraph()
	subj := core.NewURI("dblp", "Thomas")
	pred := core.NewURI("foaf", "age")
	obj := core.NewLiteral("22")
	triple := core.NewTriple(subj, pred, obj)
	graph.Add(triple)

	triples, _ := graph.Filter(subj, pred, obj)

	if len(triples) != 1 {
		t.Error("expected length == 1 but got length ==", len(triples))
	}

	if test, _ := triples[0].Equals(triple); !test {
		t.Error(triple, "not in results :", triples)
	}
}

func BenchmarkFilterIndexGraph(b *testing.B) {
    graph := NewIndexGraph()
	subj := core.NewURI("dblp", "Thomas")
	pred := core.NewURI("foaf", "age")
    for i := 0; i < 1000; 
	obj := core.NewLiteral("22")
	triple := core.NewTriple(subj, pred, obj)
	graph.Add(triple)

    for i := 0; i < b.N; i++ {
        graph.Filter(subj, pred, obj)
    }
}
