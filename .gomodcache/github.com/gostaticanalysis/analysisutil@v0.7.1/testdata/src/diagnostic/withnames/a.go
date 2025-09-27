package withnames

var a int //lint:ignore check1 reason

var b struct { // want "NG"
	N int //lint:ignore check2 reason
}

var c int // want "NG" "NG"

var d struct { // want "NG" "NG"
	N int
}
