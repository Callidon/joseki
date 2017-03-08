package sparql

import (
	"github.com/Callidon/joseki/rdf"
	"sort"
)

// selectNode represents a SELCECT operation, a projection over a set of variables
type selectNode struct {
	previous   queryNode
	projection []string
}

// newSelectNode creates a new selectNode.
func newSelectNode(previous queryNode, projection ...string) *selectNode {
	s := &selectNode{previous, projection}
	sort.Strings(s.projection)
	return s
}

// get fetch the bindings from previous node, apply the projection to them, and forward them to the operator
func (n selectNode) get() <-chan rdf.BindingsGroup {
	results := make(chan rdf.BindingsGroup, bufferSize)
	go func() {
		defer close(results)
		for b := range n.previous.get() {
			newGroup := rdf.NewBindingsGroup()
			for key, value := range b.Bindings {
        // use binary search as projection variables are always sorted
				i := sort.Search(len(n.projection), func(i int) bool { return n.projection[i] >= key })
				if i < len(n.projection) && n.projection[i] == key {
					newGroup.Bindings[key] = value
				}
			}
			results <- newGroup
		}
	}()
	return results
}
