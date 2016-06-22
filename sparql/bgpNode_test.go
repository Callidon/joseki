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
	triple := rdf.NewTriple(rdf.NewURI("http://www.w3.org/2001/sw/RDFCore/ntriples/"),
		rdf.NewURI("http://purl.org/dc/terms/title"),
		rdf.NewBlankNode("v1"))
	node := newBgpNode(triple, graph)

	datas := []rdf.Binding{
		rdf.NewBinding("v1", rdf.NewLangLiteral("N-Triples", "en")),
		rdf.NewBinding("v1", rdf.NewTypedLiteral("My Typed Literal", "<http://www.w3.org/2001/XMLSchema#string>")),
	}
	cpt := 0

	for binding := range node.execute() {
		testA, errA := binding.Equals(datas[0])
		testB, errB := binding.Equals(datas[1])
		if (!testA || errA != nil) || (!testB || errB != nil) {
			t.Error(binding, "should be one of", datas)
		}
		cpt++
	}

	if cpt != len(datas) {
		t.Error("read", cpt, "bindings found instead of", len(datas))
	}
}

func TestExecuteWithBGPNode(t *testing.T) {
	var graph = graph.NewHDTGraph()
	graph.LoadFromFile("../parser/datas/test.nt", "nt")
	triple := rdf.NewTriple(rdf.NewURI("http://www.w3.org/2001/sw/RDFCore/ntriples/"),
		rdf.NewBlankNode("v2"),
		rdf.NewBlankNode("v1"))
	node := newBgpNode(triple, graph)
	binding := rdf.NewBinding("v2", rdf.NewURI("http://purl.org/dc/terms/title"))

	datas := []rdf.Binding{
		rdf.NewBinding("v1", rdf.NewLangLiteral("N-Triples", "en")),
		rdf.NewBinding("v1", rdf.NewTypedLiteral("My Typed Literal", "<http://www.w3.org/2001/XMLSchema#string>")),
	}
	cpt := 0

	for binding := range node.executeWith(binding) {
		testA, errA := binding.Equals(datas[0])
		testB, errB := binding.Equals(datas[1])
		if (!testA || errA != nil) || (!testB || errB != nil) {
			t.Error(binding, "should be one of", datas)
		}
		cpt++
	}

	if cpt != len(datas) {
		t.Error("read", cpt, "bindings found instead of", len(datas))
	}
}
