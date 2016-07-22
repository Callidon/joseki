// Copyright (c) 2016 Thomas Minier. All rights reserved.
// Use of this source code is governed by a MIT License
// license that can be found in the LICENSE file.

package sparql

import (
	"github.com/Callidon/joseki/graph"
	"github.com/Callidon/joseki/rdf"
	"sort"
)

// tripleNode is the lowest level of SPARQL query execution plan.
// Its role is to retrieve bindings according to a triple pattern from a graph.
type tripleNode struct {
	pattern       rdf.Triple
	graph         graph.Graph
	limit, offset int
	bNames        []string
	rebuildNames  bool
}

// newTripleNode creates a new tripleNode.
func newTripleNode(pattern rdf.Triple, graph graph.Graph, limit int, offset int) *tripleNode {
	return &tripleNode{pattern, graph, limit, offset, nil, true}
}

// execute retrieves bindings from a graph that match a triple pattern.
func (n tripleNode) execute() <-chan rdf.BindingsGroup {
	out := make(chan rdf.BindingsGroup, bufferSize)

	// find free vars in triple pattern
	subject, freeSubject := n.pattern.Subject.(rdf.Variable)
	predicate, freePredicate := n.pattern.Predicate.(rdf.Variable)
	object, freeObject := n.pattern.Object.(rdf.Variable)

	// retrieves triples & form bindings to send
	go func() {
		defer close(out)
		for triple := range n.graph.FilterSubset(n.pattern.Subject, n.pattern.Predicate, n.pattern.Object, n.limit, n.offset) {
			group := rdf.NewBindingsGroup()
			if freeSubject {
				group.Bindings[subject.Value] = triple.Subject
			}
			if freePredicate {
				group.Bindings[predicate.Value] = triple.Predicate
			}
			if freeObject {
				group.Bindings[object.Value] = triple.Object
			}
			out <- group
		}
	}()
	return out
}

// executeWith retrieves bindings from a graph that match a triple pattern, completed by a given binding.
func (n tripleNode) executeWith(group rdf.BindingsGroup) <-chan rdf.BindingsGroup {
	out := make(chan rdf.BindingsGroup, bufferSize)

	// complete the triple attern using a groupf of bindings, then find its free vars
	newTriple := n.pattern.Complete(group)
	subject, freeSubject := newTriple.Subject.(rdf.Variable)
	predicate, freePredicate := newTriple.Predicate.(rdf.Variable)
	object, freeObject := newTriple.Object.(rdf.Variable)

	// retrieves triples & form bindings to send
	go func() {
		defer close(out)
		for triple := range n.graph.FilterSubset(newTriple.Subject, newTriple.Predicate, newTriple.Object, n.limit, n.offset) {
			newGroup := group.Clone()
			if freeSubject {
				newGroup.Bindings[subject.Value] = triple.Subject
			}
			if freePredicate {
				newGroup.Bindings[predicate.Value] = triple.Predicate
			}
			if freeObject {
				newGroup.Bindings[object.Value] = triple.Object
			}
			out <- newGroup
		}
	}()
	return out
}

// bindingNames returns the names of the bindings produced.
// Those names are stored in cache after the first call of this method, in
// order to speed up later calls of the method.
func (n tripleNode) bindingNames() []string {
	if n.rebuildNames {
		n.bNames = make([]string, 0, 3)
		// find free vars in triple pattern
		subject, freeSubject := n.pattern.Subject.(rdf.Variable)
		predicate, freePredicate := n.pattern.Predicate.(rdf.Variable)
		object, freeObject := n.pattern.Object.(rdf.Variable)
		if freeSubject {
			n.bNames = append(n.bNames, subject.Value)
		}
		if freePredicate {
			n.bNames = append(n.bNames, predicate.Value)
		}
		if freeObject {
			n.bNames = append(n.bNames, object.Value)
		}
		sort.Strings(n.bNames)
		n.rebuildNames = false
	}
	return n.bNames
}

// Equals test if two Triple nodes are equals.
func (n tripleNode) Equals(other sparqlNode) bool {
	tripleN, isTriple := other.(*tripleNode)
	if !isTriple {
		return false
	}
	test, _ := n.pattern.Equals(tripleN.pattern)
	return test
}

// String serialize the node in string format.
func (n tripleNode) String() string {
	return "Triple(" + n.pattern.Subject.String() + " " + n.pattern.Predicate.String() + " " + n.pattern.Object.String() + ")"
}
