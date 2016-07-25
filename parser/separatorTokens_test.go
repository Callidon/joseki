// Copyright (c) 2016 Thomas Minier. All rights reserved.
// Use of this source code is governed by a MIT License
// license that can be found in the LICENSE file.

package parser

import (
	"github.com/Callidon/joseki/rdf"
	"testing"
)

func TestInterpretTokenEnd(t *testing.T) {
	token := newTokenEnd(1, 1)
	stack := newStack()
	out := make(chan rdf.Triple, 1)
	subject := rdf.NewURI("http://example.org/subject")
	predicate := rdf.NewURI("http://example.org/predicate")
	object := rdf.NewURI("http://example.org/object")
	expectedTriple := rdf.NewTriple(subject, predicate, object)

	// Test for correct interpretation of the token
	stack.Push(subject)
	stack.Push(predicate)
	stack.Push(object)

	if err := token.Interpret(stack, nil, out); err != nil {
		t.Error("interpretation of a correct TokenEnd shouldn't produce the error :", err)
	}
	triple := <-out
	if test, err := expectedTriple.Equals(triple); !test || err != nil {
		t.Error(triple, "produced by tokenEnd.Interpret should be equals to", expectedTriple)
	}
	if stack.Len() > 0 {
		t.Error("after an interpretation of the token, this stack should be empty")
	}
}

func TestInterpretErrorsTokenEnd(t *testing.T) {
	token := newTokenEnd(1, 1)
	stack := newStack()
	out := make(chan rdf.Triple, 1)

	// test with not enough nodes in the stack
	stack.Push(rdf.NewURI("http://example.org/subject"))
	stack.Push(rdf.NewURI("http://example.org/predicate"))

	if err := token.Interpret(stack, nil, out); err == nil {
		t.Error("interpretation of a TokenEnd with not enough tokens in the stack should produce an error")
	}

	// test with an incorrect element in the stack
	stack.Push("incorrect node")
	if err := token.Interpret(stack, nil, out); err == nil {
		t.Error("interpretation of a TokenEnd with an incorrect element in the stack should produce an error")
	}
}

func TestInterpretTokenSep(t *testing.T) {
	var top rdf.Node
	var triple rdf.Triple
	token := newTokenSep(";", 1, 1)
	stack := newStack()
	out := make(chan rdf.Triple, 1)
	subject := rdf.NewURI("http://example.org/subject")
	predicate := rdf.NewURI("http://example.org/predicate")
	object := rdf.NewURI("http://example.org/object")
	expectedTriple := rdf.NewTriple(subject, predicate, object)

	// Test for a "," separator
	stack.Push(subject)
	stack.Push(predicate)
	stack.Push(object)

	if err := token.Interpret(stack, nil, out); err != nil {
		t.Error("interpretation of a correct TokenEnd shouldn't produce the error :", err)
	}

	triple = <-out
	if test, err := expectedTriple.Equals(triple); !test || err != nil {
		t.Error(triple, "produced by tokenSep.Interpret should be equals to", expectedTriple)
	}

	top, _ = stack.Pop().(rdf.Node)
	if test, err := top.Equals(subject); !test || err != nil {
		t.Error("tokenSep.Interpret with the ';' separator value should have put back the subject on top of the stack")
	}
	if stack.Len()+1 != 1 {
		t.Error("after an interpretation of the token, the stack's size should be exactly equals to 1'")
	}

	// Test for a "," separator
	token.value = ","
	stack.Push(subject)
	stack.Push(predicate)
	stack.Push(object)

	if err := token.Interpret(stack, nil, out); err != nil {
		t.Error("interpretation of a correct TokenEnd shouldn't produce the error :", err)
	}

	triple = <-out
	if test, err := expectedTriple.Equals(triple); !test || err != nil {
		t.Error(triple, "produced by tokenSep.Interpret should be equals to", expectedTriple)
	}

	top, _ = stack.Pop().(rdf.Node)
	if test, err := top.Equals(predicate); !test || err != nil {
		t.Error("tokenSep.Interpret with the ',' separator value should have put back the predicate on top of the stack")
	}
	top, _ = stack.Pop().(rdf.Node)
	if test, err := top.Equals(subject); !test || err != nil {
		t.Error("tokenSep.Interpret with the ',' separator value should have put back the subject in the 2nd position of the stack")
	}

	if stack.Len()+2 != 2 {
		t.Error("after an interpretation of the token, the stack's size should be exactly equals to 2'")
	}
}

func TestInterpretErrorsTokenSep(t *testing.T) {
	token := newTokenSep("l", 1, 1)
	stack := newStack()
	out := make(chan rdf.Triple, 1)

	// test with incorrect separator value
	if err := token.Interpret(stack, nil, out); err == nil {
		t.Error("interpretation of a TokenSep with an incorrect separator value should produce an error")
	}

	// test with the separator "[" and not enough nodes in the stack
	token.value = "["
	stack.Push(rdf.NewURI("http://example.org/subject"))

	if err := token.Interpret(stack, nil, out); err == nil {
		t.Error("interpretation of a TokenSep with the separator '[' & not enough tokens in the stack should produce an error")
	}

	// test with the separator "[" and not enough nodes in the stack
	stack.Push("incorrect node")
	if err := token.Interpret(stack, nil, out); err == nil {
		t.Error("interpretation of a TokenSep with the separator '[' & an incorrect element in the stack should produce an error")
	}

	// test with the separator ";" and not enough nodes in the stack
	token.value = ","
	stack.Push(rdf.NewURI("http://example.org/subject"))
	stack.Push(rdf.NewURI("http://example.org/predicate"))

	if err := token.Interpret(stack, nil, out); err == nil {
		t.Error("interpretation of a TokenSep not enough tokens in the stack should produce an error")
	}

	// test with the separator ";" and not enough nodes in the stack
	stack.Push("incorrect node")
	if err := token.Interpret(stack, nil, out); err == nil {
		t.Error("interpretation of a TokenSep with an incorrect element in the stack should produce an error")
	}
}
