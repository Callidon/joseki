package graph

import (
	"github.com/Callidon/joseki/rdf"
	"math/rand"
	"testing"
)

func TestAddListGraph(t *testing.T) {
	graph := NewListGraph()
	subj := rdf.NewURI("dblp:Thomas")
	pred := rdf.NewURI("foaf:age")
	obj := rdf.NewLiteral("22")
	triple := rdf.NewTriple(subj, pred, obj)
	graph.Add(triple)

	if test, err := graph.triples[0].Equals(triple); !test && (err != nil) {
		t.Error(triple, "hasn't been inserted into the graph")
	}
}

func TestSimpleFilterListGraph(t *testing.T) {
	graph := NewListGraph()
	subj := rdf.NewURI("dblp:Thomas")
	pred := rdf.NewURI("foaf:age")
	obj := rdf.NewLiteral("22")
	triple := rdf.NewTriple(subj, pred, obj)
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
	subj := rdf.NewURI("dblp:foo")

	// insert random triples in the graph
	for i := 0; i < nbDatas; i++ {
		triple := rdf.NewTriple(subj, rdf.NewURI(string(rand.Intn(nbDatas))), rdf.NewLiteral(string(rand.Intn(nbDatas))))
		graph.Add(triple)
	}

	// select all triple of the graph
	for _ = range graph.Filter(subj, rdf.NewBlankNode("v"), rdf.NewBlankNode("w")) {
		cpt += 1
	}

	if cpt != nbDatas {
		t.Error("expected ", nbDatas, "results but got ", cpt, "results")
	}
}

func BenchmarkAddListGraph(b *testing.B) {
	graph := NewListGraph()
	nbDatas := 1000
	datas := make([]rdf.Triple, 0)

	// create triples to be inserted
	for i := 0; i < nbDatas; i++ {
		triple := rdf.NewTriple(rdf.NewURI(string(rand.Intn(nbDatas))), rdf.NewURI(string(rand.Intn(nbDatas))), rdf.NewLiteral(string(rand.Intn(nbDatas))))
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
	subj := rdf.NewURI("dblp:foo")

	// insert random triples in the graph
	for i := 0; i < nbDatas; i++ {
		triple := rdf.NewTriple(subj, rdf.NewURI(string(rand.Intn(nbDatas))), rdf.NewLiteral(string(rand.Intn(nbDatas))))
		graph.Add(triple)
	}

	for i := 0; i < b.N; i++ {
		// select all triple of the graph
		for _ = range graph.Filter(subj, rdf.NewBlankNode("v"), rdf.NewBlankNode("w")) {
			cpt += 1
		}
	}
}

func BenchmarkSpecificFilterListGraph(b *testing.B) {
	graph := NewListGraph()
	nbDatas := 1000
	cpt := 0
	subj := rdf.NewURI("dblp:foo")
	pred := rdf.NewURI("foaf:age")
	obj := rdf.NewURI("22")

	// insert random triples in the graph
	for i := 0; i < nbDatas; i++ {
		triple := rdf.NewTriple(subj, rdf.NewURI(string(rand.Intn(nbDatas))), rdf.NewLiteral(string(rand.Intn(nbDatas))))
		graph.Add(triple)
	}
	// insert a specific triple at the end
	triple := rdf.NewTriple(subj, pred, obj)
	graph.Add(triple)

	for i := 0; i < b.N; i++ {
		// fetch the last inserted triple into the graph
		for _ = range graph.Filter(subj, pred, obj) {
			cpt += 1
		}
	}
}