// Copyright (c) 2016 Thomas Minier. All rights reserved.
// Use of this source code is governed by a MIT License
// license that can be found in the LICENSE file.

package sparql

import (
    "github.com/Callidon/joseki/rdf"
    "github.com/Callidon/joseki/graph"
)

// BGPNode is the lowest level of SPARQL query execution plan.
// Its role is to retrieve bindings according to a triple pattern from a graph.
type BGPNode struct {
    pattern rdf.Triple
    graph graph.Graph
}

// NewBGPNode creates a new BGPNode.
func NewBGPNode(pattern rdf.Triple, graph graph.Graph) *BGPNode {
    return &BGPNode{pattern, graph}
}

// Execute retrieves bindings from a graph that match a triple pattern.
func (n *BGPNode) Execute() chan rdf.Binding {
    out := make(chan rdf.Binding, bufferSize)
    // find free vars in triple pattern
    subject, freeSubject := n.pattern.Subject.(rdf.BlankNode)
    predicate, freePredicate := n.pattern.Predicate.(rdf.BlankNode)
    object, freeObject := n.pattern.Object.(rdf.BlankNode)

    // retrieves triples & form bindings to send
    go func() {
        defer close(out)
        for triple := range n.graph.Filter(n.pattern.Subject, n.pattern.Predicate, n.pattern.Object) {
            if freeSubject {
                out <- rdf.NewBinding(subject.Variable, triple.Subject)
            }
            if freePredicate {
                out <- rdf.NewBinding(predicate.Variable, triple.Predicate)
            }
            if freeObject {
                out <- rdf.NewBinding(object.Variable, triple.Object)
            }
        }
    }()
    return out
}

// ExecuteWith retrieves bindings from a graph that match a triple pattern, completed by a given binding.
func (n *BGPNode) ExecuteWith(binding rdf.Binding) chan rdf.Binding {
    var querySubj, queryPred, queryObj rdf.Node
    out := make(chan rdf.Binding, bufferSize)
    // find free vars in triple pattern
    subject, freeSubject := n.pattern.Subject.(rdf.BlankNode)
    predicate, freePredicate := n.pattern.Predicate.(rdf.BlankNode)
    object, freeObject := n.pattern.Object.(rdf.BlankNode)

    // complete triple pattern using the binding given in parameter
    if freeSubject && subject.Variable == binding.Variable {
        querySubj = binding.Value
        freeSubject = false
    } else {
        querySubj = n.pattern.Subject
    }
    if freePredicate && predicate.Variable == binding.Variable {
        queryPred = binding.Value
        freePredicate = false
    } else {
        queryPred = n.pattern.Predicate
    }
    if freeObject && object.Variable == binding.Variable {
        queryObj = binding.Value
        freeObject = false
    } else {
        queryObj = n.pattern.Object
    }

    // retrieves triples & form bindings to send
    go func() {
        defer close(out)
        for triple := range n.graph.Filter(querySubj, queryPred, queryObj) {
            if freeSubject {
                out <- rdf.NewBinding(subject.Variable, triple.Subject)
            }
            if freePredicate {
                out <- rdf.NewBinding(predicate.Variable, triple.Predicate)
            }
            if freeObject {
                out <- rdf.NewBinding(object.Variable, triple.Object)
            }
        }
    }()
    return out
}
