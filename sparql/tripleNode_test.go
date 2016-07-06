// Copyright (c) 2016 Thomas Minier. All rights reserved.
// Use of this source code is governed by a MIT License
// license that can be found in the LICENSE file.

package sparql

import (
	"github.com/Callidon/joseki/graph"
	"github.com/Callidon/joseki/rdf"
	"testing"
)

func TestExecuteBGPNode(t *testing.T) {
	var graph = graph.NewHDTGraph()
	graph.LoadFromFile("../parser/datas/test.nt", "nt")
	triple := rdf.NewTriple(rdf.NewURI("http://www.w3.org/2001/sw/RDFCore/ntriples"),
		rdf.NewURI("http://purl.org/dc/terms/title"),
		rdf.NewBlankNode("v1"))
	node := newTripleNode(triple, graph, -1, -1)
	cpt := 0

	datas := []rdf.BindingsGroup{
		rdf.NewBindingsGroup(),
		rdf.NewBindingsGroup(),
	}
	datas[0].Bindings["v1"] = rdf.NewLangLiteral("N-Triples", "en")
	datas[1].Bindings["v1"] = rdf.NewTypedLiteral("My Typed Literal", "<http://www.w3.org/2001/XMLSchema#string>")

	for bindings := range node.execute() {
		testA, errA := bindings.Equals(datas[0])
		testB, errB := bindings.Equals(datas[1])
		if (!testA || (errA != nil)) && (!testB || (errB != nil)) {
			t.Error(bindings, "should be one of", datas)
		}
		cpt++
	}

	if cpt != len(datas) {
		t.Error(cpt, "bindings found instead of", len(datas))
	}
}

func TestExecuteNoResultBGPNode(t *testing.T) {
	var graph = graph.NewHDTGraph()
	graph.LoadFromFile("../parser/datas/test.nt", "nt")
	triple := rdf.NewTriple(rdf.NewURI("http://www.w3.org/2001/sw/RDFCore/ntriples"),
		rdf.NewURI("http://example.org/funny-predicate"),
		rdf.NewBlankNode("v1"))
	node := newTripleNode(triple, graph, -1, -1)
	cpt := 0

	for _ = range node.execute() {
		cpt++
	}

	if cpt > 0 {
		t.Error("should not found any bindings, but instead found", cpt, "bindings")
	}
}

func TestExecuteWithBGPNode(t *testing.T) {
	var graph = graph.NewHDTGraph()
	graph.LoadFromFile("../parser/datas/test.nt", "nt")
	triple := rdf.NewTriple(rdf.NewURI("http://www.w3.org/2001/sw/RDFCore/ntriples"),
		rdf.NewBlankNode("v2"),
		rdf.NewBlankNode("v1"))
	node := newTripleNode(triple, graph, -1, -1)
	group := rdf.NewBindingsGroup()
	group.Bindings["v2"] = rdf.NewURI("http://purl.org/dc/terms/title")
	cpt := 0

	datas := []rdf.BindingsGroup{
		rdf.NewBindingsGroup(),
		rdf.NewBindingsGroup(),
	}
	datas[0].Bindings["v1"] = rdf.NewLangLiteral("N-Triples", "en")
	datas[0].Bindings["v2"] = rdf.NewURI("http://purl.org/dc/terms/title")
	datas[1].Bindings["v1"] = rdf.NewTypedLiteral("My Typed Literal", "<http://www.w3.org/2001/XMLSchema#string>")
	datas[1].Bindings["v2"] = rdf.NewURI("http://purl.org/dc/terms/title")

	for bindings := range node.executeWith(group) {
		testA, errA := bindings.Equals(datas[0])
		testB, errB := bindings.Equals(datas[1])
		if (!testA || (errA != nil)) && (!testB || (errB != nil)) {
			t.Error(bindings, "should be one of", datas)
		}
		cpt++
	}

	if cpt != len(datas) {
		t.Error(cpt, "bindings found instead of", len(datas))
	}
}
