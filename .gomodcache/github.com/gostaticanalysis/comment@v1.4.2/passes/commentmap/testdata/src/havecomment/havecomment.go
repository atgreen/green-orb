-- a.go --
package havecomment

func _() {
	// var is no-op
	//lint:ignore check havecomment
	var _ = ""
}

-- b.go --
package havecomment

func _() {
	//lint:ignore check havecomment
	var _ = "" // var is no-op
}

-- c.go --
package havecomment

func _() {
	// var is no-op
	var _ = "" //lint:ignore check havecomment
}
