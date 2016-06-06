package core

import "testing"

func TestURIEquals(t *testing.T) {
	uri := NewURI("", "toto.org")
	literal := NewLiteral("Toto")
	bnode := NewBlankNode("v")

	if test, err := uri.Equals(uri); (!test) && (err != nil) {
		t.Error("Error : an URI must be equals to itself")
	}

	if test, err := uri.Equals(literal); test && (err == nil) {
		t.Error("Error : an URI and a Literal cannot be equal")
	}

	if test, err := uri.Equals(bnode); test && (err == nil) {
		t.Error("Error : an URI and a Blank Node cannot be equal")
	}
}

func TestLiteralEquals(t *testing.T) {
	uri := NewURI("", "toto.org")
	literal := NewLiteral("Toto")
	bnode := NewBlankNode("v")

	if test, err := literal.Equals(literal); (!test) && (err != nil) {
		t.Error("Error : a Literal must be equals to itself")
	}

	if test, err := literal.Equals(uri); test && (err == nil) {
		t.Error("Error : a Literal and a URI cannot be equal")
	}

	if test, err := literal.Equals(bnode); test && (err == nil) {
		t.Error("Error : an URI and a Blank Node cannot be equal")
	}
}

func TestBlankNodeEquals(t *testing.T) {
	uri := NewURI("", "toto.org")
	literal := NewLiteral("Toto")
	bnode := NewBlankNode("v")

	if test, err := bnode.Equals(bnode); (!test) && (err != nil) {
		t.Error("Error : a Blank Node must be equals to itself")
	}

	if test, err := bnode.Equals(uri); test && (err == nil) {
		t.Error("Error : a Blank and a URI cannot be equal")
	}

	if test, err := bnode.Equals(literal); test && (err == nil) {
		t.Error("Error : a Blank Node and a Literal cannot be equal")
	}
}
