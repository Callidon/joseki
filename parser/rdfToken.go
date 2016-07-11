// Copyright (c) 2016 Thomas Minier. All rights reserved.
// Use of this source code is governed by a MIT License
// license that can be found in the LICENSE file.

// Package parser provides utilities to work with RDF based languages
package parser

import (
	"errors"
	"github.com/Callidon/joseki/rdf"
	"strconv"
)

// rdfToken represent a token in a RDF based language
//
// It follows the Interpretor pattern (https://en.wikipedia.org/wiki/Interpreter_pattern)
// and can be used to extract triple pattersn& prefiex when reading a file
type rdfToken interface {
	// Interpret evaluate the token & produce an action
	Interpret(nodeStack *stack, prefixes *map[string]string, out chan rdf.Triple) error
}

// tokenPosition represent the position of a token
type tokenPosition struct {
	lineNumber int
	rowNumber  int
}

// newTokenPosition creates a new tokenPosition
func newTokenPosition(line, row int) *tokenPosition {
	return &tokenPosition{line, row}
}

// position returns a string representation of the token's position
func (t tokenPosition) position() string {
	return "line : " + strconv.Itoa(t.lineNumber) + " row : " + strconv.Itoa(t.rowNumber)
}

// tokenIllegal is an illegal token in the RDF syntax
type tokenIllegal struct {
	errMsg string
	*tokenPosition
}

// NewTokenIllegal crates a new tokenIllegal
func NewTokenIllegal(err string, line int, row int) *tokenIllegal {
	return &tokenIllegal{err, newTokenPosition(line, row)}
}

// Interpret evaluate the token & produce an action. In the case of a tokenIllegal, it causes a panic.
func (t tokenIllegal) Interpret(nodeStack *stack, prefixes *map[string]string, out chan rdf.Triple) error {
	return errors.New(t.errMsg + " at " + t.position())
}
