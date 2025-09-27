package conventional

import (
	"context"
	"testing"
)

// TODO: add *testing.F once we switch to Go 1.18

func Variable() {
	var (
		ctx context.Context
		b   *testing.B
		m   *testing.M
		pb  *testing.PB
		t   *testing.T
		tb  testing.TB
	)
	// fill
	// fill
	// fill
	// fill
	// fill
	_ = ctx
	_ = b
	_ = m
	_ = pb
	_ = t
	_ = tb
}

func Param(ctx context.Context,
	b *testing.B,
	m *testing.M,
	pb *testing.PB,
	t *testing.T,
	tb testing.TB) {
	// fill
	// fill
	// fill
	// fill
	// fill
	_ = ctx
	_ = b
	_ = m
	_ = pb
	_ = t
	_ = tb
}
