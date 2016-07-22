// Copyright (c) 2016 Thomas Minier. All rights reserved.
// Use of this source code is governed by a MIT License
// license that can be found in the LICENSE file.

package sparql

import (
	"github.com/Callidon/joseki/graph"
	"github.com/Callidon/joseki/rdf"
	"testing"
)

func TestExecuteUnionNode(t *testing.T) {
	g := graph.NewHDTGraph()
	g.LoadFromFile("../parser/datas/test.nt", "nt")
	tripleA := rdf.NewTriple(rdf.NewVariable("v1"),
		rdf.NewURI("http://purl.org/dc/terms/title"),
		rdf.NewVariable("v2"))
	tripleB := rdf.NewTriple(rdf.NewVariable("v3"),
		rdf.NewURI("http://www.w3.org/1999/02/22-rdf-syntax-ns#type"),
		rdf.NewURI("http://xmlns.com/foaf/0.1/Document"))
	nodeA := newTripleNode(tripleA, g, -1, 0)
	nodeB := newTripleNode(tripleB, g, -1, 0)
	union := newUnionNode(nodeA, nodeB)

	datas := []rdf.BindingsGroup{
		rdf.NewBindingsGroup(),
		rdf.NewBindingsGroup(),
		rdf.NewBindingsGroup(),
	}
	datas[0].Bindings["v1"] = rdf.NewURI("http://www.w3.org/2001/sw/RDFCore/ntriples")
	datas[0].Bindings["v2"] = rdf.NewLangLiteral("N-Triples", "en")
	datas[1].Bindings["v1"] = rdf.NewURI("http://www.w3.org/2001/sw/RDFCore/ntriples")
	datas[1].Bindings["v2"] = rdf.NewTypedLiteral("My Typed Literal", "<http://www.w3.org/2001/XMLSchema#string>")
	datas[2].Bindings["v3"] = rdf.NewURI("http://www.w3.org/2001/sw/RDFCore/ntriples")
	cpt := 0

	for bindings := range union.execute() {
		testA, errA := bindings.Equals(datas[0])
		testB, errB := bindings.Equals(datas[1])
		testC, errC := bindings.Equals(datas[2])
		if (!testA || (errA != nil)) && (!testB || (errB != nil)) && (!testC || (errC != nil)) {
			t.Error(bindings, "should be one of", datas)
		}
		cpt++
	}

	if cpt != len(datas) {
		t.Error(cpt, "bindings found instead of", len(datas))
	}
}

func TestExecuteNoResultUnionNode(t *testing.T) {
	g := graph.NewHDTGraph()
	g.LoadFromFile("../parser/datas/test.nt", "nt")
	tripleA := rdf.NewTriple(rdf.NewVariable("v1"),
		rdf.NewURI("http://example.org/funny-predicate"),
		rdf.NewVariable("v2"))
	tripleB := rdf.NewTriple(rdf.NewVariable("v1"),
		rdf.NewURI("http://www.w3.org/1999/02/22-rdf-syntax-ns#type"),
		rdf.NewURI("https://schema.org/Book"))
	nodeA := newTripleNode(tripleA, g, -1, 0)
	nodeB := newTripleNode(tripleB, g, -1, 0)
	union := newUnionNode(nodeA, nodeB)
	cpt := 0

	for _ = range union.execute() {
		cpt++
	}

	if cpt > 0 {
		t.Error("should not found any bindings, but instead found", cpt, "bindings")
	}
}

func TestBindingNamesUnionNode(t *testing.T) {
	g := graph.NewHDTGraph()
	tripleA := rdf.NewTriple(rdf.NewVariable("v1"),
		rdf.NewURI("http://purl.org/dc/terms/title"),
		rdf.NewVariable("v2"))
	tripleB := rdf.NewTriple(rdf.NewVariable("v1"),
		rdf.NewURI("http://www.w3.org/1999/02/22-rdf-syntax-ns#type"),
		rdf.NewURI("http://xmlns.com/foaf/0.1/Document"))
	nodeA := newTripleNode(tripleA, g, -1, 0)
	nodeB := newTripleNode(tripleB, g, -1, 0)
	union := newUnionNode(nodeA, nodeB)

	expected := []string{"v1", "v2"}
	cpt := 0

	for _, bindingName := range union.bindingNames() {
		if expected[cpt] != bindingName {
			t.Error("expected", expected[cpt], "but instead got", bindingName)
		}
		cpt++
	}
}
