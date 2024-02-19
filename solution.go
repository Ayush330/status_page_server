package main

import "fmt"

type stack struct {
	el      []string
	pointer int
}

func (s *stack) back() {
	currPointer := s.pointer
	s.pointer = max(currPointer-1, 0)
}

func (t *stack) forward() {
	currVal := t.pointer
	t.pointer = min(len(t.el), currVal+1)
}

func (t *stack) new(element string) {x``
	// remove all things
	currIndexTemp := t.pointer
	currIndex := min(currIndexTemp, 1)
	t.el = t.el[:currIndex]
	t.el = append(t.el, element)
	t.pointer = len(t.el)
}

func (t *stack) size() int {
	return len(t.el)
}

func (t *stack) stringify() string {
	result := "/home"
	for _, element := range t.el {
		result = fmt.Sprintf("%s/%s", result, element)
	}
	return fmt.Sprintf("%s", result)
}

func solution(N int, S []string) string {
	Stack := &stack{}
	for _, action := range S {
		if action == "back" {
			Stack.back()
		} else if action == "forward" {
			Stack.forward()
		} else {
			Stack.new(action)
		}
	}
	return Stack.stringify()
}
