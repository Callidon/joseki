package core

import "testing"

// Test the Equals operator of the URI struct
func TestURIEquals(t *testing.T) {
	uri := NewURI("dblp:foo")
	other_uri := NewURI("foaf:hasFriend")
	literal := NewLiteral("Toto")
	bnode := NewBlankNode("v")

	if test, err := uri.Equals(uri); (!test) && (err != nil) {
		t.Error("a URI must be equals to itself")
	}

	if test, err := uri.Equals(other_uri); test && (err == nil) {
		t.Error(uri, "must be different of", other_uri)
	}

	if test, err := uri.Equals(literal); test && (err == nil) {
		t.Error("a URI and a Literal cannot be equal")
	}

	if test, err := uri.Equals(bnode); test && (err == nil) {
		t.Error("a URI and a Blank Node cannot be equal")
	}
}

// Test the Equivalent operator of the URI struct
func TestURIEquivalent(t *testing.T) {
	uri := NewURI("dblp:foo")
	other_uri := NewURI("foaf:hasFriend")
	literal := NewLiteral("Toto")
	bnode := NewBlankNode("v")

	if test, err := uri.Equivalent(uri); (!test) && (err != nil) {
		t.Error(" when comparing a URI with itself, the result should be true")
	}

	if test, err := uri.Equivalent(other_uri); test && (err == nil) {
		t.Error(uri, "must be different of", other_uri)
	}

	if test, err := uri.Equivalent(bnode); (!test) && (err != nil) {
		t.Error("when comparing a URI to a Blank Node, the result should be true")
	}

	if test, err := uri.Equivalent(literal); test && (err == nil) {
		t.Error("a URI and a Literal cannot be compare")
	}
}

// Test the Equals operator of the Literal struct
func TestLiteralEquals(t *testing.T) {
	uri := NewURI("dblp:foo")
	literal := NewLiteral("Toto")
	other_literal := NewLiteral("20")
	bnode := NewBlankNode("v")

	if test, err := literal.Equals(literal); (!test) && (err != nil) {
		t.Error("a Literal must be equals to itself")
	}

	if test, err := literal.Equals(other_literal); test && (err == nil) {
		t.Error(literal, "must be different of", other_literal)
	}

	if test, err := literal.Equals(uri); test && (err == nil) {
		t.Error("a Literal and a URI cannot be equal")
	}

	if test, err := literal.Equals(bnode); test && (err == nil) {
		t.Error("a URI and a Blank Node cannot be equal")
	}
}

// Test the Equivalent operator of the Literal struct
func TestLiteralEquivalent(t *testing.T) {
	uri := NewURI("dblp:foo")
	literal := NewLiteral("Toto")
	other_literal := NewLiteral("20")
	bnode := NewBlankNode("v")

	if test, err := literal.Equivalent(literal); (!test) && (err != nil) {
		t.Error("when comparing a Literal with itself, the result should be true")
	}

	if test, err := literal.Equivalent(other_literal); test && (err == nil) {
		t.Error(literal, "must be different of", other_literal)
	}

	if test, err := literal.Equivalent(bnode); (!test) && (err != nil) {
		t.Error("when comparing a Literal to a Blank Node, the result should be true")
	}

	if test, err := literal.Equivalent(uri); test && (err == nil) {
		t.Error("a Literal and a URI cannot be compare")
	}
}

// Test the Equals operator of the BlankNode struct
func TestBlankNodeEquals(t *testing.T) {
	uri := NewURI("dblp:foo")
	literal := NewLiteral("Toto")
	bnode := NewBlankNode("v")
	other_bnode := NewBlankNode("w")

	if test, err := bnode.Equals(bnode); (!test) && (err != nil) {
		t.Error("a Blank Node must be equals to itself")
	}

	if test, err := bnode.Equals(other_bnode); test && (err == nil) {
		t.Error(bnode, "must be different of", other_bnode)
	}

	if test, err := bnode.Equals(uri); test && (err == nil) {
		t.Error("a Blank and a URI cannot be equal")
	}

	if test, err := bnode.Equals(literal); test && (err == nil) {
		t.Error("a Blank Node and a Literal cannot be equal")
	}
}

// Test the Equivalent operator of the BlankNode struct
func TestBlankNodeEquivalent(t *testing.T) {
	uri := NewURI("dblp:foo")
	literal := NewLiteral("Toto")
	bnode := NewBlankNode("v")
	other_bnode := NewBlankNode("w")

	if test, err := bnode.Equivalent(bnode); (!test) && (err != nil) {
		t.Error("when comparing two Blank Node, the result should be true")
	}

	if test, err := bnode.Equivalent(other_bnode); (!test) && (err != nil) {
		t.Error("when comparing two Blank Node, the result should be true")
	}

	if test, err := bnode.Equivalent(uri); (!test) && (err != nil) {
		t.Error("when comparing a Blank Node with a URI, the result should be true")
	}

	if test, err := bnode.Equivalent(literal); (!test) && (err != nil) {
		t.Error("when comparing a Blank Node with a Literal, the result should be true")
	}
}

// Test the Equals operator for Triple struct
func TestTripleEquals(t *testing.T) {
	tripleA := NewTriple(NewURI("foaf:foo"), NewURI("schema:bar"), NewLiteral("22"))
	tripleB := NewTriple(NewURI("schema:bar"), NewURI("foaf:foo"), NewLiteral("22"))
	tripleC := NewTriple(NewBlankNode("v"), NewURI("schema:bar"), NewLiteral("22"))

	if test, err := tripleA.Equals(tripleA); (err != nil) && !test {
		t.Error("a triple should be equals to itself")
	}
	if test, err := tripleA.Equals(tripleB); (err != nil) && test {
		t.Error(tripleA, "cannot be equals to", tripleB)
	}
	if _, err := tripleA.Equals(tripleC); err == nil {
		t.Error("cannot compare two triples with blank nodes in one of them")
	}
}

// Test the Equals operator for Triple struct
func TestTripleEquivalent(t *testing.T) {
	tripleA := NewTriple(NewURI("foaf:foo"), NewURI("schema:bar"), NewLiteral("22"))
	tripleB := NewTriple(NewURI("schema:bar"), NewURI("foaf:foo"), NewLiteral("22"))
	tripleC := NewTriple(NewBlankNode("v"), NewURI("schema:bar"), NewLiteral("22"))

	if test, err := tripleA.Equivalent(tripleA); (err != nil) && !test {
		t.Error("a triple should be equivalent to itself")
	}
	if test, err := tripleA.Equivalent(tripleB); (err != nil) && test {
		t.Error(tripleA, "cannot be equivalent to", tripleB)
	}
	if test, err := tripleA.Equivalent(tripleC); (err != nil) && !test {
		t.Error(tripleA, "should be equivalent to", tripleC)
	}
}
