package core

import "errors"

// Interface which represent a generic node in a RDF Graph
type Node interface {
	Equals(n Node) (bool, error)
	String() string
}

// Type which represent an URI Node in a RDF Graph
// RDF URI reference : https://www.w3.org/TR/2004/REC-rdf-concepts-20040210/#section-Graph-URIref
type URI struct {
	value  string
	prefix string
}

// Type which represent a Literal Node in a RDF Graph
// RDF Literal reference : https://www.w3.org/TR/2004/REC-rdf-concepts-20040210/#section-Graph-Literal
type Literal struct {
	value string
}

// Type which represent an Blank Node in a RDF Graph
// RDF Blank Node reference : https://www.w3.org/TR/2004/REC-rdf-concepts-20040210/#section-blank-nodes
type BlankNode struct {
	variable string
}

// Type which represent a RDF triple
// RDF Triple reference : https://www.w3.org/TR/2004/REC-rdf-concepts-20040210/#section-triples
type Triple struct {
	subject   Node
	predicate Node
	object    Node
}

// Expand an URI with it's prefix
func (u URI) expandName() string {
	if u.prefix != "" {
		return "<" + u.prefix + ":" + u.value + ">"
	} else {
		return "<" + u.value + ">"
	}
}

// Return True if Two URIs are equals, False if not
func (u URI) Equals(n Node) (bool, error) {
	other, ok := n.(URI)
	if ok {
		return u.expandName() == other.expandName(), nil
	} else {
		return false, errors.New("Error : mismatch type, can only compare two URIs")
	}
}

// Serialize an URI to string and return it
func (u URI) String() string {
	return u.expandName()
}

// Create an new URI
func NewURI(prefix, value string) URI {
	return URI{value, prefix}
}

// Return True if Two Literals are equals, False if not
func (l Literal) Equals(n Node) (bool, error) {
	other, ok := n.(Literal)
	if ok {
		return l.value == other.value, nil
	} else {
		return false, errors.New("Error : mismatch type, can only compare two Literals")
	}
}

// Serialize a Literal to string and return it
func (l Literal) String() string {
	return "\"" + l.value + "\""
}

// Create a new Literal
func NewLiteral(value string) Literal {
	return Literal{value}
}

// Return True if Two Blank Node are equals, False if not
func (b BlankNode) Equals(n Node) (bool, error) {
	other, ok := n.(BlankNode)
	if ok {
		return b.variable == other.variable, nil
	} else {
		return false, errors.New("Error : mismatch type, can only compare two Blank Nodes")
	}
}

// Serialize an Blank Node to string and return it
func (b BlankNode) String() string {
	return "?" + b.variable
}

// Create a new Literal
func NewBlankNode(variable string) BlankNode {
	return BlankNode{variable}
}

// Return True if two triples are equals, False if not
func (t Triple) Equals(other Triple) (bool, error) {
	test_subj, err := t.subject.Equals(other.subject)
	if err != nil {
		return false, err
	}
	test_pred, err := t.predicate.Equals(other.predicate)
	if err != nil {
		return false, err
	}
	test_obj, err := t.object.Equals(other.object)
	if err != nil {
		return false, err
	}
	return test_subj && test_pred && test_obj, nil
}

// Create a new Triple
func NewTriple(subject Node, predicate Node, object Node) Triple {
	return (Triple{subject, predicate, object})
}
