package disabledDecNumCheck

func ca() {}

type cc int // want "type must not be placed after func \\(desired order: type,const,var,func\\)"

const cd = 1 // want "const must not be placed after func \\(desired order: type,const,var,func\\)"

var ce = 1 // want "var must not be placed after func \\(desired order: type,const,var,func\\)"
