// Copyright (c) 2016 Thomas Minier. All rights reserved.
// Use of this source code is governed by a MIT License
// license that can be found in the LICENSE file.

// Package tokens provides utilities to work with RDF based languages
package tokens

import (
	"errors"
	"github.com/Callidon/joseki/rdf"
	"strconv"
)

// RDFToken represent a Token in a RDF based language
//
// It follows the Interpretor pattern (https://en.wikipedia.org/wiki/Interpreter_pattern)
// and can be used to extract triple pattersn& prefiex when reading a file
type RDFToken interface {
	// Interpret evaluate the token & produce an action
	Interpret(nodeStack *Stack, prefixes *map[string]string, out chan rdf.Triple) error
}

// tokenPosition represent the position of a Token
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

// TokenIllegal is an illegal Token in the RDF syntax
type TokenIllegal struct {
	errMsg string
	*tokenPosition
}

// NewTokenIllegal crates a new TokenIllegal
func NewTokenIllegal(err string, line int, row int) *TokenIllegal {
	return &TokenIllegal{err, newTokenPosition(line, row)}
}

// Interpret evaluate the token & produce an action. In the case of a TokenIllegal, it causes a panic.
func (t TokenIllegal) Interpret(nodeStack *Stack, prefixes *map[string]string, out chan rdf.Triple) error {
	return errors.New(t.errMsg + " - at " + t.position())
}
