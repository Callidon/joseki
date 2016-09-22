// Copyright (c) 2016 Thomas Minier. All rights reserved.
// Use of this source code is governed by a MIT License
// license that can be found in the LICENSE file.

package graph

import (
	"github.com/Callidon/joseki/rdf"
	"sync"
)

// TreeGraph is a implementation of a RDF Graph based on the HDT-MR model proposed by Giménez-García et al.
//
// For more details, see http://dataweb.infor.uva.es/projects/hdt-mr/
type TreeGraph struct {
	dictionnary *bimap
	root        *bitmapNode
	nextID      int
	triples     map[string][]rdf.Triple
	*sync.RWMutex
	*rdfReader
}

// NewTreeGraph creates a new empty Tree Graph.
func NewTreeGraph() *TreeGraph {
	reader := newRDFReader()
	g := &TreeGraph{newBimap(), newBitmapNode(-1), 0, make(map[string][]rdf.Triple), &sync.RWMutex{}, reader}
	reader.graph = g
	return g
}

// Register a new Node in the graph dictionnary, then return its unique ID.
func (g *TreeGraph) registerNode(node rdf.Node) int {
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
func (g *TreeGraph) removeNodes(root, previous *bitmapNode, datas []*rdf.Node) {
	if root != nil {
		node := (*datas[0])
		_, isVar := node.(rdf.Variable)
		id, inDict := g.dictionnary.locate(node)
		// delegate operation to root's sons if it's a Variable or if the root match the current citeria
		if isVar || (inDict && root.id == id) {
			for _, son := range root.sons {
				g.removeNodes(son, root, datas[1:])
			}
			// if root doesn't have any sons after the operation, delete it
			if len(root.sons) == 0 {
				delete(previous.sons, root.id)
			}
		}
	}
}

// sendTriples sends triples collected by another process, and respect the limit & offset of a query
func sendTriple(input <-chan bitmapTriple, out chan<- rdf.Triple, dict *bimap, limit, offset int) {
	defer close(out)
	cpt := 0
	for bTriple := range input {
		// send the triple if the offset threshold has been reached but not the limit threashold
		if cpt >= offset && (cpt-offset <= limit || limit == -1) {
			triple, err := bTriple.Triple(dict)
			if err != nil {
				panic(err)
			}
			out <- triple
		}
		cpt++
	}
}

// Recursively collect data from the graph in order to form triple pattern matching criterias.
// The graph can be query with a Limit (the max number of rsults to send in the output channel)
// and an Offset (the number of results to skip before sending them in the output channel).
// These two parameters can be set to -1 to be ignored.
func (g *TreeGraph) findNodes(root *bitmapNode, datas []*rdf.Node, triple []int, out chan<- bitmapTriple, wg *sync.WaitGroup) {
	defer wg.Done()
	// utilitary function to update WaitGroup counter when skipping sons
	skipSons := func(wg *sync.WaitGroup) {
		for _, son := range root.sons {
			son.updateCounter(wg)
		}
	}

	// skip the node if the limit have a default value or has been reached
	node := (*datas[0])
	_, isVar := node.(rdf.Variable)
	id, inDict := g.dictionnary.locate(node)
	// when the root is a variable or the value we need, save it & delegate the operation to its sons
	if isVar || (inDict && root.id == id) {
		if len(root.sons) == 0 {
			out <- newBitmapTriple(triple[0], triple[1], root.id)
			//sendValue(append(triple, root.id), out, g.dictionnary, limit, offset)
		} else {
			for _, son := range root.sons {
				go g.findNodes(son, datas[1:], append(triple, root.id), out, wg)
			}
		}
	} else {
		// the node doesn't match our query, so there's no need to visit its sons
		skipSons(wg)
	}
}

// Add a new Triple pattern to the graph.
func (g *TreeGraph) Add(triple rdf.Triple) {
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
func (g *TreeGraph) Delete(subject, predicate, object rdf.Node) {
	g.Lock()
	defer g.Unlock()
	for _, son := range g.root.sons {
		g.removeNodes(son, g.root, []*rdf.Node{&subject, &predicate, &object})
	}
}

// FilterSubset fetch triples form the graph that match a BGP given in parameters.
// It impose a Limit(the max number of results to be send in the output channel)
// and an Offset (the number of results to skip before sending them in the output channel) to the nodes requested.
func (g *TreeGraph) FilterSubset(subject rdf.Node, predicate rdf.Node, object rdf.Node, limit int, offset int) <-chan rdf.Triple {
	var wg sync.WaitGroup
	bitmapResults := make(chan bitmapTriple, bufferSize)
	results := make(chan rdf.Triple, bufferSize)

	// fetch data in the tree & wait for the operation to be complete before closing the pipeline
	g.RLock()
	go sendTriple(bitmapResults, results, g.dictionnary, limit, offset)
	for _, son := range g.root.sons {
		wg.Add(son.length() + 1)
		go g.findNodes(son, []*rdf.Node{&subject, &predicate, &object}, make([]int, 0, 3), bitmapResults, &wg)
	}
	// use a daemon to wait for the end of all related goroutines before closing the channel
	go func() {
		defer close(bitmapResults)
		defer g.RUnlock()
		wg.Wait()
	}()
	return results
}

// Filter fetch triples form the graph that match a BGP given in parameters.
func (g *TreeGraph) Filter(subject, predicate, object rdf.Node) <-chan rdf.Triple {
	return g.FilterSubset(subject, predicate, object, -1, 0)
}
