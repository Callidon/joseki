package rdf

// Type which represent a RDF triple.
//
// RDF Triple reference : https://www.w3.org/TR/2004/REC-rdf-concepts-20040210/#section-triples
type Triple struct {
	Subject   Node
	Predicate Node
	Object    Node
}

// Create a new Triple.
func NewTriple(subject, predicate, object Node) Triple {
	return (Triple{subject, predicate, object})
}

// Return True if two triples are strictly equals, False if not.
func (t Triple) Equals(other Triple) (bool, error) {
	test_subj, err := t.Subject.Equals(other.Subject)
	if err != nil {
		return false, err
	}
	test_pred, err := t.Predicate.Equals(other.Predicate)
	if err != nil {
		return false, err
	}
	test_obj, err := t.Object.Equals(other.Object)
	if err != nil {
		return false, err
	}
	return test_subj && test_pred && test_obj, nil
}

// Test if a Triple is equivalent to another triple, assuming that blank node are equals to any other node types.
//
// Return True if the two triples are equivalent with this criteria, False if not.
func (t Triple) Equivalent(other Triple) (bool, error) {
	test_subj, err := t.Subject.Equivalent(other.Subject)
	if err != nil {
		return false, err
	}
	test_pred, err := t.Predicate.Equivalent(other.Predicate)
	if err != nil {
		return false, err
	}
	test_obj, err := t.Object.Equivalent(other.Object)
	if err != nil {
		return false, err
	}
	return test_subj && test_pred && test_obj, nil
}
