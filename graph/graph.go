// Package graph provides various implementation for a RDF Graph
package graph

import (
	"errors"
	"github.com/Callidon/joseki/parser"
	"github.com/Callidon/joseki/rdf"
	"strings"
)

// Graph represent a generic RDF Graph
//
// Package graph provides several implementations for this interface.
type Graph interface {
	// Add a new Triple pattern to the graph.
	Add(triple rdf.Triple)
	// Delete triples from the graph that match a BGP given in parameters.
	Delete(subject, object, predicate rdf.Node)
	// Fetch triples form the graph that match a BGP given in parameters.
	Filter(subject, predicate, object rdf.Node) chan rdf.Triple
	// Serialize the graph into a given format and return it as a string.
	Serialize(format string) string
}

// rdfReader represent a reader capable of reading RDF data encoded in various format.
//
// This structure is designed to be embedded into types which implement the Graph interface
type rdfReader struct {
	graph Graph
	// list of prefixes used in some RDF formats (Turtle, JSON-LD, ...)
	prefixes map[string]string
}

// newRDFReader creates a new rdfReader
func newRDFReader() *rdfReader {
	return &rdfReader{nil, nil}
}

// Generic function for loading triples from a file into a graph, with a given format
func (r *rdfReader) LoadFromFile(filename string, format string) {
	var p parser.Parser
	hasPrefixes := false
	// determine which parser to use depending on the format
	switch strings.ToLower(format) {
	case "nt", "n-triples":
		p = parser.NewNTParser()
	case "ttl", "turtle":
		p = parser.NewTurtleParser()
		hasPrefixes = true
	default:
		panic(errors.New("Error : " + format + " is not a supported format." +
			"Please see the documentation at https://godoc.org/github.com/Callidon/joseki/parser to see the available parsers."))
	}
	// read triples from file, then load prefixes if necessary
	for triple := range p.Read(filename) {
		r.graph.Add(triple)
	}
	if hasPrefixes {
		r.prefixes = p.Prefixes()
	}
}
