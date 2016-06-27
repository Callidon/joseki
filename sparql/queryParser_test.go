// Copyright (c) 2016 Thomas Minier. All rights reserved.
// Use of this source code is governed by a MIT License
// license that can be found in the LICENSE file.

package sparql

import (
	"github.com/Callidon/joseki/graph"
	"testing"
)

func TestQueryParser(t *testing.T) {
	graph := graph.NewHDTGraph()
	_ = LoadQuery("select * where { ?s <http://xmlns.com/foaf/spec/name> 'England'@en . }", graph)
}
