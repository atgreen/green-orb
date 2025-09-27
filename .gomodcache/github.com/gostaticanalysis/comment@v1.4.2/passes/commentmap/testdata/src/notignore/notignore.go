-- a.go --
package ignore

func _() {
	//lint:ignore check1 test
	var _ = ""
}
-- b.go --
package ignore

func _() {
	//lint:ignore check1
	var _ = ""
}
