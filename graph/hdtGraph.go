// Copyright (c) 2016 Thomas Minier. All rights reserved.
// Use of this source code is governed by a MIT License
// license that can be found in the LICENSE file.

package graph

import (
	"github.com/Callidon/joseki/rdf"
	"sync"
)

// atomicCounter respresent a synchronized counter
type atomicCounter struct {
	cpt       int
	threshold int
	*sync.Mutex
}

// newAtomicCounter creates a new Atomic Counter
func newAtomicCounter(cpt, limit int) *atomicCounter {
	return &atomicCounter{cpt, limit, &sync.Mutex{}}
}

// HDTGraph is a implementation of a RDF Graph based on the HDT-MR model proposed by Giménez-García et al.
//
// For more details, see http://dataweb.infor.uva.es/projects/hdt-mr/
type HDTGraph struct {
	dictionnary *bimap
	root        *bitmapNode
	nextID      int
	triples     map[string][]rdf.Triple
	*sync.Mutex
	*rdfReader
}

// NewHDTGraph creates a new empty HDT Graph.
func NewHDTGraph() *HDTGraph {
	reader := newRDFReader()
	g := &HDTGraph{newBimap(), newBitmapNode(-1), 0, make(map[string][]rdf.Triple), &sync.Mutex{}, reader}
	reader.graph = g
	return g
}

// Register a new Node in the graph dictionnary, then return its unique ID.
func (g *HDTGraph) registerNode(node rdf.Node) int {
	// insert the node in dictionnary only if it's not in
	key, inDict := g.dictionnary.locate(node)
	if !inDict {
		g.dictionnary.push(g.nextID, node)
		g.nextID++
		return g.nextID - 1
	}
	return key
}

// Recursively remove nodes that match criteria
func (g *HDTGraph) removeNodes(root *bitmapNode, datas []*rdf.Node) {
	// it's a blank node, delete all his sons
	node := (*datas[0])
	if _, isVar := node.(rdf.Variable); isVar {
		root.removeSons()
	} else {
		// search for the specific node in the root's sons
		refNodeID, inDict := g.dictionnary.locate(node)
		if inDict {
			son, inSons := root.sons[refNodeID]
			if inSons {
				// delete his sons that match the next criteria
				g.removeNodes(son, datas[1:])
			}
		}
	}
}

// Recursively collect data from the graph in order to form triple pattern matching criterias.
// The graph can be query with a Limit (the max number of rsults to send in the output channel)
// and an Offset (the number of results to skip before sending them in the output channel).
// These two parameters can be set to -1 to be ignored.
func (g *HDTGraph) queryNodes(root *bitmapNode, datas []*rdf.Node, triple []int, out chan<- rdf.Triple, wg *sync.WaitGroup, limit *atomicCounter, offset *atomicCounter) {
	limit.Lock()
	defer wg.Done()
	defer limit.Unlock()
	// utilitary function to update WaitGroup when skipping sons
	updateSons := func(skipKey int, wg *sync.WaitGroup) {
		for key, son := range root.sons {
			if key != skipKey {
				son.updateCounter(wg)
			}
		}
	}

	// skip the node if the limit have a default value or has been reached
	if (limit.threshold != -1) && (limit.cpt >= limit.threshold) {
		updateSons(-1, wg)
	} else if len(triple) >= 3 {
		offset.Lock()
		defer offset.Unlock()
		// skip result and update offset if its threashold hasn't been reached
		if offset.cpt < offset.threshold {
			offset.cpt++
		} else {
			// when possible, create a new triple pattern & send it into the output pipeline
			bitmapTriple := newBitmapTriple(triple[0], triple[1], triple[2])
			triple, err := bitmapTriple.Triple(g.dictionnary)
			if err != nil {
				panic(err)
			}
			out <- triple
			limit.cpt++
		}
	} else {
		node := (*datas[0])
		switch node.(type) {
		case rdf.Variable:
			// search in every sons if the current node is a variable
			for _, son := range root.sons {
				go g.queryNodes(son, datas[1:], append(triple, son.id), out, wg, limit, offset)
			}
		default:
			// search for a specific node
			id, inDict := g.dictionnary.locate(node)
			son, inSons := root.sons[id]
			if inDict && inSons {
				go g.queryNodes(son, datas[1:], append(triple, id), out, wg, limit, offset)
				updateSons(id, wg)
			} else {
				updateSons(-1, wg)
			}
		}
	}
}

// Add a new Triple pattern to the graph.
func (g *HDTGraph) Add(triple rdf.Triple) {
	defer g.Unlock()
	// add each node of the triple to the dictionnary & then update the graph
	subjID, predID, objID := g.registerNode(triple.Subject), g.registerNode(triple.Predicate), g.registerNode(triple.Object)
	datas := []int{subjID, predID, objID}
	currentNode := g.root
	g.Lock()
	// insert each data in the graph
	for _, nodeID := range datas {
		node, inSons := currentNode.sons[nodeID]
		if inSons {
			// skip to next node if the current data is the same as the current node
			currentNode = node
		} else {
			// add the new node, then use it for the next data ton insert
			currentNode.sons[nodeID] = newBitmapNode(nodeID)
			currentNode = currentNode.sons[nodeID]
		}
	}
}

// Delete triples from the graph that match a BGP given in parameters.
func (g *HDTGraph) Delete(subject, predicate, object rdf.Node) {
	g.Lock()
	defer g.Unlock()
	g.removeNodes(g.root, []*rdf.Node{&subject, &predicate, &object})
}

// FilterSubset fetch triples form the graph that match a BGP given in parameters.
// It impose a Limit(the max number of results to be send in the output channel)
// and an Offset (the number of results to skip before sending them in the output channel) to the nodes requested.
func (g *HDTGraph) FilterSubset(subject rdf.Node, predicate rdf.Node, object rdf.Node, limit int, offset int) <-chan rdf.Triple {
	var wg sync.WaitGroup
	results := make(chan rdf.Triple, bufferSize)
	limitCpt, offsetCpt := newAtomicCounter(0, limit), newAtomicCounter(0, offset)
	// fetch data in the tree & wait for the operation to be complete before closing the pipeline
	g.Lock()
	wg.Add(g.root.length() + 1)
	go g.queryNodes(g.root, []*rdf.Node{&subject, &predicate, &object}, make([]int, 0, 3), results, &wg, limitCpt, offsetCpt)
	// use a daemon to wait for the end of all related goroutines before closing the channel
	go func() {
		defer close(results)
		defer g.Unlock()
		wg.Wait()
	}()
	return results
}

// Filter fetch triples form the graph that match a BGP given in parameters.
func (g *HDTGraph) Filter(subject, predicate, object rdf.Node) <-chan rdf.Triple {
	return g.FilterSubset(subject, predicate, object, -1, 0)
}
