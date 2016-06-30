// Copyright (c) 2016 Thomas Minier. All rights reserved.
// Use of this source code is governed by a MIT License
// license that can be found in the LICENSE file.

package sparql

import (
	"github.com/Callidon/joseki/rdf"
	"testing"
)

func TestSimpleFindJoin(t *testing.T) {
	var joinFound sparqlNode
	var ind int
	tripleA := rdf.NewTriple(rdf.NewURI("example.org"), rdf.NewURI("foaf:friendOf"), rdf.NewBlankNode("v1"))
	tripleB := rdf.NewTriple(rdf.NewBlankNode("v1"), rdf.NewURI("rdf:name"), rdf.NewBlankNode("v2"))
	nodeA := newTripleNode(tripleA, nil)
	nodeB := newTripleNode(tripleB, nil)

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
	tripleA := rdf.NewTriple(rdf.NewURI("example.org"), rdf.NewURI("foaf:friendOf"), rdf.NewBlankNode("v1"))
	tripleB := rdf.NewTriple(rdf.NewBlankNode("v2"), rdf.NewURI("rdf:type"), rdf.NewURI("schema.org#People"))
	tripleC := rdf.NewTriple(rdf.NewBlankNode("v1"), rdf.NewURI("rdf:name"), rdf.NewBlankNode("v2"))
	nodeA := newTripleNode(tripleA, nil)
	nodeB := newTripleNode(tripleB, nil)
	nodeC := newTripleNode(tripleC, nil)

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
