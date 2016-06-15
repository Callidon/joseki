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
	prefixesOut chan string
}

// NewTurtleParser creates a new TurtleParser
func NewTurtleParser() TurtleParser {
	return TurtleParser{}
}

// Prefixes returns a channel which will be used to fetch prefixes during the parsing of a file.
func (p TurtleParser) Prefixes() chan string {
	p.prefixesOut = make(chan string)
	return p.prefixesOut
}

// Read a file containg RDF triples in Turtle format & convert them in triples.
//
// Triples generated are send throught a channel, which is closed when the parsing of the file has been completed.
func (p TurtleParser) Read(filename string) chan rdf.Triple {
	out := make(chan rdf.Triple)
	// walk through the file using a goroutine
	go func() {
		var subject, predicate, object rdf.Node
		var prefixName, prefixValue string
		var scanPrefixesDone bool
		var err error

		f, err := os.Open(filename)
		check(err)
		defer f.Close()

		scanner := bufio.NewScanner(bufio.NewReader(f))

		lineNumber := 0
		bnodeCpt := 0
		for scanner.Scan() {
			line := extractSegments(scanner.Text())
			scanPrefixesDone = (line[0] != "@prefix") || (p.prefixesOut == nil)
			for _, elt := range line {
				// scan for prefixes until they have been all found, then scan for triples
				if !scanPrefixesDone {
					// skip to next element if reading the @prefix keyword
					if elt == "@prefix" {
						continue
					} else if elt == "." {
						// when hitting the separator, send the prefix (name & value)
						p.prefixesOut <- prefixName
						p.prefixesOut <- prefixValue
						prefixName, prefixValue = "", ""
					} else if prefixName == "" {
						prefixName = elt
					} else if prefixValue == "" {
						prefixValue = elt
					} else {
						err = errors.New("Error at line " + string(lineNumber) + " of file during prefixes scan : bad syntax")
					}
				} else {
					// when hitting the separator, send triple into channel
					if (elt == ".") || (elt == "]") {
						sendTriple(subject, predicate, object, out)
						// reset the value
						subject, predicate, object = nil, nil, nil
					} else if elt == ";" {
						// send previous value & keep subject for the next triple
						sendTriple(subject, predicate, object, out)
						predicate, object = nil, nil
					} else if elt == "," {
						// send previous value & keep subject and predicate ofr the next triple
						sendTriple(subject, predicate, object, out)
						object = nil
					} else if elt == "[" {
						// generate a new objectn send triple and then use the new blank Node as the new subject
						object = rdf.NewBlankNode("v" + strconv.Itoa(bnodeCpt))
						sendTriple(subject, predicate, object, out)
						subject = object
						predicate, object = nil, nil
						bnodeCpt++
					} else if subject == nil {
						subject, err = parseNode(elt)
					} else if predicate == nil {
						predicate, err = parseNode(elt)
					} else if object == nil {
						object, err = parseNode(elt)
					} else {
						err = errors.New("Error at line " + string(lineNumber) + " of file : bad syntax")
					}
				}
				// check for error during the parsing
				check(err)
				lineNumber++
			}
		}
		close(out)
        if p.prefixesOut != nil {
            close(p.prefixesOut)
        }
	}()
	return out
}
