package test

import "testing"

func testFn1(t *testing.T) {
	t.Parallel()
	call("Named1_Sub")
}

func testFn2(t *testing.T) {
	call("Named2_Sub")
}

func testFn3(t *testing.T) {
	t.Parallel()
	call("Named3_Sub")
}

func testFn4(t *testing.T) {
	call("Named4_Sub")
}

func Test_Named1(t *testing.T) { // want "Test_Named1 should call t.Parallel on the top level as well as its subtests" "Test_Named1 should use t.Cleanup"
	teardown := setup("Test_Named1")
	defer teardown()

	fn := testFn1

	t.Run("Named1_Sub1", fn)

	t.Run("Named1_Sub2", fn)
}

func Test_Named2(t *testing.T) { // want "Test_Named2's subtests should call t.Parallel"
	teardown := setup("Test_Named1")
	defer teardown()
	t.Parallel()

	fn := testFn2

	t.Run("Named2_Sub1", fn)

	t.Run("Named2_Sub2", fn)
}

func Test_Named3(t *testing.T) { // OK
	teardown := setup("Test_Named3")
	t.Cleanup(teardown)
	t.Parallel()

	fn := testFn3

	t.Run("Named3_Sub1", fn)

	t.Run("Named3_Sub2", fn)
}

func Test_Named4(t *testing.T) { // want "Test_Named4's subtests should call t.Parallel"
	teardown := setup("Test_Named4")
	t.Cleanup(teardown)
	t.Parallel()

	tests := []struct {
		name string
	}{
		{
			name: "Named4_Sub1",
		},
		{
			name: "Named4_Sub2",
		},
	}

	fn := testFn4
	for _, tt := range tests {
		t.Run(tt.name, fn)
	}
}
