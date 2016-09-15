// Copyright (c) 2016 Thomas Minier. All rights reserved.
// Use of this source code is governed by a MIT License
// license that can be found in the LICENSE file.

package graph

import (
	"github.com/Callidon/joseki/rdf"
	"math/rand"
	"os"
	"testing"
)

// skip a test if a file does'nt exist
func skipTest(file string, t *testing.T) {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		t.Skip(file, "doesn't exist, skipping test")
	}
}

func TestAddTreeGraph(t *testing.T) {
	var node rdf.Node
	graph := NewTreeGraph()

	subj := rdf.NewURI("http://dbpl.org#Thomas")
	predA := rdf.NewURI("http://foaf.com/age")
	predB := rdf.NewURI("http://Schema.org#livesIn")
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
		t.Error("expected <http://foaf.com/age> to be the first predicate of <http://dbpl.org#Thomas> but found", node)
	}
	node, _ = graph.dictionnary.extract(graph.root.sons[0].sons[3].id)
	if test, _ := node.Equals(predB); !test {
		t.Error("expected <http://Schema.org#livesIn> to be the second predicate of <http://dbpl.org#Thomas> but found", node)
	}
	node, _ = graph.dictionnary.extract(graph.root.sons[0].sons[1].sons[2].id)
	if test, _ := node.Equals(objA); !test {
		t.Error("expected \"20\" to be the object of <http://dbpl.org#Thomas> <http://foaf.com/age> but found", node)
	}
	node, _ = graph.dictionnary.extract(graph.root.sons[0].sons[3].sons[4].id)
	if test, _ := node.Equals(objB); !test {
		t.Error("expected \"Nantes\" to be the object of <http://dbpl.org#Thomas> <http://Schema.org#livesIn> but found", node)
	}
}

func TestFilterTreeGraph(t *testing.T) {
	skipTest("./watdiv.30k.nt", t)
	graph := NewTreeGraph()
	graph.LoadFromFile("./watdiv.30k.nt", "nt")
	subj := rdf.NewURI("http://db.uwaterloo.ca/~galuc/wsdbm/Offer10001")
	pred := rdf.NewURI("http://schema.org/eligibleRegion")
	obj := rdf.NewURI("http://db.uwaterloo.ca/~galuc/wsdbm/Country0")
	triple := rdf.NewTriple(subj, pred, obj)
	cpt := 0

	// select one triple specific triple pattern
	for result := range graph.Filter(subj, pred, obj) {
		if test, err := result.Equals(triple); !test || (err != nil) {
			t.Error("expected", triple, "but instead got", result)
		}
		cpt++
	}

	if cpt != 1 {
		t.Error("expected 1 result but instead got", cpt, "results")
	}

	// select all triples
	cpt = 0
	for _ = range graph.Filter(rdf.NewVariable("v"), rdf.NewVariable("w"), rdf.NewVariable("z")) {
		cpt++
	}
	if cpt != 30000 {
		t.Error("expected 30000 results but instead got", cpt, "results")
	}

	// select multiple triples with the same subject
	cpt = 0
	for _ = range graph.Filter(rdf.NewURI("http://db.uwaterloo.ca/~galuc/wsdbm/Offer1375"), rdf.NewVariable("v"), rdf.NewVariable("w")) {
		cpt++
	}
	if cpt != 9 {
		t.Error("expected 9 results but instead got", cpt, "results")
	}

	// select multiple triples with the same predicate
	cpt = 0
	for _ = range graph.Filter(rdf.NewVariable("v"), rdf.NewURI("http://www.geonames.org/ontology#parentCountry"), rdf.NewVariable("w")) {
		cpt++
	}
	if cpt != 240 {
		t.Error("expected 240 results but instead got", cpt, "results")
	}

	// select multiple triples with the same object
	cpt = 0
	for _ = range graph.Filter(rdf.NewVariable("v"), rdf.NewVariable("w"), rdf.NewLiteral("673")) {
		cpt++
	}
	if cpt != 6 {
		t.Error("expected 6 results but instead got", cpt, "results")
	}

	// select a triple that doesn't exist in the graph
	cpt = 0
	for _ = range graph.Filter(rdf.NewURI("http://example.org"), rdf.NewVariable("v1"), rdf.NewVariable("v2")) {
		cpt++
	}

	if cpt > 0 {
		t.Error("expected no result but instead found", cpt, "results")
	}
}

func TestFilterSubsetTreeGraph(t *testing.T) {
	skipTest("./watdiv.30k.nt", t)
	graph := NewTreeGraph()
	graph.LoadFromFile("./watdiv.30k.nt", "nt")
	nbDatas, limit, offset := 30000, 600, 800
	cpt := 0

	// test a FilterSubset with a simple Limit
	for _ = range graph.FilterSubset(rdf.NewVariable("x"), rdf.NewVariable("v"), rdf.NewVariable("w"), limit, -1) {
		cpt++
	}

	if cpt != limit {
		t.Error("expected ", limit, "results but instead found ", cpt, "results")
	}

	// test a FilterSubset with a simple offset
	cpt = 0
	for _ = range graph.FilterSubset(rdf.NewVariable("x"), rdf.NewVariable("v"), rdf.NewVariable("w"), -1, offset) {
		cpt++
	}

	if cpt != nbDatas-offset {
		t.Error("expected ", nbDatas-offset, "results but instead found ", cpt, "results")
	}

	// test with a offset than doesn't allow enough results to reach the limit
	cpt = 0
	offset = nbDatas - 10
	for _ = range graph.FilterSubset(rdf.NewVariable("x"), rdf.NewVariable("v"), rdf.NewVariable("w"), limit, offset) {
		cpt++
	}

	if cpt != nbDatas-offset {
		t.Error("expected ", nbDatas-offset, "results but instead found ", cpt, "results")
	}
}

func TestDeleteTreeGraph(t *testing.T) {
	var triple rdf.Triple
	graph := NewTreeGraph()
	nbDatas := 1000
	cpt := 0
	subj := rdf.NewURI("http://dblp.org#foo")

	// insert random triples in the graph
	for i := 0; i < nbDatas; i++ {
		triple = rdf.NewTriple(subj, rdf.NewURI(string(rand.Intn(nbDatas))), rdf.NewLiteral(string(rand.Intn(nbDatas))))
		graph.Add(triple)
	}

	// remove the last triple pattern inserted
	graph.Delete(triple.Subject, triple.Predicate, triple.Object)

	// check for the absence of the triple
	for _ = range graph.Filter(triple.Subject, triple.Predicate, triple.Object) {
		cpt++
	}
	if cpt > 0 {
		t.Error("the graph shouldn't contains the triple", triple)
	}

	// check for the size of the graph
	cpt = 0
	for _ = range graph.Filter(rdf.NewVariable("v"), rdf.NewVariable("w"), rdf.NewVariable("y")) {
		cpt++
	}
	if cpt != 999 {
		t.Error("the graph should contains 999 triples, but instead it contains", cpt, "triples")
	}

	// remove all triple with a given subject
	graph.Delete(subj, rdf.NewVariable("v"), rdf.NewVariable("w"))

	// select all triple of the graph
	cpt = 0
	for _ = range graph.Filter(subj, rdf.NewVariable("v"), rdf.NewVariable("w")) {
		cpt++
	}

	if cpt > 0 {
		t.Error("Error : the graph should be empty")
	}
}

func TestLoadFromFileTreeGraph(t *testing.T) {
	graph := NewTreeGraph()
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

func BenchmarkAddTreeGraph(b *testing.B) {
	b.Skip("skipped because it's currently not accurate")
	graph := NewTreeGraph()
	graph.LoadFromFile("./watdiv.30k.nt", "nt")
	triple := rdf.NewTriple(rdf.NewURI("http://example.org/subject"), rdf.NewURI("http://example.org/predicate"), rdf.NewURI("http://example.org/object"))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		graph.Add(triple)
	}
}

func BenchmarkLoadFromFileTreeGraph(b *testing.B) {
	for i := 0; i < b.N; i++ {
		graph := NewTreeGraph()
		graph.LoadFromFile("./watdiv.30k.nt", "nt")
	}
}

func BenchmarkAllFilterTreeGraph(b *testing.B) {
	graph := NewTreeGraph()
	graph.LoadFromFile("./watdiv.30k.nt", "nt")
	cpt := 0
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// select all triple of the graph
		for _ = range graph.Filter(rdf.NewVariable("v"), rdf.NewVariable("w"), rdf.NewVariable("z")) {
			cpt++
		}
	}
}

func BenchmarkSpecificFilterTreeGraph(b *testing.B) {
	graph := NewTreeGraph()
	graph.LoadFromFile("./watdiv.30k.nt", "nt")
	subj := rdf.NewURI("http://db.uwaterloo.ca/~galuc/wsdbm/Offer10001")
	pred := rdf.NewURI("http://schema.org/eligibleRegion")
	obj := rdf.NewURI("http://db.uwaterloo.ca/~galuc/wsdbm/Country0")
	cpt := 0
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// fetch the last inserted triple into the graph
		for _ = range graph.Filter(subj, pred, obj) {
			cpt++
		}
	}
}

func BenchmarkAllFilterSubsetTreeGraph(b *testing.B) {
	graph := NewTreeGraph()
	graph.LoadFromFile("./watdiv.30k.nt", "nt")
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
