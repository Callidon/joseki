package graph

import (
    "testing"
    "github.com/Callidon/joseki/core"
)

func TestBimap(t *testing.T) {
    var value core.Node
    bimap := newBimap()
    nodeA := core.NewLiteral("foo")
    nodeB := core.NewLiteral("bar")
    nodeC := core.NewLiteral("test")

    // test with classic insertion
    bimap.push(1, nodeA)
    bimap.push(2, nodeB)

    if key, _ := bimap.locate(nodeA); key != 1 {
        t.Error("expected key = 1 but got", key)
    }
    value, _ = bimap.extract(1)
    if test, _ := value.Equals(nodeA); !test {
        t.Error("expected value = 'foo' but got", value)
    }
    if key, _ := bimap.locate(nodeB); key != 2 {
        t.Error("expected key = 2 but got", key)
    }
    value, _ = bimap.extract(2)
    if test, _ := value.Equals(nodeB); !test {
        t.Error("expected value = 'bar' but got", value)
    }

    // test with non-existent key/value
    if _, test := bimap.locate(nodeC); test {
        t.Error("cannot return true when locating a non-existent value")
    }
    if _, test := bimap.extract(5); test {
        t.Error("cannot return true when extracting a non-existent key")
    }

    // test with value override
    bimap.push(1, nodeB)
    if key, _ := bimap.locate(nodeB); key != 1 {
        t.Error("expected key = 1 but got", key)
    }
    value, _ = bimap.extract(1)
    if test, _ := value.Equals(nodeB); !test {
        t.Error("expected value = 'bar' but got", value)
    }
    if _, test := bimap.locate(nodeA); test {
        t.Error("cannot return true when locating a non-existent value")
    }
    if _, test := bimap.extract(2); test {
        t.Error("cannot return true when extracting a non-existent key")
    }
}
