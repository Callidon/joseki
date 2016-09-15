// Copyright (c) 2016 Thomas Minier. All rights reserved.
// Use of this source code is governed by a MIT License
// license that can be found in the LICENSE file.

package sparql

import (
	"github.com/Callidon/joseki/graph"
	"github.com/Callidon/joseki/rdf"
	"testing"
)

func TestSimpleFindJoin(t *testing.T) {
	var joinFound sparqlNode
	var ind int
	tripleA := rdf.NewTriple(rdf.NewURI("example.org"), rdf.NewURI("foaf:friendOf"), rdf.NewVariable("v1"))
	tripleB := rdf.NewTriple(rdf.NewVariable("v1"), rdf.NewURI("rdf:name"), rdf.NewVariable("v2"))
	nodeA := newTripleNode(tripleA, nil, -1, -1)
	nodeB := newTripleNode(tripleB, nil, -1, -1)

	join := newJoinNode(nodeA, nodeB)
	joinFound, ind = findJoin(nodeA, nodeB)

	if joinFound == nil || (ind == -1) {
		t.Error("expected to find a join between", tripleA, "and", tripleB)
	}

	if !joinFound.Equals(join) || ind != 0 {
		t.Error("expected", join, "but instead found", joinFound)
	}
}

func TestComplexFindJoin(t *testing.T) {
	var joinFound sparqlNode
	var ind int
	tripleA := rdf.NewTriple(rdf.NewURI("example.org"), rdf.NewURI("foaf:friendOf"), rdf.NewVariable("v1"))
	tripleB := rdf.NewTriple(rdf.NewVariable("v2"), rdf.NewURI("rdf:type"), rdf.NewURI("schema.org#People"))
	tripleC := rdf.NewTriple(rdf.NewVariable("v1"), rdf.NewURI("rdf:name"), rdf.NewVariable("v2"))
	nodeA := newTripleNode(tripleA, nil, -1, -1)
	nodeB := newTripleNode(tripleB, nil, -1, -1)
	nodeC := newTripleNode(tripleC, nil, -1, -1)

	firstJoin := newJoinNode(nodeA, nodeC)
	secondJoin := newJoinNode(firstJoin, nodeB)

	joinFound, ind = findJoin(nodeA, nodeB, nodeC)

	if joinFound == nil || (ind == -1) {
		t.Error("expected to find a join between", tripleA, "and", tripleB)
	}

	if !joinFound.Equals(firstJoin) || ind != 1 {
		t.Error("expected", firstJoin, "but instead found", joinFound)
	}

	joinFound, ind = findJoin(joinFound, nodeB)

	if joinFound == nil || (ind == -1) {
		t.Error("expected to find a join between", tripleA, "and", tripleB)
	}

	if !joinFound.Equals(secondJoin) || ind != 0 {
		t.Error("expected", secondJoin, "but instead found", joinFound)
	}
}

func TestBuildQueryDescriptor(t *testing.T) {
	var expectedRoot sparqlNode
	tripleA := rdf.NewTriple(rdf.NewURI("example.org"), rdf.NewURI("foaf:friendOf"), rdf.NewVariable("v1"))
	tripleB := rdf.NewTriple(rdf.NewVariable("v2"), rdf.NewURI("rdf:type"), rdf.NewURI("schema.org#People"))
	tripleC := rdf.NewTriple(rdf.NewVariable("v1"), rdf.NewURI("rdf:name"), rdf.NewVariable("v2"))
	tripleD := rdf.NewTriple(rdf.NewURI("example.org"), rdf.NewURI("rdf:name"), rdf.NewVariable("v5"))

	nodeA := newTripleNode(tripleA, nil, -1, -1)
	nodeB := newTripleNode(tripleB, nil, -1, -1)
	nodeC := newTripleNode(tripleC, nil, -1, -1)
	nodeD := newTripleNode(tripleD, nil, -1, -1)
	expectedRoot = newUnionNode(newJoinNode(newJoinNode(nodeA, nodeC), nodeB), nodeD)

	qd := newQueryDescriptor(nil, selectQuery)
	qd.Where(tripleA, tripleB, tripleC, tripleD)
	root := qd.build()

	if !expectedRoot.Equals(root) {
		t.Error("expected", root, "to be equals to", expectedRoot)
	}
}

func TestLimitQueryDescriptor(t *testing.T) {
	graph := graph.NewTreeGraph()
	graph.LoadFromFile("../parser/datas/test.nt", "nt")
	triple := rdf.NewTriple(rdf.NewURI("http://www.w3.org/2001/sw/RDFCore/ntriples"), rdf.NewVariable("v"), rdf.NewVariable("w"))
	cpt := 0

	qd := newQueryDescriptor(graph, selectQuery)
	qd.Where(triple)
	qd.Limit(2)
	root := qd.build()

	for _ = range root.execute() {
		cpt++
	}

	if cpt != 2 {
		t.Error("expected two results but instead found", cpt)
	}

}
