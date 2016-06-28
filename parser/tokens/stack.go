// Based on Douglas Hall Stack implementation in Go : https://gist.github.com/bemasher/1777766
// Modified by Thomas Minier to add synchronized access

package tokens

import "sync"

// Stack is a synchronized LIFO structure implementated as a Linked-List
type Stack struct {
	top  *StackElement
	size int
	*sync.Mutex
}

// StackElement is a element in a Stack
type StackElement struct {
	value interface{} // All types satisfy the empty interface, so we can store anything here.
	next  *StackElement
}

// NewStack creates a new Stack
func NewStack() *Stack {
	return &Stack{nil, 0, &sync.Mutex{}}
}

// Return the Stack's length
func (s *Stack) Len() int {
	s.Lock()
	defer s.Unlock()
	return s.size
}

// Push a new element onto the Stack
func (s *Stack) Push(value interface{}) {
	s.Lock()
	s.top = &StackElement{value, s.top}
	s.size++
	s.Unlock()
}

// Get the next element without popping the Stack
func (s *Stack) Peek() interface{} {
	s.Lock()
	defer s.Unlock()
	if s.size > 0 {
		return s.top.value
	}
	return nil
}

// Remove the top element from the Stack and return it's value
// If the Stack is empty, return nil
func (s *Stack) Pop() (value interface{}) {
	s.Lock()
	defer s.Unlock()
	if s.size > 0 {
		value, s.top = s.top.value, s.top.next
		s.size--
		return
	}
	return nil
}

// clear empty the Stack
func (s *Stack) Clear() {
	for s.Len() > 0 {
		s.Pop()
	}
}
