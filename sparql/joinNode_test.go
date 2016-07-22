// Copyright (c) 2016 Thomas Minier. All rights reserved.
// Use of this source code is governed by a MIT License
// license that can be found in the LICENSE file.

package sparql

import (
	"github.com/Callidon/joseki/rdf"
	"testing"
)

func TestSimpleExecuteJoinNode(t *testing.T) {
	tripleA := rdf.NewTriple(rdf.NewVariable("v1"),
		rdf.NewURI("http://purl.org/dc/terms/title"),
		rdf.NewVariable("v2"))
	tripleB := rdf.NewTriple(rdf.NewVariable("v1"),
		rdf.NewURI("http://www.w3.org/1999/02/22-rdf-syntax-ns#type"),
		rdf.NewURI("http://xmlns.com/foaf/0.1/Document"))
	nodeA := newTripleNode(tripleA, smallGraph, -1, 0)
	nodeB := newTripleNode(tripleB, smallGraph, -1, 0)
	join := newJoinNode(nodeA, nodeB)

	datas := []rdf.BindingsGroup{
		rdf.NewBindingsGroup(),
		rdf.NewBindingsGroup(),
	}
	datas[0].Bindings["v1"] = rdf.NewURI("http://www.w3.org/2001/sw/RDFCore/ntriples")
	datas[0].Bindings["v2"] = rdf.NewLangLiteral("N-Triples", "en")
	datas[1].Bindings["v1"] = rdf.NewURI("http://www.w3.org/2001/sw/RDFCore/ntriples")
	datas[1].Bindings["v2"] = rdf.NewTypedLiteral("My Typed Literal", "<http://www.w3.org/2001/XMLSchema#string>")
	cpt := 0

	for bindings := range join.execute() {
		testA, errA := bindings.Equals(datas[0])
		testB, errB := bindings.Equals(datas[1])
		if (!testA || (errA != nil)) && (!testB || (errB != nil)) {
			t.Error(bindings, "should be one of", datas)
		}
		cpt++
	}

	if cpt != len(datas) {
		t.Error("expected", len(datas), "bindings but instead found", cpt, "bindings")
	}
}

func TestComplexExecuteJoinNode(t *testing.T) {
	tripleA := rdf.NewTriple(rdf.NewVariable("v1"),
		rdf.NewURI("http://purl.org/goodrelations/price"),
		rdf.NewVariable("v2"))
	tripleB := rdf.NewTriple(rdf.NewVariable("v1"),
		rdf.NewURI("http://schema.org/eligibleQuantity"),
		rdf.NewVariable("v3"))
	nodeA := newTripleNode(tripleA, bigGraph, -1, 0)
	nodeB := newTripleNode(tripleB, bigGraph, -1, 0)
	join := newJoinNode(nodeA, nodeB)

	expected := 1378
	cpt := 0
	for _ = range join.execute() {
		cpt++
	}

	if cpt != expected {
		t.Error("expected", expected, "bindings but instead got", cpt, "bindings")
	}

	// test if the join operation is commutative
	join = newJoinNode(nodeB, nodeA)
	cpt = 0
	for _ = range join.execute() {
		cpt++
	}

	if cpt != expected {
		t.Error("join operation should be commutative : expected", expected, "bindings but instead got", cpt, "bindings")
	}
}

func TestComplexExecuteWithJoinNode(t *testing.T) {
	tripleA := rdf.NewTriple(rdf.NewVariable("v1"),
		rdf.NewURI("http://purl.org/goodrelations/price"),
		rdf.NewVariable("v2"))
	tripleB := rdf.NewTriple(rdf.NewVariable("v1"),
		rdf.NewURI("http://schema.org/eligibleQuantity"),
		rdf.NewVariable("v3"))
	nodeA := newTripleNode(tripleA, bigGraph, -1, 0)
	nodeB := newTripleNode(tripleB, bigGraph, -1, 0)
	join := newJoinNode(nodeA, nodeB)

	group := rdf.NewBindingsGroup()
	group.Bindings["v1"] = rdf.NewURI("http://db.uwaterloo.ca/~galuc/wsdbm/Offer10071")

	expected := 1
	cpt := 0
	for _ = range join.executeWith(group) {
		cpt++
	}

	if cpt != expected {
		t.Error("expected", expected, "bindings but instead got", cpt, "bindings")
	}

	// test if the join operation isn't commutative where there is no join
	join = newJoinNode(nodeB, nodeA)
	expected = 1378
	cpt = 0
	for _ = range join.execute() {
		cpt++
	}

	if cpt != expected {
		t.Error("expected", expected, "bindings but instead got", cpt, "bindings")
	}
}

func TestExecuteNoResultJoinNode(t *testing.T) {
	tripleA := rdf.NewTriple(rdf.NewVariable("v1"),
		rdf.NewURI("http://example.org/funny-predicate"),
		rdf.NewVariable("v2"))
	tripleB := rdf.NewTriple(rdf.NewVariable("v1"),
		rdf.NewURI("http://www.w3.org/1999/02/22-rdf-syntax-ns#type"),
		rdf.NewURI("http://xmlns.com/foaf/0.1/Document"))
	nodeA := newTripleNode(tripleA, smallGraph, -1, 0)
	nodeB := newTripleNode(tripleB, smallGraph, -1, 0)
	join := newJoinNode(nodeA, nodeB)
	cpt := 0

	for _ = range join.execute() {
		cpt++
	}

	if cpt > 0 {
		t.Error("shouldn't find any bindings, but instead got", cpt, "bindings")
	}
}

func TestBindingNamesJoinNode(t *testing.T) {
	tripleA := rdf.NewTriple(rdf.NewVariable("v1"),
		rdf.NewURI("http://purl.org/dc/terms/title"),
		rdf.NewVariable("v2"))
	tripleB := rdf.NewTriple(rdf.NewVariable("v1"),
		rdf.NewURI("http://www.w3.org/1999/02/22-rdf-syntax-ns#type"),
		rdf.NewURI("http://xmlns.com/foaf/0.1/Document"))
	nodeA := newTripleNode(tripleA, smallGraph, -1, 0)
	nodeB := newTripleNode(tripleB, smallGraph, -1, 0)
	join := newJoinNode(nodeA, nodeB)

	expected := []string{"v1", "v2"}
	cpt := 0

	for _, bindingName := range join.bindingNames() {
		if expected[cpt] != bindingName {
			t.Error("expected", expected[cpt], "but instead got", bindingName)
		}
		cpt++
	}
}

func TestEqualsJoinNode(t *testing.T) {
	tripleA := rdf.NewTriple(rdf.NewVariable("v1"),
		rdf.NewURI("http://purl.org/dc/terms/title"),
		rdf.NewVariable("v2"))
	tripleB := rdf.NewTriple(rdf.NewVariable("v1"),
		rdf.NewURI("http://www.w3.org/1999/02/22-rdf-syntax-ns#type"),
		rdf.NewURI("http://xmlns.com/foaf/0.1/Document"))
	nodeA := newTripleNode(tripleA, nil, -1, 0)
	nodeB := newTripleNode(tripleB, nil, -1, 0)
	join := newJoinNode(nodeA, nodeB)
	otherJoin := newJoinNode(nodeB, nodeA)

	if !join.Equals(join) {
		t.Error(join, "should be equal to itself")
	}
	if join.Equals(otherJoin) {
		t.Error(join, "shouldn't be equals to", otherJoin)
	}
	if join.Equals(nodeA) {
		t.Error(join, "shouldn't be equals to", nodeA)
	}
}

func TestStringJoinNode(t *testing.T) {
	tripleA := rdf.NewTriple(rdf.NewVariable("v1"),
		rdf.NewURI("http://purl.org/dc/terms/title"),
		rdf.NewVariable("v2"))
	tripleB := rdf.NewTriple(rdf.NewVariable("v1"),
		rdf.NewURI("http://www.w3.org/1999/02/22-rdf-syntax-ns#type"),
		rdf.NewURI("http://xmlns.com/foaf/0.1/Document"))
	nodeA := newTripleNode(tripleA, nil, -1, 0)
	nodeB := newTripleNode(tripleB, nil, -1, 0)
	join := newJoinNode(nodeA, nodeB)
	expected := "JOIN (" + nodeA.String() + ", " + nodeB.String() + ")"

	if join.String() != expected {
		t.Error(join.String(), "should be equals to", expected)
	}
}
