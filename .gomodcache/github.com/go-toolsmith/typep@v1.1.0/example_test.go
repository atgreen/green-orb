package typep_test

import (
	"fmt"
	"go/types"

	"github.com/go-toolsmith/typep"
)

func Example() {
	floatTyp := types.Typ[types.Float32]
	intTyp := types.Typ[types.Int]
	ptr := types.NewPointer(intTyp)
	arr := types.NewArray(intTyp, 64)

	fmt.Println(typep.HasFloatProp(floatTyp)) // => true
	fmt.Println(typep.HasFloatProp(intTyp))   // => false
	fmt.Println(typep.IsPointer(ptr))         // => true
	fmt.Println(typep.IsArray(arr))           // => true

	// Output:
	// true
	// false
	// true
	// true
}
