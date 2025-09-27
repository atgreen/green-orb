package regular

import "io"

var (
	ErrA = newSomeTypeWithPtr()
	ErrB = newSomeTypeWithPtr2()
	ErrC = newSomeTypeWithoutPtr()
	ErrD = newSomeTypeWithoutPtr2()
	ErrE = new(SomeTypeWithPtr)
	ErrF = &SomeTypeWithPtr{}
	ErrG = SomeTypeWithoutPtr{}

	AErr = newSomeTypeWithPtr()     // want "the variable name `AErr` should conform to the `ErrXxx` format"
	BErr = newSomeTypeWithPtr2()    // want "the variable name `BErr` should conform to the `ErrXxx` format"
	CErr = newSomeTypeWithoutPtr()  // want "the variable name `CErr` should conform to the `ErrXxx` format"
	DErr = newSomeTypeWithoutPtr2() // want "the variable name `DErr` should conform to the `ErrXxx` format"
	EErr = new(SomeTypeWithPtr)     // want "the variable name `EErr` should conform to the `ErrXxx` format"
	FErr = &SomeTypeWithPtr{}       // want "the variable name `FErr` should conform to the `ErrXxx` format"
	GErr = SomeTypeWithoutPtr{}     // want "the variable name `GErr` should conform to the `ErrXxx` format"

	AErrr error = newSomeTypeWithPtr2()   // want "the variable name `AErrr` should conform to the `ErrXxx` format"
	BErrr error = newSomeTypeWithoutPtr() // want "the variable name `BErrr` should conform to the `ErrXxx` format"
	CErrr error = new(SomeTypeWithPtr)    // want "the variable name `CErrr` should conform to the `ErrXxx` format"
	DErrr error = &SomeTypeWithPtr{}      // want "the variable name `DErrr` should conform to the `ErrXxx` format"
	EErrr error = SomeTypeWithoutPtr{}    // want "the variable name `EErrr` should conform to the `ErrXxx` format"

	ErrByAnonymousFunc = func() error { return nil }
	ByAnonymousFuncErr = func() error { return io.EOF }() // want "the variable name `ByAnonymousFuncErr` should conform to the `ErrXxx` format"
)

var (
	InitializedLaterError            error              // want "the variable name `InitializedLaterError` should conform to the `ErrXxx` format"
	InitializedLaterImplicitError    SomeTypeWithoutPtr // want "the variable name `InitializedLaterImplicitError` should conform to the `ErrXxx` format"
	InitializedLaterImplicitPtrError *SomeTypeWithPtr   // want "the variable name `InitializedLaterImplicitPtrError` should conform to the `ErrXxx` format"
)

func init() {
	InitializedLaterError = newSomeTypeWithPtr()
	InitializedLaterImplicitError = newSomeTypeWithoutPtr()
	InitializedLaterImplicitPtrError = newSomeTypeWithPtr2()
}

func newSomeTypeWithPtr() error {
	return new(SomeTypeWithPtr)
}

func newSomeTypeWithPtr2() *SomeTypeWithPtr {
	return nil
}

func newSomeTypeWithoutPtr() SomeTypeWithoutPtr {
	return SomeTypeWithoutPtr{}
}

func newSomeTypeWithoutPtr2() error {
	return SomeTypeWithoutPtr{}
}
