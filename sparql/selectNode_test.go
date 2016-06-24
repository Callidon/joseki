// Copyright (c) 2016 Thomas Minier. All rights reserved.
// Use of this source code is governed by a MIT License
// license that can be found in the LICENSE file.

package sparql

import (
	"github.com/Callidon/joseki/graph"
	"github.com/Callidon/joseki/rdf"
	"testing"
)

func TestExecuteSelectNode(t *testing.T) {
	var graph = graph.NewHDTGraph()
	graph.LoadFromFile("../parser/datas/test.nt", "nt")
	triple := rdf.NewTriple(rdf.NewBlankNode("v1"),
		rdf.NewURI("http://purl.org/dc/terms/title"),
		rdf.NewBlankNode("v2"))
	node := newBgpNode(triple, graph)
	selectNode := newSelectNode(node, []string{"v1"}...)
	cpt := 0

	datas := []rdf.BindingsGroup{
		rdf.NewBindingsGroup(),
		rdf.NewBindingsGroup(),
	}
	datas[0].Bindings["v1"] = rdf.NewURI("http://www.w3.org/2001/sw/RDFCore/ntriples")
	datas[1].Bindings["v1"] = rdf.NewURI("http://www.w3.org/2001/sw/RDFCore/ntriples")

	for bindings := range selectNode.execute() {
		testA, errA := bindings.Equals(datas[0])
		testB, errB := bindings.Equals(datas[1])
		if (!testA || (errA != nil)) && (!testB || (errB != nil)) {
			t.Error(bindings, "should be one of", datas)
		}
		cpt++
	}

	if cpt != len(datas) {
		t.Error(cpt, "bindings found instead of", len(datas))
	}
}

// No need to test joinNode.executeWith(), since it's equivalent to a call to joinNode.execute()
