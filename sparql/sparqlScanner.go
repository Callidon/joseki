// Copyright (c) 2016 Thomas Minier. All rights reserved.
// Use of this source code is governed by a MIT License
// license that can be found in the LICENSE file.

package sparql

const (
	// Max size for the buffer of this package.
	bufferSize = 100
)

// Token is the type for a SPARQL token read by a scanner.
type SparqlToken float64

const (
    _ = iota
    // TokenIllegal is an illegal token in the SPARQL syntax
    IllegalToken SparqlToken = 1 << (10 * iota)
    // TokenSelect is a SPARQL SELECT keyword
    TokenSelect
    // TokenConstruct is a SPARQL CONSTRUCT keyword
    TokenConstruct
    // TokenDescribe is a SPARQL DESCRIBE keyword
    TokenDescribe
    // TokenAsk is a SPARQL ASK keyword
    TokenAsk
    // TokenCount is a SPARQL COUNT keyword
    TokenCount
    // TokenWhere is a SPARQL WHERE keyword
    TokenWhere
    // TokenOptionnal is a SPARQL OPTIONNAL keyword
    TokenOptionnal
    // TokenUnion is a SPARQL UNION keyword
    TokenUnion
    // TokenOrderBy is a SPARQL ORDER BY keyword
    TokenOrderBy
    // TokenLimit is a SPARQL LIMIT keyword
    TokenLimit
    // TokenOffset is a SPARQL OFFSET keyword
    TokenOffset
    // TokenEnd ends a triple declaration
	TokenEnd
	// TokenSep is a RDF separator (for object/literal list, etc)
	TokenSep
	// TokenURI is a RDF URI
	TokenURI
	// TokenLiteral is a RDF Literal
	TokenLiteral
	// TokenTypedLiteral is a RDF typed Literal
	TokenTypedLiteral
	// TokenLangLiteral is a RDF Literal with lang informations
	TokenLangLiteral
	// TokenVariable is a SPARQL variable
	TokenVariable
)

// SparqlScanner is a scanner for reading a SPARQL request & extracting tokens.
type SparqlScanner struct {
}

// NewSparqlScanner creates a new SparqlScanner.
func NewSparqlScanner() *SparqlScanner {
    return &SparqlScanner{}
}

// Scan analyze a SPARQL request in string format and decompose it into tokens.
//
// Tokens are send through a channel, which is closed when the scan of the request is finished.
func (s *SparqlScanner) Scan(request string) chan SparqlToken {
    out := make(chan SparqlToken, bufferSize)
    // scan the request using a goroutine
    go func() {
        defer close(out)
        // TODO
    }()
    return out
}
