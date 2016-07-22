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
	bNames              []string
	rebuildNames        bool
}

// newUnionNode creates a new Union Node.
func newUnionNode(left, right sparqlNode) *unionNode {
	return &unionNode{left, right, nil, true}
}

// performUnion perform the Union between two streams of groups of bindings.
func (n unionNode) performUnion(leftIn, rightIn <-chan rdf.BindingsGroup, out chan<- rdf.BindingsGroup) {
	defer close(out)
	var wg sync.WaitGroup

	fetchBindings := func(input <-chan rdf.BindingsGroup, out chan<- rdf.BindingsGroup, wg *sync.WaitGroup) {
		defer wg.Done()
		for bindings := range input {
			out <- bindings
		}
	}

	// fetch the bindings from the left & the right nodes in parallel
	wg.Add(2)
	go fetchBindings(leftIn, out, &wg)
	go fetchBindings(rightIn, out, &wg)
	// wait for the completion of the previous operations before closing the channel
	wg.Wait()
}

// execute perform the Union between the two nodes of the Union Operator
// with the execute() operation apply to each of them.
//
// SPARQl 1.1 UNION reference : https://www.w3.org/TR/2013/REC-sparql11-query-20130321/#alternatives
func (n unionNode) execute() <-chan rdf.BindingsGroup {
	out := make(chan rdf.BindingsGroup, bufferSize)
	go n.performUnion(n.leftNode.execute(), n.rightNode.execute(), out)
	return out
}

//  execute perform the Union between the two nodes of the Union Operator
// with the executeWith() operation apply to each of them.
//
// SPARQl 1.1 UNION reference : https://www.w3.org/TR/2013/REC-sparql11-query-20130321/#alternatives
func (n unionNode) executeWith(bindings rdf.BindingsGroup) <-chan rdf.BindingsGroup {
	out := make(chan rdf.BindingsGroup, bufferSize)
	go n.performUnion(n.leftNode.executeWith(bindings), n.rightNode.executeWith(bindings), out)
	return out
}

// bindingNames returns the names of the bindings produced by this operation.
// Those names are stored in cache after the first call of this method, in
// order to speed up later calls of the method.
func (n unionNode) bindingNames() []string {
	if n.rebuildNames {
		n.bNames = n.leftNode.bindingNames()
		for _, name := range n.rightNode.bindingNames() {
			if !containsString(n.bNames, name) {
				n.bNames = append(n.bNames, name)
			}
		}
		sort.Strings(n.bNames)
		n.rebuildNames = false
	}
	return n.bNames
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
	return "UNION (" + n.leftNode.String() + ", " + n.rightNode.String() + ")"
}
