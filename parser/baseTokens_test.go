// Copyright (c) 2016 Thomas Minier. All rights reserved.
// Use of this source code is governed by a MIT License
// license that can be found in the LICENSE file.

package parser

import (
	"github.com/Callidon/joseki/rdf"
	"testing"
)

func TestInterpretTokenURI(t *testing.T) {
	token := newTokenURI("http://example.org/subject")
	stack := newStack()
	expectedNode := rdf.NewURI("http://example.org/subject")

	// Test for correct interpretation of the token
	if err := token.Interpret(stack, nil, nil); err != nil {
		t.Error("interpretation of a correct tokenURI shouldn't produce the error :", err)
	}
	node, _ := stack.Pop().(rdf.Node)
	if test, err := expectedNode.Equals(node); !test || err != nil {
		t.Error(expectedNode, "produced by tokenURI.Interpret should be equals to", node)
	}
	if stack.Len()+1 != 1 {
		t.Error("after an interpretation of the token, the stack's size should be exactly equals to 1'")
	}
}

func TestInterpretTokenLiteral(t *testing.T) {
	token := newTokenLiteral("Harry Potter")
	stack := newStack()
	expectedNode := rdf.NewLiteral("Harry Potter")

	// Test for correct interpretation of the token
	if err := token.Interpret(stack, nil, nil); err != nil {
		t.Error("interpretation of a correct tokenLiteral shouldn't produce the error :", err)
	}
	node, _ := stack.Pop().(rdf.Node)
	if test, err := expectedNode.Equals(node); !test || err != nil {
		t.Error(expectedNode, "produced by tokenLiteral.Interpret should be equals to", node)
	}
	if stack.Len()+1 != 1 {
		t.Error("after an interpretation of the token, the stack's size should be exactly equals to 1'")
	}
}

func TestInterpretTokenType(t *testing.T) {
	token := newTokenType("http://www.w3.org/2001/XMLSchema#string", 1, 1)
	stack := newStack()
	expectedNode := rdf.NewTypedLiteral("Harry Potter", "http://www.w3.org/2001/XMLSchema#string")

	// Test for correct interpretation of the token
	stack.Push(rdf.NewLiteral("Harry Potter"))
	if err := token.Interpret(stack, nil, nil); err != nil {
		t.Error("interpretation of a correct tokenType shouldn't produce the error :", err)
	}
	node, _ := stack.Pop().(rdf.Node)
	if test, err := expectedNode.Equals(node); !test || err != nil {
		t.Error(expectedNode, "produced by tokenType.Interpret should be equals to", node)
	}
	if stack.Len()+1 != 1 {
		t.Error("after an interpretation of the token, the stack's size should be exactly equals to 1'")
	}
}

func TestInterpretErrorsTokenType(t *testing.T) {
	token := newTokenType("http://www.w3.org/2001/XMLSchema#string", 1, 1)
	stack := newStack()

	// Test for incorrect interpretation of the token
	if err := token.Interpret(stack, nil, nil); err == nil {
		t.Error("interpretation of a tokenType when they aren't anough nodes in the stack should produce an error")
	}

	stack.Push(rdf.NewURI("http://example.org/subject"))
	if err := token.Interpret(stack, nil, nil); err == nil {
		t.Error("interpretation of a tokenType when the top of the stack is a non-Literal node should produce an error")
	}
}

func TestInterpretTokenLang(t *testing.T) {
	token := newTokenType("en", 1, 1)
	stack := newStack()
	expectedNode := rdf.NewTypedLiteral("Harry Potter", "en")

	// Test for correct interpretation of the token
	stack.Push(rdf.NewLiteral("Harry Potter"))
	if err := token.Interpret(stack, nil, nil); err != nil {
		t.Error("interpretation of a correct tokenLang shouldn't produce the error :", err)
	}
	node, _ := stack.Pop().(rdf.Node)
	if test, err := expectedNode.Equals(node); !test || err != nil {
		t.Error(expectedNode, "produced by tokenLang.Interpret should be equals to", node)
	}
	if stack.Len()+1 != 1 {
		t.Error("after an interpretation of the token, the stack's size should be exactly equals to 1'")
	}
}

func TestInterpretErrorsTokenLang(t *testing.T) {
	token := newTokenLang("en", 1, 1)
	stack := newStack()

	// Test for incorrect interpretation of the token
	if err := token.Interpret(stack, nil, nil); err == nil {
		t.Error("interpretation of a tokenLang when they aren't anough nodes in the stack should produce an error")
	}

	stack.Push(rdf.NewURI("http://example.org/subject"))
	if err := token.Interpret(stack, nil, nil); err == nil {
		t.Error("interpretation of a tokenLang when the top of the stack is a non-Literal node should produce an error")
	}
}

func TestInterpretTokenBlankNode(t *testing.T) {
	token := newTokenBlankNode("v")
	stack := newStack()
	expectedNode := rdf.NewBlankNode("v")

	// Test for correct interpretation of the token
	if err := token.Interpret(stack, nil, nil); err != nil {
		t.Error("interpretation of a correct tokenBlankNode shouldn't produce the error :", err)
	}
	node, _ := stack.Pop().(rdf.Node)
	if test, err := expectedNode.Equals(node); !test || err != nil {
		t.Error(expectedNode, "produced by tokenBlankNode.Interpret should be equals to", node)
	}
	if stack.Len()+1 != 1 {
		t.Error("after an interpretation of the token, the stack's size should be exactly equals to 1'")
	}
}
