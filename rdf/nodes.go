// Copyright (c) 2016 Thomas Minier. All rights reserved.
// Use of this source code is governed by a MIT License
// license that can be found in the LICENSE file.

// Package rdf provides primitives to work with RDF
package rdf

import "errors"

// Node represent a generic node in a RDF Grapg
//
// RDF Graph reference : https://www.w3.org/TR/2004/REC-rdf-concepts-20040210/#section-data-model
type Node interface {
	Equals(n Node) (bool, error)
	Equivalent(n Node) (bool, error)
	String() string
}

// URI represent a URI node in a RDF Graph
//
// RDF URI reference : https://www.w3.org/TR/2004/REC-rdf-concepts-20040210/#section-Graph-URIref
type URI struct {
	Value string
}

// Literal represent a Literal node in a RDF Graph.
//
// RDF Literal reference : https://www.w3.org/TR/2004/REC-rdf-concepts-20040210/#section-Graph-Literal
type Literal struct {
	Value string
}

// BlankNode represent a Blank Node in a RDF Graph.
//
// RDF Blank Node reference : https://www.w3.org/TR/2004/REC-rdf-concepts-20040210/#section-blank-nodes
type BlankNode struct {
	Variable string
}

// Equals is a function that compare two URIs and return True if they are equals, False otherwise.
func (u URI) Equals(n Node) (bool, error) {
	other, ok := n.(URI)
	if ok {
		return u.Value == other.Value, nil
	}
	return false, errors.New("Error : mismatch type, can only compare two URIs")
}

// Equivalent is a function that determine if a URI is equivalent to another RDF Node.
// Two URIs are equivalents if they are equals, and a URI is always equivalent to a Blank Node.
// Otherwise, the result is always False.
func (u URI) Equivalent(n Node) (bool, error) {
	equality, err := u.Equals(n)
	if err != nil {
		_, ok := n.(BlankNode)
		if ok {
			return true, nil
		}
		return false, errors.New("Error : can only compare a URI with another URI or a Blank Node")
	}
	return equality, nil
}

// Serialize a URI to string and return it.
func (u URI) String() string {
	return "<" + u.Value + ">"
}

// NewURI creates a new URI.
func NewURI(value string) URI {
	return URI{value}
}

// Equals is a function that compare two Literals and return True if they are equals, False otherwise.
func (l Literal) Equals(n Node) (bool, error) {
	other, ok := n.(Literal)
	if ok {
		return l.Value == other.Value, nil
	}
	return false, errors.New("Error : mismatch type, can only compare two Literals")
}

// Equivalent is a function that determine if a Literals is equivalent to another RDF Node.
// Two Literals are equivalents if they are equals, and a Literal is always equivalent to a Blank Node.
// Otherwise, the result is always False.
func (l Literal) Equivalent(n Node) (bool, error) {
	equality, err := l.Equals(n)
	if err != nil {
		_, ok := n.(BlankNode)
		if ok {
			return true, nil
		}
		return false, errors.New("Error : can only compare a Literal with another Literal or a Blank Node")
	}
	return equality, nil
}

// Serialize a Literal to string and return it.
func (l Literal) String() string {
	return "\"" + l.Value + "\""
}

// NewLiteral creates a new Literal.
func NewLiteral(value string) Literal {
	return Literal{value}
}

// Equals is a function that compare two Blank Nodes and return True if they are equals, False otherwise.
func (b BlankNode) Equals(n Node) (bool, error) {
	other, ok := n.(BlankNode)
	if ok {
		return b.Variable == other.Variable, nil
	}
	return false, errors.New("Error : mismatch type, can only compare two Blank Nodes")
}

// Equivalent is a function that determine if a Blank Node is equivalent to another RDF Node.
// Since a Blank Node is always equivalent to any RDF Node, this function always return True.
func (b BlankNode) Equivalent(n Node) (bool, error) {
	return true, nil
}

// Serialize a Blank Node to string and return it.
func (b BlankNode) String() string {
	return "_:" + b.Variable
}

// NewBlankNode creates a new Literal.
func NewBlankNode(variable string) BlankNode {
	return BlankNode{variable}
}
