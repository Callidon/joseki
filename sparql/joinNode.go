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
	innerNode sparqlNode
	outerNode sparqlNode
}

// newJoinNode creates a new Join Node.
func newJoinNode(inner, outer sparqlNode) *joinNode {
	return &joinNode{inner, outer}
}

// execute perform the join between the two nodes of the Join Operator.
// It use a parallelized Nested Loop Join algorithm, as described by
// ÖZSU, M. Tamer et VALDURIEZ, Patrick. In Principles of distributed database systems. Springer Science & Business Media, 2011
func (n joinNode) execute() chan rdf.BindingsGroup {
	var wg sync.WaitGroup
	var page []rdf.BindingsGroup
	out := make(chan rdf.BindingsGroup, bufferSize)

	// execute the outer loop with a page of bindings from the inner loop
	executeOuterLoop := func(outerNode sparqlNode, page []rdf.BindingsGroup, out chan rdf.BindingsGroup, wg *sync.WaitGroup) {
		for _, bindingGroup := range page {
			for outerBindings := range outerNode.executeWith(bindingGroup) {
				out <- outerBindings
			}
		}
		wg.Done()
	}

	go func() {
		defer close(out)
		cpt := 0
		// execute the inner loop, then the outer loop using each group of bindings previously retrieved
		for innerBindings := range n.innerNode.execute() {
			// accumule group of bindings to form pages, send them when they are completed
			if cpt < pageSize {
				page = append(page, innerBindings)
			} else {
				// execute the outer loop for the current page, then prepare the next one
				wg.Add(1)
				go executeOuterLoop(n.outerNode, page, out, &wg)
				cpt = 0
				page = nil
			}
		}
		// process the last page if it's not empty
		if len(page) > 0 {
			wg.Add(1)
			go executeOuterLoop(n.outerNode, page, out, &wg)
		}
		// wait for all process to finish their jobs
		wg.Wait()
	}()
	return out
}

// This operation has no particular meaning in the case of a joinNode, so it's equivalent to the execute method.
func (n joinNode) executeWith(binding rdf.BindingsGroup) chan rdf.BindingsGroup {
	return n.execute()
}

// bindingNames returns the names of the bindings produced by this operation.
func (n joinNode) bindingNames() (bindingNames []string) {
	bindingNames = n.innerNode.bindingNames()
	for _, name := range n.outerNode.bindingNames() {
		if !containsString(bindingNames, name) {
			bindingNames = append(bindingNames, name)
		}
	}
	sort.Strings(bindingNames)
	return
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
	return "JOIN (" + n.innerNode.String() + ", " + n.outerNode.String() + ")"
}
