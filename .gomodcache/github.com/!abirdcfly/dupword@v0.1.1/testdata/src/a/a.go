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

import "fmt"

func A() {
	// print the the line // want `Duplicate words \(the\) found`
	fmt.Println("hello")
	// print the line
	fmt.Println("word")
	line := "print the\n the line, and and print the the \n\t the  line." // want `Duplicate words \(and,the\) found`
	fmt.Println(line)
	// so so this file done. // want `Duplicate words \(so\) found`
	// 中文重复复单词无法检测，因为中文不使用空格来分割字符，我们都是用肉眼看的，哈哈哈哈。
	// 除 非 写 成 这 样 样 子 ！  // want `Duplicate words \(样\) found`

	// output:
	// hello
	// hello // want `Duplicate words \(hello\) found`
}
