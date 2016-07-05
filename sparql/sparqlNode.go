// Copyright (c) 2016 Thomas Minier. All rights reserved.
// Use of this source code is governed by a MIT License
// license that can be found in the LICENSE file.

package sparql

import "github.com/Callidon/joseki/rdf"

// sparqlNode represents a node in a SPARQL query execution plan.
// Each implementation of this interface represents a type of operation executed during a SPARQL request.
//
// When all nodes of SPARQL query execution plan are executed in the correct order, a response to the corresponding request will be produced.
// Package sparql provides several implementations for this interface.
type sparqlNode interface {
	execute() <-chan rdf.BindingsGroup
	executeWith(group rdf.BindingsGroup) <-chan rdf.BindingsGroup
	bindingNames() []string
	Equals(other sparqlNode) bool
	String() string
}
