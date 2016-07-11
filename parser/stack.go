// Based on Douglas Hall stack implementation in Go : https://gist.github.com/bemasher/1777766
// Modified by Thomas Minier to add synchronized access

package parser

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

// NewStack creates a new stack
func NewStack() *stack {
	return &stack{nil, 0, &sync.Mutex{}}
}

// Return the stack's length
func (s *stack) Len() int {
	s.Lock()
	defer s.Unlock()
	return s.size
}

// Push a new element onto the stack
func (s *stack) Push(value interface{}) {
	s.Lock()
	s.top = &stackElement{value, s.top}
	s.size++
	s.Unlock()
}

// Get the next element without popping the stack
func (s *stack) Peek() interface{} {
	s.Lock()
	defer s.Unlock()
	if s.size > 0 {
		return s.top.value
	}
	return nil
}

// Remove the top element from the stack and return it's value
// If the stack is empty, return nil
func (s *stack) Pop() (value interface{}) {
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
func (s *stack) Clear() {
	for s.Len() > 0 {
		s.Pop()
	}
}
