// Copyright (c) 2016 Thomas Minier. All rights reserved.
// Use of this source code is governed by a MIT License
// license that can be found in the LICENSE file.

package graph

import (
	"github.com/Callidon/joseki/rdf"
	"math/rand"
	"testing"
)

func TestAddHDTGraph(t *testing.T) {
	var node rdf.Node
	graph := NewHDTGraph()

	subj := rdf.NewURI("dblp:Thomas")
	predA := rdf.NewURI("foaf:age")
	predB := rdf.NewURI("schema:livesIn")
	objA := rdf.NewLiteral("22")
	objB := rdf.NewLiteral("Nantes")
	tripleA := rdf.NewTriple(subj, predA, objA)
	tripleB := rdf.NewTriple(subj, predB, objB)
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
	if test, _ := node.Equals(subj); !test {
		t.Error("expected <dbpl:Thomas> to be the only subject node but found", node)
	}
	node, _ = graph.dictionnary.extract(graph.root.sons[0].sons[1].id)
	if test, _ := node.Equals(predA); !test {
		t.Error("expected <foaf:age> to be the first predicate of <dblp:Thomas> but found", node)
	}
	node, _ = graph.dictionnary.extract(graph.root.sons[0].sons[3].id)
	if test, _ := node.Equals(predB); !test {
		t.Error("expected <schema:livesIn> to be the second predicate of <dblp:Thomas> but found", node)
	}
	node, _ = graph.dictionnary.extract(graph.root.sons[0].sons[1].sons[2].id)
	if test, _ := node.Equals(objA); !test {
		t.Error("expected \"20\" to be the object of <dblp:Thomas> <foaf:age> but found", node)
	}
	node, _ = graph.dictionnary.extract(graph.root.sons[0].sons[3].sons[4].id)
	if test, _ := node.Equals(objB); !test {
		t.Error("expected \"Nantes\" to be the object of <dblp:Thomas> <schema:livesIn> but found", node)
	}
}

func TestSimpleFilterHDTGraph(t *testing.T) {
	graph := NewHDTGraph()
	subj := rdf.NewURI("dblp:Thomas")
	predA := rdf.NewURI("foaf:age")
	predB := rdf.NewURI("schema:livesIn")
	objA := rdf.NewLiteral("22")
	objB := rdf.NewURI("dbpedia:Nantes")
	tripleA := rdf.NewTriple(subj, predA, objA)
	tripleB := rdf.NewTriple(subj, predB, objB)
	graph.Add(tripleA)
	graph.Add(tripleB)

	// select one triple
	for result := range graph.Filter(subj, predA, objA) {
		if test, err := result.Equals(tripleA); test && (err != nil) {
			t.Error(tripleA, "not in results :", result)
		}
	}

	// select multiple triples using Blank Nodes
	cpt := 0
	for _ = range graph.Filter(subj, rdf.NewVariable("v"), rdf.NewVariable("w")) {
		cpt++
	}

}

func TestFilterNoResultHDTGraph(t *testing.T) {
	graph := NewHDTGraph()
	subj := rdf.NewURI("dblp:Thomas")
	predA := rdf.NewURI("foaf:age")
	predB := rdf.NewURI("schema:livesIn")
	objA := rdf.NewLiteral("22")
	objB := rdf.NewURI("dbpedia:Nantes")
	tripleA := rdf.NewTriple(subj, predA, objA)
	tripleB := rdf.NewTriple(subj, predB, objB)
	graph.Add(tripleA)
	graph.Add(tripleB)

	// select a triple that doesn't exist in the graph
	cpt := 0
	for _ = range graph.Filter(rdf.NewURI("<htt://example.org>"), rdf.NewVariable("v1"), rdf.NewVariable("v2")) {
		cpt++
	}

	if cpt > 0 {
		t.Error("expected no result but instead found", cpt, "results")
	}

}

func TestComplexFilterHDTGraph(t *testing.T) {
	graph := NewHDTGraph()
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
		t.Error("expected ", nbDatas, "results but instead found ", cpt, "results")
	}
}

func TestComplexFilterSubsetHDTGraph(t *testing.T) {
	graph := NewHDTGraph()
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

func TestDeleteHDTGraph(t *testing.T) {
	graph := NewHDTGraph()
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
		t.Error("Error : the graph should be empty")
	}
}

func TestLoadFromFileHDTGraph(t *testing.T) {
	graph := NewHDTGraph()
	cpt := 0
	graph.LoadFromFile("../parser/datas/test.nt", "nt")

	// select all triple of the graph
	for _ = range graph.Filter(rdf.NewVariable("y"), rdf.NewVariable("v"), rdf.NewVariable("w")) {
		cpt++
	}

	if cpt != 4 {
		t.Error("the graph should contains 4 triples, but it contains", cpt, "triples")
	}
}

// Benchmarking

func BenchmarkAddHDTGraph(b *testing.B) {
	b.Skip("skipped because it's currently not accurate")
	graph := NewHDTGraph()
	graph.LoadFromFile("../parser/datas/watdiv1k.nt", "nt")
	triple := rdf.NewTriple(rdf.NewURI("http://example.org/subject"), rdf.NewURI("http://example.org/predicate"), rdf.NewURI("http://example.org/object"))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		graph.Add(triple)
	}
}

func BenchmarkLoadFromFileHDTGraph(b *testing.B) {

	for i := 0; i < b.N; i++ {
		graph := NewHDTGraph()
		graph.LoadFromFile("../parser/datas/watdiv1k.nt", "nt")
	}
}

func BenchmarkAllFilterHDTGraph(b *testing.B) {
	graph := NewHDTGraph()
	graph.LoadFromFile("../parser/datas/watdiv1k.nt", "nt")
	cpt := 0
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// select all triple of the graph
		for _ = range graph.Filter(rdf.NewVariable("v"), rdf.NewVariable("w"), rdf.NewVariable("z")) {
			cpt++
		}
	}
}

func BenchmarkSpecificFilterHDTGraph(b *testing.B) {
	graph := NewHDTGraph()
	graph.LoadFromFile("../parser/datas/watdiv1k.nt", "nt")
	subj := rdf.NewURI("http://db.uwaterloo.ca/~galuc/wsdbm/User999")
	pred := rdf.NewURI("http://xmlns.com/foaf/age")
	obj := rdf.NewURI("http://db.uwaterloo.ca/~galuc/wsdbm/AgeGroup2")
	cpt := 0
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// fetch the last inserted triple into the graph
		for _ = range graph.Filter(subj, pred, obj) {
			cpt++
		}
	}
}

func BenchmarkAllFilterSubsetHDTGraph(b *testing.B) {
	graph := NewHDTGraph()
	graph.LoadFromFile("../parser/datas/watdiv1k.nt", "nt")
	limit, offset := 600, 200
	cpt := 0
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// select all triple of the graph
		for _ = range graph.FilterSubset(rdf.NewVariable("v"), rdf.NewVariable("w"), rdf.NewVariable("z"), limit, offset) {
			cpt++
		}
	}
}
