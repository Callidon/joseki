// Based on Douglas Hall stack implementation in Go : https://gist.github.com/bemasher/1777766
// Modified by Thomas Minier to add synchronized access

package sparql

import "sync"

// stack is a synchronized LIFO structure implementated as a Linked-List
type stack struct {
	top  *stackElement
	size int
	*sync.Mutex
}

// stackElement is a element in a stack
type stackElement struct {
	value interface{} // All types satisfy the empty interface, so we can store anything here.
	next  *stackElement
}

// newStack creates a new stack
func newStack() *stack {
	return &stack{nil, 0, &sync.Mutex{}}
}

// Return the stack's length
func (s *stack) len() int {
	s.Lock()
	defer s.Unlock()
	return s.size
}

// Push a new element onto the stack
func (s *stack) push(value interface{}) {
	s.Lock()
	s.top = &stackElement{value, s.top}
	s.size++
	s.Unlock()
}

// Remove the top element from the stack and return it's value
// If the stack is empty, return nil
func (s *stack) pop() (value interface{}) {
	s.Lock()
	defer s.Unlock()
	if s.size > 0 {
		value, s.top = s.top.value, s.top.next
		s.size--
		return
	}
	return nil
}

// clear empty the stack
func (s *stack) clear() {
	for s.len() > 0 {
		s.pop()
	}
}
