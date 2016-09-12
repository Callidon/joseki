// Copyright (c) 2016 Thomas Minier. All rights reserved.
// Use of this source code is governed by a MIT License
// license that can be found in the LICENSE file.

// Package parser provides parser to work with several RDF formats (N-Triples, Turtles, JSON-LD, ...)
package parser

import (
	"github.com/Callidon/joseki/rdf"
	"regexp"
)

const (
	// Max size for the buffer of this package
	bufferSize = 100
	// Regexp used to isolate triples and their elements
	wordRegexp = "'.*?'|\".*?\"|\\S+"
)

// Parser represent a generic interface for parsing every RDF format.
//
// Package parser provides several implementations for this interface.
type Parser interface {
	Read(filename string) chan rdf.Triple
	Prefixes() map[string]string
}

// lineCutter wraps up the regexp used isolate triples and their elements in the RDF standard
// It's main purpose is to ensure that the regexp is compiled only once, since it's a high cost operation.
type lineCutter struct {
	*regexp.Regexp
}

// newLineCutter creates a new lineCutter
func newLineCutter(reg string) *lineCutter {
	return &lineCutter{regexp.MustCompile(reg)}
}

// extractSegments parse a string and split the segments into a slice.
// A segment is a string quoted or separated from the other by whitespaces.
func (l lineCutter) extractSegments(line string) []string {
	return l.FindAllString(line, -1)
}

// Utility function for checking errors
func check(err error) {
	if err != nil {
		panic(err)
	}
}
