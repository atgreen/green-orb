package warnings

import (
	"context"
	"testing"
)

// TODO: add *testing.F once we switch to Go 1.18

func Variable_Conventional() {
	var (
		c2 context.Context // want `variable name 'c2' is too short for the scope of its usage`
		b2 *testing.B      // want `variable name 'b2' is too short for the scope of its usage`
		m2 *testing.M      // want `variable name 'm2' is too short for the scope of its usage`
		p2 *testing.PB     // want `variable name 'p2' is too short for the scope of its usage`
		t2 *testing.T      // want `variable name 't2' is too short for the scope of its usage`
		t3 testing.TB      // want `variable name 't3' is too short for the scope of its usage`
	)
	// fill
	// fill
	// fill
	// fill
	// fill
	_ = c2
	_ = b2
	_ = m2
	_ = p2
	_ = t2
	_ = t3
}

func Param_Conventional(c2 context.Context, // want `parameter name 'c2' is too short for the scope of its usage`
	b2 *testing.B, // want `parameter name 'b2' is too short for the scope of its usage`
	m2 *testing.M, // want `parameter name 'm2' is too short for the scope of its usage`
	p2 *testing.PB, // want `parameter name 'p2' is too short for the scope of its usage`
	t2 *testing.T, // want `parameter name 't2' is too short for the scope of its usage`
	t3 testing.TB) { // want `parameter name 't3' is too short for the scope of its usage`
	// fill
	// fill
	// fill
	// fill
	// fill
	_ = c2
	_ = b2
	_ = m2
	_ = p2
	_ = t2
	_ = t3
}
