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

// Test the Equals operator for Triple struct
func TestTripleEquivalent(t *testing.T) {
	tripleA := NewTriple(NewURI("foaf:foo"), NewURI("schema:bar"), NewLiteral("22"))
	tripleB := NewTriple(NewURI("schema:bar"), NewURI("foaf:foo"), NewLiteral("22"))
	tripleC := NewTriple(NewBlankNode("v"), NewURI("schema:bar"), NewLiteral("22"))

	if test, err := tripleA.Equivalent(tripleA); !test || (err != nil) {
		t.Error("a triple should be equivalent to itself")
	}
	if test, err := tripleA.Equivalent(tripleB); (err != nil) && test {
		t.Error(tripleA, "cannot be equivalent to", tripleB)
	}
	if test, err := tripleA.Equivalent(tripleC); !test || (err != nil) {
		t.Error(tripleA, "should be equivalent to", tripleC)
	}
}
