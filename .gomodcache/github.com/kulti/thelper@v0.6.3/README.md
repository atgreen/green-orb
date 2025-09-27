# thelper

![CI](https://github.com/kulti/thelper/workflows/CI/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/kulti/thelper)](https://goreportcard.com/report/github.com/kulti/thelper)
[![Coverage](https://coveralls.io/repos/github/kulti/thelper/badge.svg?branch=master)](https://coveralls.io/github/kulti/thelper?branch=master)
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENSE)

thelper detects golang test helpers without `t.Helper()` call. Also, it checks the consistency of test helpers and has similar checks for benchmarks and TB interface. Read further to learn more.

## Why?
### Why I need to add t.Helper() into my helpers functions?

[With this call go test prints correct lines of code for failed tests](https://golang.org/pkg/testing/#T.Helper
). Otherwise, printed lines will be inside helpers functions.

Compare (https://goplay.space/#I_pAJ-vWcRq):
```
func TestBad(t *testing.T) {
	testBad(t)
}

func testBad(t *testing.T) {
	t.Fatal("unhelpful line") // <--- output point this line
}

func TestGood(t *testing.T) {
	testGood(t) // <--- output point this line
}

func testGood(t *testing.T) {
	t.Helper()
	t.Fatal("clear output")
}
```

### Why should I call `t.Helper()` as the first statement in my test helper functions?

Because `t.Helper()` only affects asserts that are called _after_ the `t.Helper()` call, so requiring it to be the first statement helps ensure
all assertions in the helper function are affected.

### Why do I need to name `*testing.T` as `t`? Why it should be the first parameter?

It adds more consistency into your code. When common variables have the same name and are placed into the same position, it makes it easier to understand and read.

Note that you can also have it as the second parameter - this is primarily intended to allow compatibility with `context.Context`.

## Installation

```
go get github.com/kulti/thelper/cmd/thelper
```

## Usage

### golangci-lint

[golangci-lint](https://golangci-lint.run/) supports thelper, so you can enable this linter and use it.

### Shell

Check everything:
```
thelper ./...
```

### Enable only required checks

If you run via golangci-lint look at [.golangci.example.yml](https://golangci-lint.run/usage/configuration/#config-file) for an example of the configuration.

Otherwise you can run thelper with `--checks` command line argument. E.g. the following command checks that `*testing.T` is the first param and `*testing.B` is named `b`:
```
thelper --checks=t_first,b_name ./...
```

### Available checks
* t_begin - `t.Helper()` should begin helper function.
* t_name - `*testing.T` should be named `t`.
* t_first - `*testing.T` should be the first param of helper function.

The same for fuzzing, benchmarks and TB interface:
* f_begin - `f.Helper()` should begin helper function.
* f_name - `*testing.F` should be named `f`.
* f_first - `*testing.F` should be the first param of helper function.
* b_begin - `b.Helper()` should begin helper function.
* b_name - `*testing.B` should be named `b`.
* b_first - `*testing.B` should be the first param of helper function.
* tb_begin - `tb.Helper()` should begin helper function.
* tb_name - `*testing.TB` should be named `tb`.
* tb_first - `*testing.TB` should be the first param of helper function.

### Exceptions

* t_name allows using `_` name.
* t_begin allows subtests and fuzz tests not begin from `t.Helper()`.
* t_first allows to be the second after `context.Context` param.

## Examples

### Without t.Helper()

```go
func TestSmth(t *testing.T) {
	checkSmth(t)
}

// Bad
func checkSmth(t *testing.T) {
	t.Fatal("check failed")
}

// Good
func checkSmth(t *testing.T) {
	t.Helper()
	t.Fatal("check failed")
}
```

### With invalid name

```go
// Bad
func checkSmth(tt *testing.T) {
    // ...
}

// Good
func checkSmth(t *testing.T) {
    // ...
}
```

### With t as not the first param

```go
// Bad
func checkSmth(doneCh <-chan struct{}, t *testing.T) {
    // ...
}

// Good
func checkSmth(t *testing.T, doneCh <-chan struct{}) {
    // ...
}
```
