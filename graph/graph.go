// Package graph provides various implementation for a RDF Graph
package graph

import "github.com/Callidon/joseki/rdf"

// Graph represent a generic RDF Graph
//
// Package graph provides several implementations for this interface.
type Graph interface {
	// Load the content of a RDF graph stored in a file into the current graph.
	LoadFromFile(filename, format string)
	// Add a new Triple pattern to the graph.
	Add(triple rdf.Triple)
	// Fetch triples form the graph that match a BGP given in parameters.
	Filter(subject, predicate, object rdf.Node) chan rdf.Triple
	// Serialize the graph into a given format and return it as a string.
	Serialize(format string) string
}
