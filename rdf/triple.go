// Copyright (c) 2016 Thomas Minier. All rights reserved.
// Use of this source code is governed by a MIT License
// license that can be found in the LICENSE file.

package rdf

// Triple represents a RDF Triple
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

// Complete use a group of bindings to complete the variable in the triple pattern
// and then return a new completed Triple pattern
func (t Triple) Complete(group BindingsGroup) Triple {
	newSubj, newPred, newObj := t.Subject, t.Predicate, t.Object
	// find the nodes of the triple wich can be completed
	subject, freeSubject := t.Subject.(Variable)
	predicate, freePredicate := t.Predicate.(Variable)
	object, freeObject := t.Object.(Variable)
	if !freeSubject {
		newSubj = t.Subject
	}
	if !freePredicate {
		newPred = t.Predicate
	}
	if !freeObject {
		newObj = t.Object
	}

	for key, binding := range group.Bindings {
		// try to complete any node of the triple using the current binding
		if freeSubject && subject.Value == key {
			newSubj = binding
		}
		if freePredicate && predicate.Value == key {
			newPred = binding
		}
		if freeObject && object.Value == key {
			newObj = binding
		}
	}
	return NewTriple(newSubj, newPred, newObj)
}
