package test

import "testing"

func tRun1(t *testing.T) {
	t.Parallel()

	t.Run("Named_tRun1_Sub1", func(t *testing.T) {
		call("Named_tRun1_Sub")
	})

	t.Run("Named_tRun1_Sub2", func(t *testing.T) {
		call("Named_tRun1_Sub2")
	})
}

func tRun2(t *testing.T) {
	t.Run("Named_tRun2_Sub", func(t *testing.T) {
		call("Named_tRun2_Sub")
	})
}

func tRun3(t *testing.T) {
	t.Run("Named_tRun3_Sub", func(t *testing.T) {
		t.Parallel()
		call("Named_tRun3_Sub")
	})
}

func tRun4(t *testing.T) {
	t.Run("Named_tRun4_Sub", func(t *testing.T) {
		t.Parallel()
		call("Named_tRun4_Sub")
	})
}

func tRun5(t *testing.T) {
	t.Parallel()

	t.Run("Named_tRun5_Sub1", func(t *testing.T) {
		t.Parallel()
		call("Named_tRun5_Sub1")
	})

	t.Run("Named_tRun5_Sub2", func(t *testing.T) {
		t.Parallel()
		call("Named_tRun5_Sub2")
	})
}

func tRun6(t *testing.T) {
	t.Run("Named_tRun6_Sub1", func(t *testing.T) {
		t.Parallel()
		call("Named_tRun6_Sub1")
	})
}

func Test_Named_tRun1(t *testing.T) { // want "Test_Named_tRun1's subtests should call t.Parallel"
	call("Test_Named_tRun1")

	tRun1(t)
}

func Test_Named_tRun2(t *testing.T) { // want "Test_Named_tRun2's subtests should call t.Parallel"
	t.Parallel()

	call("Test_Named_tRun2")

	tRun2(t)
}

func Test_Named_tRun3(t *testing.T) { // want "Test_Named_tRun3 should call t.Parallel on the top level as well as its subtests"
	call("Test_Named_tRun3")

	tRun3(t)
}

func Test_Named_tRun4(t *testing.T) { // OK
	t.Parallel()

	call("Test_Named_tRun4")

	tRun4(t)
}

func Test_Named_tRun5(t *testing.T) { // OK
	call("Test_Named_tRun5")

	tRun5(t)
}

func Test_Named_tRun6(t *testing.T) { // OK
	t.Parallel()
	call("Test_Named_tRun6")

	tRun6(t)

	t.Run("Named_tRun6_Sub2", func(t *testing.T) {
		t.Parallel()
		call("Named_tRun6_Sub2")
	})

}
