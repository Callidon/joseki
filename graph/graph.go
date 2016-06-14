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
	// Load the content of a RDF graph stored in a file into the current graph.
	LoadFromFile(filename, format string)
	// Add a new Triple pattern to the graph.
	Add(triple rdf.Triple)
	// Delete triples from the graph that match a BGP given in parameters.
	Delete(subject, object, predicate rdf.Node)
	// Fetch triples form the graph that match a BGP given in parameters.
	Filter(subject, predicate, object rdf.Node) chan rdf.Triple
	// Serialize the graph into a given format and return it as a string.
	Serialize(format string) string
}

// Generic function for loading triples from a file into a graph, with a given format
func loadFromFile(g Graph, filename string, format string) {
	var p parser.Parser
	// determine which parser to use depending on the format
	switch strings.ToLower(format) {
	case "nt", "n-triples":
		p = parser.NewNTParser()
	case "ttl", "turtle":
		p = parser.NewTurtleParser()
	default:
		panic(errors.New("Error : " + format + " is not a supported format." +
			"Please see the documentation at https://godoc.org/github.com/Callidon/joseki/parser to see the available parsers."))
	}
	for triple := range p.Read(filename) {
		g.Add(triple)
	}
}
