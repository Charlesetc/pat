// stack.go

package stack

import "errors"

type stackJoint struct {
	previous *stackJoint
	value    interface{}
}

type Stack struct {
	start *stackJoint
	saved *stackJoint
}

func newJoint(val interface{}) *stackJoint {
	return &stackJoint{value: val}
}

func (s *Stack) Add(val interface{}) {
	joint := newJoint(val)
	joint.previous = s.start
	s.start = joint
	s.saved = nil
}

func (s *Stack) Pop() (interface{}, error) {
	if s.start == nil {
		return nil, errors.New("Stack is empty")
	}
	saved := s.saved
	s.saved = s.start

	out := s.start.value
	s.start = s.start.previous
	s.saved.previous = saved
	return out, nil
}

func (s *Stack) UnPop() (interface{}, error) {
	if s.saved == nil {
		return nil, errors.New("Saved Stack is empty")
	}
	start := s.start
	s.start = s.saved

	out := s.saved.value
	s.saved = s.saved.previous
	s.start.previous = start
	return out, nil
}

func New() *Stack {
	return &Stack{nil, nil}
}
