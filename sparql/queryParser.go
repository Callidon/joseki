// Copyright (c) 2016 Thomas Minier. All rights reserved.
// Use of this source code is governed by a MIT License
// license that can be found in the LICENSE file.

package sparql

import (
	"errors"
	"fmt"
	"github.com/Callidon/joseki/graph"
	"github.com/Callidon/joseki/rdf"
	"strconv"
)

// Query represent classic SPARQL 1.1 query.
//
// SPARQL 1.1 reference : https://www.w3.org/TR/sparql11-overview/
type Query struct {
	root sparqlNode
}

// Execute run the query against a graph and return the results as a channel of groups of binding
func (q Query) Execute() chan rdf.BindingsGroup {
	return q.root.execute()
}

// checkRoot assert that the root node of the query (ie it's type : SELECT, ASK, etc) has not yet been found
func checkRoot(root sparqlNode) {
	if root != nil {
		panic(errors.New("A query must contain only one of the following keywords : SELECT, DESCRIBE, ASK, CONSTRUCT"))
	}
}

// saveTriple add a new Triple Node to the BGP currently evaluated
func saveTriple(bgp []sparqlNode, graph graph.Graph, subject rdf.Node, predicate rdf.Node, object rdf.Node) []sparqlNode {
	return append(bgp, newTripleNode(rdf.NewTriple(subject, predicate, object), graph))
}

// LoadQuery creates a new query which will be executed agianst a given RDF Graph
func LoadQuery(query string, graph graph.Graph) *Query {
	var root sparqlNode
	var subject, object, predicate rdf.Node
	var literalValue string
	bgps := make(map[int][]sparqlNode)
	readingBGP := false
	nextBGPID := 0
	bnodeCpt := 0
	scanner := newSparqlScanner()
	//stack := newStack()

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

	// extract elements from the query
	// TODO : complete with all SPARQL keywords
	for token := range scanner.scan(query) {
		switch token.Type {
		case tokenSelect:
			checkRoot(root)
			root = &selectNode{}
		case tokenBGPBegin:
			readingBGP = true
			bgps[nextBGPID] = make([]sparqlNode, 0)
		case tokenBGPEnd:
			if !readingBGP {
				panic(errors.New("found '}' but no matching '{' has been found"))
			}
			readingBGP = false
			nextBGPID++
		case tokenEnd:
			// add a new triple to the BGP currently evaluated
			bgps[nextBGPID] = saveTriple(bgps[nextBGPID], graph, subject, predicate, object)
			subject, predicate, object = nil, nil, nil
		case tokenSep:
			switch token.Value {
			case ";":
				// save previous triple & keep subject for the next triple
				bgps[nextBGPID] = saveTriple(bgps[nextBGPID], graph, subject, predicate, object)
				predicate, object = nil, nil
			case ",":
				// save previous triple & keep subject and predicate ofr the next triple
				bgps[nextBGPID] = saveTriple(bgps[nextBGPID], graph, subject, predicate, object)
				object = nil
			case "[":
				// generate a new object, save triple and then use the new blank Node as the new subject
				object = rdf.NewBlankNode("v" + strconv.Itoa(bnodeCpt))
				bgps[nextBGPID] = saveTriple(bgps[nextBGPID], graph, subject, predicate, object)
				subject = object
				predicate, object = nil, nil
				bnodeCpt++
			default:
				panic(errors.New("Unexpected separator token " + token.Value))
			}
		case tokenURI:
			assignNode(rdf.NewURI(token.Value))
		case tokenPrefixedURI:
			/*sepIndex := strings.Index(token.Value, ":")
			prefixURI, knownPrefix := p.prefixes[token.Value[0:sepIndex]]
			if knownPrefix {
				assignNode(rdf.NewURI(prefixURI + token.Value[sepIndex+1:]))
			} else {
				panic(errors.New("Unknown prefix " + token.Value[0:sepIndex] + " found"))
			}*/
		case tokenVariable:
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
		case illegalToken:
			panic(token.Value)
		default:
			//panic(errors.New("Unsupported token type for " + token.Value))
		}
	}

	// analyse the BGPs found to determine the join ordering
	// TODO

	return &Query{root}
}
