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
	"strings"
)

// TurtleParser is a parser for reading & loading triples in Turtle format.
//
// Turtle reference : https://www.w3.org/TR/turtle/
type TurtleParser struct {
	prefixes map[string]string
}

// rdfToken is a scanner for reading triples in Turtle format.
type turtleScanner struct {
}

// newTurtleScanner creates a new rdfToken
func newTurtleScanner() *turtleScanner {
	return &turtleScanner{}
}

// scan read a file in Turtle format, identify and extract token with their values.
//
// The results are sent through a channel, which is closed when the scan of the file has been completed.
func (s *turtleScanner) scan(filename string) chan rdfToken {
	out := make(chan rdfToken, bufferSize)
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
						out <- newRDFToken(tokenPrefixName, prefixName)
						out <- newRDFToken(tokenPrefixValue, prefixValue)
						prefixName, prefixValue = "", ""
					} else if prefixName == "" {
						if string(elt[len(elt)-1]) != ":" {
							out <- newRDFToken(tokenIllegal, "Unexpected token at line "+string(lineNumber)+" : "+elt)
							return
						}
						prefixName = elt[0 : len(elt)-1]
					} else if prefixValue == "" {
						if (string(elt[0]) != "<") && (string(elt[len(elt)-1]) != ">") {
							out <- newRDFToken(tokenIllegal, "Unexpected token at line "+string(lineNumber)+" : "+elt)
							return
						}
						prefixValue = elt[1 : len(elt)-1]
					}
				} else {
					// when hitting the separator, send triple into channel
					if (elt == ".") || (elt == "]") {
						out <- newRDFToken(tokenEnd, elt)
					} else if (elt == ";") || (elt == ",") || (elt == "[") {
						out <- newRDFToken(tokenSep, elt)
					} else if (string(elt[0]) == "<") && (string(elt[len(elt)-1]) == ">") {
						out <- newRDFToken(tokenURI, elt[1:len(elt)-1])
					} else if (string(elt[0]) == "\"") && (string(elt[len(elt)-1]) == "\"") {
						out <- newRDFToken(tokenLiteral, elt[1:len(elt)-1])
					} else if elt[0:2] == "^^" {
						out <- newRDFToken(tokenTypedLiteral, elt[2:])
					} else if string(elt[0]) == "@" {
						out <- newRDFToken(tokenLangLiteral, elt[1:])
					} else if (string(elt[0]) == "_") && (string(elt[1]) == ":") {
						out <- newRDFToken(tokenBlankNode, elt[2:])
					} else if strings.Index(elt, ":") > -1 {
						out <- newRDFToken(tokenPrefixedURI, elt)
					} else {
						out <- newRDFToken(tokenIllegal, "Unexpected token at line "+string(lineNumber)+" of file : bad syntax")
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
		bnodeCpt := 0
		scanner := newTurtleScanner()
		for token := range scanner.scan(filename) {
			switch token.Type {
			case tokenEnd:
				sendTriple(subject, predicate, object, out)
				subject, predicate, object = nil, nil, nil
			case tokenSep:
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
					// generate a new object & send triple and then use the new blank Node as the new subject
					object = rdf.NewBlankNode("v" + strconv.Itoa(bnodeCpt))
					sendTriple(subject, predicate, object, out)
					subject = object
					predicate, object = nil, nil
					bnodeCpt++
				default:
					panic(errors.New("Unexpected separator token " + token.Value))
				}
			case tokenPrefixName:
				prefixName = token.Value
			case tokenPrefixValue:
				p.prefixes[prefixName] = token.Value
			case tokenURI:
				assignNode(rdf.NewURI(token.Value))
			case tokenPrefixedURI:
				sepIndex := strings.Index(token.Value, ":")
				prefixURI, knownPrefix := p.prefixes[token.Value[0:sepIndex]]
				if knownPrefix {
					assignNode(rdf.NewURI(prefixURI + token.Value[sepIndex+1:]))
				} else {
					panic(errors.New("Unknown prefix " + token.Value[0:sepIndex] + " found"))
				}
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
