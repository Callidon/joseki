package parser

import "github.com/Callidon/joseki/rdf"

type Parser interface {
	Read(filename string) chan rdf.Triple
}
