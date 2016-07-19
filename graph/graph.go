// Copyright (c) 2016 Thomas Minier. All rights reserved.
// Use of this source code is governed by a MIT License
// license that can be found in the LICENSE file.

// Package graph provides various implementation for a RDF Graph
package graph

import (
	"errors"
	"github.com/Callidon/joseki/parser"
	"github.com/Callidon/joseki/rdf"
	"strings"
)

// Graph represents a generic RDF Graph
//
// Package graph provides several implementations for this interface.
// RDF Graph reference : https://www.w3.org/TR/rdf11-concepts/#section-rdf-graph
type Graph interface {
	// Add a new Triple pattern to the graph.
	Add(triple rdf.Triple)
	// Delete triples from the graph that match a BGP given in parameters.
	Delete(subject, predicate, object rdf.Node)
	// Fetch triples form the graph that match a BGP given in parameters.
	Filter(subject, predicate, object rdf.Node) <-chan rdf.Triple
	// Same as Filter, but with a Limit and an Offset
	FilterSubset(subject rdf.Node, predicate rdf.Node, object rdf.Node, limit int, offset int) <-chan rdf.Triple
}

// rdfReader represents a reader capable of reading RDF data encoded in various format.
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

// LoadFromFile loads triples from a file into a graph, with a given format
// In the desired format isn't supported or doesn't exist, no new triples will
// be inserted into the graph and an error will be returned.
func (r *rdfReader) LoadFromFile(filename string, format string) error {
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
		return errors.New("Error : " + format + " is not a supported format." +
			"Please see the documentation at https://godoc.org/github.com/Callidon/joseki/parser to see the available parsers.")
	}
	// read triples from file, then load prefixes if necessary
	for triple := range p.Read(filename) {
		r.graph.Add(triple)
	}
	if hasPrefixes {
		r.prefixes = p.Prefixes()
	}
	return nil
}

// Utility function for checking errors
func check(err error) {
	if err != nil {
		panic(err)
	}
}
