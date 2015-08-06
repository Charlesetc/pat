// stack_test.go

package stack

import "testing"

var stack *Stack

func init() {
	stack = New()
}

func TestAdd(t *testing.T) {
	stack.Add("a")
	stack.Add("c")
	stack.Add("d")
	ret, err := stack.Pop()
	if err != nil || ret.(string) != "d" {
		t.Errorf("Expected 'd', got '%s' instead", ret.(string))
	}
	ret, err = stack.Pop()
	if err != nil || ret.(string) != "c" {
		t.Errorf("Expected 'c', got '%s' instead", ret.(string))
	}
	ret, err = stack.Pop()
	if err != nil || ret.(string) != "a" {
		t.Errorf("Expected 'a', got '%s' instead", ret.(string))
	}
	ret, err = stack.Pop()
	if err == nil {
		t.Errorf("Expected popping empty stack to error")
	}
	ret, err = stack.UnPop()
	if err != nil || ret.(string) != "a" {
		t.Errorf("Expected 'a', got '%s' instead", ret.(string))
	}
	ret, err = stack.UnPop()
	if err != nil || ret.(string) != "c" {
		t.Errorf("Expected 'b', got '%s' instead", ret.(string))
	}
	ret, err = stack.UnPop()
	if err != nil || ret.(string) != "d" {
		t.Errorf("Expected 'c', got '%s' instead", ret.(string))
	}
}
