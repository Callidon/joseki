// Copyright (c) 2016 Thomas Minier. All rights reserved.
// Use of this source code is governed by a MIT License
// license that can be found in the LICENSE file.

package sparql

import (
	"github.com/Callidon/joseki/graph"
	"github.com/Callidon/joseki/rdf"
)

// tripleNode is the lowest level of SPARQL query execution plan.
// Its role is to retrieve bindings according to a triple pattern from a graph.
type tripleNode struct {
	pattern rdf.Triple
	graph   graph.Graph
}

// newTripleNode creates a new tripleNode.
func newTripleNode(pattern rdf.Triple, graph graph.Graph) *tripleNode {
	return &tripleNode{pattern, graph}
}

// execute retrieves bindings from a graph that match a triple pattern.
func (n *tripleNode) execute() chan rdf.BindingsGroup {
	out := make(chan rdf.BindingsGroup, bufferSize)
	// find free vars in triple pattern
	subject, freeSubject := n.pattern.Subject.(rdf.BlankNode)
	predicate, freePredicate := n.pattern.Predicate.(rdf.BlankNode)
	object, freeObject := n.pattern.Object.(rdf.BlankNode)

	// retrieves triples & form bindings to send
	go func() {
		defer close(out)
		for triple := range n.graph.Filter(n.pattern.Subject, n.pattern.Predicate, n.pattern.Object) {
			group := rdf.NewBindingsGroup()
			if freeSubject {
				group.Bindings[subject.Variable] = triple.Subject
			}
			if freePredicate {
				group.Bindings[predicate.Variable] = triple.Predicate
			}
			if freeObject {
				group.Bindings[object.Variable] = triple.Object
			}
			out <- group
		}
	}()
	return out
}

// executeWith retrieves bindings from a graph that match a triple pattern, completed by a given binding.
func (n *tripleNode) executeWith(group rdf.BindingsGroup) chan rdf.BindingsGroup {
	var querySubj, queryPred, queryObj rdf.Node
	out := make(chan rdf.BindingsGroup, bufferSize)
	// find free vars in triple pattern
	subject, freeSubject := n.pattern.Subject.(rdf.BlankNode)
	predicate, freePredicate := n.pattern.Predicate.(rdf.BlankNode)
	object, freeObject := n.pattern.Object.(rdf.BlankNode)

	// complete triple pattern using the group of bindings given in parameter
	for bindingKey, bindingValue := range group.Bindings {
		if freeSubject && subject.Variable == bindingKey {
			querySubj = bindingValue
			freeSubject = false
		} else {
			querySubj = n.pattern.Subject
		}
		if freePredicate && predicate.Variable == bindingKey {
			queryPred = bindingValue
			freePredicate = false
		} else {
			queryPred = n.pattern.Predicate
		}
		if freeObject && object.Variable == bindingKey {
			queryObj = bindingValue
			freeObject = false
		} else {
			queryObj = n.pattern.Object
		}
	}

	// retrieves triples & form bindings to send
	go func() {
		defer close(out)
		for triple := range n.graph.Filter(querySubj, queryPred, queryObj) {
			newGroup := group.Clone()
			if freeSubject {
				newGroup.Bindings[subject.Variable] = triple.Subject
			}
			if freePredicate {
				newGroup.Bindings[predicate.Variable] = triple.Predicate
			}
			if freeObject {
				newGroup.Bindings[object.Variable] = triple.Object
			}
			out <- newGroup
		}
	}()
	return out
}

// bindingNames returns the names of the bindings produced
func (n *tripleNode) bindingNames() []string {
	var bindings []string
	// find free vars in triple pattern
	subject, freeSubject := n.pattern.Subject.(rdf.BlankNode)
	predicate, freePredicate := n.pattern.Predicate.(rdf.BlankNode)
	object, freeObject := n.pattern.Object.(rdf.BlankNode)
	if freeSubject {
		bindings = append(bindings, subject.Variable)
	}
	if freePredicate {
		bindings = append(bindings, predicate.Variable)
	}
	if freeObject {
		bindings = append(bindings, object.Variable)
	}
	return bindings
}
