// Package joseki/parser provides parser to work with several RDF formats (N-Triples, Turtles, JSON-LD, ...)
package parser

import (
    "github.com/Callidon/joseki/rdf"
	"errors"
    "regexp"
)

// Generic parser interface for every RDF format.
//
// Package joseki/parser provides several implementations for these parsers.
type Parser interface {
	Read(filename string) chan rdf.Triple
}

// RDF file reader wich can be used with multiples formats.
//
// Works using different separators depending of the format used.
type rdfReader struct {
    separator string
    predListSeparator string
    objectListSeparator string
    nestingBegin string
    nestingEnd string
}

// Utility function for checking errors
func check(err error) {
	if err != nil {
		panic(err)
	}
}

// Send RDF Node of a triple pattern throught a channel
func sendTriple(subject rdf.Node, predicate rdf.Node, object rdf.Node, out chan rdf.Triple) {
    out <- rdf.NewTriple(subject, predicate, object)
}

// Parse a string node to find its type & return the corresponding RDF Node
func parseNode(elt string) (rdf.Node, error) {
	var node rdf.Node
	var err error
	if (string(elt[0]) == "<") && (string(elt[len(elt)-1]) == ">") {
		node = rdf.NewURI(elt[1 : len(elt)-1])
	} else if (string(elt[0]) == "\"") && (string(elt[len(elt)-1]) == "\"") {
		// TODO add a security when a xml type is given with the literal
		node = rdf.NewLiteral(elt[1 : len(elt)-1])
	} else if (string(elt[0]) == "_") && (string(elt[1]) == ":") {
		node = rdf.NewBlankNode(elt[2:])
	} else {
		err = errors.New("Error : cannot parse " + elt + " into a RDF Node")
	}
	return node, err
}

func extractSegments(line string) []string {
    r := regexp.MustCompile("'.*?'|\".*?\"|\\S+")
    return r.FindAllString(line, -1)
}
