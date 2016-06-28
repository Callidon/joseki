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
)

// Parser represent a generic interface for parsing every RDF format.
//
// Package parser provides several implementations for this interface.
type Parser interface {
	Read(filename string) chan rdf.Triple
	Prefixes() map[string]string
}

// Utility function for checking errors
func check(err error) {
	if err != nil {
		panic(err)
	}
}

// extractSegments parse a string and split the segments into a slice.
// A segment is a string quoted or separated from the other by whitespaces.
func extractSegments(line string) []string {
	r := regexp.MustCompile("'.*?'|\".*?\"|\\S+")
	return r.FindAllString(line, -1)
}
