// Copyright (c) 2016 Thomas Minier. All rights reserved.
// Use of this source code is governed by a MIT License
// license that can be found in the LICENSE file.

package sparql

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
	// tokenLiteral is a RDF Literal
	tokenLiteral
	// tokenTypedLiteral is a RDF typed Literal
	tokenTypedLiteral
	// tokenLangLiteral is a RDF Literal with lang informations
	tokenLangLiteral
	// tokenVariable is a SPARQL variable
	tokenVariable
)

// sparqlScanner is a scanner for reading a SPARQL request & extracting tokens.
type sparqlScanner struct {
}

// NewsparqlScanner creates a new sparqlScanner.
func newSparqlScanner() *sparqlScanner {
	return &sparqlScanner{}
}

// Analyze a SPARQL request in string format and decompose it into tokens.
//
// Tokens are send through a channel, which is closed when the scan of the request is finished.
func (s *sparqlScanner) scan(request string) chan sparqlToken {
	out := make(chan sparqlToken, bufferSize)
	// scan the request using a goroutine
	go func() {
		defer close(out)
		// TODO
	}()
	return out
}
