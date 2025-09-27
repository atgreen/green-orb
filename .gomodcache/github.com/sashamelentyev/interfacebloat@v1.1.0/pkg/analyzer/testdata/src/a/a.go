package a

type Example01 interface { // want "the interface has more than 10 methods: 11"
	a01()
	a02()
	a03()
	a04()
	a05()
	a06()
	a07()
	a08()
	a09()
	a10()
	a11()
}

func Example02() {
	var _ interface { // want "the interface has more than 10 methods: 11"
		a01()
		a02()
		a03()
		a04()
		a05()
		a06()
		a07()
		a08()
		a09()
		a10()
		a11()
	}
}

func Example03() interface { // want "the interface has more than 10 methods: 11"
	a01()
	a02()
	a03()
	a04()
	a05()
	a06()
	a07()
	a08()
	a09()
	a10()
	a11()
} {
	return nil
}

type Example04 struct {
	Foo interface { // want "the interface has more than 10 methods: 11"
		a01()
		a02()
		a03()
		a04()
		a05()
		a06()
		a07()
		a08()
		a09()
		a10()
		a11()
	}
}

type Small01 interface {
	a01()
	a02()
	a03()
	a04()
	a05()
}

type Small02 interface {
	a06()
	a07()
	a08()
	a09()
	a10()
	a11()
}

type Example05 interface {
	Small01
	Small02
}

type Example06 interface {
	interface { // want "the interface has more than 10 methods: 11"
		a01()
		a02()
		a03()
		a04()
		a05()
		a06()
		a07()
		a08()
		a09()
		a10()
		a11()
	}
}

type TypeGeneric interface {
	~uint8 | ~uint16 | ~uint32 | ~uint64 | uint |
		~int8 | ~int16 | ~int32 | ~int64 | int |
		~float32 | ~float64 |
		~string
}

func ExampleNoProblem() interface {
	a01()
	a02()
	a03()
	a04()
	a05()
	a06()
	a07()
	a08()
	a09()
	a10()
} {
	return nil
}
