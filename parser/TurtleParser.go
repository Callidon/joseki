// Copyright (c) 2016 Thomas Minier. All rights reserved.
// Use of this source code is governed by a MIT License
// license that can be found in the LICENSE file.

package parser

import (
	"bufio"
	"github.com/Callidon/joseki/parser/tokens"
	"github.com/Callidon/joseki/rdf"
	"io"
	"os"
	"strings"
)

// TurtleParser is a parser for reading & loading triples in Turtle format.
//
// Turtle reference : https://www.w3.org/TR/turtle/
type TurtleParser struct {
	prefixes map[string]string
}

// scanTurtle read a file in Turtle format, identify and extract token with their values.
//
// The results are sent through a channel, which is closed when the scan of the file has been completed.
func scanTurtle(reader io.Reader) chan tokens.RDFToken {
	out := make(chan tokens.RDFToken, bufferSize)
	// walk through the file using a goroutine
	go func() {
		defer close(out)
		var prefixName, prefixValue string
		var scanPrefixesDone bool

		scanner := bufio.NewScanner(reader)
		lineNumber, rowNumber := 1, 1
		for scanner.Scan() {
			line := extractSegments(scanner.Text())
			// skip blank lines & comments
			if (len(line) == 0) || (line[0] == "#") {
				continue
			}
			scanPrefixesDone = (line[0] != "@prefix")
			// scan elements of the line
			for _, elt := range line {
				// skip comments
				if string(elt[0]) == "#" {
					break
				}
				if !scanPrefixesDone {
					if (elt == "@prefix") || (elt == ":") {
						continue
					} else if elt == "." {
						out <- tokens.NewTokenPrefix(prefixName, prefixValue)
						prefixName, prefixValue = "", ""
					} else if prefixName == "" {
						if string(elt[len(elt)-1]) != ":" {
							out <- tokens.NewTokenIllegal("Unexpected token "+elt, lineNumber, rowNumber)
							return
						}
						prefixName = elt[0 : len(elt)-1]
					} else if prefixValue == "" {
						if (string(elt[0]) != "<") && (string(elt[len(elt)-1]) != ">") {
							out <- tokens.NewTokenIllegal("Unexpected token "+elt, lineNumber, rowNumber)
							return
						}
						prefixValue = elt[1 : len(elt)-1]
					}
				} else {
					// when hitting the separator, send triple into channel
					if (elt == ".") || (elt == "]") {
						out <- tokens.NewTokenEnd(lineNumber, rowNumber)
					} else if (elt == ";") || (elt == ",") || (elt == "[") {
						out <- tokens.NewTokenSep(elt, lineNumber, rowNumber)
					} else if (string(elt[0]) == "<") && (string(elt[len(elt)-1]) == ">") {
						out <- tokens.NewTokenURI(elt[1 : len(elt)-1])
					} else if ((string(elt[0]) == "\"") && (string(elt[len(elt)-1]) == "\"")) || ((string(elt[0]) == "'") && (string(elt[len(elt)-1]) == "'")) {
						out <- tokens.NewTokenLiteral(elt[1 : len(elt)-1])
					} else if elt[0:2] == "^^" {
						out <- tokens.NewTokenType(elt[2:], lineNumber, rowNumber)
					} else if string(elt[0]) == "@" {
						out <- tokens.NewTokenLang(elt[1:], lineNumber, rowNumber)
					} else if (string(elt[0]) == "_") && (string(elt[1]) == ":") {
						out <- tokens.NewTokenBlankNode(elt[2:])
					} else if string(elt[0]) == "?" {
						out <- tokens.NewTokenBlankNode(elt[1:])
					} else if strings.Index(elt, ":") > -1 {
						out <- tokens.NewTokenPrefixedURI(elt, lineNumber, rowNumber)
					} else {
						out <- tokens.NewTokenIllegal("Unexpected token when scanning "+elt, lineNumber, rowNumber)
					}
				}
				rowNumber += len(elt) + 1
			}
			lineNumber++
		}
	}()
	return out
}

// NewTurtleParser creates a new TurtleParser
func NewTurtleParser() *TurtleParser {
	return &TurtleParser{make(map[string]string)}
}

// Prefixes returns the prefixes read by the parser during the last parsing.
func (p TurtleParser) Prefixes() map[string]string {
	return p.prefixes
}

// Read a file containg RDF triples in Turtle format & convert them in triples.
//
// Triples generated are send throught a channel, which is closed when the parsing of the file has been completed.
func (p *TurtleParser) Read(filename string) chan rdf.Triple {
	out := make(chan rdf.Triple, bufferSize)
	stack := tokens.NewStack()

	// scan the file & analyse the tokens using a goroutine
	go func() {
		defer close(out)
		f, err := os.Open(filename)
		check(err)
		defer f.Close()

		for token := range scanTurtle(bufio.NewReader(f)) {
			err = token.Interpret(stack, &p.prefixes, out)
			check(err)
		}
	}()
	return out
}
