// Copyright (c) 2016 Thomas Minier. All rights reserved.
// Use of this source code is governed by a MIT License
// license that can be found in the LICENSE file.

// Package joseki is a pure Go library for working with RDF, a powerful framework for representing informations as graphs.
// For more informations about RDF itself, please see https://www.w3.org/TR/rdf11-concepts
//
// Features
//
// Joseki provides the following features to work with RDF :
//
// * Structures to represent and manipulate the RDF model (URIs, Literals, Blank Nodes, Triples, etc).
//
// * RDF Graphs to store data, with several implentations provided.
//
// * A Low level API to query data stored in graphs.
//
// * A High level API to query data using the SPARQL 1.1 query language.
//
// * Query processing using modern techniques such as join ordering or optimized query execution plans.
//
// * Load RDF data stored in files in various formats (N-Triples, Turtle, etc) into any graph.
//
// * Serialize a RDF Graph into various formats.
//
// Getting Started
//
// This package aims to work with RDF graphs, which are composed of RDF Triple {Subject Object Predicate}.
// Using joseki, you can represent an RDF Triple as followed :
//
//  import (
//      "github.com/Callidon/joseki/rdf"
//      "fmt"
//  )
//  subject := rdf.NewURI("http://example.org/book/book1")
//  predicate := rdf.NewURI("http://purl.org/dc/terms/title")
//  object := rdf.NewLiteral("Harry Potter and the Order of the Phoenix")
//  triple := rdf.NewTriple(subject, predicate, object)
//  fmt.Println(triple)
//  // Output : {<http://example.org/book/book1> <http://purl.org/dc/terms/title> "Harry Potter and the Order of the Phoenix"}
//
// You can also store your RDF Triples in a RDF Graph, using various type of graphs.
// Here, we use a Tree Graph to store our triple :
//
//  import (
//      "github.com/Callidon/joseki/rdf"
//      "github.com/Callidon/joseki/graph"
//  )
//  subject := rdf.NewURI("http://example.org/book/book1")
//  predicate := rdf.NewURI("http://purl.org/dc/terms/title")
//  object := rdf.NewLiteral("Harry Potter and the Order of the Phoenix")
//  graph := graph.NewTreeGraph()
//  graph.Add(rdf.NewTriple(subject, predicate, object))
//
// You can also query any triple from a RDF Graph, using a low level API or a SPARQL query.
//
//  import (
//      "github.com/Callidon/joseki/rdf"
//      "github.com/Callidon/joseki/graph"
//      "fmt"
//  )
//  graph := graph.NewTreeGraph()
//  // Datas stored in a file can be easily loaded into a graph
//  graph.LoadFromFile("datas/awesome-books.ttl", "turtle")
//  // Let's fetch the titles of all the books in our graph !
//  subject := rdf.NewVariable("title")
//  predicate := rdf.NewURI("http://www.w3.org/1999/02/22-rdf-syntax-ns#type")
//  object := rdf.NewURI("https://schema.org/Book")
//  for bindings := range graph.Filter(subject, predicate, object) {
//      fmt.Println(bindings)
//  }
//
// For more informations about specific features, see the documentation of each subpackage.
//
// Author : Thomas Minier
package joseki

import "fmt"

func ComingSoon() {
	fmt.Println("More to come ;)")
}
