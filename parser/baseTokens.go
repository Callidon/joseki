// Copyright (c) 2016 Thomas Minier. All rights reserved.
// Use of this source code is governed by a MIT License
// license that can be found in the LICENSE file.

package parser

import (
	"errors"
	"github.com/Callidon/joseki/rdf"
)

// tokenURI represent a RDF URI
type tokenURI struct {
	value string
}

// NewTokenURI creates a new tokenURI
func NewTokenURI(value string) *tokenURI {
	return &tokenURI{value}
}

// Interpret evaluate the token & produce an action.
// In the case of a tokenURI, it push a URI on top of the stack
func (t tokenURI) Interpret(nodeStack *stack, prefixes *map[string]string, out chan rdf.Triple) error {
	nodeStack.Push(rdf.NewURI(t.value))
	return nil
}

// tokenLiteral represent a RDF Literal
type tokenLiteral struct {
	value string
}

// NewTokenLiteral creates a new tokenLiteral
func NewTokenLiteral(value string) *tokenLiteral {
	return &tokenLiteral{value}
}

// Interpret evaluate the token & produce an action.
// In the case of a tokenLiteral, it push a Literal on top of the stack
func (t tokenLiteral) Interpret(nodeStack *stack, prefixes *map[string]string, out chan rdf.Triple) error {
	nodeStack.Push(rdf.NewLiteral(t.value))
	return nil
}

// tokenType represent a type for a RDF Literal
type tokenType struct {
	value string
	*tokenPosition
}

// NewTokenType creates a new tokenType.
// Since this token can produce an error, its position is needed for a better error handling
func NewTokenType(value string, line int, row int) *tokenType {
	return &tokenType{value, newTokenPosition(line, row)}
}

// Interpret evaluate the token & produce an action.
// In the case of a tokenType, it push a typed Literal on top of the stack
func (t tokenType) Interpret(nodeStack *stack, prefixes *map[string]string, out chan rdf.Triple) error {
	literal, isLiteral := nodeStack.Pop().(rdf.Literal)
	if !isLiteral {
		return errors.New("A XML type can only be associated with a RDF Literal, at " + t.position())
	}
	nodeStack.Push(rdf.NewTypedLiteral(literal.Value, t.value))
	return nil
}

// tokenLang represent a localization information about a RDF Literal
type tokenLang struct {
	value string
	*tokenPosition
}

// NewTokenLang creates a new tokenLang.
// Since this token can produce an error, its position is needed for a better error handling
func NewTokenLang(value string, line int, row int) *tokenLang {
	return &tokenLang{value, newTokenPosition(line, row)}
}

// Interpret evaluate the token & produce an action.
// In the case of a tokenLang, it push a Literal with its associated language on top of the stack
func (t tokenLang) Interpret(nodeStack *stack, prefixes *map[string]string, out chan rdf.Triple) error {
	literal, isLiteral := nodeStack.Pop().(rdf.Literal)
	if !isLiteral {
		return errors.New("A localization information can only be associated with a RDF Literal, at " + t.position())
	}
	nodeStack.Push(rdf.NewLangLiteral(literal.Value, t.value))
	return nil
}

// tokenBlankNode represent a RDF Blank Node
type tokenBlankNode struct {
	value string
}

// NewTokenBlankNode creates a new tokenBlankNode
func NewTokenBlankNode(value string) *tokenBlankNode {
	return &tokenBlankNode{value}
}

// Interpret evaluate the token & produce an action.
// In the case of a tokenBlankNode, it push a Blank Node on top of the stack
func (t tokenBlankNode) Interpret(nodeStack *stack, prefixes *map[string]string, out chan rdf.Triple) error {
	nodeStack.Push(rdf.NewBlankNode(t.value))
	return nil
}
