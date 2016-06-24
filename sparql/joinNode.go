// Copyright (c) 2016 Thomas Minier. All rights reserved.
// Use of this source code is governed by a MIT License
// license that can be found in the LICENSE file.

package sparql

import "github.com/Callidon/joseki/rdf"

// joinNode represent a Join Operator in a SPARQL query execution plan
type joinNode struct {
	innerNode sparqlNode
	outerNode sparqlNode
}

// newJoinNode creates a new Join Node
func newJoinNode(inner, outer sparqlNode) *joinNode {
	return &joinNode{inner, outer}
}

// execute perform the join between the two nodes of the Join Operator
func (n *joinNode) execute() chan rdf.BindingsGroup {
	out := make(chan rdf.BindingsGroup, bufferSize)

	go func() {
		defer close(out)
		// execute the inner join, then the outer join using each group of bindings retrieved from the inner join
		for innerBindings := range n.innerNode.execute() {
			// TODO : use parallelization for the outer join ?
			for outerBindings := range n.outerNode.executeWith(innerBindings) {
				out <- outerBindings
			}
		}
	}()
	return out
}

// This operation has no particular meaning in the case of a joinNode, so it's equivalent to the execute method
func (n *joinNode) executeWith(binding rdf.BindingsGroup) chan rdf.BindingsGroup {
	return n.execute()
}
