package ignoredecls

import (
	"bytes"
	strs "strings"
)

const I = 123

func Const() {
	// fill
	// fill
	// fill
	// fill
	// fill
	_ = I
}

func Type_Int() {
	i := 123
	// fill
	// fill
	// fill
	// fill
	// fill
	_ = i
}

func Type_Map() {
	i := map[string]string{}
	// fill
	// fill
	// fill
	// fill
	// fill
	_ = i["foo"]
}

func Type_Pointer() {
	i := &bytes.Buffer{}
	// fill
	// fill
	// fill
	// fill
	// fill
	_ = i
}

func Type_ImportAlias() {
	i := &strs.Builder{}
	// fill
	// fill
	// fill
	// fill
	// fill
	_ = i
}

func Param_Int(i int) {
	// fill
	// fill
	// fill
	// fill
	// fill
	_ = i
}

func Param_Pointer(i *bytes.Buffer) {
	// fill
	// fill
	// fill
	// fill
	// fill
	_ = i
}

func Param_ImportAlias(i *strs.Builder) {
	// fill
	// fill
	// fill
	// fill
	// fill
	_ = i
}

func Return() (i int) {
	// fill
	// fill
	// fill
	// fill
	// fill
	i = 123
	return
}
