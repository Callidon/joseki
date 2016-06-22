// Copyright (c) 2016 Thomas Minier. All rights reserved.
// Use of this source code is governed by a MIT License
// license that can be found in the LICENSE file.

package sparql

// SparqlParser is a parser for reading & decomposing a SPARQL request in a query execution plan.
//
// SPARQL 1.1 reference : https://www.w3.org/TR/sparql11-overview/
type SparqlParser struct {
    *SparqlScanner
}

// NewSparqlParser creates a new SparqlParser.
func NewSparqlParser() *SparqlParser {
    return &SparqlParser{NewSparqlScanner()}
}