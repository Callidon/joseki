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

// Test the Equals operator for Triple struct
func TestTripleComplete(t *testing.T) {
	var completed Triple
	datas := []Triple{
		NewTriple(NewVariable("x"), NewURI("example.org#pred"), NewLiteral("22")),
		NewTriple(NewVariable("x"), NewVariable("y"), NewLiteral("22")),
		NewTriple(NewVariable("x"), NewVariable("y"), NewVariable("z")),
		NewTriple(NewURI("example.org#subj"), NewVariable("y"), NewLiteral("22")),
		NewTriple(NewURI("example.org#subj"), NewVariable("y"), NewVariable("z")),
		NewTriple(NewVariable("w"), NewURI("example.org#pred"), NewLiteral("22")),
		NewTriple(NewVariable("x"), NewVariable("w"), NewLiteral("22")),
		NewTriple(NewVariable("x"), NewVariable("y"), NewVariable("w")),
		NewTriple(NewURI("example.org#subj"), NewURI("example.org#pred"), NewLiteral("22")),
	}
	expected := []Triple{
		NewTriple(NewURI("example.org#subj"), NewURI("example.org#pred"), NewLiteral("22")),
		NewTriple(NewURI("example.org#subj"), NewURI("example.org#pred"), NewLiteral("22")),
		NewTriple(NewURI("example.org#subj"), NewURI("example.org#pred"), NewLiteral("22")),
		NewTriple(NewURI("example.org#subj"), NewURI("example.org#pred"), NewLiteral("22")),
		NewTriple(NewURI("example.org#subj"), NewURI("example.org#pred"), NewLiteral("22")),
		NewTriple(NewVariable("w"), NewURI("example.org#pred"), NewLiteral("22")),
		NewTriple(NewURI("example.org#subj"), NewVariable("w"), NewLiteral("22")),
		NewTriple(NewURI("example.org#subj"), NewURI("example.org#pred"), NewVariable("w")),
		NewTriple(NewURI("example.org#subj"), NewURI("example.org#pred"), NewLiteral("22")),
	}
	cpt := 0

	group := NewBindingsGroup()
	group.Bindings["x"] = NewURI("example.org#subj")
	group.Bindings["y"] = NewURI("example.org#pred")
	group.Bindings["z"] = NewLiteral("22")

	for _, data := range datas {
		completed = data.Complete(group)
		switch {
		case completed.Subject == nil:
			t.Error("complete", data, "with", group, "shouldn't produce result with nil subject")
		case completed.Predicate == nil:
			t.Error("complete", data, "with", group, "shouldn't produce result with nil predicate")
		case completed.Object == nil:
			t.Error("complete", data, "with", group, "shouldn't produce result with nil object")
		}
		if test, err := expected[cpt].Equals(completed); !test || err != nil {
			t.Error(completed, "should be equal to", expected[cpt])
		}
		cpt++
	}
}
