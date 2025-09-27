-- a.go --
package ignore

func _() {
	//lint:ignore check test
	var _ = ""
}

-- b.go --
package ignore

func _() {
	//lint:ignore check1,check,check2 test
	var _ = ""
}

-- c.go --
package ignore

func _() {
	// lint:ignore check test
	var _ = ""
}

-- d.go --
package ignore

func _() {
	//lint:ignore check multiple words in reason
	var _ = ""
}
