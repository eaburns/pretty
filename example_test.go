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

func ExamplePrint_emptyStruct() {
	type T struct{}
	Print(T{})
	// Output: T{}
}

func ExamplePrint_unexportedStruct() {
	type T struct{ a int }
	Print(T{})
	// Output: T{…}
}

func ExamplePrint_exportedAndUnexportedStruct() {
	type T struct{ A, b int }
	Print(T{})
	// Output: T{
	//	A: 0
	// 	…
	// }
}

func ExamplePrint_Indent() {
	type T struct{ A, b int }
	orig := Indent
	Indent = "----"
	Print(T{})
	Indent = orig
	// Output: T{
	// ----A: 0
	// ----…
	// }
}

func ExamplePrint_map() {
	type T map[string]int
	Print(T{
		"a": 5,
		"b": 6,
		"α": 7,
	})
	// Output: T{
	// 	"a": 5
	// 	"b": 6
	//	"α": 7
	// }
}
