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

package a

import (
	"testing"
)

func TestA(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add the the test cases. // want `Duplicate words \(the\) found`
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			A()
		})
	}
}

func ExampleA() {
	// duplicate words in the example output below should be ignored

	// Output:
	// hello
	// hello
	// hello

	// output:
	// hello
	// hello
	// hello

	// Unordered output:
	// hello
	// hello
	// hello

	// unordered output:
	// hello
	// hello
	// hello

	// this comment block _doesn't_ start with 'output:'
	// output:
	// hello
	// hello // want `Duplicate words \(hello\) found`

	// this also isn't an output block
	// Unordered output:
	// hello
	// hello // want `Duplicate words \(hello\) found`
}
