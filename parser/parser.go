// Copyright (c) 2016 Thomas Minier. All rights reserved.
// Use of this source code is governed by a MIT License
// license that can be found in the LICENSE file.

// Package parser provides parser to work with several RDF formats (N-Triples, Turtles, JSON-LD, ...)
package parser

import (
	"errors"
	"github.com/Callidon/joseki/rdf"
	"regexp"
)

// Token is the type for a token read by a scanner
type Token float64

const (
	_ = iota
	// TokenIllegal is an illegal token in a RDF syntax
	TokenIllegal Token = 1 << (10 * iota)
	// TokenEnd ends a triple declaration
	TokenEnd
	// TokenSep is a RDF separator (for object/literal list, etc)
	TokenSep
	// TokenURI is a RDF URI
	TokenURI
	// TokenLiteral is a RDF Literal
	TokenLiteral
	// TokenTypedLiteral is a RDF typed Literal
	TokenTypedLiteral
	// TokenLangLiteral is a RDF Literal with lang informations
	TokenLangLiteral
	// TokenBlankNode is a RDF Blank Node
	TokenBlankNode
)

// Parser represent a generic interface for parsing every RDF format.
//
// Package parser provides several implementations for this interface.
type Parser interface {
	Read(filename string) chan rdf.Triple
	Prefixes() map[string]string
}

// RDFScanner represent a generic interface for scanning an RDF file in every format.
// This interface act as a Lexer during the parsing process.
//
// Package parser provides several implementations for this interface.
type RDFScanner interface {
	Scan(filename string) chan RDFToken
}

// RDFToken is a token extracted during the scan of a RDF file.
// It's meant to be used by a Parser implementation during the parsing phase.
type RDFToken struct {
	Type  Token
	Value string
}

// NewRDFToken creates a new RDFToken
func NewRDFToken(tokType Token, value string) RDFToken {
	return RDFToken{tokType, value}
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
