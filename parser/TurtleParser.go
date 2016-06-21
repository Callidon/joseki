// Copyright (c) 2016 Thomas Minier. All rights reserved.
// Use of this source code is governed by a MIT License
// license that can be found in the LICENSE file.

package parser

import (
	"bufio"
	"errors"
	"github.com/Callidon/joseki/rdf"
	"os"
	"strconv"
)

// TurtleParser is a parser for reading & loading triples in Turtle format.
//
// Turtle reference : https://www.w3.org/TR/turtle/
type TurtleParser struct {
	prefixes map[string]string
}

// TurtleScanner is a scanner for reading triples in Turtle format.
type TurtleScanner struct {
}

// NewTurtleScanner creates a new TurtleScanner
func NewTurtleScanner() *TurtleScanner {
	return &TurtleScanner{}
}

// Scan read a file in Turtle format, identify and extract token with their values.
//
// The results are sent through a channel, which is closed when the scan of the file has been completed.
func (s *TurtleScanner) Scan(filename string) chan RDFToken {
	out := make(chan RDFToken, bufferSize)
	// walk through the file using a goroutine
	go func() {
		defer close(out)
		var prefixName, prefixValue string
		var scanPrefixesDone bool

		f, err := os.Open(filename)
		check(err)
		defer f.Close()

		scanner := bufio.NewScanner(bufio.NewReader(f))
		lineNumber := 0
		for scanner.Scan() {
			line := extractSegments(scanner.Text())
			scanPrefixesDone = (line[0] != "@prefix")
			for _, elt := range line {
				if !scanPrefixesDone {
					if (elt == "@prefix") || (elt == ":") {
						continue
					} else if elt == "." {
						out <- NewRDFToken(TokenPrefixName, prefixName)
						out <- NewRDFToken(TokenPrefixValue, prefixValue)
						prefixName, prefixValue = "", ""
					} else if prefixName == "" {
						if string(elt[len(elt)-1]) != ":" {
							out <- NewRDFToken(TokenIllegal, "Unexpected token at line "+string(lineNumber)+" : "+elt)
							return
						}
						prefixName = elt[0 : len(elt)-1]
					} else if prefixValue == "" {
						if (string(elt[0]) != "<") && (string(elt[len(elt)-1]) != ">") {
							out <- NewRDFToken(TokenIllegal, "Unexpected token at line "+string(lineNumber)+" : "+elt)
							return
						}
						prefixValue = elt[1 : len(elt)-1]
					}
				} else {
					// when hitting the separator, send triple into channel
					if (elt == ".") || (elt == "]") {
						out <- NewRDFToken(TokenEnd, elt)
					} else if (elt == ";") || (elt == ",") || (elt == "[") {
						out <- NewRDFToken(TokenSep, elt)
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
				}
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
	var subject, predicate, object rdf.Node
	var prefixName string
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
		bnodeCpt := 0
		scanner := NewTurtleScanner()
		for token := range scanner.Scan(filename) {
			switch token.Type {
			case TokenEnd:
				sendTriple(subject, predicate, object, out)
				subject, predicate, object = nil, nil, nil
			case TokenSep:
				switch token.Value {
				case ";":
					// send previous value & keep subject for the next triple
					sendTriple(subject, predicate, object, out)
					predicate, object = nil, nil
				case ",":
					// send previous value & keep subject and predicate ofr the next triple
					sendTriple(subject, predicate, object, out)
					object = nil
				case "[":
					// generate a new objectn send triple and then use the new blank Node as the new subject
					object = rdf.NewBlankNode("v" + strconv.Itoa(bnodeCpt))
					sendTriple(subject, predicate, object, out)
					subject = object
					predicate, object = nil, nil
					bnodeCpt++
				default:
					panic(errors.New("Unexpected separator token " + token.Value))
				}
			case TokenPrefixName:
				prefixName = token.Value
			case TokenPrefixValue:
				p.prefixes[prefixName] = token.Value
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
	}()
	return out
}
