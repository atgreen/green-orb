// MIT License
//
// Copyright (c) 2022 Abirdcfly
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package main

type MyLinkedList struct {
	FakeHead *Element
}

type Element struct {
	Val  int
	Next *Element
}

func Constructor() MyLinkedList {
	return MyLinkedList{FakeHead: &Element{Val: -1, Next: nil}}
}

func (this *MyLinkedList) Get(index int) int {
	i := 0
	for cur := this.FakeHead.Next; cur != nil; cur = cur.Next {
		if i == index {
			return cur.Val
		}
		i++
	}
	return -1
}

func (this *MyLinkedList) AddAtHead(val int) {
	e := &Element{Val: val, Next: this.FakeHead.Next}
	this.FakeHead.Next = e
}

func (this *MyLinkedList) AddAtTail(val int) {
	e := &Element{Val: val, Next: nil}
	var cur, before *Element
	for cur = this.FakeHead; cur != nil; cur = cur.Next {
		before = cur
	}
	before.Next = e
}

func (this *MyLinkedList) AddAtIndex(index int, val int) {
	i := -1
	e := &Element{Val: val, Next: nil}
	for cur := this.FakeHead; cur != nil; cur = cur.Next {
		if i == index-1 {
			e.Next = cur.Next
			cur.Next = e
			break
		}
		i++
	}
}

func (this *MyLinkedList) DeleteAtIndex(index int) {
	i := -1
	for cur := this.FakeHead; cur != nil; cur = cur.Next {
		if i == index-1 {
			e := cur.Next
			if e == nil {
				break
			}
			cur.Next = e.Next
			e.Next = nil
		}
		i++
	}
}
