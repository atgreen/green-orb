package test

import (
	"fmt"
	"testing"
)

func call(name string) {
	fmt.Println(name)
}

func setup(name string) func() {
	fmt.Printf("setup: %s\n", name)
	return func() {
		fmt.Println("clean up finished")
	}
}

func Test_Func1(t *testing.T) { // want "Test_Func1 should call t.Parallel on the top level as well as its subtests"
	teardown := setup("Test_Func1")
	t.Cleanup(teardown)

	t.Run("Func1_Sub1", func(t *testing.T) {
		call("Func1_Sub1")
		t.Parallel()
	})

	t.Run("Func1_Sub2", func(t *testing.T) {
		call("Func1_Sub2")
		t.Parallel()
	})
}

func Test_Func2(t *testing.T) { // want "Test_Func2's subtests should call t.Parallel"
	teardown := setup("Test_Func2")
	t.Cleanup(teardown)
	t.Parallel()

	t.Run("Func2_Sub1", func(t *testing.T) {
		call("Func2_Sub1")
	})

	t.Run("Func2_Sub2", func(t *testing.T) {
		call("Func2_Sub2")
	})
}

func Test_Func3(t *testing.T) { // OK
	teardown := setup("Test_Func3")
	t.Cleanup(teardown)
	t.Parallel()

	t.Run("Func3_Sub1", func(t *testing.T) {
		call("Func3_Sub1")
		t.Parallel()
	})

	t.Run("Func3_Sub2", func(t *testing.T) {
		call("Func3_Sub2")
		t.Parallel()
	})
}

func Test_Func4(t *testing.T) { // OK
	teardown := setup("Test_Func4")
	defer teardown()
}
