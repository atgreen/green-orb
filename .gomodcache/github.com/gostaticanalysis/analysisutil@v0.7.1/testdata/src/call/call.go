package call

var Flag bool

// res is a resource which must be closed after using.
type res struct{}

func newRes() *res {
	return &res{}
}

func (r *res) close() {}

func test1() {
	r := newRes() // want `NG`
	_ = r
}

func test2() {
	r := newRes() // OK
	r.close()
}

func test3() {
	r := newRes() // want `NG`
	for i := 0; i < 3; i++ {
	}
	_ = r
}

func test4() {
	r := newRes() // OK
	for i := 0; i < 3; i++ {
	}
	r.close()
}

func test5() {
	r := newRes() // want `NG`
	if Flag {
		// this flow does not close res
		return
	}
	r.close()
}

func test6() {
	r := newRes() // OK
	if Flag {
		r.close()
		return
	}
	r.close()
}

func test7() {
	r := newRes() // OK
	defer r.close()
	if Flag {
		return
	}
}

func test8() {
	r := newRes() // OK
	defer r.close()
	func() {
		func(r *res) {}(r)
	}()
}

func test9() {
	r := newRes() // NG
	func() {
		func(r *res) {}(r)
	}()
}

func test10() {
	r := newRes() // OK

	if Flag { // divide into multiple blocks
		println("hoge")
	}

	defer r.close()
	func() {
		func(r *res) {}(r)
	}()
}

func test11() {
	r := newRes() // NG

	if Flag { // divide into multiple blocks
		println("hoge")
	}

	func() {
		func(r *res) {}(r)
	}()
}

func test12() {
	r := newRes() // OK
	defer r.close()
	func() {
		_ = r // OK (free var)
	}()
}

func test13() {
	r := newRes() // OK
	defer r.close()
	func() {
		r2 := r // OK (free var)
		_ = r2
	}()
}

func test14() {
	r := newRes() // OK
	defer r.close()
	go func() {
		r2 := r // OK (free var)
		_ = r2
	}()
}

var pkgRes = newRes()

func test15() {
	r := pkgRes // OK (package variable)
	_ = r
}

// embedded type
type res2 struct{ *res }

func newRes2() *res2 {
	return &res2{newRes()}
}
func test16() {
	r := newRes2() // want `NG`
	_ = r
}

// embedded type
type res3 struct{ res }

func newRes3() *res3 {
	return &res3{}
}
func test17() {
	r := newRes3() // want `NG`
	_ = r
}

// embedded type
type res4 struct{ res }

func newRes4() res4 {
	return res4{}
}
func test18() {
	r := newRes4() // want `NG`
	_ = r
}

// embedded type
type res5 struct{ *res2 }

func newRes5() res5 {
	return res5{}
}
func test19() {
	r := newRes5() // want `NG`
	_ = r
}

func test20() {
	println("hello")
	r := newRes() // want `NG`
	_ = r
}

func test21() interface{} {
	println("hello")
	r := newRes()
	return r
}

func test22() interface{} {
	println("hello")
	r := newRes()
	return struct {
		v *res
	}{r}
}

func test23() interface{} {
	println("hello")
	r := newRes()
	return struct {
		v interface{}
	}{r}
}

func test24() []*res {
	r := newRes()
	return []*res{r}
}

func test25() map[int]*res {
	r := newRes()
	return map[int]*res{0: r}
}

func test26() []*res {
	r := newRes()
	s := []*res{nil}
	s[0] = r
	return s
}

func test27() [1]*res {
	r := newRes()
	return [...]*res{r}
}

func test28() interface{} {
	r := newRes()
	return []*res{0:r}
}

func test29() (*res, int) {
	r := newRes()
	return r, 1
}

func test30() (r *res) {
	r = newRes()
	return
}
