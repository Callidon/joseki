// Copyright (c) 2016 Thomas Minier. All rights reserved.
// Use of this source code is governed by a MIT License
// license that can be found in the LICENSE file.

package rdf

// Binding represent a binding, a pair (variable, RDF term), as described in SPARQL 1.1 language reference.
//
// SPARQL 1.1 reference : https://www.w3.org/TR/sparql11-query/
type Binding struct {
    Variable string
    Value Node
}

// NewBinding creates a new Binding.
func NewBinding(variable string, value Node) Binding {
    return Binding{variable, value}
}

// Equals is a function that compare two bindings and return True if they are equals, False otherwise.
func (b Binding) Equals(other Binding) (bool, error) {
    if b.Variable == other.Variable {
        return true, nil
    }
    return b.Value.Equals(other.Value)
}
