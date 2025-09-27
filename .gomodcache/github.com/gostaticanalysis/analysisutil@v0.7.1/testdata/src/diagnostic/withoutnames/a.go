package withoutnames

var a int //lint:ignore test reason

var b struct { // want "NG"
	N int //lint:ignore test reason
}

var c int // want "NG" "NG"

var d struct { // want "NG" "NG"
	N int
}
