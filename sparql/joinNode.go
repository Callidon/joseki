// Copyright (c) 2016 Thomas Minier. All rights reserved.
// Use of this source code is governed by a MIT License
// license that can be found in the LICENSE file.

package sparql

import (
	"github.com/Callidon/joseki/rdf"
	"sort"
	"sync"
)

// joinNode represent a Join Operator in a SPARQL query execution plan.
type joinNode struct {
	outerNode, innerNode sparqlNode
}

// newJoinNode creates a new Join Node.
func newJoinNode(outer, inner sparqlNode) *joinNode {
	return &joinNode{outer, inner}
}

// performJoin apply a join between a stream of groups of bindings from a node and another node.
// It use a parallelized Nested Loop Join algorithm, as described by
// ÖZSU, M. Tamer et VALDURIEZ, Patrick. In Principles of distributed database systems. Springer Science & Business Media, 2011
func (n joinNode) performJoin(in <-chan rdf.BindingsGroup, out chan<- rdf.BindingsGroup) {
	defer close(out)
	var wg sync.WaitGroup
	page := make([]rdf.BindingsGroup, 0, pageSize)
	cpt := 0

	// execute the inner loop with a page of bindings from the inner loop
	executeInnerLoop := func(innerNode sparqlNode, page []rdf.BindingsGroup, out chan<- rdf.BindingsGroup, wg *sync.WaitGroup) {
		for _, bindingGroup := range page {
			for innerBindings := range innerNode.executeWith(bindingGroup) {
				out <- innerBindings
			}
		}
		wg.Done()
	}

	// execute the outer loop, then the inner loop using each group of bindings previously retrieved
	for outerBindings := range in {
		// accumule group of bindings to form pages, send them when they are completed
		if cpt < pageSize {
			page = append(page, outerBindings)
		} else {
			// execute the inner loop for the current page, then prepare the next one
			wg.Add(1)
			go executeInnerLoop(n.innerNode, page, out, &wg)
			cpt = 0
			page = make([]rdf.BindingsGroup, 0, pageSize)
		}
	}
	// process the last page if it's not empty
	if len(page) > 0 {
		wg.Add(1)
		go executeInnerLoop(n.innerNode, page, out, &wg)
	}
	// wait for all process to finish their jobs
	wg.Wait()
}

// execute perform the join between the two nodes of the Join Operator.
func (n joinNode) execute() <-chan rdf.BindingsGroup {
	out := make(chan rdf.BindingsGroup, bufferSize)
	go n.performJoin(n.outerNode.execute(), out)
	return out
}

// execute perform the join between the two nodes of the Join Operator using a group of bindings
func (n joinNode) executeWith(bindings rdf.BindingsGroup) <-chan rdf.BindingsGroup {
	out := make(chan rdf.BindingsGroup, bufferSize)
	go n.performJoin(n.outerNode.executeWith(bindings), out)
	return out
}

// bindingNames returns the names of the bindings produced by this operation.
func (n joinNode) bindingNames() []string {
	bindingNames := n.outerNode.bindingNames()
	for _, name := range n.innerNode.bindingNames() {
		if !containsString(bindingNames, name) {
			bindingNames = append(bindingNames, name)
		}
	}
	sort.Strings(bindingNames)
	return bindingNames
}

// Equals test if two Join nodes are equals.
func (n joinNode) Equals(other sparqlNode) bool {
	join, isJoin := other.(*joinNode)
	if !isJoin {
		return false
	}
	return n.innerNode.Equals(join.innerNode) && n.outerNode.Equals(join.outerNode)
}

// String serialize the node in string format.
func (n joinNode) String() string {
	return "JOIN (" + n.outerNode.String() + ", " + n.innerNode.String() + ")"
}
