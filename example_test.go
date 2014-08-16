package pretty

import "fmt"

// Recall that if you pass a cyclic object by value then a copy is made.
// The copy is not part of the cycle.
func ExamplePrint_passValue() {
	type T struct{ X *T }
	var t T
	t.X = &t
	Print(t)
	// Output: T{
	// 	X: T{
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
	Print(&t)
	// Output: T{
	// 	X: <cycle>
	// }
}

type prettyPrinter struct {
	x, y, z int
}

func (p prettyPrinter) PrettyPrint() string {
	return fmt.Sprintf("<%d, %d, %d>", p.x, p.y, p.z)
}

func ExamplePrint_prettyPrinter() {
	Print(prettyPrinter{5, 6, 7})
	// Output: <5, 6, 7>
}

type T0 struct{}

func ExamplePrint_emptyStruct() {
	Print(T0{})
	// Output: T0{}
}

type T1 struct{ a int }

func ExamplePrint_unexportedStruct() {
	Print(T1{})
	// Output: T1{…}
}

type T2 struct{ A, b int }

func ExamplePrint_exportedAndUnexportedStruct() {
	Print(T2{})
	// Output: T2{
	//	A: 0
	// 	…
	// }
}

func ExamplePrint_Indent() {
	orig := Indent
	Indent = "----"
	Print(T2{})
	Indent = orig
	// Output: T2{
	// ----A: 0
	// ----…
	// }
}
