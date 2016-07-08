// Copyright (c) 2016 Thomas Minier. All rights reserved.
// Use of this source code is governed by a MIT License
// license that can be found in the LICENSE file.

package rdf

import "testing"

// Test the Equals operator of the URI struct
func TestURIEquals(t *testing.T) {
	uri := NewURI("dblp:foo")
	otherURI := NewURI("foaf:hasFriend")
	literal := NewLiteral("Toto")
	bnode := NewBlankNode("v")
	variable := NewVariable("x")

	if test, err := uri.Equals(uri); !test || (err != nil) {
		t.Error("a URI must be equals to itself")
	}

	if test, err := uri.Equals(otherURI); test && (err == nil) {
		t.Error(uri, "must be different of", otherURI)
	}

	if test, err := uri.Equals(literal); test && (err == nil) {
		t.Error("a URI and a Literal cannot be equal")
	}

	if test, err := uri.Equals(bnode); !test && (err == nil) {
		t.Error("a URI and a Blank Node cannot be equal")
	}

	if test, err := uri.Equals(variable); !test || (err != nil) {
		t.Error("a URI and a Variable are always equals")
	}
}

// Test the Equals operator of the Literal struct
func TestLiteralEquals(t *testing.T) {
	uri := NewURI("dblp:foo")
	literal := NewLiteral("Toto")
	otherLiteral := NewLiteral("20")
	bnode := NewBlankNode("v")
	variable := NewVariable("x")

	if test, err := literal.Equals(literal); !test || (err != nil) {
		t.Error("a Literal must be equals to itself")
	}

	if test, err := literal.Equals(otherLiteral); test && (err == nil) {
		t.Error(literal, "must be different of", otherLiteral)
	}

	if test, err := literal.Equals(uri); test && (err == nil) {
		t.Error("a Literal and a URI cannot be equal")
	}

	if test, err := literal.Equals(bnode); test && (err == nil) {
		t.Error("a Literal and a Blank Node cannot be equal")
	}

	if test, err := literal.Equals(variable); !test || (err != nil) {
		t.Error("a Literal and a Variable are always equals")
	}
}

// Test the Equals operator of the BlankNode struct
func TestBlankNodeEquals(t *testing.T) {
	uri := NewURI("dblp:foo")
	literal := NewLiteral("Toto")
	bnode := NewBlankNode("v")
	otherBnode := NewBlankNode("w")
	variable := NewVariable("x")

	if test, err := bnode.Equals(bnode); !test || (err != nil) {
		t.Error("a Blank Node must be equals to itself")
	}

	if test, err := bnode.Equals(otherBnode); test && (err == nil) {
		t.Error(bnode, "must be different of", otherBnode)
	}

	if test, err := bnode.Equals(uri); test && (err == nil) {
		t.Error("a Blank Node and a URI cannot be equal")
	}

	if test, err := bnode.Equals(literal); test && (err == nil) {
		t.Error("a Blank Node and a Literal cannot be equal")
	}

	if test, err := bnode.Equals(variable); !test || (err != nil) {
		t.Error("a Blank Node and a Variable are always equals")
	}
}

// Test the Equals operator of the Variable struct
func TestVariableEquals(t *testing.T) {
	uri := NewURI("dblp:foo")
	literal := NewLiteral("Toto")
	bnode := NewBlankNode("v")
	variable := NewVariable("x")
	otherVariable := NewVariable("w")

	if test, err := variable.Equals(variable); !test || (err != nil) {
		t.Error("a Variable must be equals to itself")
	}

	if test, err := variable.Equals(otherVariable); test && (err == nil) {
		t.Error(variable, "must be different of", otherVariable)
	}

	if test, err := variable.Equals(uri); !test || (err != nil) {
		t.Error("a Variable and a URI are always equals")
	}

	if test, err := variable.Equals(literal); !test || (err != nil) {
		t.Error("a Variable and a Literal are always equals")
	}

	if test, err := bnode.Equals(bnode); !test || (err != nil) {
		t.Error("a Variable and a Blank Node are always equals")
	}
}
