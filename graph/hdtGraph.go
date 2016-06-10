package graph

import (
	"errors"
	"github.com/Callidon/joseki/core"
	"os"
	"sync"
)

const (
	// Max number of goroutines allowed in worker pools
	MAX_GOROUTINES_WP = 5.0
)

// Node represented in the Bitmap standard, following the HDT-MR model
type bitmapNode struct {
	id   int
	sons map[int]*bitmapNode
}

// Triple represented in the Bitmap standard, following the HDT-MR model
type bitmapTriple struct {
	subject_id   int
	predicate_id int
	object_id    int
}

// Implementation of a RDF Graph based on the HDT-MR model proposed by Giménez-García et al
// For more details, see http://dataweb.infor.uva.es/projects/hdt-mr/
type HDTGraph struct {
	dictionnary bimap
	root        bitmapNode
	nextId      int
	triples     map[string][]core.Triple
}

// Return a new Bitmap Node without any son
func newBitmapNode(id int) bitmapNode {
	return bitmapNode{id, make(map[int]*bitmapNode)}
}

// Add a son to a Bitmap Node
func (n *bitmapNode) addSon(id int) {
	n.sons[id] = &bitmapNode{id, make(map[int]*bitmapNode)}
}

// Return the depth of the tree starting to this node
func (n *bitmapNode) depth() int {
	res := 0
	if len(n.sons) > 0 {
		res += len(n.sons)
		for _, son := range n.sons {
			res += son.depth()
		}
	}
	return res
}

// Update a Wait Group counter for a node & his sons recursively
func (n *bitmapNode) updateCounter(wg *sync.WaitGroup) {
	wg.Done()
	for _, son := range n.sons {
		son.updateCounter(wg)
	}
}

// Return a New Bitmap Triple
func newBitmapTriple(subj, pred, obj int) bitmapTriple {
	return bitmapTriple{subj, pred, obj}
}

// Convert a BitMap Triple to a RDF Triple
func (t *bitmapTriple) Triple(dict *bimap) (core.Triple, error) {
	var triple core.Triple
	subj, foundSubj := dict.extract(t.subject_id)
	if !foundSubj {
		return triple, errors.New("Error : cannot found the subject id in the dictionnary")
	}
	pred, foundPred := dict.extract(t.predicate_id)
	if !foundPred {
		return triple, errors.New("Error : cannot found the predicate id in the dictionnary")
	}
	obj, foundObj := dict.extract(t.object_id)
	if !foundObj {
		return triple, errors.New("Error : cannot found the object id in the dictionnary")
	}
	triple = core.NewTriple(subj, pred, obj)
	return triple, nil
}

// Return a new empty HDT Graph
func NewHDTGraph() HDTGraph {
	return HDTGraph{newBimap(), newBitmapNode(-1), 0, make(map[string][]core.Triple)}
}

// Register a new Node in the graph dictionnary, then return its unique ID
func (g *HDTGraph) registerNode(node core.Node) int {
	// insert the node in dictionnary only if it's not in
	key, inDict := g.dictionnary.locate(node)
	if !inDict {
		g.dictionnary.push(g.nextId, node)
		g.nextId += 1
		return g.nextId - 1
	} else {
		return key
	}
}

// Recursively update the nodes of the graph with new datas
func (g *HDTGraph) updateNodes(root *bitmapNode, datas []int) {
	// if they are data to insert in the graph
	if len(datas) > 0 {
		id := datas[0]
		// if the node's id in already in the root sons, continue the operation with it
		node, inSons := root.sons[id]
		if inSons {
			g.updateNodes(node, datas[1:])
		} else {
			// add the new node, then continue the operation with its sons
			root.addSon(id)
			g.updateNodes(root.sons[id], datas[1:])
		}
	}
}

// Recursively collect datas from the graph in order to form triple pattern matching criterias
func (g *HDTGraph) queryNodes(root *bitmapNode, datas []*core.Node, triple []int, out chan core.Triple, wg *sync.WaitGroup) {
	defer wg.Done()
	// when possible, create a new triple pattern & send it into the output pipeline
	if len(triple) == 3 {
		bitmapTriple := newBitmapTriple(triple[0], triple[1], triple[2])
		triple, err := bitmapTriple.Triple(&g.dictionnary)
		if err != nil {
			panic(err)
		}
		out <- triple
	} else {
		node := (*datas[0])
		// if the current node to search is a blank node, search in every sons
		_, isBnode := node.(core.BlankNode)
		if isBnode {
			// IDEA : allow a pool of workers to query datas from all the sons of the root
			go func() {
				for _, son := range root.sons {
					g.queryNodes(son, datas[1:], append(triple, son.id), out, wg)
				}
			}()
		} else {
			// search for a specific node
			id, inDict := g.dictionnary.locate(node)
			if _, inSons := root.sons[id]; inDict && (inSons || root.sons[id] == nil) {
				go g.queryNodes(root.sons[id], datas[1:], append(triple, id), out, wg)
			}
			// update the counter for the sons that will not be visited
			for key, son := range root.sons {
				if key != id {
					son.updateCounter(wg)
				}
			}
		}
	}
}

func (g *HDTGraph) LoadFromFile(file *os.File) {
	//TODO
}

// Add a new Triple pattern to the graph
func (g *HDTGraph) Add(triple core.Triple) {
	// add each node of the triple to the dictionnary & then update the graph
	subjId, predID, objId := g.registerNode(triple.Subject), g.registerNode(triple.Predicate), g.registerNode(triple.Object)
	g.updateNodes(&g.root, []int{subjId, predID, objId})
}

// Fetch triples form the graph that match a BGP given in parameters
func (g *HDTGraph) Filter(subject, predicate, object core.Node) chan core.Triple {
	var wg sync.WaitGroup
	results := make(chan core.Triple)
	// fetch data in the tree & wait for the operation to be complete before closing the pipeline
	wg.Add(g.root.depth() + 1)
	go g.queryNodes(&g.root, []*core.Node{&subject, &predicate, &object}, make([]int, 0), results, &wg)
	go func() {
		defer close(results)
		wg.Wait()
	}()
	return results
}

func (g *HDTGraph) Serialize(format string) string {
	// TODO
	return ""
}
