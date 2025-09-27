package generics

type NotErrorGeneric[T float64 | int] struct {
	Limit T
}

type SomeError[T ~string] struct{ Code T }

func (e SomeError[T]) Error() string { return string(e.Code) }

type SomePtrError[T ~string] struct{ Code T }

func (e *SomePtrError[T]) Error() string { return string(e.Code) }

type someErr[T ~string] struct{ Code T } //  want "the type name `someErr` should conform to the `xxxError` format"
func (e someErr[T]) Error() string       { return string(e.Code) }

type SomePtrErr[T ~string] struct{ Code T } //  want "the type name `SomePtrErr` should conform to the `XxxError` format"
func (e *SomePtrErr[T]) Error() string      { return string(e.Code) }

var (
	ErrOK = &SomePtrError[string]{Code: "200 OK"}
	okErr = &SomePtrError[string]{Code: "200 OK"} // want "the variable name `okErr` should conform to the `errXxx` format"

	ErrNotFound = SomeError[string]{Code: "Not Found"}
	NotFoundErr = SomeError[string]{Code: "Not Found"} // want "the variable name `NotFoundErr` should conform to the `ErrXxx` format"

	statusCodeError       = new(SomePtrError[string]) // want "the variable name `statusCodeError` should conform to the `errXxx` format"
	ExplicitError   error = new(SomePtrError[string]) // want "the variable name `ExplicitError` should conform to the `ErrXxx` format"
)
