package parser

import (
	"github.com/Callidon/joseki/rdf"
	"testing"
)

func TestReadTurtleParser(t *testing.T) {
	parser := NewTurtleParser()
	cpt := 0
	prefixes := [][]string{
		[]string{"sw", "http://www.w3.org/2001/sw/RDFCore/"},
		[]string{"foaf", "http://xmlns.com/foaf/0.1/"},
	}
	datas := []rdf.Triple{
		rdf.NewTriple(rdf.NewURI("http://www.w3.org/2001/sw/RDFCore/ntriples/"),
			rdf.NewURI("http://www.w3.org/1999/02/22-rdf-syntax-ns#type"),
			rdf.NewURI("http://xmlns.com/foaf/0.1/Document")),
		rdf.NewTriple(rdf.NewURI("http://www.w3.org/2001/sw/RDFCore/ntriples/"),
			rdf.NewURI("http://www.w3.org/1999/02/22-rdf-syntax-ns#type"),
			rdf.NewURI("http://xmlns.com/foaf/0.1/Document")),
		rdf.NewTriple(rdf.NewURI("http://www.w3.org/2001/sw/RDFCore/ntriples/"),
			rdf.NewURI("http://purl.org/dc/terms/title"),
			rdf.NewLiteral("N-Triples")),
		rdf.NewTriple(rdf.NewURI("http://www.w3.org/2001/sw/RDFCore/ntriples/"),
			rdf.NewURI("http://purl.org/dc/terms/title"),
			rdf.NewLiteral("Turtle")),
		rdf.NewTriple(rdf.NewURI("http://www.w3.org/2001/sw/RDFCore/ntriples/"),
			rdf.NewURI("http://xmlns.com/foaf/0.1/maker"),
			rdf.NewBlankNode("v0")),
		rdf.NewTriple(rdf.NewBlankNode("v0"),
			rdf.NewURI("http://purl.org/dc/terms/title"),
			rdf.NewLiteral("My Title")),
	}

	// check for triples
	for elt := range parser.Read("datas/test.ttl") {
		if test, err := elt.Equals(datas[cpt]); !test || (err != nil) {
			t.Error(elt, "should be equal to", datas[cpt])
		}
		cpt++
	}
	if cpt != len(datas) {
		t.Error("read", cpt, "nodes of the file instead of", len(datas))
	}

	// check for prefixes
	parserPrefixes := parser.Prefixes()
	for _, prefix := range prefixes {
		value, inPrefixes := parserPrefixes[prefix[0]]
		if !inPrefixes {
			t.Error("the prefix", prefix[0], "hasn't been read by the parser")
		}
		if value != prefix[1] {
			t.Error("for key", prefix[0], "expected value", prefix[1], "but got", value)
		}
	}
}
