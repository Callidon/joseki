// Copyright (c) 2016 Thomas Minier. All rights reserved.
// Use of this source code is governed by a MIT License
// license that can be found in the LICENSE file.

package sparql

import (
	"github.com/Callidon/joseki/rdf"
	"strings"
)

// selectNode represent a Select operation in a SPARQL query execution plan.
type selectNode struct {
	node         sparqlNode
	varNames     []string
	bNames       []string
	rebuildNames bool
}

// newSelectNode creates a new Select Node.
func newSelectNode(node sparqlNode, bindings ...string) *selectNode {
	return &selectNode{node, bindings, nil, true}
}

// applySelect apply a SELECT operation on a stream of groups of binding
func (n selectNode) applySelect(in <-chan rdf.BindingsGroup, out chan<- rdf.BindingsGroup) {
	defer close(out)
	// request groups of bindings from the node below & filter them
	for group := range in {
		newGroup := rdf.NewBindingsGroup()
		for _, bindingName := range n.varNames {
			value, inGroup := group.Bindings[bindingName]
			if inGroup {
				newGroup.Bindings[bindingName] = value
			}
		}
		out <- newGroup
	}
}

// execute apply a SELECT operation to the bindings produced by
// another node when applying the execute() operation to it.
//
// SPARQL 1.1 SELECT reference : https://www.w3.org/TR/2013/REC-sparql11-query-20130321/#select
func (n selectNode) execute() <-chan rdf.BindingsGroup {
	out := make(chan rdf.BindingsGroup, bufferSize)
	go n.applySelect(n.node.execute(), out)
	return out
}

// executeWith apply a SELECT operation to the bindings produced by
// another node when applying the executeWith() operation to it.
//
// SPARQL 1.1 SELECT reference : https://www.w3.org/TR/2013/REC-sparql11-query-20130321/#select
func (n selectNode) executeWith(bindings rdf.BindingsGroup) <-chan rdf.BindingsGroup {
	out := make(chan rdf.BindingsGroup, bufferSize)
	go n.applySelect(n.node.executeWith(bindings), out)
	return out
}

// bindingNames returns the names of the bindings produced by this operation.
// Those names are stored in cache after the first call of this method, in
// order to speed up later calls of the method.
func (n selectNode) bindingNames() []string {
	if n.rebuildNames {
		n.bNames = n.node.bindingNames()
		n.rebuildNames = false
	}
	return n.bNames
}

// Equals test if two Select nodes are equals.
func (n selectNode) Equals(other sparqlNode) bool {
	selectN, isSelect := other.(*selectNode)
	if !isSelect {
		return false
	}
	return n.node.Equals(selectN.node)
}

// String serialize the node in string format.
func (n selectNode) String() string {
	return "SELECT " + strings.Join(n.varNames, ",") + " (" + n.node.String() + ")"
}
