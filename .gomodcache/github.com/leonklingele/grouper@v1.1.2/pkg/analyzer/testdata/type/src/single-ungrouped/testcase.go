package testcase

type a string // want "should only use grouped global 'type' declarations"

func dummy() { type _ string }
