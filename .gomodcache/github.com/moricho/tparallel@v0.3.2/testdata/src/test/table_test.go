package test

import (
	"testing"
)

func Test_Table1(t *testing.T) { // want "Test_Table1 should call t.Parallel on the top level as well as its subtests" "Test_Table1 should use t.Cleanup"
	teardown := setup("Test_Table1")
	defer teardown()

	tests := []struct {
		name string
	}{
		{
			name: "Table1_Sub1",
		},
		{
			name: "Table1_Sub2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			call(tt.name)
		})
	}
}

func Test_Table2(t *testing.T) { // want "Test_Table2's subtests should call t.Parallel"
	teardown := setup("Test_Table2")
	t.Cleanup(teardown)

	t.Parallel()

	tests := []struct {
		name string
	}{
		{
			name: "Table2_Sub1",
		},
		{
			name: "Table2_Sub2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			call(tt.name)
		})
	}
}

func Test_Table3(t *testing.T) { // OK
	teardown := setup("Test_Table3")
	t.Cleanup(teardown)

	t.Parallel()

	tests := []struct {
		name string
	}{
		{
			name: "Table3_Sub1",
		},
		{
			name: "Table3_Sub2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			call(tt.name)
		})
	}
}
