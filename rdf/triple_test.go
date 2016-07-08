// Copyright (c) 2016 Thomas Minier. All rights reserved.
// Use of this source code is governed by a MIT License
// license that can be found in the LICENSE file.

package rdf

import "testing"

// Test the Equals operator for Triple struct
func TestTripleEquals(t *testing.T) {
	tripleA := NewTriple(NewURI("foaf:foo"), NewURI("schema:bar"), NewLiteral("22"))
	tripleB := NewTriple(NewURI("schema:bar"), NewURI("foaf:foo"), NewLiteral("22"))
	tripleC := NewTriple(NewBlankNode("v"), NewURI("schema:bar"), NewLiteral("22"))

	if test, err := tripleA.Equals(tripleA); !test || (err != nil) {
		t.Error("a triple should be equals to itself")
	}
	if test, err := tripleA.Equals(tripleB); (err != nil) && test {
		t.Error(tripleA, "cannot be equals to", tripleB)
	}
	if _, err := tripleA.Equals(tripleC); err == nil {
		t.Error("cannot compare two triples with blank nodes in one of them")
	}
}
