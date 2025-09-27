# testableexamples <br> [![CI][ci-img]][ci-url] [![Codecov][codecov-img]][codecov-url] [![Codebeat][codebeat-img]][codebeat-url] [![Maintainability][codeclimate-img]][codeclimate-url] [![Go Report Card][goreportcard-img]][goreportcard-url] [![License][license-img]][license-url] [![Go Reference][godoc-img]][godoc-url]

Linter checks if examples are testable (have an expected output).

> Author of idea is [Jamie Tanna](https://github.com/jamietanna) (see [this issue](https://github.com/golangci/golangci-lint/issues/3084)).


## Description

Example functions without output comments are compiled but not executed by `go test`, see [doc](https://pkg.go.dev/testing#hdr-Examples).  
It means that such examples are not validated. This can lead to outdated code.  
That's why the linter complains on missing output.

```go
func Example_bad() { // <- linter will complain on missing output
	fmt.Println("hello")
}

func Example_good() {
	fmt.Println("hello")
	// Output: hello
}
```


## Usage

The best way is to use [golangci-lint](https://golangci-lint.run/).  
It includes [testableexamples](https://golangci-lint.run/usage/linters/#list-item-testableexamples) and a lot of other great linters.

### Install

See [official site](https://golangci-lint.run/usage/install/).

### Enable

`testableexamples` is disabled by default.  
To enable it, add the following to your `.golangci.yml`:

```yaml
linters:
  enable:
    testableexamples
```

### Run

```shell
golangci-lint run
```

### Directive `//nolint:testableexamples`

Here is incorrect usage of `//nolint` directive.  
`golangci-lint` understands it correctly, but `godoc` will include comment to the code.

```go
// Description
func Example_nolintIncorrect() { //nolint:testableexamples // that's why
	fmt.Println("hello")
}
```

Here are two examples of correct usage of `//nolint` directive.  
`godoc` will ignore comment.

```go
//nolint:testableexamples // that's why
func Example_nolintCorrect() {
	fmt.Println("hello")
}

// Description
//
//nolint:testableexamples // that's why
func Example_nolintCorrectWithDescription() {
	fmt.Println("hello")
}
```


## Usage as standalone linter

### Install
```shell
go install github.com/maratori/testableexamples@latest
```

### Run

```shell
testableexamples ./...
```

### Nolint

Standalone linter doesn't support `//nolint` directive.  
And there is no alternative for that, please use `golangci-lint`.


## License

[MIT License][license-url]


[ci-img]: https://github.com/maratori/testableexamples/actions/workflows/ci.yml/badge.svg
[ci-url]: https://github.com/maratori/testableexamples/actions/workflows/ci.yml
[codecov-img]: https://codecov.io/gh/maratori/testableexamples/branch/main/graph/badge.svg?token=VMXc2fc7cJ
[codecov-url]: https://codecov.io/gh/maratori/testableexamples
[codebeat-img]: https://codebeat.co/badges/1b813bf1-336d-4886-b4fa-1d482bedc754
[codebeat-url]: https://codebeat.co/projects/github-com-maratori-testableexamples-main
[codeclimate-img]: https://api.codeclimate.com/v1/badges/47ed5db4a7595d4f95d5/maintainability
[codeclimate-url]: https://codeclimate.com/github/maratori/testableexamples/maintainability
[goreportcard-img]: https://goreportcard.com/badge/github.com/maratori/testableexamples
[goreportcard-url]: https://goreportcard.com/report/github.com/maratori/testableexamples
[license-img]: https://img.shields.io/github/license/maratori/testableexamples.svg
[license-url]: /LICENSE
[godoc-img]: https://pkg.go.dev/badge/github.com/maratori/testableexamples.svg
[godoc-url]: https://pkg.go.dev/github.com/maratori/testableexamples
