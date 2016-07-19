// Copyright (c) 2016 Thomas Minier. All rights reserved.
// Use of this source code is governed by a MIT License
// license that can be found in the LICENSE file.

package graph

import (
	"errors"
	"github.com/Callidon/joseki/rdf"
	"sync"
)

// Node represented in the Bitmap standard, following the HDT-MR model.
type bitmapNode struct {
	id   int
	sons map[int]*bitmapNode
}

// Triple represented in the Bitmap standard, following the HDT-MR model.
type bitmapTriple struct {
	subjectID   int
	predicateID int
	objectID    int
}

// newBitmapNode creates a new Bitmap Node without any son.
func newBitmapNode(id int) *bitmapNode {
	return &bitmapNode{id, make(map[int]*bitmapNode)}
}

// addSon add a son to a Bitmap Node.
func (n *bitmapNode) addSon(id int) {
	n.sons[id] = newBitmapNode(id)
}

// depth calculates the number of nodes in the tree, starting from this node.
func (n *bitmapNode) length() int {
	res := 0
	res += len(n.sons)
	for _, son := range n.sons {
		res += son.length()
	}
	return res
}

// updateCounter update a Wait Group counter for a node & his sons recursively.
func (n *bitmapNode) updateCounter(wg *sync.WaitGroup) {
	wg.Done()
	for _, son := range n.sons {
		son.updateCounter(wg)
	}
}

// Recursively remove the sons of a Bitmap Node
func (n *bitmapNode) removeSons() {
	for key, son := range n.sons {
		son.removeSons()
		delete(n.sons, key)
	}
}

// newBitmapTriple creates a New Bitmap Triple.
func newBitmapTriple(subj, pred, obj int) bitmapTriple {
	return bitmapTriple{subj, pred, obj}
}

// Equals returns True if two Bitmap Triple are equals, False otherwise
func (b *bitmapTriple) Equals(other bitmapTriple) bool {
	subjEq := b.subjectID == other.subjectID
	predEq := b.predicateID == other.predicateID
	objEq := b.objectID == other.objectID
	if b.subjectID < 0 || other.subjectID < 0 {
		subjEq = true
	}
	if b.predicateID < 0 || other.predicateID < 0 {
		predEq = true
	}
	if b.objectID < 0 || other.objectID < 0 {
		objEq = true
	}
	return subjEq && predEq && objEq
}

// Convert a BitMap Triple to a RDF Triple.
func (b *bitmapTriple) Triple(dict *bimap) (rdf.Triple, error) {
	var triple rdf.Triple
	subj, foundSubj := dict.extract(b.subjectID)
	if !foundSubj {
		return triple, errors.New("Error : cannot found the subject id in the dictionnary")
	}
	pred, foundPred := dict.extract(b.predicateID)
	if !foundPred {
		return triple, errors.New("Error : cannot found the predicate id in the dictionnary")
	}
	obj, foundObj := dict.extract(b.objectID)
	if !foundObj {
		return triple, errors.New("Error : cannot found the object id in the dictionnary")
	}
	triple = rdf.NewTriple(subj, pred, obj)
	return triple, nil
}
