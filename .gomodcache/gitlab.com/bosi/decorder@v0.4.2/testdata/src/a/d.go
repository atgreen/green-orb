package a

var da = 1

const db = 1 // want "const must not be placed after var"

type dd int // want "type must not be placed after const"

func init() {}

func init() {}

func dg() {}
