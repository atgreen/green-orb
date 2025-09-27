package customDecOrderAll

func ba() {}

var bc = 1

const bb = 1 // want "const must not be placed after var \\(desired order: const,var\\)"

type bd int
