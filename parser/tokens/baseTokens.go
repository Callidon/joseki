// Copyright (c) 2016 Thomas Minier. All rights reserved.
// Use of this source code is governed by a MIT License
// license that can be found in the LICENSE file.

package tokens

import (
	"errors"
	"github.com/Callidon/joseki/rdf"
)

// TokenURI represent a RDF URI
type TokenURI struct {
	value string
}

// NewTokenURI creates a new TokenURI
func NewTokenURI(value string) *TokenURI {
	return &TokenURI{value}
}

// Interpret evaluate the token & produce an action.
// In the case of a TokenURI, it push a URI on top of the stack
func (t TokenURI) Interpret(nodeStack *Stack, prefixes *map[string]string, out chan rdf.Triple) error {
	nodeStack.Push(rdf.NewURI(t.value))
	return nil
}

// TokenLiteral represent a RDF Literal
type TokenLiteral struct {
	value string
}

// NewTokenLiteral creates a new TokenLiteral
func NewTokenLiteral(value string) *TokenLiteral {
	return &TokenLiteral{value}
}

// Interpret evaluate the token & produce an action.
// In the case of a TokenLiteral, it push a Literal on top of the stack
func (t TokenLiteral) Interpret(nodeStack *Stack, prefixes *map[string]string, out chan rdf.Triple) error {
	nodeStack.Push(rdf.NewLiteral(t.value))
	return nil
}

// TokenType represent a type for a RDF Literal
type TokenType struct {
	value string
	*tokenPosition
}

// NewTokenType creates a new TokenType.
// Since this token can produce an error, its position is needed for a better error handling
func NewTokenType(value string, line int, row int) *TokenType {
	return &TokenType{value, newTokenPosition(line, row)}
}

// Interpret evaluate the token & produce an action.
// In the case of a TokenType, it push a typed Literal on top of the stack
func (t TokenType) Interpret(nodeStack *Stack, prefixes *map[string]string, out chan rdf.Triple) error {
	literal, isLiteral := nodeStack.Pop().(rdf.Literal)
	if !isLiteral {
		return errors.New("A XML type can only be associated with a RDF Literal - at " + t.position())
	}
	nodeStack.Push(rdf.NewTypedLiteral(literal.Value, t.value))
	return nil
}

// TokenLang represent a localization information about a RDF Literal
type TokenLang struct {
	value string
	*tokenPosition
}

// NewTokenLang creates a new TokenLang.
// Since this token can produce an error, its position is needed for a better error handling
func NewTokenLang(value string, line int, row int) *TokenLang {
	return &TokenLang{value, newTokenPosition(line, row)}
}

// Interpret evaluate the token & produce an action.
// In the case of a TokenLang, it push a Literal with its associated language on top of the stack
func (t TokenLang) Interpret(nodeStack *Stack, prefixes *map[string]string, out chan rdf.Triple) error {
	literal, isLiteral := nodeStack.Pop().(rdf.Literal)
	if !isLiteral {
		return errors.New("A localization information can only be associated with a RDF Literal - at " + t.position())
	}
	nodeStack.Push(rdf.NewLangLiteral(literal.Value, t.value))
	return nil
}

// TokenBlankNode represent a RDF Blank Node
type TokenBlankNode struct {
	value string
}

// NewTokenBlankNode creates a new TokenBlankNode
func NewTokenBlankNode(value string) *TokenBlankNode {
	return &TokenBlankNode{value}
}

// Interpret evaluate the token & produce an action.
// In the case of a TokenBlankNode, it push a Blank Node on top of the stack
func (t TokenBlankNode) Interpret(nodeStack *Stack, prefixes *map[string]string, out chan rdf.Triple) error {
	nodeStack.Push(rdf.NewBlankNode(t.value))
	return nil
}
