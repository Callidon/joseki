// Package joseki/graph provides various implementation for RDF Graph
package graph

import (
	"github.com/Callidon/joseki/rdf"
	"os"
)

// Generic representation of a RDF Graph
// Various implementation are proposed in the joseki/graph package
type Graph interface {
	LoadFromFile(file *os.File)
	Add(triple rdf.Triple)
	Filter(subject, predicate, object rdf.Node) chan rdf.Triple
	Serialize(format string) string
}
