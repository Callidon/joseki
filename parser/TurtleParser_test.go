// Copyright (c) 2016 Thomas Minier. All rights reserved.
// Use of this source code is governed by a MIT License
// license that can be found in the LICENSE file.

package parser

import (
	"github.com/Callidon/joseki/rdf"
	"strings"
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
		rdf.NewTriple(rdf.NewURI("http://www.w3.org/2001/sw/RDFCore/ntriples"),
			rdf.NewURI("http://www.w3.org/1999/02/22-rdf-syntax-ns#type"),
			rdf.NewURI("http://xmlns.com/foaf/0.1/Document")),
		rdf.NewTriple(rdf.NewURI("http://www.w3.org/2001/sw/RDFCore/ntriples"),
			rdf.NewURI("http://www.w3.org/1999/02/22-rdf-syntax-ns#type"),
			rdf.NewURI("http://xmlns.com/foaf/0.1/Document")),
		rdf.NewTriple(rdf.NewURI("http://www.w3.org/2001/sw/RDFCore/ntriples"),
			rdf.NewURI("http://purl.org/dc/terms/title"),
			rdf.NewLangLiteral("N-Triples", "en")),
		rdf.NewTriple(rdf.NewURI("http://www.w3.org/2001/sw/RDFCore/ntriples"),
			rdf.NewURI("http://purl.org/dc/terms/title"),
			rdf.NewTypedLiteral("Turtle", "<http://www.w3.org/2001/XMLSchema#string>")),
		rdf.NewTriple(rdf.NewURI("http://www.w3.org/2001/sw/RDFCore/ntriples"),
			rdf.NewURI("http://purl.org/dc/terms/title"),
			rdf.NewBlankNode("a")),
		rdf.NewTriple(rdf.NewURI("http://www.w3.org/2001/sw/RDFCore/ntriples"),
			rdf.NewURI("http://xmlns.com/foaf/0.1/maker"),
			rdf.NewVariable("v0")),
		rdf.NewTriple(rdf.NewVariable("v0"),
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

func TestIllegalTokenTurtleParser(t *testing.T) {
	inputs := []string{
		"@prefix incorrect_uri",
		"@prefix <http://example.org> :",
		"@prefix <http://example.org> : illegal_value",
		"illegal_token",
	}
	expectedMsg := []string{
		"Unexpected token : incorrect_uri at line : 1 row : 1",
		"Unexpected token : <http://example.org> at line : 1 row : 1",
		"Unexpected token : <http://example.org> at line : 1 row : 1",
		"Unexpected token when scanning 'illegal_token' at line : 1 row : 1",
	}
	cpt := 0

	for _, input := range inputs {
		token := <-scanTurtle(strings.NewReader(input))
		tokenErr := token.Interpret(nil, nil, nil).Error()
		if tokenErr != expectedMsg[cpt] {
			t.Error("expected illegal token", expectedMsg[cpt], "but instead got", tokenErr)
		}
		cpt++
	}
}
