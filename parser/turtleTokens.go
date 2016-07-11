// Copyright (c) 2016 Thomas Minier. All rights reserved.
// Use of this source code is governed by a MIT License
// license that can be found in the LICENSE file.

package parser

import (
	"errors"
	"github.com/Callidon/joseki/rdf"
	"strings"
)

// tokenPrefixedURI represent a prefixed RDF URI
type tokenPrefixedURI struct {
	value string
	*tokenPosition
}

// newTokenPrefixedURI creates a new tokenPrefixedURI
func newTokenPrefixedURI(value string, line int, row int) *tokenPrefixedURI {
	return &tokenPrefixedURI{value, newTokenPosition(line, row)}
}

// Interpret evaluate the token & produce an action.
// In the case of a tokenPrefixedURI, it push a URI to the stack
func (t tokenPrefixedURI) Interpret(nodeStack *stack, prefixes *map[string]string, out chan rdf.Triple) error {
	sepIndex := strings.Index(t.value, ":")
	prefixValue := string(t.value[0:sepIndex])
	prefixURI, inPrefixes := (*prefixes)[prefixValue]
	if !inPrefixes {
		return errors.New("unkown prefix " + prefixValue + " at " + t.position())
	}
	nodeStack.Push(rdf.NewURI(prefixURI + t.value[sepIndex+1:]))
	return nil
}

// tokenPrefix represent a prefix
type tokenPrefix struct {
	name  string
	value string
}

// newTokenPrefix creates a new tokenPrefix
func newTokenPrefix(name, value string) *tokenPrefix {
	return &tokenPrefix{name, value}
}

// Interpret evaluate the token & produce an action.
// In the case of a tokenPrefix, it register a new prefix
func (t tokenPrefix) Interpret(nodeStack *stack, prefixes *map[string]string, out chan rdf.Triple) error {
	(*prefixes)[t.name] = t.value
	return nil
}
