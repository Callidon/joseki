// Copyright (c) 2016 Thomas Minier. All rights reserved.
// Use of this source code is governed by a MIT License
// license that can be found in the LICENSE file.

package sparql

import (
	"github.com/Callidon/joseki/graph"
	"github.com/Callidon/joseki/rdf"
)

const (
	// Max size for the buffer of this package
	bufferSize = 100
	// Max size for the pages of group of bindings
	pageSize = 15
)

// queryType is the type of a SPARQL query
type queryType float64

const (
	_ = iota
	// selectQuery is a SPARQL SELECT query
	selectQuery queryType = 1 << (10 * iota)
	// askQuery is a SPARQL ASK query
	askQuery
	// describeQuery is a SPARQL DESCRIBE query
	describeQuery
	// constructQuery is a SPARQL CONSTRUCT query
	constructQuery
)

// queryDescriptor is a virtual description of a SPARQL query.
// It holds informations about the query during its construction process
// and can be converted into a real query.
type queryDescriptor struct {
	graph graph.Graph
	qType queryType
	bgps  [][]sparqlNode
}

// newQueryDescriptor creates a new queryDescriptor.
func newQueryDescriptor(graph graph.Graph, qType queryType) *queryDescriptor {
	return &queryDescriptor{graph, qType, make([][]sparqlNode, 0)}
}

// From set the source's graph for the query.
// Multiple calls of this method will override the previous source each time.
func (q *queryDescriptor) From(graph graph.Graph) {
	q.graph = graph
}

// Where add multiples triples pattern as a new BGP evualuted by the query
// Multiple calls of this method will each time add a new BGP to the query.
func (q *queryDescriptor) Where(triples ...rdf.Triple) {
	var nodes []sparqlNode
	for _, triple := range triples {
		nodes = append(nodes, newTripleNode(triple, q.graph))
	}
	q.bgps = append(q.bgps, nodes)
}

// joinNode try to find the first join possible between a node and a list of other nodes/
// It returns the joinNode & the index of the right
func findJoin(leftNode sparqlNode, otherNodes ...sparqlNode) (sparqlNode, int) {
	leftVariables := leftNode.bindingNames()
	leftSelectivity := len(leftVariables)
	// search for possible joins between the current node & the others
	for rightInd, rightNode := range otherNodes {
		rightVariables := rightNode.bindingNames()
		rightSelectivity := len(rightVariables)
		for _, variable := range rightVariables {
			// a join is possible between the two nodes
			if containsString(leftVariables, variable) {
				// pre-optimization : put the most selective node at the left of the join
				// apply only it if the two nodes are Triple Nodes
				_, isLeftTriple := leftNode.(*tripleNode)
				_, isRightTriple := rightNode.(*tripleNode)
				if isLeftTriple && isRightTriple && (leftSelectivity >= rightSelectivity) {
					return newJoinNode(rightNode, leftNode), rightInd
				}
				return newJoinNode(leftNode, rightNode), rightInd
			}
		}
	}
	// return nil if we doesn't find any joins
	return nil, -1
}

// build analyse the query execution plan and return its first node
func (q *queryDescriptor) build() sparqlNode {
	var root, joinNode sparqlNode
	var currentBGP, bgpRoots []sparqlNode
	var rightInd int
	var processNodes []int
	var joinFound bool

	// find the possible joins for each BGP
	for _, bgp := range q.bgps {
		nextBGP := make([]sparqlNode, len(bgp))
		copy(nextBGP, bgp)
		// look for joins until all nodes have been processed
		for len(currentBGP) != len(nextBGP) {
			currentBGP = make([]sparqlNode, len(nextBGP))
			copy(currentBGP, nextBGP)
			nextBGP = nil
			processNodes = nil
			joinFound = false

			for leftInd, leftNode := range currentBGP {
				// once we found one join, fill the set in order to skip to the next step
				if joinFound {
					if !containsInt(processNodes, leftInd) {
						nextBGP = append(nextBGP, leftNode)
					}
					continue
				}
				// search for a possible join between the current nodes & the next nodes
				joinNode, rightInd = findJoin(leftNode, currentBGP[leftInd+1:]...)
				if (joinNode == nil) && (rightInd == -1) {
					// save the node since no join is currently possible with it
					nextBGP = append(nextBGP, leftNode)
				} else {
					// save the new join we were able to find
					nextBGP = append(nextBGP, joinNode)
					joinFound = true
					processNodes = append(processNodes, rightInd+leftInd+1)
				}
			}
		}
		// join the remaining nodes (joins & triples) between them using UNION operators
		root = nextBGP[0]
		for _, bgpNode := range nextBGP[1:] {
			root = newUnionNode(root, bgpNode)
		}
		bgpRoots = append(bgpRoots, root)
	}
	// assemble the BGPs together with UNION operators
	root = bgpRoots[0]
	for _, bgpNode := range bgpRoots[1:] {
		root = newUnionNode(root, bgpNode)
	}
	return root
}
