// Copyright (c) 2016 Thomas Minier. All rights reserved.
// Use of this source code is governed by a MIT License
// license that can be found in the LICENSE file.

package rdf

// BindingsGroup represents a group of bindings, a pair (variable, RDF term), as described in SPARQL 1.1 language reference.
//
// SPARQL 1.1 reference : https://www.w3.org/TR/sparql11-query/
type BindingsGroup struct {
	Bindings map[string]Node
}

// NewBindingsGroup creates a new BindingsGroup.
func NewBindingsGroup() BindingsGroup {
	return BindingsGroup{make(map[string]Node)}
}

// Equals is a function that compare two group of bindings and return True if they are equals, False otherwise.
func (b BindingsGroup) Equals(other BindingsGroup) (bool, error) {
	for key, value := range b.Bindings {
		if _, inOther := other.Bindings[key]; !inOther {
			return false, nil
		}
		otherValue, _ := other.Bindings[key]
		test, err := value.Equals(otherValue)
		if !test && (err != nil) {
			return false, err
		}
	}
	return true, nil
}

// Clone creates a duplicate of the group of bindings
func (b BindingsGroup) Clone() BindingsGroup {
	newGroup := NewBindingsGroup()
	for key, value := range b.Bindings {
		newGroup.Bindings[key] = value
	}
	return newGroup
}
