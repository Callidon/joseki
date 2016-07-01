// Copyright (c) 2016 Thomas Minier. All rights reserved.
// Use of this source code is governed by a MIT License
// license that can be found in the LICENSE file.

// Package sparql provides support for requesting RDF Graphs using SPARQL query language
package sparql

import "github.com/Callidon/joseki/rdf"

// SelectQuery is a SPARQL SELECT query.
//
// The following example shows how to build a simple SELECT query :
//  graph := graph.NewHDTGraph().LoadFromFile("datas.nt", "nt")
//  triples := []rdf.Triple {
//    // let's initialize some triples here ...
//  }
//  query := NewSelectQuery("?s", "?p")
//  query.From(graph)
//  query.Where(triples...)
//  for bindings := range query.Execute() {
//		fmt.Println(bindings)
//	}
//
type SelectQuery struct {
	variables []string
	*queryDescriptor
}

// NewSelectQuery creates a new SPARQL SELECT query.
func NewSelectQuery(variables ...string) *SelectQuery {
	return &SelectQuery{variables, newQueryDescriptor(nil, selectQuery)}
}

// Execute run the Select query.
// The group of bindings which answers the query are send through a channel.
func (q SelectQuery) Execute() chan rdf.BindingsGroup {
	// get the query execution plan & add the SELECT modifier
	root := newSelectNode(q.build(), q.variables...)
	// TODO : apply optimization heuristic to the plan
	return root.execute()
}

// AskQuery is a SPARQL ASK query.
//
// The following example shows how to build a simple ASK query :
//  graph := graph.NewHDTGraph().LoadFromFile("datas.nt", "nt")
//  triples := []rdf.Triple {
//    // let's initialize some triples here ...
//  }
//  query := NewAskQuery()
//  query.From(graph)
//  query.Where(triples...)
//  fmt.Println(query.Execute()) // will display "true" or "false"
//
type AskQuery struct {
	*queryDescriptor
}

// NewAskQuery creates a new SPARQL ASK query.
func NewAskQuery() *AskQuery {
	return &AskQuery{newQueryDescriptor(nil, askQuery)}
}

type DescribeQuery struct {
	variables []string
	*queryDescriptor
}

// NewDescribeQuery creates a new SPARQL DESCRIBE query.
func NewDescribeQuery(variables ...string) *DescribeQuery {
	return &DescribeQuery{variables, newQueryDescriptor(nil, describeQuery)}
}

type ConstructQuery struct {
	triples []rdf.Triple
	*queryDescriptor
}

// NewConstructQuery creates a new SPARQL CONSTRUCT query.
func NewConstructQuery(triples ...rdf.Triple) *ConstructQuery {
	return &ConstructQuery{triples, newQueryDescriptor(nil, constructQuery)}
}
