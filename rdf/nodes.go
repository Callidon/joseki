// Copyright (c) 2016 Thomas Minier. All rights reserved.
// Use of this source code is governed by a MIT License
// license that can be found in the LICENSE file.

// Package rdf provides primitives to work with RDF
package rdf

import "errors"

// Node represents a generic node in a RDF Grapg
//
// RDF Graph reference : https://www.w3.org/TR/rdf11-concepts/#section-rdf-graph
type Node interface {
	Equals(n Node) (bool, error)
	String() string
}

// URI represents a URI node in a RDF Graph
//
// RDF URI reference : https://www.w3.org/TR/rdf11-concepts/#section-IRIs
type URI struct {
	Value string
}

// Literal represents a Literal node in a RDF Graph.
//
// RDF Literal reference : https://www.w3.org/TR/rdf11-concepts/#section-Graph-Literal
type Literal struct {
	Value string
	Type  string
	Lang  string
}

// BlankNode represents a Blank Node in a RDF Graph.
//
// RDF Blank Node reference : https://www.w3.org/TR/rdf11-concepts/#section-blank-nodes
type BlankNode struct {
	Value string
}

// Variable represents a SPARQL variable used when querying data in a RDF graph
type Variable struct {
	Value string
}

// Equals is a function that compare a URI with another RDF Node and return True if they are equals, False otherwise.
func (u URI) Equals(n Node) (bool, error) {
	other, ok := n.(URI)
	if ok {
		return u.Value == other.Value, nil
	} else if _, isVar := n.(Variable); isVar {
		return true, nil
	}
	return false, errors.New("Error : mismatch type, can only compare two URIs")
}

// Serialize a URI to string and return it.
func (u URI) String() string {
	return "<" + u.Value + ">"
}

// NewURI creates a new URI.
func NewURI(value string) URI {
	return URI{value}
}

// Equals is a function that compare a Literal with another RDF Node and return True if they are equals, False otherwise.
func (l Literal) Equals(n Node) (bool, error) {
	other, ok := n.(Literal)
	if ok {
		return (l.Value == other.Value) && (l.Type == other.Type) && (l.Lang == other.Lang), nil
	} else if _, isVar := n.(Variable); isVar {
		return true, nil
	}
	return false, errors.New("Error : mismatch type, can only compare two Literals")
}

// Serialize a Literal to string and return it.
func (l Literal) String() string {
	if l.Type != "" {
		return "\"" + l.Value + "\"^^" + l.Type
	} else if l.Lang != "" {
		return "\"" + l.Value + "\"@" + l.Lang
	}
	return "\"" + l.Value + "\""
}

// NewLiteral creates a new Literal.
func NewLiteral(value string) Literal {
	return Literal{value, "", ""}
}

// NewTypedLiteral returns a new Literal with a type.
func NewTypedLiteral(value, xmlType string) Literal {
	return Literal{value, xmlType, ""}
}

// NewLangLiteral returns a new Literal with a language.
func NewLangLiteral(value, lang string) Literal {
	return Literal{value, "", lang}
}

// Equals is a function that compare a Blank Node with another RDF Node and return True if they are equals, False otherwise.
func (b BlankNode) Equals(n Node) (bool, error) {
	other, ok := n.(BlankNode)
	if ok {
		return b.Value == other.Value, nil
	} else if _, isVar := n.(Variable); isVar {
		return true, nil
	}
	return false, errors.New("Error : mismatch type, can only compare two Blank Nodes")
}

// Serialize a Blank Node to string and return it.
func (b BlankNode) String() string {
	return "_:" + b.Value
}

// NewBlankNode creates a new Blank Node.
func NewBlankNode(variable string) BlankNode {
	return BlankNode{variable}
}

// Equals is a function that compare a Variable with another RDF Node and return True if they are equals, False otherwise.
// Two variables are equals if they have the same value, and a variable is always equals to any other RDF node.
func (v Variable) Equals(n Node) (bool, error) {
	other, ok := n.(Variable)
	if ok {
		return v.Value == other.Value, nil
	}
	return true, nil
}

// Serialize a Variable to string and return it.
func (v Variable) String() string {
	return "?" + v.Value
}

// NewVariable creates a new Variable.
func NewVariable(value string) Variable {
	return Variable{value}
}
