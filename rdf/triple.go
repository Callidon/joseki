// Copyright (c) 2016 Thomas Minier. All rights reserved.
// Use of this source code is governed by a MIT License
// license that can be found in the LICENSE file.

package rdf

// Triple represent a RDF Triple
//
// RDF Triple reference : https://www.w3.org/TR/rdf11-concepts/#section-triples
type Triple struct {
	Subject   Node
	Predicate Node
	Object    Node
}

// NewTriple creates a new Triple.
func NewTriple(subject, predicate, object Node) Triple {
	return (Triple{subject, predicate, object})
}

// Equals is a function that compare two Triples and return True if they are equals, False otherwise.
func (t Triple) Equals(other Triple) (bool, error) {
	testSubj, err := t.Subject.Equals(other.Subject)
	if err != nil {
		return false, err
	}
	testPred, err := t.Predicate.Equals(other.Predicate)
	if err != nil {
		return false, err
	}
	testObj, err := t.Object.Equals(other.Object)
	if err != nil {
		return false, err
	}
	return testSubj && testPred && testObj, nil
}

// Equivalent is a function that determine if two Triple are equivalent.
// Two triples are equivalents if their nodes are equivalents between them.
// Otherwise, the result is always False.
func (t Triple) Equivalent(other Triple) (bool, error) {
	testSubj, err := t.Subject.Equivalent(other.Subject)
	if err != nil {
		return false, err
	}
	testPred, err := t.Predicate.Equivalent(other.Predicate)
	if err != nil {
		return false, err
	}
	testObj, err := t.Object.Equivalent(other.Object)
	if err != nil {
		return false, err
	}
	return testSubj && testPred && testObj, nil
}
