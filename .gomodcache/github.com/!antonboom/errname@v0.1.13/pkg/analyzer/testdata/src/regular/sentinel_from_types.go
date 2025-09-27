package regular

import "net"

var (
	InvalidAddrError       = new(net.AddrError) // unsupported
	InvalidAddrErr   error = new(net.AddrError) // want "the variable name `InvalidAddrErr` should conform to the `ErrXxx` format"
	NotErr                 = new(NotErrorType)

	Aa = new(someTypeWithPtr) // want "the variable name `Aa` should conform to the `ErrXxx` format"
	Bb = someTypeWithoutPtr{} // want "the variable name `Bb` should conform to the `ErrXxx` format"

	cC error = new(someTypeWithPtr) // want "the variable name `cC` should conform to the `errXxx` format"
	dD error = someTypeWithoutPtr{} // want "the variable name `dD` should conform to the `errXxx` format"
)
