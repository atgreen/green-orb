package regular

import "strings"

type ValidationErrors []string

func (ve ValidationErrors) Error() string { return strings.Join(ve, "\n") }

type validationErrors []string

func (ve validationErrors) Error() string { return strings.Join(ve, "\n") }

type TenErrors [10]string

func (te TenErrors) Error() string { return strings.Join(te[:], "\n") }

type tenErrors [10]string

func (te tenErrors) Error() string { return strings.Join(te[:], "\n") }

type MultiError []error             // want "the type name `MultiError` should conform to the `XxxErrors` format"
func (me MultiError) Error() string { return "" }

type multiError []error             // want "the type name `multiError` should conform to the `xxxErrors` format"
func (me multiError) Error() string { return "" }

type TwoError [2]error            // want "the type name `TwoError` should conform to the `XxxErrors` format"
func (te TwoError) Error() string { return te[0].Error() + "\n" + te[1].Error() }

type twoError [2]error            // want "the type name `twoError` should conform to the `xxxErrors` format"
func (te twoError) Error() string { return te[0].Error() + "\n" + te[1].Error() }
