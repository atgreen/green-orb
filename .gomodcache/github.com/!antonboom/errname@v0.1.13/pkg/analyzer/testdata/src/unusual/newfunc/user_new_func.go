package newfunc

var FromUserNewError = new() // want "the variable name `FromUserNewError` should conform to the `ErrXxx` format"

func new() error {
	return nil
}
