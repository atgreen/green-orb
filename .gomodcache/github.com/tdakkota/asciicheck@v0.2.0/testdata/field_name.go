package testdata

type Field struct{}

type JustStruct struct {
	Tеst Field // want `identifier "Tеst" contain non-ASCII character: U\+0435 'е'`
}
