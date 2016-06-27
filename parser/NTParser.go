// Copyright (c) 2016 Thomas Minier. All rights reserved.
// Use of this source code is governed by a MIT License
// license that can be found in the LICENSE file.

package parser

import (
	"bufio"
	"errors"
	"github.com/Callidon/joseki/rdf"
	"os"
)

// NTParser is a parser for reading & loading triples in N-Triples format.
//
// N-Triples reference : https://www.w3.org/2011/rdf-wg/wiki/N-Triples-Format
type NTParser struct {
}

// ntScanner is a scanner for reading triples in N-Triples format.
type ntScanner struct {
}

// newNTScanner creates a new ntScanner
func newNTScanner() *ntScanner {
	return &ntScanner{}
}

// Scan read a file in N-Triples format, identify and extract token with their values.
//
// The results are sent through a channel, which is closed when the scan of the file has been completed.
func (s *ntScanner) scan(filename string) chan rdfToken {
	out := make(chan rdfToken, bufferSize)
	// walk through the file using a goroutine
	go func() {
		defer close(out)
		f, err := os.Open(filename)
		check(err)
		defer f.Close()

		scanner := bufio.NewScanner(bufio.NewReader(f))
		lineNumber := 0
		for scanner.Scan() {
			line := extractSegments(scanner.Text())
			// skip blank lines & comments
			if (len(line) == 0) || (line[0] == "#") {
				continue
			}
			// scan elements of the line
			for _, elt := range line {
				if string(elt[0]) == "#" {
					break
				} else if elt == "." {
					out <- newRDFToken(tokenEnd, ".")
				} else if (string(elt[0]) == "<") && (string(elt[len(elt)-1]) == ">") {
					out <- newRDFToken(tokenURI, elt[1:len(elt)-1])
				} else if (string(elt[0]) == "_") && (string(elt[1]) == ":") {
					out <- newRDFToken(tokenBlankNode, elt[2:])
				} else if ((string(elt[0]) == "\"") && (string(elt[len(elt)-1]) == "\"")) || ((string(elt[0]) == "'") && (string(elt[len(elt)-1]) == "'")) {
					out <- newRDFToken(tokenLiteral, elt[1:len(elt)-1])
				} else if elt[0:2] == "^^" {
					out <- newRDFToken(tokenTypedLiteral, elt[2:])
				} else if string(elt[0]) == "@" {
					out <- newRDFToken(tokenLangLiteral, elt[1:])
				} else {
					out <- newRDFToken(tokenIllegal, "Unexpected token when scanning "+elt)
				}
			}
			lineNumber++
		}
	}()
	return out
}

// NewNTParser creates a new NTParser
func NewNTParser() *NTParser {
	return &NTParser{}
}

// Prefixes returns the prefixes read by the parser during the last parsing.
// Since N-Triples format doesn't use prefixes, this function always return nil.
func (p NTParser) Prefixes() map[string]string {
	return nil
}

// Read a file containg RDF triples in N-Triples format & convert them in triples.
//
// Triples generated are send through a channel, which is closed when the parsing of the file has been completed.
func (p NTParser) Read(filename string) chan rdf.Triple {
	var subject, predicate, object rdf.Node
	var literalValue string
	out := make(chan rdf.Triple, bufferSize)
	// utility function for assigning a value to the first available node
	assignNode := func(value rdf.Node) {
		if subject == nil {
			subject = value
		} else if predicate == nil {
			predicate = value
		} else if object == nil {
			object = value
		}
	}
	// scan the file & analyse the tokens using a goroutine
	go func() {
		defer close(out)
		scanner := newNTScanner()
		for token := range scanner.scan(filename) {
			switch token.Type {
			case tokenEnd:
				sendTriple(subject, predicate, object, out)
				// reset the values
				subject, predicate, object = nil, nil, nil
			case tokenURI:
				assignNode(rdf.NewURI(token.Value))
			case tokenBlankNode:
				assignNode(rdf.NewBlankNode(token.Value))
			case tokenLiteral:
				assignNode(rdf.NewLiteral(token.Value))
				literalValue = token.Value
			case tokenTypedLiteral:
				_, ok := object.(rdf.Literal)
				if ok {
					object = rdf.NewTypedLiteral(literalValue, token.Value)
				} else {
					panic(errors.New("Trying to assign a type to a non literal object"))
				}
			case tokenLangLiteral:
				_, ok := object.(rdf.Literal)
				if ok {
					object = rdf.NewLangLiteral(literalValue, token.Value)
				} else {
					panic(errors.New("Trying to assign a language to a non literal object"))
				}
			case tokenIllegal:
				panic(token.Value)
			default:
				panic(errors.New("Unexpected token " + token.Value))
			}
		}
	}()
	return out
}
