# grouper — a Go linter to analyze expression groups

## Installation

```sh
go install github.com/leonklingele/grouper@latest
grouper -help
```

## Run analyzer

```sh
grouper -import-require-single-import -import-require-grouping ./...

# Example output:
GOPATH/src/github.com/leonklingele/grouper/pkg/analyzer/analyzer.go:8:1: should only use a single 'import' declaration, 2 found
GOPATH/src/github.com/leonklingele/grouper/pkg/analyzer/flags.go:3:1: should only use grouped 'import' declarations
```

### Available flags

```
  -const-require-grouping
    	require the use of grouped global 'const' declarations
  -const-require-single-const
    	require the use of a single global 'const' declaration only

  -import-require-grouping
    	require the use of grouped 'import' declarations
  -import-require-single-import
    	require the use of a single 'import' declaration only

  -type-require-grouping
    	require the use of grouped global 'type' declarations
  -type-require-single-type
    	require the use of a single global 'type' declaration only

  -var-require-grouping
    	require the use of grouped global 'var' declarations
  -var-require-single-var
    	require the use of a single global 'var' declaration only
```
