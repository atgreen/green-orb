# Contribution Guide

## Prerequisites

Make sure following installed on your machine:
 - Go
 - Make
 - Docker

## Before Create PR

### Write Tests

All new code must be covered with unit tests.
Please make sure you've added all necessary tests.

### Run Tests

```shell
make test
```

### Run Linter

```shell
make lint
```

## Development inside container

Docker container contains all necessary tools for development. Just run bash in the dev container.

```shell
make bash
```

## Makefile

To see all available make commands run

```shell
make help
```
