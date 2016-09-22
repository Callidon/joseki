// Copyright (c) 2016 Thomas Minier. All rights reserved.
// Use of this source code is governed by a MIT License
// license that can be found in the LICENSE file.

package graph

import (
	"github.com/Callidon/joseki/rdf"
	"sync"
)

// ListGraph is implementation of a RDF Graph, using a slice to store RDF Triples.
type ListGraph struct {
	dictionnary *bimap
	triples     []bitmapTriple
	nextID      int
	*sync.RWMutex
	*rdfReader
}

// NewListGraph creates a new List Graph.
func NewListGraph() *ListGraph {
	reader := newRDFReader()
	g := &ListGraph{newBimap(), make([]bitmapTriple, 0), 0, &sync.RWMutex{}, reader}
	reader.graph = g
	return g
}

// Register a new Node in the graph dictionnary, then return its unique ID.
func (g *ListGraph) registerNode(node rdf.Node) int {
	// insert the node in dictionnary only if it's not in
	key, inDict := g.dictionnary.locate(node)
	if !inDict {
		g.dictionnary.push(g.nextID, node)
		g.nextID++
		return g.nextID - 1
	}
	return key
}

// identifyNode returns the id of a node in the graph dictionnary, and a boolean to
// indicate if the node is known by the dictionnary
func (g *ListGraph) identifyNode(node rdf.Node) (int, bool) {
	if _, isVar := node.(rdf.Variable); isVar {
		return -1, true
	}
	return g.dictionnary.locate(node)
}

// Add a new Triple pattern to the graph.
func (g *ListGraph) Add(triple rdf.Triple) {
	g.Lock()
	defer g.Unlock()
	// add each node of the triple to the dictionnary & then update the slice
	subjID, predID, objID := g.registerNode(triple.Subject), g.registerNode(triple.Predicate), g.registerNode(triple.Object)
	g.triples = append(g.triples, newBitmapTriple(subjID, predID, objID))
}

// Delete triples from the graph that match a BGP given in parameters.
func (g *ListGraph) Delete(subject, predicate, object rdf.Node) {
	g.Lock()
	defer g.Unlock()
	newTriples := make([]bitmapTriple, 0, len(g.triples))
	subjID, subjKnown := g.identifyNode(subject)
	predID, predKnown := g.identifyNode(predicate)
	objID, objKnown := g.identifyNode(object)
	// continue onyl if we know all elements of the pattern
	if subjKnown && predKnown && objKnown {
		refTriple := newBitmapTriple(subjID, predID, objID)
		// resinsert into the graph the elements we doesn't want to delete
		for _, triple := range g.triples {
			if test := triple.Equals(refTriple); !test {
				newTriples = append(newTriples, triple)
			}
		}
		// update the slice
		g.triples = make([]bitmapTriple, len(newTriples))
		copy(g.triples, newTriples)
	}
}

// FilterSubset fetch triples form the graph that match a BGP given in parameters.
// It impose a Limit(the max number of results to be send in the output channel)
// and an Offset (the number of results to skip before sending them in the output channel) to the nodes requested.
func (g *ListGraph) FilterSubset(subject rdf.Node, predicate rdf.Node, object rdf.Node, limit int, offset int) <-chan rdf.Triple {
	results := make(chan rdf.Triple, bufferSize)
	// search for matching triple pattern in graph
	go func() {
		g.RLock()
		defer g.RUnlock()
		defer close(results)
		subjID, subjKnown := g.identifyNode(subject)
		predID, predKnown := g.identifyNode(predicate)
		objID, objKnown := g.identifyNode(object)
		// continue onyl if we know all elements of the pattern
		if subjKnown && predKnown && objKnown {
			refTriple := newBitmapTriple(subjID, predID, objID)
			cpt := 0
			for _, triple := range g.triples {
				if test := refTriple.Equals(triple); test {
					// send the result only if the offset has been reached
					if cpt >= offset {
						value, err := triple.Triple(g.dictionnary)
						check(err)
						results <- value
					}
					cpt++
				}
				// terminate the loop when the limit has been reached
				if (limit != -1) && (cpt-offset > limit) {
					break
				}
			}
		}
	}()
	return results
}

// Filter fetch triples form the graph that match a BGP given in parameters.
func (g *ListGraph) Filter(subject, predicate, object rdf.Node) <-chan rdf.Triple {
	return g.FilterSubset(subject, predicate, object, -1, 0)
}
