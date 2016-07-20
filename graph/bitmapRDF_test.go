// Copyright (c) 2016 Thomas Minier. All rights reserved.
// Use of this source code is governed by a MIT License
// license that can be found in the LICENSE file.

package graph

import (
	"github.com/Callidon/joseki/rdf"
	"testing"
)

// Test the Equals operator for bitmapTriple struct
func TestBitmapTripleEquals(t *testing.T) {
	tripleA := newBitmapTriple(1, 2, 3)
	tripleB := newBitmapTriple(2, 1, 3)
	tripleC := newBitmapTriple(-1, 2, 3)

	if test := tripleA.Equals(tripleA); !test {
		t.Error("a triple should be equals to itself")
	}
	if test := tripleA.Equals(tripleB); test {
		t.Error(tripleA, "cannot be equals to", tripleB)
	}
	if test := tripleA.Equals(tripleC); !test {
		t.Error(tripleA, "should be equals to", tripleC)
	}
}

// Test the Triple method for bitmapTriple struct
func TestBitmapTripleCast(t *testing.T) {
	subject := rdf.NewURI("http://example.org#subject")
	predicate := rdf.NewURI("http://example.org#predicate")
	object := rdf.NewURI("http://example.org#object")
	refTriple := rdf.NewTriple(subject, predicate, object)

	bimap := newBimap()
	bimap.push(1, subject)
	bimap.push(2, predicate)
	bimap.push(3, object)

	tripleA := newBitmapTriple(1, 2, 3)
	tripleB := newBitmapTriple(4, 2, 3)
	tripleC := newBitmapTriple(1, 4, 3)
	tripleD := newBitmapTriple(1, 2, 4)

	if _, err := tripleA.Triple(bimap); err != nil {
		t.Error(tripleA, "shouldn't throw an error when cast to rdf.Triple")
	} else {
		triple, _ := tripleA.Triple(bimap)
		if test, err := triple.Equals(refTriple); !test && err != nil {
			t.Error(triple, "should be equals to", refTriple)
		}
	}

	if _, err := tripleB.Triple(bimap); err == nil {
		t.Error(tripleB, "should throw an error when casted to rdf.triple")
	}

	if _, err := tripleC.Triple(bimap); err == nil {
		t.Error(tripleC, "should throw an error when casted to rdf.triple")
	}

	if _, err := tripleD.Triple(bimap); err == nil {
		t.Error(tripleD, "should throw an error when casted to rdf.triple")
	}
}
