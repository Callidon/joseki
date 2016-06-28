// Copyright (c) 2016 Thomas Minier. All rights reserved.
// Use of this source code is governed by a MIT License
// license that can be found in the LICENSE file.

package parser

import (
	"bufio"
	"github.com/Callidon/joseki/parser/tokens"
	"github.com/Callidon/joseki/rdf"
	"os"
)

// NTParser is a parser for reading & loading triples in N-Triples format.
//
// N-Triples reference : https://www.w3.org/2011/rdf-wg/wiki/N-Triples-Format
type NTParser struct {
}

// scanNtriples read a file in N-Triples format, identify and extract token with their values.
//
// The results are sent through a channel, which is closed when the scan of the file has been completed.
func scanNtriples(filename string) chan tokens.RDFToken {
	out := make(chan tokens.RDFToken, bufferSize)
	// walk through the file using a goroutine
	go func() {
		defer close(out)
		f, err := os.Open(filename)
		check(err)
		defer f.Close()

		scanner := bufio.NewScanner(bufio.NewReader(f))
		lineNumber, rowNumber := 1, 1
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
					out <- tokens.NewTokenEnd(lineNumber, rowNumber)
				} else if (string(elt[0]) == "<") && (string(elt[len(elt)-1]) == ">") {
					out <- tokens.NewTokenURI(elt[1 : len(elt)-1])
				} else if (string(elt[0]) == "_") && (string(elt[1]) == ":") {
					out <- tokens.NewTokenBlankNode(elt[2:])
				} else if ((string(elt[0]) == "\"") && (string(elt[len(elt)-1]) == "\"")) || ((string(elt[0]) == "'") && (string(elt[len(elt)-1]) == "'")) {
					out <- tokens.NewTokenLiteral(elt[1 : len(elt)-1])
				} else if elt[0:2] == "^^" {
					out <- tokens.NewTokenType(elt[2:], lineNumber, rowNumber)
				} else if string(elt[0]) == "@" {
					out <- tokens.NewTokenLang(elt[1:], lineNumber, rowNumber)
				} else {
					out <- tokens.NewTokenIllegal("Unexpected token when scanning "+elt, lineNumber, rowNumber)
				}
				rowNumber += len(elt) + 1
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
	var err error
	out := make(chan rdf.Triple, bufferSize)
	stack := tokens.NewStack()

	// scan the file & analyse the tokens using a goroutine
	go func() {
		defer close(out)
		for token := range scanNtriples(filename) {
			err = token.Interpret(stack, nil, out)
			check(err)
		}
	}()
	return out
}
