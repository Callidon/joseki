// Package joseki/parser provides parser to work with several RDF formats (N-Triples, Turtles, JSON-LD, ...)
package parser

import "github.com/Callidon/joseki/rdf"

// Generic parser interface for every RDF format.
//
// Package joseki/parser provides several implementations for these parsers.
type Parser interface {
	Read(filename string) chan rdf.Triple
}
