// Copyright (c) 2016 Thomas Minier. All rights reserved.
// Use of this source code is governed by a MIT License
// license that can be found in the LICENSE file.

package parser

import (
	"github.com/Callidon/joseki/rdf"
	"testing"
)

func TestInterpretTokenPrefixedURI(t *testing.T) {
	token := newTokenPrefixedURI("example:subject", 1, 1)
	stack := newStack()
	prefixes := make(map[string]string)
	prefixes["example"] = "http://example.org/"
	expectedNode := rdf.NewURI("http://example.org/subject")

	// Test for correct interpretation of the token
	if err := token.Interpret(stack, &prefixes, nil); err != nil {
		t.Error("interpretation of a correct tokenPrefixedURI shouldn't produce the error :", err)
	}
	node, _ := stack.Pop().(rdf.Node)
	if test, err := expectedNode.Equals(node); !test || err != nil {
		t.Error(expectedNode, "produced by tokenEnd.Interpret should be equals to", node)
	}
	if stack.Len()+1 != 1 {
		t.Error("after an interpretation of the token, the stack's size should be exactly equals to 1'")
	}
}

func TestInterpretErrorsTokenPrefixedURI(t *testing.T) {
	token := newTokenPrefixedURI("example:subject", 1, 1)
	stack := newStack()
	prefixes := make(map[string]string)

	// Test for incorrect interpretation of the token
	if err := token.Interpret(stack, &prefixes, nil); err == nil {
		t.Error("interpretation of a tokenPrefixedURI with an unknown prefix should produce an error")
	}

	if stack.Len() > 0 {
		t.Error("after an incorrect interpretation of the token, the stack should be empty'")
	}
}

func TestInterpretTokenPrefix(t *testing.T) {
	key, expectedValue := "example", "http://example.org/"
	token := newTokenPrefix(key, expectedValue)
	prefixes := make(map[string]string)

	// Test for correct interpretation of the token
	if err := token.Interpret(nil, &prefixes, nil); err != nil {
		t.Error("interpretation of a correct tokenPrefix shouldn't produce the error :", err)
	}
	if value, inPrefixes := prefixes[key]; !inPrefixes || value != expectedValue {
		t.Error("after tokenPrefix.Interpet, the prefix", key, "should exist and have", expectedValue, "as value")
	}
}
