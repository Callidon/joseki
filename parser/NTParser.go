// Copyright (c) 2016 Thomas Minier. All rights reserved.
// Use of this source code is governed by a MIT License
// license that can be found in the LICENSE file.

package parser

import (
	"bufio"
	"github.com/Callidon/joseki/rdf"
	"io"
	"os"
)

// NTParser is a parser for reading & loading triples in N-Triples format.
//
// N-Triples reference : https://www.w3.org/2011/rdf-wg/wiki/N-Triples-Format
type NTParser struct {
	cutter *lineCutter
}

// scanNtriples read a file in N-Triples format, identify and extract token with their values.
//
// The results are sent through a channel, which is closed when the scan of the file has been completed.
func scanNtriples(reader io.Reader, out chan<- rdfToken, l *lineCutter) {
	// walk through the file using a goroutine
	go func() {
		defer close(out)

		scanner := bufio.NewScanner(reader)
		lineNumber := 1
		for scanner.Scan() {
			line := l.extractSegments(scanner.Text())
			rowNumber := 1
			// skip blank lines & comments
			if len(line) == 0 || line[0] == "#" {
				lineNumber++
				continue
			}
			// scan elements of the line
			for _, elt := range line {
				// skip to next line when a comment is detect
				if string(elt[0]) == "#" {
					break
				}
				switch {
				case elt == ".":
					out <- newTokenEnd(lineNumber, rowNumber)
				case string(elt[0]) == "<" && string(elt[len(elt)-1]) == ">":
					out <- newTokenURI(elt[1 : len(elt)-1])
				case string(elt[0]) == "_" && string(elt[1]) == ":":
					out <- newTokenBlankNode(elt[2:])
				case string(elt[0]) == "\"" && string(elt[len(elt)-1]) == "\"", string(elt[0]) == "'" && string(elt[len(elt)-1]) == "'":
					out <- newTokenLiteral(elt[1 : len(elt)-1])
				case len(elt) >= 2 && elt[0:2] == "^^":
					out <- newTokenType(elt[2:], lineNumber, rowNumber)
				case string(elt[0]) == "@":
					out <- newTokenLang(elt[1:], lineNumber, rowNumber)
				default:
					out <- newTokenIllegal("Unexpected token when scanning '"+elt+"'", lineNumber, rowNumber)
				}
				rowNumber += len(elt) + 1
			}
			lineNumber++
		}
	}()
}

// NewNTParser creates a new NTParser
func NewNTParser() *NTParser {
	return &NTParser{newLineCutter(wordRegexp)}
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
	tokenPipe := make(chan rdfToken, bufferSize)
	out := make(chan rdf.Triple, bufferSize)
	stack := newStack()

	// scan the file & analyse the tokens using a goroutine
	go func() {
		defer close(out)
		f, err := os.Open(filename)
		check(err)
		defer f.Close()
		// launch the scan, then interpret each token produced
		go scanNtriples(bufio.NewReader(f), tokenPipe, p.cutter)
		for token := range tokenPipe {
			err = token.Interpret(stack, nil, out)
			check(err)
		}
	}()
	return out
}
