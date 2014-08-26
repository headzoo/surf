package element

import (
	ut "github.com/headzoo/surf/unittest"
	"testing"
)

func TestPageStack(t *testing.T) {
	ut.Run(t)
	stack := NewPageStack()

	page1 := &Page{}
	stack.Push(page1)
	ut.AssertEquals(1, stack.Len())
	ut.AssertEquals(page1, stack.Top())

	page2 := &Page{}
	stack.Push(page2)
	ut.AssertEquals(2, stack.Len())
	ut.AssertEquals(page2, stack.Top())

	page := stack.Pop()
	ut.AssertEquals(page, page2)
	ut.AssertEquals(1, stack.Len())
	ut.AssertEquals(page1, stack.Top())

	page = stack.Pop()
	ut.AssertEquals(page, page1)
	ut.AssertEquals(0, stack.Len())
}
