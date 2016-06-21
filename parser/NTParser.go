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

// NTScanner is a scanner for reading triples in N-Triples format.
type NTScanner struct {
}

// NewNTScanner creates a new NTScanner
func NewNTScanner() *NTScanner {
	return &NTScanner{}
}

// Scan read a file in N-Triples format, identify and extract token with their values.
//
// The results are sent through a channel, which is closed when the scan of the file has been completed.
func (s *NTScanner) Scan(filename string) chan RDFToken {
	out := make(chan RDFToken)
	// walk through the file using a goroutine
	go func() {
		f, err := os.Open(filename)
		check(err)
		defer f.Close()

		scanner := bufio.NewScanner(bufio.NewReader(f))
		for scanner.Scan() {
			lineNumber := 0
			line := extractSegments(scanner.Text())
			for _, elt := range line {
				if elt == "." {
					out <- NewRDFToken(TokenEnd, ".")
				} else if (string(elt[0]) == "<") && (string(elt[len(elt)-1]) == ">") {
					out <- NewRDFToken(TokenURI, elt[1:len(elt)-1])
				} else if (string(elt[0]) == "\"") && (string(elt[len(elt)-1]) == "\"") {
					// TODO add a security when a xml type or a lang is given with the literal
					out <- NewRDFToken(TokenLiteral, elt[1:len(elt)-1])
				} else if (string(elt[0]) == "_") && (string(elt[1]) == ":") {
					out <- NewRDFToken(TokenBlankNode, elt[2:])
				} else {
					out <- NewRDFToken(TokenIllegal, "Unexpected token at line "+string(lineNumber)+" of file : bad syntax")
				}
				lineNumber++
			}
		}
		close(out)
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
	out := make(chan rdf.Triple)

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
		scanner := NewNTScanner()
		for token := range scanner.Scan(filename) {
			switch token.Type {
			case TokenEnd:
				sendTriple(subject, predicate, object, out)
				// reset the values
				subject, predicate, object = nil, nil, nil
			case TokenURI:
				assignNode(rdf.NewURI(token.Value))
			case TokenBlankNode:
				assignNode(rdf.NewBlankNode(token.Value))
			case TokenLiteral:
				assignNode(rdf.NewLiteral(token.Value))
			case TokenIllegal:
				panic(token.Value)
			default:
				panic(errors.New("Unexpected token " + token.Value))
			}
		}
		close(out)
	}()
	return out
}
