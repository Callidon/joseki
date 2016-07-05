// Copyright (c) 2016 Thomas Minier. All rights reserved.
// Use of this source code is governed by a MIT License
// license that can be found in the LICENSE file.

package sparql

import (
	"github.com/Callidon/joseki/rdf"
	"sort"
	"sync"
)

// unionNode represent a Union Operator in a SPARQL query execution plan.
type unionNode struct {
	leftNode, rightNode sparqlNode
}

// newUnionNode creates a new Union Node.
func newUnionNode(left, right sparqlNode) *unionNode {
	return &unionNode{left, right}
}

// execute perform the Union between the two nodes of the Union Operator.
func (n unionNode) execute() <-chan rdf.BindingsGroup {
	var wg sync.WaitGroup
	out := make(chan rdf.BindingsGroup, bufferSize)

	fetchBindings := func(node sparqlNode, out chan rdf.BindingsGroup, wg *sync.WaitGroup) {
		defer wg.Done()
		for bindings := range node.execute() {
			out <- bindings
		}
	}

	// fetch the bindings from the left & the right nodes in parallel
	wg.Add(2)
	go fetchBindings(n.leftNode, out, &wg)
	go fetchBindings(n.rightNode, out, &wg)
	// wait for the completion of the previous operations before closing the channel
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

// This operation has no particular meaning in the case of a unionNode, so it's equivalent to the execute method.
func (n unionNode) executeWith(binding rdf.BindingsGroup) <-chan rdf.BindingsGroup {
	return n.execute()
}

// bindingNames returns the names of the bindings produced by this operation
func (n unionNode) bindingNames() (bindingNames []string) {
	bindingNames = n.leftNode.bindingNames()
	for _, name := range n.rightNode.bindingNames() {
		if !containsString(bindingNames, name) {
			bindingNames = append(bindingNames, name)
		}
	}
	sort.Strings(bindingNames)
	return
}

// Equals test if two Union nodes are equals.
func (n unionNode) Equals(other sparqlNode) bool {
	union, isUnion := other.(*unionNode)
	if !isUnion {
		return false
	}
	return n.leftNode.Equals(union.leftNode) && n.rightNode.Equals(union.rightNode)
}

// String serialize the node in string format.
func (n unionNode) String() string {
	return "Union (" + n.leftNode.String() + ", " + n.rightNode.String() + ")"
}
