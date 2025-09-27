#!/bin/bash

my_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$my_dir/.."

coverdir="$(mktemp -d)"

die () {
    echo fatal: "$@" >&2
    exit 1
}

sed-i () {
    # what if we're running on a BSD sed (like on macOS)?
    local is_gnu_sed=false
    sed --version 2>&1 | grep 'GNU sed' > /dev/null && is_gnu_sed=true
    if "$is_gnu_sed"; then
        sed -i "$@"
    else
        # assume everything else is BSD-like
        sed -i '' "$@"
    fi
}

check () {
    local pkg="$1"  # xxx (transformed into ./testdata/xxx)
    local shouldSucceed="$2"  # true / false
    local expectedFixture="$3"  # expected1.txt
    shift
    shift
    shift

    local pkgdir="./testdata/$pkg"
    local stdout="$(mktemp)"
    GOCOVERDIR="$coverdir" ./gosmopolitan "$@" "$pkgdir" > "$stdout" 2>&1
    local ret=$?

    if "$shouldSucceed"; then
        [[ $ret -eq 0 ]] || die "return code should be zero, but is $ret"
    else
        [[ $ret -ne 0 ]] || die "return code should be non-zero"
    fi

    sed-i "s@^.*/testdata/$pkg/@ROOT/@" "$stdout"
    diff -u "$pkgdir/$expectedFixture" "$stdout" || die "unexpected linter output"
    rm "$stdout"
}

check pkgFoo false expected0.txt

check pkgFoo false expected1.txt \
    -escapehatches '(github.com/xen0n/gosmopolitan/testdata/pkgFoo).escapeHatch,(github.com/xen0n/gosmopolitan/testdata/pkgFoo).pri18ntln,(github.com/xen0n/gosmopolitan/testdata/pkgFoo).i18nMessage'

# edge case: unmatched parens
check pkgFoo false expected1.txt \
    -escapehatches '(github.com/xen0n/gosmopolitan/testdata/pkgFoo.escapeHatch,(github.com/xen0n/gosmopolitan/testdata/pkgFoo).pri18ntln,(github.com/xen0n/gosmopolitan/testdata/pkgFoo).i18nMessage'

check pkgFoo false expected2.txt \
    -allowtimelocal \
    -escapehatches 'github.com/xen0n/gosmopolitan/testdata/pkgFoo.escapeHatch,github.com/xen0n/gosmopolitan/testdata/pkgFoo.pri18ntln,github.com/xen0n/gosmopolitan/testdata/pkgFoo.i18nMessage'

check pkgFoo true expected3.txt \
    -allowtimelocal \
    -watchforscripts Arabic

check pkgFoo false expected4.txt \
    -allowtimelocal \
    -watchforscripts Arabic,Devanagari

check pkgFoo false expected5.txt \
    -lookattests

go tool covdata textfmt -i="$coverdir" -o ./coverage.txt
rm -rf "$coverdir"

# report to Codecov if running in CI
if [[ -n $CI ]]; then
    "$my_dir"/codecov.sh
fi
