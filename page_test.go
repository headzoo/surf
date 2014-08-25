package gosurf

import (
	"github.com/headzoo/gosurf/unittest"
	"testing"
)

func TestPageStack(t *testing.T) {
	unittest.Run(t)
	stack := NewPageStack()

	page1 := &Page{}
	stack.Push(page1)
	unittest.AssertEquals(1, stack.Len())
	unittest.AssertEquals(page1, stack.Top())

	page2 := &Page{}
	stack.Push(page2)
	unittest.AssertEquals(2, stack.Len())
	unittest.AssertEquals(page2, stack.Top())

	page := stack.Pop()
	unittest.AssertEquals(page, page2)
	unittest.AssertEquals(1, stack.Len())
	unittest.AssertEquals(page1, stack.Top())

	page = stack.Pop()
	unittest.AssertEquals(page, page1)
	unittest.AssertEquals(0, stack.Len())
}
