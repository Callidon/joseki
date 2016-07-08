// Copyright (c) 2016 Thomas Minier. All rights reserved.
// Use of this source code is governed by a MIT License
// license that can be found in the LICENSE file.

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
	for _ = range graph.Filter(subj, rdf.NewVariable("v"), rdf.NewVariable("w")) {
		cpt++
	}

	if cpt != nbDatas {
		t.Error("expected", nbDatas, "results but inftsead got", cpt, "results")
	}
}

func TestComplexFilterSubsetListGraph(t *testing.T) {
	graph := NewListGraph()
	nbDatas, limit, offset := 1000, 600, 800
	cpt := 0
	subj := rdf.NewURI("dblp:foo")

	// insert random triples in the graph
	for i := 0; i < nbDatas; i++ {
		triple := rdf.NewTriple(subj, rdf.NewURI(string(rand.Intn(nbDatas))), rdf.NewLiteral(string(rand.Intn(nbDatas))))
		graph.Add(triple)
	}

	// test a FilterSubset with a simple Limit
	for _ = range graph.FilterSubset(subj, rdf.NewVariable("v"), rdf.NewVariable("w"), limit, -1) {
		cpt++
	}

	if cpt != limit {
		t.Error("expected ", limit, "results but instead found ", cpt, "results")
	}

	// test a FilterSubset with a simple offset
	cpt = 0
	for _ = range graph.FilterSubset(subj, rdf.NewVariable("v"), rdf.NewVariable("w"), -1, offset) {
		cpt++
	}

	if cpt != nbDatas-offset {
		t.Error("expected ", nbDatas-offset, "results but instead found ", cpt, "results")
	}

	// test with a offset than doesn't allow enough results to reach the limit
	cpt = 0
	for _ = range graph.FilterSubset(subj, rdf.NewVariable("v"), rdf.NewVariable("w"), limit, offset) {
		cpt++
	}

	if cpt != nbDatas-offset {
		t.Error("expected ", nbDatas-offset, "results but instead found ", cpt, "results")
	}
}

func TestDeleteListGraph(t *testing.T) {
	graph := NewListGraph()
	nbDatas := 1000
	cpt := 0
	subj := rdf.NewURI("dblp:foo")

	// insert random triples in the graph
	for i := 0; i < nbDatas; i++ {
		triple := rdf.NewTriple(subj, rdf.NewURI(string(rand.Intn(nbDatas))), rdf.NewLiteral(string(rand.Intn(nbDatas))))
		graph.Add(triple)
	}

	// remove all triple with a given subject
	graph.Delete(subj, rdf.NewVariable("v"), rdf.NewVariable("w"))

	// select all triple of the graph
	for _ = range graph.Filter(subj, rdf.NewVariable("v"), rdf.NewVariable("w")) {
		cpt++
	}

	if cpt > 0 {
		t.Error("the graph should be empty")
	}
}

func TestLoadFromFileListGraph(t *testing.T) {
	graph := NewListGraph()
	cpt := 0
	graph.LoadFromFile("../parser/datas/test.nt", "nt")

	// select all triple of the graph
	for _ = range graph.Filter(rdf.NewVariable("y"), rdf.NewVariable("v"), rdf.NewVariable("w")) {
		cpt++
	}

	if cpt != 5 {
		t.Error("the graph should contains 4 triples, but it contains", cpt, "triples")
	}
}

// Benchmarking with WatDiv 1K

func BenchmarkAddListGraph(b *testing.B) {
	graph := NewListGraph()

	for i := 0; i < b.N; i++ {
		graph.LoadFromFile("../parser/datas/watdiv1k.nt", "nt")
	}
}

func BenchmarkDeleteAllListGraph(b *testing.B) {
	graph := NewListGraph()
	graph.LoadFromFile("../parser/datas/watdiv1k.nt", "nt")
	pred := rdf.NewURI("http://purl.org/goodrelations/price")

	for i := 0; i < b.N; i++ {
		graph.Delete(rdf.NewVariable("v"), pred, rdf.NewVariable("w"))
	}
}

func BenchmarkAllFilterListGraph(b *testing.B) {
	graph := NewListGraph()
	graph.LoadFromFile("../parser/datas/watdiv1k.nt", "nt")
	cpt := 0

	for i := 0; i < b.N; i++ {
		// select all triple of the graph
		for _ = range graph.Filter(rdf.NewVariable("v"), rdf.NewVariable("w"), rdf.NewVariable("z")) {
			cpt++
		}
	}
}

func BenchmarkSpecificFilterListGraph(b *testing.B) {
	graph := NewListGraph()
	graph.LoadFromFile("../parser/datas/watdiv1k.nt", "nt")
	subj := rdf.NewURI("http://db.uwaterloo.ca/~galuc/wsdbm/Product43")
	pred := rdf.NewURI("http://purl.org/stuff/rev#hasReview")
	obj := rdf.NewURI("http://db.uwaterloo.ca/~galuc/wsdbm/Review864")
	cpt := 0

	// insert a specific triple at the end
	triple := rdf.NewTriple(subj, pred, obj)
	graph.Add(triple)

	for i := 0; i < b.N; i++ {
		// fetch the last inserted triple into the graph
		for _ = range graph.Filter(subj, pred, obj) {
			cpt++
		}
	}
}

func BenchmarkAllFilterSubsetListGraph(b *testing.B) {
	graph := NewListGraph()
	graph.LoadFromFile("../parser/datas/watdiv1k.nt", "nt")
	limit, offset := 600, 200
	cpt := 0

	for i := 0; i < b.N; i++ {
		// select all triple of the graph
		for _ = range graph.FilterSubset(rdf.NewVariable("v"), rdf.NewVariable("w"), rdf.NewVariable("z"), limit, offset) {
			cpt++
		}
	}
}
