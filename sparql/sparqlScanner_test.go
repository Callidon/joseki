// Copyright (c) 2016 Thomas Minier. All rights reserved.
// Use of this source code is governed by a MIT License
// license that can be found in the LICENSE file.

package sparql

import "testing"

func TestScanSimpleSparqlScanner(t *testing.T) {
	scanner := newSparqlScanner()
	datas := []valuedToken{
		newValuedToken(tokenSelect, "select"),
		newValuedToken(tokenAll, "*"),
		newValuedToken(tokenWhere, "where"),
		newValuedToken(tokenBGPBegin, "{"),
		newValuedToken(tokenVariable, "s"),
		newValuedToken(tokenURI, "http://xmlns.com/foaf/spec/name"),
		newValuedToken(tokenLiteral, "England"),
		newValuedToken(tokenLangLiteral, "en"),
		newValuedToken(tokenEnd, "."),
		newValuedToken(tokenBGPEnd, "}"),
	}
	cpt := 0

	for token := range scanner.scan("select * where { ?s <http://xmlns.com/foaf/spec/name> 'England'@en . }") {
		if (token.Value != datas[cpt].Value) || (token.Type != datas[cpt].Type) {
			t.Error("expected", datas[cpt], "but instead got", token)
		}
		cpt++
	}

	if cpt != len(datas) {
		t.Error("expected", len(datas), "tokens but instead found", cpt, "tokens")
	}
}
