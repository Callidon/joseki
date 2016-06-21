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
	for _ = range graph.Filter(subj, rdf.NewBlankNode("v"), rdf.NewBlankNode("w")) {
		cpt++
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
	for _ = range graph.Filter(subj, rdf.NewBlankNode("v"), rdf.NewBlankNode("w")) {
		cpt++
	}

	if cpt != nbDatas {
		t.Error("expected ", nbDatas, "results but got ", cpt, "results")
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
	graph.Delete(subj, rdf.NewBlankNode("v"), rdf.NewBlankNode("w"))

	// select all triple of the graph
	for _ = range graph.Filter(subj, rdf.NewBlankNode("v"), rdf.NewBlankNode("w")) {
		cpt++
	}

	if cpt > 0 {
		t.Error("Error : the graph should be empty")
	}
}

func TestLoadFromFileHDTGraph(t *testing.T) {
	graph := NewListGraph()
	cpt := 0
	graph.LoadFromFile("../parser/datas/test.nt", "nt")

	// select all triple of the graph
	for _ = range graph.Filter(rdf.NewBlankNode("y"), rdf.NewBlankNode("v"), rdf.NewBlankNode("w")) {
		cpt++
	}

	if cpt != 4 {
		t.Error("the graph should contains 4 triples, but it contains", cpt, "triples")
	}
}

// Benchmarking

func BenchmarkAddHDTGraph(b *testing.B) {
	graph := NewHDTGraph()
	nbDatas := 1000
	var datas []rdf.Triple

	// create triples to be inserted
	for i := 0; i < nbDatas; i++ {
		triple := rdf.NewTriple(rdf.NewURI(string(rand.Int())), rdf.NewURI(string(rand.Int())), rdf.NewLiteral(string(rand.Int())))
		datas = append(datas, triple)
	}

	for i := 0; i < b.N; i++ {
		for _, triple := range datas {
			graph.Add(triple)
		}
	}
}

func BenchmarkDeleteAllHDTGraph(b *testing.B) {
	graph := NewHDTGraph()
	nbDatas := 1000
	subj := rdf.NewURI("dblp:foo")

	// insert random triples in the graph
	for i := 0; i < nbDatas; i++ {
		triple := rdf.NewTriple(subj, rdf.NewURI(string(rand.Intn(nbDatas))), rdf.NewLiteral(string(rand.Intn(nbDatas))))
		graph.Add(triple)
	}

	for i := 0; i < b.N; i++ {
		graph.Delete(subj, rdf.NewBlankNode("v"), rdf.NewBlankNode("w"))
	}
}

func BenchmarkAllFilterHDTGraph(b *testing.B) {
	graph := NewHDTGraph()
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
			cpt++
		}
	}
}

func BenchmarkSpecificFilterHDTGraph(b *testing.B) {
	graph := NewHDTGraph()
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
			cpt++
		}
	}
}
