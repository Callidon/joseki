// Copyright (c) 2016 Thomas Minier. All rights reserved.
// Use of this source code is governed by a MIT License
// license that can be found in the LICENSE file.

package sparql

import (
	"github.com/Callidon/joseki/graph"
	"github.com/Callidon/joseki/rdf"
	"testing"
)

func TestExecuteSelectNode(t *testing.T) {
	g := graph.NewHDTGraph()
	g.LoadFromFile("../parser/datas/test.nt", "nt")
	triple := rdf.NewTriple(rdf.NewVariable("v1"),
		rdf.NewURI("http://purl.org/dc/terms/title"),
		rdf.NewVariable("v2"))
	node := newTripleNode(triple, g, -1, -1)
	selectNode := newSelectNode(node, []string{"v1"}...)

	datas := []rdf.BindingsGroup{
		rdf.NewBindingsGroup(),
		rdf.NewBindingsGroup(),
	}
	datas[0].Bindings["v1"] = rdf.NewURI("http://www.w3.org/2001/sw/RDFCore/ntriples")
	datas[1].Bindings["v1"] = rdf.NewURI("http://www.w3.org/2001/sw/RDFCore/ntriples")
	cpt := 0

	for bindings := range selectNode.execute() {
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

func TestExecuteWithSelectNode(t *testing.T) {
	g := graph.NewHDTGraph()
	g.LoadFromFile("../parser/datas/test.nt", "nt")
	triple := rdf.NewTriple(rdf.NewURI("http://www.w3.org/2001/sw/RDFCore/ntriples"),
		rdf.NewVariable("v2"),
		rdf.NewVariable("v1"))
	node := newTripleNode(triple, g, -1, -1)
	selectNode := newSelectNode(node, []string{"v1"}...)
	group := rdf.NewBindingsGroup()
	group.Bindings["v2"] = rdf.NewURI("http://purl.org/dc/terms/title")

	datas := []rdf.BindingsGroup{
		rdf.NewBindingsGroup(),
		rdf.NewBindingsGroup(),
	}
	datas[0].Bindings["v1"] = rdf.NewLangLiteral("N-Triples", "en")
	datas[1].Bindings["v1"] = rdf.NewTypedLiteral("My Typed Literal", "<http://www.w3.org/2001/XMLSchema#string>")
	cpt := 0

	for bindings := range selectNode.executeWith(group) {
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

func TestBindingNamesSelectNode(t *testing.T) {
	triple := rdf.NewTriple(rdf.NewVariable("v1"),
		rdf.NewURI("http://purl.org/dc/terms/title"),
		rdf.NewVariable("v2"))
	node := newTripleNode(triple, nil, -1, -1)
	selectNode := newSelectNode(node, []string{"v1"}...)

	expectedNames := []string{"v1", "v2"}
	cpt := 0

	for _, name := range selectNode.bindingNames() {
		if name != expectedNames[cpt] {
			t.Error("expected", expectedNames[cpt], "but instead got", name)
		}
		cpt++
	}
}

func TestEqualsSelectNode(t *testing.T) {
	tripleA := rdf.NewTriple(rdf.NewURI("http://example.org#subjectA"),
		rdf.NewURI("http://purl.org/dc/terms/title"),
		rdf.NewLiteral("Harry Potter 2"))
	tripleB := rdf.NewTriple(rdf.NewURI("http://example.org#subjectB"),
		rdf.NewURI("http://purl.org/dc/terms/title"),
		rdf.NewLiteral("World of Warcraft"))
	tnodeA := newTripleNode(tripleA, nil, -1, 0)
	tnodeB := newTripleNode(tripleB, nil, -1, 0)
	snode := newSelectNode(tnodeA, []string{"v1"}...)
	otherSnode := newSelectNode(tnodeB, []string{"v1"}...)

	if !snode.Equals(snode) {
		t.Error(snode, "should be equal to itself")
	}
	if snode.Equals(otherSnode) {
		t.Error(snode, "shouldn't be equals to", otherSnode)
	}
	if snode.Equals(tnodeA) {
		t.Error(snode, "shouldn't be equals to", tnodeA)
	}
}

func TestStringSelectNode(t *testing.T) {
	triple := rdf.NewTriple(rdf.NewVariable("v1"),
		rdf.NewURI("http://purl.org/dc/terms/title"),
		rdf.NewVariable("v2"))
	tnode := newTripleNode(triple, nil, -1, 0)
	snode := newSelectNode(tnode, []string{"v1", "v2"}...)
	expected := "SELECT v1,v2 (" + tnode.String() + ")"

	if snode.String() != expected {
		t.Error(snode.String(), "should be equals to", expected)
	}
}
