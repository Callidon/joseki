package graph

import "github.com/Callidon/joseki/rdf"

// Data structure that represents bidirectional relations between elements of two collections.
// Used as a dictionnary in the HDT-MR Graph implementation
type bimap struct {
	keyToValue map[int]rdf.Node
	valueToKey map[rdf.Node]int
}

// Return a new empty Bimap.
func newBimap() bimap {
	return bimap{make(map[int]rdf.Node), make(map[rdf.Node]int)}
}

// Add a (key, value) to the Bimap.
// If a key or a value already exist in the bimap, erase the associate relation.
func (b *bimap) push(key int, value rdf.Node) {
	// insert value in map, but check if it's already present before
	previousValue, inMap := b.keyToValue[key]
	if inMap {
		b.keyToValue[key] = value
		// remove association in other map before updating it
		delete(b.valueToKey, previousValue)
	} else {
		b.keyToValue[key] = value
	}
	// same thing for the key
	previousKey, inMap := b.valueToKey[value]
	if inMap {
		b.valueToKey[value] = key
		// remove association in other map before updating it
		delete(b.keyToValue, previousKey)
	} else {
		b.valueToKey[value] = key
	}
}

// Return the key associated to a value in the Bimap.
func (b *bimap) locate(value rdf.Node) (int, bool) {
	key, test := b.valueToKey[value]
	return key, test
}

// Return the value associated to a key in the Bimap.
func (b *bimap) extract(key int) (rdf.Node, bool) {
	value, test := b.keyToValue[key]
	return value, test
}
