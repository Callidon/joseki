// Copyright (c) 2016 Thomas Minier. All rights reserved.
// Use of this source code is governed by a MIT License
// license that can be found in the LICENSE file.

package parser

import (
	"errors"
	"github.com/Callidon/joseki/rdf"
	"math/rand"
	"strconv"
)

// tokenEnd represent a RDF URI
type tokenEnd struct {
	*tokenPosition
}

// newTokenEnd creates a new tokenEnd
func newTokenEnd(line, row int) *tokenEnd {
	return &tokenEnd{newTokenPosition(line, row)}
}

// Interpret evaluate the token & produce an action.
// In the case of a tokenEnd, it form a new triple using the nodes in the stack
func (t tokenEnd) Interpret(nodeStack *stack, prefixes *map[string]string, out chan rdf.Triple) error {
	if nodeStack.Len() < 3 {
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

// tokenSep represent a Turtle separator
type tokenSep struct {
	value string
	*tokenPosition
}

// newTokenSep creates a new tokenSep
func newTokenSep(value string, line int, row int) *tokenSep {
	return &tokenSep{value, newTokenPosition(line, row)}
}

// Interpret evaluate the token & produce an action.
// In the case of a tokenSep, it form a new triple based on the separator, using the nodes in the stack
func (t tokenSep) Interpret(nodeStack *stack, prefixes *map[string]string, out chan rdf.Triple) error {
	// case of a object separator
	if t.value == "[" {
		if nodeStack.Len() < 2 {
			return errors.New("encountered a malformed triple pattern at " + t.position())
		}
		predicate, predIsNode := nodeStack.Pop().(rdf.Node)
		subject, subjIsNode := nodeStack.Pop().(rdf.Node)
		object := rdf.NewBlankNode("v" + strconv.Itoa(rand.Int()))
		if !predIsNode || !subjIsNode {
			return errors.New("expected a Node in stack but doesn't found it")
		}
		out <- rdf.NewTriple(subject, predicate, object)
		nodeStack.Push(object)
	} else {
		if nodeStack.Len() < 3 {
			return errors.New("encountered a malformed triple pattern at " + t.position())
		}
		object, objIsNode := nodeStack.Pop().(rdf.Node)
		predicate, predIsNode := nodeStack.Pop().(rdf.Node)
		subject, subjIsNode := nodeStack.Pop().(rdf.Node)
		if !objIsNode || !predIsNode || !subjIsNode {
			return errors.New("expected a Node in stack but doesn't found it")
		}
		out <- rdf.NewTriple(subject, predicate, object)

		switch t.value {
		case ";":
			// push back the subject into the stack
			nodeStack.Push(subject)
		case ",":
			// push back the subject & the predicate into the stack
			nodeStack.Push(subject)
			nodeStack.Push(predicate)
		default:
			return errors.New("Unexpected separator token " + t.value + " - at " + t.position())
		}
	}
	return nil
}
