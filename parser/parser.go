// Package parser provides parser to work with several RDF formats (N-Triples, Turtles, JSON-LD, ...)
package parser

import (
	"errors"
	"github.com/Callidon/joseki/rdf"
	"regexp"
)

// Parser represent a generic interface for parsing every RDF format.
//
// Package parser provides several implementations for this interface.
type Parser interface {
	Read(filename string) chan rdf.Triple
	Prefixes() map[string]string
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

// extractSegments parse a string and split the segments into a slice.
// A segment is a string quoted or separated from the other by whitespaces.
func extractSegments(line string) []string {
	r := regexp.MustCompile("'.*?'|\".*?\"|\\S+")
	return r.FindAllString(line, -1)
}
