// Copyright (c) 2016 Thomas Minier. All rights reserved.
// Use of this source code is governed by a MIT License
// license that can be found in the LICENSE file.

package sparql

import "sort"

// containsInt use a binary search to find if a slice contains a Int
func containsInt(data []int, elt int) bool {
	// use binary search on the data
	pos := sort.SearchInts(data, elt)
	if pos == len(data) {
		return false
	} else if pos == 0 {
		return data[0] == elt
	}
	return (data[pos-1] == elt) || (data[pos] == elt)
}

// containsString use a binary search to find if a slice contains a String
func containsString(data []string, elt string) bool {
	// use binary search on the data
	pos := sort.SearchStrings(data, elt)
	if pos == len(data) {
		return false
	} else if pos == 0 {
		return data[0] == elt
	}
	return (data[pos-1] == elt) || (data[pos] == elt)
}
