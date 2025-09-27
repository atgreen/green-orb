# interfacebloat

Interface bloat (anti-pattern, also called fat interface) is when an interface incorporates too many operations on some data.

A linter that checks length of interface.

The bigger the interface, the weaker the abstraction. (C) Go Proverbs

## Install

```bash
go install github.com/sashamelentyev/interfacebloat
```

## Examples

```bash
interfacebloat ./...
```

## Links

- https://en.wikipedia.org/wiki/Interface_bloat
- https://en.wikipedia.org/wiki/Interface_segregation_principle
