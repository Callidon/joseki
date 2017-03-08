package sparql

import "github.com/Callidon/joseki/rdf"

const (
	// Max size for the buffer of this package
	bufferSize = 100
)

// queryNode represents a generic node in a query execution plan.
type queryNode interface {
	get() <-chan rdf.BindingsGroup
}
