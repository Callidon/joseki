// Copyright (c) 2016 Thomas Minier. All rights reserved.
// Use of this source code is governed by a MIT License
// license that can be found in the LICENSE file.

package tokens

import (
	"errors"
	"github.com/Callidon/joseki/rdf"
	"strings"
)

// TokenPrefixedURI represent a prefixed RDF URI
type TokenPrefixedURI struct {
	value string
	*tokenPosition
}

// NewTokenPrefixedURI creates a new TokenPrefixedURI
func NewTokenPrefixedURI(value string, line int, row int) *TokenPrefixedURI {
	return &TokenPrefixedURI{value, newTokenPosition(line, row)}
}

// Interpret evaluate the token & produce an action.
// In the case of a TokenPrefixedURI, it push a URI to the stack
func (t TokenPrefixedURI) Interpret(nodeStack *Stack, prefixes *map[string]string, out chan rdf.Triple) error {
	sepIndex := strings.Index(t.value, ":")
	prefixValue := string(t.value[0:sepIndex])
	prefixURI, inPrefixes := (*prefixes)[prefixValue]
	if !inPrefixes {
		return errors.New("unkown prefix " + prefixValue + " at " + t.position())
	}
	nodeStack.Push(rdf.NewURI(prefixURI + t.value[sepIndex+1:]))
	return nil
}

// TokenPrefix represent a prefix
type TokenPrefix struct {
	name  string
	value string
}

// NewTokenPrefix creates a new TokenPrefix
func NewTokenPrefix(name, value string) *TokenPrefix {
	return &TokenPrefix{name, value}
}

// Interpret evaluate the token & produce an action.
// In the case of a TokenPrefix, it register a new prefix
func (t TokenPrefix) Interpret(nodeStack *Stack, prefixes *map[string]string, out chan rdf.Triple) error {
	(*prefixes)[t.name] = t.value
	return nil
}
