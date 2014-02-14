// Package pp provides a pretty-printer for Go types.  It produces a
// lightweight, Go-syntax-like output, but it is not intended to produce
// valid Go syntax.  It elides a bunch of type information, leaving only
// struct names.  It omits commas, and possibly much much more.
package pp

import (
	"fmt"
	"io"
	"reflect"
	"strconv"
)

// Print pretty-prints the value to the given writer.
func Print(out io.Writer, v interface{}) (err error) {
	defer func() {	
		if r := recover(); r == nil {
			return
		} else if e, ok := r.(error); ok {
			err = e
		}
	}()
	print(out, "\n", reflect.ValueOf(v))
	return err
}

func print(out io.Writer, indent string, v reflect.Value) {
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

	case reflect.Array:
	case reflect.Slice:
		pr(out, "[")
		indent2 := indent + "\t"
		for i := 0; i < v.Len(); i++ {
			pr(out, indent2)
			print(out, indent2, v.Index(i))
		}
		pr(out, indent+"]")

	case reflect.Interface, reflect.Ptr:
		print(out, indent, v.Elem())

	case reflect.String:
		pr(out, strconv.Quote(v.String()))

	case reflect.Struct:
		t := v.Type()
		pr(out, "%s {", t.Name())
		indent2 := indent + "\t"
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			pr(out, "%s%s: ", indent2, f.Name)
			print(out, indent2, v.Field(i))
		}
		pr(out, "%s}", indent)

	case reflect.Chan:
		pr(out, "<chan>")
	case reflect.Func:
		pr(out, "<function>")
	case reflect.Map:
		pr(out, "<map>")
	case reflect.UnsafePointer:
		pr(out, "<unsafe pointer>")
	}
}

func pr(out io.Writer, f string, args ...interface{}) {
	if _, err := fmt.Fprintf(out, f, args...); err != nil {
		panic(err)
	}
}