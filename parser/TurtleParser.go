package parser

import (
	"bufio"
	"github.com/Callidon/joseki/rdf"
	"os"
    "errors"
    "strconv"
)

// Parser for reading & loading triples in Turtle format.
//
// Turtle reference : https://www.w3.org/TR/turtle/
type TurtleParser struct {
}

// Read a file containg RDF triples in Turtle format & convert them in triples.
//
// Triples generated are send throught a channel, which is closed when the parsing of the file has been completed.
func (p *TurtleParser) Read(filename string) chan rdf.Triple {
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
            bnodeCpt := 0
			line := extractSegments(scanner.Text())
            for _, elt := range line {
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
                    bnodeCpt += 1
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
                lineNumber += 1
            }
		}
		close(out)
	}()
	return out
}
