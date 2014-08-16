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

func ExamplePrint_unexportedStructFields() {
	type T struct{ a int }
	Print(T{})
	// Output: T{…}
}

func ExamplePrint_exportedAndUnexportedStructFields() {
	type T struct{ A, b int }
	Print(T{})
	// Output: T{
	//	A: 0
	// 	…
	// }
}

func ExamplePrint_invalid() {
	Print(nil)
	// Output: nil
}

func ExamplePrint_nil() {
	var i *int
	Print(i)
	// Output: nil
}

func ExamplePrint_complex() {
	Print(3 + 5i)
	// Output: (3.000000+5.000000i)
}

func ExamplePrint_boolMap() {
	type T map[bool]int
	Print(T{
		true:  5,
		false: 6,
	})
	// Output: T{
	// 	false: 6
	// 	true: 5
	// }
}

func ExamplePrint_intMap() {
	type T map[int]int
	Print(T{
		4: 7,
		1: 5,
		2: 6,
	})
	// Output: T{
	// 	1: 5
	// 	2: 6
	// 	4: 7
	// }
}

func ExamplePrint_uintMap() {
	type T map[uint]int
	Print(T{
		4: 7,
		1: 5,
		2: 6,
	})
	// Output: T{
	// 	1: 5
	// 	2: 6
	// 	4: 7
	// }
}

func ExamplePrint_floatMap() {
	type T map[float32]int
	Print(T{
		4.3: 7,
		1.1: 5,
		1.2: 6,
	})
	// Output: T{
	// 	1.100000: 5
	// 	1.200000: 6
	// 	4.300000: 7
	// }
}

func ExamplePrint_stringMap() {
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

func ExamplePrint_array() {
	type T [5]int
	Print(T{5, 6, 7, 8, 9})
	// Output: [
	// 	5
	// 	6
	//	7
	//	8
	//	9
	// ]
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
