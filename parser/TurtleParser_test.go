package parser

import (
	"github.com/Callidon/joseki/rdf"
	"testing"
)

func TestReadTurtleParser(t *testing.T) {
	parser := TurtleParser{}
	cpt := 0
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

	for elt := range parser.Read("datas/test.ttl") {
		if test, err := elt.Equals(datas[cpt]); !test || (err != nil) {
			t.Error(elt, "should be equal to", datas[cpt])
		}
		cpt++
	}

	if cpt != len(datas) {
		t.Error("read", cpt, "nodes of the file instead of", len(datas))
	}
}
