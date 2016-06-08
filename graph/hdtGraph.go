package graph

import (
	"github.com/Callidon/joseki/core"
	"os"
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
    root    bitmapNode
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
func (g *HDTGraph) queryNodes(root *bitmapNode, datas []*core.Node, triple []int) {
    // if a triple pattern can be formed using the datas collected
    if len(triple) == 3 {

    } else {

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

func (g *HDTGraph) Filter(subject, predicate, object core.Node) chan core.Triple {
	results := make(chan core.Triple)
	// TODO
	return results
}

func (g *HDTGraph) Serialize(format string) string {
	// TODO
	return ""
}
