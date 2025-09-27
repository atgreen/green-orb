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

package dupword_test

import (
	"testing"

	"github.com/Abirdcfly/dupword"
	"golang.org/x/tools/go/analysis/analysistest"
)

func Test(t *testing.T) {
	analyzer := dupword.NewAnalyzer()
	tests := []string{"a", "good"}
	analysistest.Run(t, analysistest.TestData(), analyzer, tests...)

	analyzer1 := dupword.NewAnalyzer()
	_ = analyzer1.Flags.Set("ignore", "the")
	defer dupword.ClearIgnoreWord()
	analysistest.Run(t, analysistest.TestData(), analyzer1, "a_ignore_the")

	analyzer2 := dupword.NewAnalyzer()
	_ = analyzer.Flags.Set("ignore", "the,and")
	analysistest.Run(t, analysistest.TestData(), analyzer2, "a_ignore_the_and")
}

func Test_checkOneKey(t *testing.T) {
	type args struct {
		raw string
		key string
	}
	tests := []struct {
		name     string
		args     args
		wantNew  string
		wantKey  string
		wantFind bool
	}{
		{
			name: "one word",
			args: args{
				raw: "Done",
				key: "the",
			},
			wantNew:  "",
			wantFind: false,
		},
		{
			name: "one word with space",
			args: args{
				raw: " Done \n \t",
				key: "the",
			},
			wantNew:  "",
			wantFind: false,
		},
		{
			name: "one line without key word",
			args: args{
				raw: "hello word",
				key: "the",
			},
			wantNew:  "",
			wantFind: false,
		},
		{
			name: "one line with key word only once",
			args: args{
				raw: "hello the word",
				key: "the",
			},
			wantNew:  "",
			wantFind: false,
		},
		{
			name: "one line with key word twice",
			args: args{
				raw: "hello the the world",
				key: "the",
			},
			wantNew:  "hello the world",
			wantKey:  "the",
			wantFind: true,
		},
		{
			name: "one line with key word multi times",
			args: args{
				raw: "hello the the the world",
				key: "the",
			},
			wantNew:  "hello the world",
			wantKey:  "the",
			wantFind: true,
		},
		{
			name: "multi line with key word once",
			args: args{
				raw: "hello \t the \nworld",
				key: "the",
			},
			wantNew:  "",
			wantFind: false,
		},
		{
			name: "multi line with key word twice",
			args: args{
				raw: "// NewDBStats returns a new DBStats object initialized using the\n the new function from each component.",
				key: "the",
			},
			wantNew:  "// NewDBStats returns a new DBStats object initialized using the\n new function from each component.",
			wantKey:  "the",
			wantFind: true,
		},
		{
			name: "multi line with key word multi times",
			args: args{
				raw: "hello \t the \n   the the world",
				key: "the",
			},
			wantNew:  "hello \t the \n   world",
			wantKey:  "the",
			wantFind: true,
		},
		{
			name: "multi line with key word multi times",
			args: args{
				raw: "print the\nthe line, print the\n\t the line.",
				key: "the",
			},
			wantNew:  "print the\nline, print the\n\t line.",
			wantKey:  "the",
			wantFind: true,
		},
		{
			name: "multi line with key word multi times",
			args: args{
				raw: "print the\nthe line, print the\n\t the line.",
			},
			wantNew:  "print the\nline, print the\n\t line.",
			wantKey:  "the",
			wantFind: true,
		},
		{
			name: "duplicate word is last word",
			args: args{
				raw: "this line include duplicate word	and\nand",
			},
			wantNew:  "this line include duplicate word	and",
			wantKey:  "and",
			wantFind: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotNew, gotKey, gotFind := dupword.CheckOneKey(tt.args.raw, tt.args.key)
			if gotNew != tt.wantNew {
				t.Errorf("CheckOneKey() gotNew = %q, want %q", gotNew, tt.wantNew)
			}
			if gotFind != tt.wantFind {
				t.Errorf("CheckOneKey() gotFind = %v, want %v", gotFind, tt.wantFind)
			}
			if gotKey != tt.wantKey {
				t.Errorf("CheckOneKey() gotKey = %v, want %v", gotKey, tt.wantKey)
			}
		})
	}
}

func TestExcludeWords(t *testing.T) {
	type args struct {
		word string
	}
	tests := []struct {
		name        string
		args        args
		wantExclude bool
	}{
		{
			name:        "normal words should not exclude",
			args:        args{word: "normal"},
			wantExclude: false,
		},
		{
			name:        "number should not exclude",
			args:        args{word: "3"},
			wantExclude: true,
		},
		{
			name:        "number,number should exclude",
			args:        args{word: "3,3 3,3"},
			wantExclude: true,
		},
		{
			name:        "%s should exclude",
			args:        args{word: "%s"},
			wantExclude: true,
		},
		{
			name:        "</div> should exclude",
			args:        args{word: "</div>"},
			wantExclude: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotExclude := dupword.ExcludeWords(tt.args.word); gotExclude != tt.wantExclude {
				t.Errorf("excludeWords() = %v, want %v", gotExclude, tt.wantExclude)
			}
		})
	}
}
