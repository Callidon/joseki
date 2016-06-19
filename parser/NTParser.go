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
// Triples generated are send throught a channel, which is closed when the parsing of the file has been completed.
func (p NTParser) Read(filename string) chan rdf.Triple {
	var subject, predicate, object rdf.Node
	out := make(chan rdf.Triple)
	// walk through the file using a goroutine
	go func() {
		f, err := os.Open(filename)
		check(err)
		defer f.Close()

		scanner := bufio.NewScanner(bufio.NewReader(f))
		for scanner.Scan() {
			var err error
			lineNumber := 0
			line := extractSegments(scanner.Text())
			for _, elt := range line {
				// when hitting the separator, send triple into channel
				if elt == "." {
					sendTriple(subject, predicate, object, out)
					// reset the value
					subject, predicate, object = nil, nil, nil
				} else if subject == nil {
					subject, err = parseNode(elt)
				} else if predicate == nil {
					predicate, err = parseNode(elt)
				} else if object == nil {
					object, err = parseNode(elt)
				} else {
					err = errors.New("Error at line " + string(lineNumber) + " of file : bad syntax")
				}
				// check for error during the parsing
				check(err)
				lineNumber++
			}
		}
		close(out)
	}()
	return out
}
