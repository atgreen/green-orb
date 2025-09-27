# gocheckcompilerdirectives

Check that go compiler directives (`//go:` comments) are valid and catch easy
mistakes.

Compiler directives are comments like `//go:generate`, `//go:embed`,
`//go:build`, etc.

## Why

Go compiler directives are comments in the form of `//go:` that provide an
instruction to the compiler.

The directives are easy to make mistakes with. The linter will detect the
following mistakes:

1. Adding a space in between the comment bars and the first character, e.g. `//
go:`, will cause the compiler to silently ignore the comment.

2. Mistyping a directives name, e.g. `//go:embod`, will cause the compiler to silently ignore the comment.

## Install

```
go install 4d63.com/gocheckcompilerdirectives@latest
```

## Usage

```
gocheckcompilerdirectives [package]
```

```
gocheckcompilerdirectives ./...
```
