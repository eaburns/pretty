package pp

import (
	"os"
)

// Recall that if you pass a cyclic object by value then a copy is made.
// The copy is not part of the cycle.
func ExamplePrint_passValue() {
	type T struct{ X *T }
	var t T
	t.X = &t
	Print(os.Stdout, t)
	// Output: T {
	// 	X: T {
	// 		X: <cycle>
	//	}
	// }
}

// Recall that if you pass a cyclic object as a value, a copy is made.
// The copy is not part of the cycle.  But, if you pass a pointer to the
// value then the argument will be on the cycle.
func ExamplePrint_passPointer() {
	type T struct{ X *T }
	var t T
	t.X = &t
	Print(os.Stdout, &t)
	// Output: T {
	// 	X: <cycle>
	// }
}
