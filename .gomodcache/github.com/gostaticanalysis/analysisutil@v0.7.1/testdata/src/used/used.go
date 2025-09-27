package used

var (
	N int
	V interface{}
)

func f1(v interface{}) { // want "used"
	println(v)
}

func f2(v interface{}) {} // unsed

func f3(v interface{}) { // unsed
	{
		v := 100
		println(v)
	}
}

func f4(v interface{}) { // want "used"
	V = v
}

func f5(v interface{}) { // want "used"
	if N == 0 {
		return
	}
	V = v
}

func f6(v interface{}) { // want "used"
	func() {
		println(v)
	}()
}

func f7(v interface{}) { // want "used"
	go func() {
		println(v)
	}()
}

func f8(v interface{}) { // want "used"
	defer func() {
		println(v)
	}()
}

func f9(v interface{}) { // want "used"
	func(v interface{}) { // want "used"
		println(v)
	}(v)
}
