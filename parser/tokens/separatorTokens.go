// Copyright (c) 2016 Thomas Minier. All rights reserved.
// Use of this source code is governed by a MIT License
// license that can be found in the LICENSE file.

package tokens

import (
	"errors"
	"github.com/Callidon/joseki/rdf"
)

// TokenEnd represent a RDF URI
type TokenEnd struct {
  *tokenPosition
}

// NewTokenEnd creates a new TokenEnd
func NewTokenEnd(line, row int) *TokenEnd {
	return &TokenEnd{newTokenPosition(line, row)}
}

// Interpret evaluate the token & produce an action.
// In the case of a TokenEnd, it form a new triples using the nodes in the stack
func (t TokenEnd) Interpret(nodeStack *Stack, prefixes *map[string]string, out chan rdf.Triple) error {
  if nodeStack.Len() > 3 {
    return errors.New("encountered a malformed triple pattern at " + t.position())
  }
	object, objIsNode := nodeStack.Pop().(rdf.Node)
  predicate, predIsNode := nodeStack.Pop().(rdf.Node)
  subject, subjIsNode := nodeStack.Pop().(rdf.Node)
  if !objIsNode || !predIsNode || !subjIsNode {
    return errors.New("expected a Node in stack but doesn't found it")
  }
  out <- rdf.NewTriple(subject, predicate, object)
	return nil
}
