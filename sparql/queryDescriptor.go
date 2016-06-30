// Copyright (c) 2016 Thomas Minier. All rights reserved.
// Use of this source code is governed by a MIT License
// license that can be found in the LICENSE file.

package sparql

import (
	"github.com/Callidon/joseki/graph"
	"github.com/Callidon/joseki/rdf"
	"sort"
)

const (
	// Max size for the buffer of this package
	bufferSize = 100
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
			// we find an possible join
			if sort.SearchStrings(leftVariables, variable) != len(leftVariables) {
				// pre-optimization : put the most selective node at the left of the join
				// apply only it if the two nodes are Triple Nodes
				// based on "The SPARQL Query Graph Model for Query Optimization" (Olaf Hartig, Ralf Hees)
				_, isLeftTriple := leftNode.(*tripleNode)
				_, isRightTriple := rightNode.(*tripleNode)
				if isLeftTriple && isRightTriple {
					if leftSelectivity < rightSelectivity {
						return newJoinNode(leftNode, rightNode), rightInd
					}
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
	//var bgpRoots []sparqlNode
	var rightInd int
	// finds the possible join for each BGP
	for _, bgp := range q.bgps {
		processNodes := 0
		// look for joins until all nodes have been processed
		for processNodes < len(bgp) {
			for leftInd, leftNode := range bgp {
				if leftInd == 0 {
					joinNode, rightInd = findJoin(leftNode, bgp[1:]...)
				} else {
					joinNode, rightInd = findJoin(leftNode, bgp[0:leftInd]...)
					// look in the other part of the slice if first search fails
					if (joinNode == nil) && (rightInd == -1) {
						joinNode, rightInd = findJoin(leftNode, bgp[leftInd+1:]...)
					}
				}
				// if a join has been found
				if (joinNode != nil) && (rightInd >= 0) {
					// TODO
					// update counter
					if processNodes == 0 {
						processNodes = 2
					} else {
						processNodes++
					}
				} else {
					// TODO
				}
			}
		}
		// join the remaining nodes (joins & triples) between them using UNION operators
		// ...
	}
	// assemble the BGPs together with UNION operators
	/*root = bgpRoots[0]
	  for _, bgpNode := range bgpRoots[1:] {
	    root = newUnionNode(root, bgpNode)
	  }*/
	return root
}
