package testcase

const a = "a" // want "should only use grouped global 'const' declarations"

func dummy() { const _ = "ignore" }
