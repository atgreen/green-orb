package testcase

var a = "a" // want "should only use grouped global 'var' declarations"

func dummy() { var _ = "ignore"; println(a) }
