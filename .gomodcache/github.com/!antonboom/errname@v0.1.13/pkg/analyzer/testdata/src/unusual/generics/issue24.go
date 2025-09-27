package generics

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
)

type (
	Req  any
	Resp any
)

type timeoutErr[REQ Req, RESP Resp] struct { //  want "the type name `timeoutErr` should conform to the `xxxError` format"
	err     error
	sending bool
}

func (e *timeoutErr[REQ, RESP]) Error() string {
	var req REQ
	var resp RESP

	direction := "sending"
	if !e.sending {
		direction = "receiving"
	}

	return fmt.Sprintf("deferred call %T->%T timeout %s: %s",
		reflect.TypeOf(req), reflect.TypeOf(resp), direction, e.err.Error())
}

func (e *timeoutErr[REQ, RESP]) Unwrap() error {
	return e.err
}

type TimeoutError[REQ Req, RESP Resp] struct{} //
func (TimeoutError[REQ, RESP]) Error() string  { return "timeouted" }

type ValErr[A, B, C, D, E, F any] struct{}     //  want "the type name `ValErr` should conform to the `XxxError` format"
func (ValErr[A, B, C, D, E, F]) Error() string { return "boom!" }

var (
	ErrTimeout error = &timeoutErr[*http.Request, *http.Response]{err: context.DeadlineExceeded, sending: false}
	tErr       error = &timeoutErr[string, string]{err: context.DeadlineExceeded, sending: true} // want "the variable name `tErr` should conform to the `errXxx` format"
)
