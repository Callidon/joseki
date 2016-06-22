// Copyright (c) 2016 Thomas Minier. All rights reserved.
// Use of this source code is governed by a MIT License
// license that can be found in the LICENSE file.

// Package parser provides parser to work with several RDF formats (N-Triples, Turtles, JSON-LD, ...)
package parser

import (
	"github.com/Callidon/joseki/rdf"
	"regexp"
)

const (
	// Max size for the buffer of this package
	bufferSize = 100
)

// token is the type for a token read by a scanner
type token float64

const (
	_ = iota
	// tokenIllegal is an illegal token in the RDF syntax
	tokenIllegal token = 1 << (10 * iota)
	// tokenEnd ends a triple declaration
	tokenEnd
	// tokenSep is a RDF separator (for object/literal list, etc)
	tokenSep
	// tokenURI is a RDF URI
	tokenURI
	// tokenPrefixedURI is a RDF URI with a prefix
	tokenPrefixedURI
	// tokenLiteral is a RDF Literal
	tokenLiteral
	// tokenTypedLiteral is a RDF typed Literal
	tokenTypedLiteral
	// tokenLangLiteral is a RDF Literal with lang informations
	tokenLangLiteral
	// tokenBlankNode is a RDF Blank Node
	tokenBlankNode
	// tokenPrefixName is the name of a prefix
	tokenPrefixName
	// tokenPrefixValue is the value of a prefix
	tokenPrefixValue
)

// Parser represent a generic interface for parsing every RDF format.
//
// Package parser provides several implementations for this interface.
type Parser interface {
	Read(filename string) chan rdf.Triple
	Prefixes() map[string]string
}

// rdfScanner represent a generic interface for scanning an RDF file in every format.
// This interface act as a Lexer during the parsing process.
//
// Package parser provides several implementations for this interface.
type rdfScanner interface {
	Scan(filename string) chan rdfToken
}

// rdfToken is a token extracted during the scan of a RDF file.
// It's meant to be used by a Parser implementation during the parsing phase.
type rdfToken struct {
	Type  token
	Value string
}

// newRDFToken creates a new rdfToken
func newRDFToken(tokType token, value string) rdfToken {
	return rdfToken{tokType, value}
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

// extractSegments parse a string and split the segments into a slice.
// A segment is a string quoted or separated from the other by whitespaces.
func extractSegments(line string) []string {
	r := regexp.MustCompile("'.*?'|\".*?\"|\\S+")
	return r.FindAllString(line, -1)
}
