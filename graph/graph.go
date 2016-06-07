package graph

import (
	"github.com/Callidon/joseki/core"
	"os"
)

// Generic representation of a RDF Graph
// Various implementation are proposed in the joseki/graph package
type Graph interface {
	LoadFromFile(file *os.File)
	Add(triple core.Triple)
	Filter(subject, predicate, object core.Node) []core.Triple
	Serialize(format string) string
}
