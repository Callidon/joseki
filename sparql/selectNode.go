// Copyright (c) 2016 Thomas Minier. All rights reserved.
// Use of this source code is governed by a MIT License
// license that can be found in the LICENSE file.

package sparql

import "github.com/Callidon/joseki/rdf"

// selectNode represent a Select operation in a SPARQL query execution plan
type selectNode struct {
	node  sparqlNode
	names []string
}

// newSelectNode creates a new Select Node
func newSelectNode(node sparqlNode, bindings ...string) *selectNode {
	return &selectNode{node, bindings}
}

// execute apply a Select operation to the bindings produced by another node
func (n *selectNode) execute() chan rdf.BindingsGroup {
	out := make(chan rdf.BindingsGroup, bufferSize)

	go func() {
		defer close(out)
		// request groups of bindings from the node below & filter them
		for group := range n.node.execute() {
			newGroup := rdf.NewBindingsGroup()
			for _, bindingName := range n.names {
				value, inGroup := group.Bindings[bindingName]
				if inGroup {
					newGroup.Bindings[bindingName] = value
				}
			}
			out <- newGroup
		}
	}()
	return out
}

// This operation has no particular meaning in the case of a selectNode, so it's equivalent to the execute method
func (n *selectNode) executeWith(binding rdf.BindingsGroup) chan rdf.BindingsGroup {
	return n.execute()
}

// bindingNames returns the names of the bindings produced by this operation
func (n *selectNode) bindingNames() []string {
	return n.node.bindingNames()
}
