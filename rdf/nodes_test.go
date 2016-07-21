// Copyright (c) 2016 Thomas Minier. All rights reserved.
// Use of this source code is governed by a MIT License
// license that can be found in the LICENSE file.

package rdf

import "testing"

// Test the Equals operator of the URI struct
func TestURIEquals(t *testing.T) {
	uri := NewURI("http://dblp.org#foo")
	otherURI := NewURI("http://foaf.com/hasFriend")
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

func TestURIString(t *testing.T) {
	uri := NewURI("http://dblp.org#foo")
	expected := "<http://dblp.org#foo>"

	if uri.String() != expected {
		t.Error(uri.String(), "should be equals to", expected)
	}
}

// Test the Equals operator of the Literal struct
func TestLiteralEquals(t *testing.T) {
	uri := NewURI("http://dblp.org#foo")
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

func TestLiteralString(t *testing.T) {
	literal := NewLiteral("22")
	typedLiteral := NewTypedLiteral("22", "http://www.w3.org/2001/XMLSchema#string")
	langLiteral := NewLangLiteral("World of Warcraft", "en")
	expectedLiteral := "\"22\""
	expectedTypedLiteral := "\"22\"^^<http://www.w3.org/2001/XMLSchema#string>"
	expectedLangLiteral := "\"World of Warcraft\"@en"

	if literal.String() != expectedLiteral {
		t.Error(literal.String(), "should be equals to", expectedLiteral)
	}

	if typedLiteral.String() != expectedTypedLiteral {
		t.Error(typedLiteral.String(), "should be equals to", expectedTypedLiteral)
	}

	if langLiteral.String() != expectedLangLiteral {
		t.Error(langLiteral.String(), "should be equals to", expectedLangLiteral)
	}
}

// Test the Equals operator of the BlankNode struct
func TestBlankNodeEquals(t *testing.T) {
	uri := NewURI("http://dblp.org#foo")
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

func TestBlankNodeString(t *testing.T) {
	bnode := NewBlankNode("a")
	expected := "_:a"

	if bnode.String() != expected {
		t.Error(bnode.String(), "should be equals to", expected)
	}
}

// Test the Equals operator of the Variable struct
func TestVariableEquals(t *testing.T) {
	uri := NewURI("http://dblp.org#foo")
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

func TestVariableString(t *testing.T) {
	variable := NewVariable("v")
	expected := "?v"

	if variable.String() != expected {
		t.Error(variable.String(), "should be equals to", expected)
	}
}
