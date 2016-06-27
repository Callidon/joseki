// Copyright (c) 2016 Thomas Minier. All rights reserved.
// Use of this source code is governed by a MIT License
// license that can be found in the LICENSE file.

package sparql

import (
	"fmt"
	"regexp"
	"strings"
)

var _ = fmt.Println

const (
	// Max size for the buffer of this package.
	bufferSize = 100
)

// Token is the type for a SPARQL token read by a scanner.
type sparqlToken float64

const (
	_ = iota
	// illegalToken is an illegal token in the SPARQL syntax
	illegalToken sparqlToken = 1 << (10 * iota)
	// tokenSelect is a SPARQL SELECT keyword
	tokenSelect
	// tokenConstruct is a SPARQL CONSTRUCT keyword
	tokenConstruct
	// tokenDescribe is a SPARQL DESCRIBE keyword
	tokenDescribe
	// tokenAsk is a SPARQL ASK keyword
	tokenAsk
	// tokenCount is a SPARQL COUNT keyword
	tokenCount
	// tokenWhere is a SPARQL WHERE keyword
	tokenWhere
	// tokenOptionnal is a SPARQL OPTIONNAL keyword
	tokenOptionnal
	// tokenUnion is a SPARQL UNION keyword
	tokenUnion
	// tokenOrderBy is a SPARQL ORDER BY keyword
	tokenOrderBy
	// tokenLimit is a SPARQL LIMIT keyword
	tokenLimit
	// tokenOffset is a SPARQL OFFSET keyword
	tokenOffset
	// tokenEnd ends a triple declaration
	tokenEnd
	// tokenSep is a RDF separator (for object/literal list, etc)
	tokenSep
	// tokenURI is a RDF URI
	tokenURI
	// tokenPrefixedURI is a RDF URI with a prefix
	tokenPrefixedURI
	// tokenLiteral is a RDF Literal
	tokenLiteral
	// tokenTypedLiteral is a RDF typed Literal
	tokenTypedLiteral
	// tokenLangLiteral is a RDF Literal with lang informations
	tokenLangLiteral
	// tokenVariable is a SPARQL variable
	tokenVariable
	// tokenAll represent the wild card symbol *
	tokenAll
	// tokenBGPBegin represent the begin of a BGP statement
	tokenBGPBegin
	// tokenBGPEnd represent the end of a BGP statement
	tokenBGPEnd
)

// valuedToken is a SparqlkToken with an associated value
type valuedToken struct {
	Type  sparqlToken
	Value string
}

// newValuedToken creates a new valuedToken
func newValuedToken(tokenType sparqlToken, value string) valuedToken {
	return valuedToken{tokenType, value}
}

// sparqlScanner is a scanner for reading a SPARQL query & extracting tokens.
type sparqlScanner struct {
}

// NewsparqlScanner creates a new sparqlScanner.
func newSparqlScanner() *sparqlScanner {
	return &sparqlScanner{}
}

// extractSegments parse a string and split the segments into a slice.
// A segment is a string quoted or separated from the other by whitespaces.
func extractSegments(line string) []string {
	r := regexp.MustCompile("'.*?'|\".*?\"|\\S+")
	return r.FindAllString(line, -1)
}

// Analyze a SPARQL query in string format and decompose it into tokens.
//
// Tokens are send through a channel, which is closed when the scan of the query is finished.
func (s sparqlScanner) scan(query string) chan valuedToken {
	out := make(chan valuedToken, bufferSize)
	// scan the query using a goroutine
	go func() {
		defer close(out)
		expectBy := false
		lineNumber := 0

		for _, line := range strings.Split(query, "\n") {
			for _, elt := range extractSegments(line) {
				if expectBy && (elt != "by") {
					out <- newValuedToken(illegalToken, "Syntax error : expected a 'by' keyword but instead found "+elt)
					continue
				}
				// first, process basic SPARQL keyword cases
				switch strings.ToLower(elt) {
				case "select":
					out <- newValuedToken(tokenSelect, elt)
				case "describe":
					out <- newValuedToken(tokenDescribe, elt)
				case "ask":
					out <- newValuedToken(tokenAsk, elt)
				case "construct":
					out <- newValuedToken(tokenConstruct, elt)
				case "where":
					out <- newValuedToken(tokenWhere, elt)
				case "optionnal":
					out <- newValuedToken(tokenOptionnal, elt)
				case "order":
					expectBy = true
				case "by":
					if expectBy {
						out <- newValuedToken(tokenOrderBy, "order by")
						expectBy = false
					} else {
						out <- newValuedToken(illegalToken, "Syntax error : a 'by' must be precede by a 'order'. Instead found "+elt)
					}
				case "count":
					out <- newValuedToken(tokenCount, elt)
				case "union":
					out <- newValuedToken(tokenUnion, elt)
				case "limit":
					out <- newValuedToken(tokenLimit, elt)
				case "offset":
					out <- newValuedToken(tokenOffset, elt)
				case "*":
					out <- newValuedToken(tokenAll, elt)
				case ".", "]":
					out <- newValuedToken(tokenEnd, elt)
				case ",", ";", "[":
					out <- newValuedToken(tokenSep, elt)
				case "{":
					out <- newValuedToken(tokenBGPBegin, elt)
				case "}":
					out <- newValuedToken(tokenBGPEnd, elt)
				default:
					// process more specific cases (URIs, Literals, Variables, separators, etc)
					if (string(elt[0]) == "<") && (string(elt[len(elt)-1]) == ">") {
						out <- newValuedToken(tokenURI, elt[1:len(elt)-1])
					} else if ((string(elt[0]) == "\"") && (string(elt[len(elt)-1]) == "\"")) || ((string(elt[0]) == "'") && (string(elt[len(elt)-1]) == "'")) {
						out <- newValuedToken(tokenLiteral, elt[1:len(elt)-1])
					} else if elt[0:2] == "^^" {
						out <- newValuedToken(tokenTypedLiteral, elt[2:])
					} else if string(elt[0]) == "@" {
						out <- newValuedToken(tokenLangLiteral, elt[1:])
					} else if string(elt[0]) == "?" {
						out <- newValuedToken(tokenVariable, elt[1:])
					} else if strings.Index(elt, ":") > -1 {
						out <- newValuedToken(tokenPrefixedURI, elt)
					} else {
						out <- newValuedToken(illegalToken, "Unexpected token at line "+string(lineNumber)+" : bad syntax")
					}
				}
			}
			lineNumber++
		}
	}()
	return out
}
