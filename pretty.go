// Package pretty provides a pretty-printer for Go types. It produces a
// lightweight, Go-syntax-like output. It elides some type information
// and syntactic details. The intent is to show a data structure, such
// as an abstract syntax tree, without much clutter.
package pretty

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"reflect"
	"strconv"
)

// A PrettyPrinter implements the PrettyPrint method.
type PrettyPrinter interface {
	// PrettyPrint returns a string, overriding the output of Print.
	PrettyPrint() string
}

// Fprint pretty-prints a value to the given writer.
// If a type implementing PrettyPrinter is encountered, its PrettyPrint
// method is used to print it. Print prunes cycles.
//
// Recall that if you pass a cyclic object as a
// value, a copy is made. The copy is not part of the cycle.
func Fprint(out io.Writer, v interface{}) (err error) {
	defer func() {
		if r := recover(); r == nil {
			return
		} else if e, ok := r.(error); ok {
			err = e
		} else {
			panic(err)
		}
	}()
	print(out, make(map[reflect.Value]bool), "\n", reflect.ValueOf(v))
	return err
}

// Print pretty-prints a value to os.Stdout.
func Print(v interface{}) error {
	return Fprint(os.Stdout, v)
}

// String pretty-prints a value, returning it as a string.
func String(v interface{}) string {
	buf := bytes.NewBuffer(nil)
	if err := Fprint(buf, v); err != nil {
		panic(err)
	}
	return buf.String()
}

func print(out io.Writer, path map[reflect.Value]bool, indent string, v reflect.Value) {
	if !v.IsValid() {
		pr(out, "nil")
		return
	}
	if path[v] {
		pr(out, "<cycle>")
		return
	}
	path[v] = true
	defer func() { path[v] = false }()
	if pper, ok := v.Interface().(PrettyPrinter); ok {
		pr(out, "%s", pper.PrettyPrint())
		return
	}
	switch v.Kind() {
	case reflect.Bool:
		pr(out, "%t", v.Bool())

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		pr(out, "%d", v.Int())

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		pr(out, "%d", v.Uint())

	case reflect.Float32, reflect.Float64:
		pr(out, "%f", v.Float())

	case reflect.Complex64, reflect.Complex128:
		pr(out, "%f", v.Complex())

	case reflect.Array, reflect.Slice:
		pr(out, "[")
		indent2 := indent + "\t"
		for i := 0; i < v.Len(); i++ {
			pr(out, indent2)
			print(out, path, indent2, v.Index(i))
		}
		pr(out, indent+"]")

	case reflect.Interface, reflect.Ptr:
		if v.IsNil() {
			pr(out, "nil")
		} else {
			print(out, path, indent, v.Elem())
		}

	case reflect.String:
		pr(out, strconv.Quote(v.String()))

	case reflect.Struct:
		printStruct(out, path, indent, v)

	case reflect.Chan:
		pr(out, "<chan>")
	case reflect.Func:
		pr(out, "<function>")
	case reflect.Map:
		pr(out, "<map>")
	case reflect.UnsafePointer:
		pr(out, "<unsafe pointer>")
	case reflect.Invalid:
		pr(out, "<invalid>")
	}
}

func printStruct(out io.Writer, path map[reflect.Value]bool, indent string, v reflect.Value) {
	t := v.Type()
	pr(out, "%s {", t.Name())
	indent2 := indent + "\t"

	var u, e bool
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if !exported(&f) {
			u = true
			continue
		}
		e = true
		pr(out, "%s%s: ", indent2, f.Name)
		print(out, path, indent2, v.Field(i))
	}
	if !e {
		// No exported fields, so don't put '}' on a new line.
		indent = ""
		indent2 = ""
	}
	if u {
		pr(out, "%s…", indent2)
	}
	pr(out, "%s}", indent)
}

func pr(out io.Writer, f string, args ...interface{}) {
	if _, err := fmt.Fprintf(out, f, args...); err != nil {
		panic(err)
	}
}

func exported(f *reflect.StructField) bool {
	return len(f.PkgPath) == 0
}
