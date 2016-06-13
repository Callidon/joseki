package parser

import (
	"github.com/Callidon/joseki/rdf"
	"testing"
)

func TestReadNTParser(t *testing.T) {
	parser := NTParser{}
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
			rdf.NewURI("http://xmlns.com/foaf/0.1/maker"),
			rdf.NewBlankNode("art")),
	}

	for elt := range parser.Read("datas/test.nt") {
		if test, err := elt.Equals(datas[cpt]); !test || (err != nil) {
			t.Error(datas[cpt], "should be equal to", elt)
		}
		cpt++
	}

	if cpt != len(datas) {
		t.Error("read", cpt, "nodes of the file instead of", len(datas))
	}
}
