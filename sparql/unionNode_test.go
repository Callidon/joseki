// Copyright (c) 2016 Thomas Minier. All rights reserved.
// Use of this source code is governed by a MIT License
// license that can be found in the LICENSE file.

package sparql

import (
	"github.com/Callidon/joseki/rdf"
	"testing"
)

func TestExecuteUnionNode(t *testing.T) {
	tripleA := rdf.NewTriple(rdf.NewVariable("v1"),
		rdf.NewURI("http://purl.org/dc/terms/title"),
		rdf.NewVariable("v2"))
	tripleB := rdf.NewTriple(rdf.NewVariable("v3"),
		rdf.NewURI("http://www.w3.org/1999/02/22-rdf-syntax-ns#type"),
		rdf.NewURI("http://xmlns.com/foaf/0.1/Document"))
	nodeA := newTripleNode(tripleA, smallGraph, -1, 0)
	nodeB := newTripleNode(tripleB, smallGraph, -1, 0)
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
		t.Error("expected", len(datas), "bindings but instead found", cpt, "bindings")
	}
}

func TestComplexExecuteUnionNode(t *testing.T) {
	tripleA := rdf.NewTriple(rdf.NewVariable("v1"),
		rdf.NewURI("http://www.geonames.org/ontology#parentCountry"),
		rdf.NewVariable("v2"))
	tripleB := rdf.NewTriple(rdf.NewVariable("v1"),
		rdf.NewURI("http://schema.org/eligibleQuantity"),
		rdf.NewVariable("v2"))
	nodeA := newTripleNode(tripleA, bigGraph, -1, 0)
	nodeB := newTripleNode(tripleB, bigGraph, -1, 0)
	union := newUnionNode(nodeA, nodeB)
	expected := 1618
	cpt := 0

	for _ = range union.execute() {
		cpt++
	}

	if cpt != expected {
		t.Error("expected", expected, "bindings but instead got", cpt, "bindings")
	}

	// test if the union is commutative
	union = newUnionNode(nodeB, nodeA)
	cpt = 0
	for _ = range union.execute() {
		cpt++
	}

	if cpt != expected {
		t.Error("union operation should be commutative : expected", expected, "bindings but instead got", cpt, "bindings")
	}
}

func TestComplexExecuteWithUnionNode(t *testing.T) {
	tripleA := rdf.NewTriple(rdf.NewVariable("v1"),
		rdf.NewURI("http://www.geonames.org/ontology#parentCountry"),
		rdf.NewVariable("v2"))
	tripleB := rdf.NewTriple(rdf.NewVariable("v1"),
		rdf.NewURI("http://schema.org/eligibleQuantity"),
		rdf.NewVariable("v2"))
	nodeA := newTripleNode(tripleA, bigGraph, -1, 0)
	nodeB := newTripleNode(tripleB, bigGraph, -1, 0)
	union := newUnionNode(nodeA, nodeB)
	expected := 1
	cpt := 0

	group := rdf.NewBindingsGroup()
	group.Bindings["v1"] = rdf.NewURI("http://db.uwaterloo.ca/~galuc/wsdbm/City125")

	for _ = range union.executeWith(group) {
		cpt++
	}

	if cpt != expected {
		t.Error("expected", expected, "bindings but instead got", cpt, "bindings")
	}

	// test if the union is commutative
	union = newUnionNode(nodeB, nodeA)
	cpt = 0
	for _ = range union.executeWith(group) {
		cpt++
	}

	if cpt != expected {
		t.Error("union operation should be commutative : expected", expected, "bindings but instead got", cpt, "bindings")
	}
}

func TestExecuteNoResultUnionNode(t *testing.T) {
	tripleA := rdf.NewTriple(rdf.NewVariable("v1"),
		rdf.NewURI("http://example.org/funny-predicate"),
		rdf.NewVariable("v2"))
	tripleB := rdf.NewTriple(rdf.NewVariable("v1"),
		rdf.NewURI("http://www.w3.org/1999/02/22-rdf-syntax-ns#type"),
		rdf.NewURI("https://schema.org/Book"))
	nodeA := newTripleNode(tripleA, smallGraph, -1, 0)
	nodeB := newTripleNode(tripleB, smallGraph, -1, 0)
	union := newUnionNode(nodeA, nodeB)
	cpt := 0

	for _ = range union.execute() {
		cpt++
	}

	if cpt > 0 {
		t.Error("shouldn't find any bindings, but instead got", cpt, "bindings")
	}
}

func TestBindingNamesUnionNode(t *testing.T) {
	tripleA := rdf.NewTriple(rdf.NewVariable("v1"),
		rdf.NewURI("http://purl.org/dc/terms/title"),
		rdf.NewVariable("v2"))
	tripleB := rdf.NewTriple(rdf.NewVariable("v1"),
		rdf.NewURI("http://www.w3.org/1999/02/22-rdf-syntax-ns#type"),
		rdf.NewVariable("v3"))
	nodeA := newTripleNode(tripleA, nil, -1, 0)
	nodeB := newTripleNode(tripleB, nil, -1, 0)
	union := newUnionNode(nodeA, nodeB)

	expected := []string{"v1", "v2", "v3"}
	cpt := 0

	for _, bindingName := range union.bindingNames() {
		if expected[cpt] != bindingName {
			t.Error("expected", expected[cpt], "but instead got", bindingName)
		}
		cpt++
	}
}

func TestEqualsUnionNode(t *testing.T) {
	tripleA := rdf.NewTriple(rdf.NewVariable("v1"),
		rdf.NewURI("http://purl.org/dc/terms/title"),
		rdf.NewVariable("v2"))
	tripleB := rdf.NewTriple(rdf.NewVariable("v1"),
		rdf.NewURI("http://www.w3.org/1999/02/22-rdf-syntax-ns#type"),
		rdf.NewURI("http://xmlns.com/foaf/0.1/Document"))
	nodeA := newTripleNode(tripleA, nil, -1, 0)
	nodeB := newTripleNode(tripleB, nil, -1, 0)
	union := newUnionNode(nodeA, nodeB)
	otherUnion := newUnionNode(nodeB, nodeA)

	if !union.Equals(union) {
		t.Error(union, "should be equal to itself")
	}
	if union.Equals(otherUnion) {
		t.Error(union, "shouldn't be equals to", otherUnion)
	}
	if union.Equals(nodeA) {
		t.Error(union, "shouldn't be equals to", nodeA)
	}
}

func TestStringUnionNode(t *testing.T) {
	tripleA := rdf.NewTriple(rdf.NewVariable("v1"),
		rdf.NewURI("http://purl.org/dc/terms/title"),
		rdf.NewVariable("v2"))
	tripleB := rdf.NewTriple(rdf.NewVariable("v1"),
		rdf.NewURI("http://www.w3.org/1999/02/22-rdf-syntax-ns#type"),
		rdf.NewURI("http://xmlns.com/foaf/0.1/Document"))
	nodeA := newTripleNode(tripleA, nil, -1, 0)
	nodeB := newTripleNode(tripleB, nil, -1, 0)
	union := newUnionNode(nodeA, nodeB)
	expected := "UNION (" + nodeA.String() + ", " + nodeB.String() + ")"

	if union.String() != expected {
		t.Error(union.String(), "should be equals to", expected)
	}
}
