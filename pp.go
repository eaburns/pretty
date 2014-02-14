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
	"unicode"
	"unicode/utf8"
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

	case reflect.Array, reflect.Slice:
		pr(out, "[")
		indent2 := indent + "\t"
		for i := 0; i < v.Len(); i++ {
			pr(out, indent2)
			print(out, indent2, v.Index(i))
		}
		pr(out, indent+"]")

	case reflect.Interface, reflect.Ptr:
		if v.IsNil() {
			pr(out, "nil")
		} else {
			print(out, indent, v.Elem())
		}

	case reflect.String:
		pr(out, strconv.Quote(v.String()))

	case reflect.Struct:
		t := v.Type()
		pr(out, "%s {", t.Name())
		indent2 := indent + "\t"
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			if !exported(f.Name) {
				// Don't output unexported fields.
				continue
			}
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
	case reflect.Invalid:
		pr(out, "<invalid>")
	}
}

func pr(out io.Writer, f string, args ...interface{}) {
	if _, err := fmt.Fprintf(out, f, args...); err != nil {
		panic(err)
	}
}

// Print the value to the writer using the dot language of graphviz.
func Dot(out io.Writer, v interface{}) (err error) {
	defer func() {	
		if r := recover(); r == nil {
			return
		} else if e, ok := r.(error); ok {
			err = e
		}
	}()
	if _, err = io.WriteString(out, "digraph {\n"); err != nil {
		return err
	}
	dot(out, 0, reflect.ValueOf(v))
	_, err = io.WriteString(out, "}")
	return err
}

func dot(out io.Writer, n int, v reflect.Value) int {
	switch v.Kind() {
	case reflect.Bool:
		return node(out, n, "%t", v.Bool())

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return node(out, n, "%d", v.Int())

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return node(out, n, "%d", v.Uint())

	case reflect.Float32, reflect.Float64:
		return node(out, n, "%f", v.Float())

	case reflect.Complex64, reflect.Complex128:
		return node(out, n, "%f", v.Complex())

	case reflect.Array, reflect.Slice:
		m := node(out, n, "%s[]", v.Type().Elem().Name())
		for i := 0; i < v.Len(); i++ {
			arc(out, n, m, "")
			m = dot(out, m, v.Index(i))
		}
		return m

	case reflect.Interface, reflect.Ptr:
		if v.IsNil() {
			return node(out, n, "nil")
		}
		return dot(out, n, v.Elem())

	case reflect.String:
		return node(out, n, "%s", strconv.Quote(v.String()))

	case reflect.Struct:
		t := v.Type()
		m := node(out, n, t.Name())
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			if !exported(f.Name) {
				// Don't output unexported fields.
				continue
			}
			arc(out, n, m, f.Name)
			m = dot(out, m, v.Field(i))
		}
		return m

	case reflect.Chan:
		return node(out, n, "<chan>")
	case reflect.Func:
		return node(out, n, "<function>")
	case reflect.Map:
		return node(out, n, "<map>")
	case reflect.UnsafePointer:
		return node(out, n, "<unsafe pointer>")
	case reflect.Invalid:
		return node(out, n, "<invalid>")
	}
	panic("unreachable: " + v.Kind().String())
}

func node(out io.Writer, n int, f string, args ...interface{}) int {
	s := fmt.Sprintf(f, args...)
	_, err := fmt.Fprintf(out, "\tn%d [label=%s]\n", n, strconv.Quote(s))
	if err != nil {
		panic(err)
	}
	return n+1
}

func arc(out io.Writer, src, dst int, label string) {
	if label == "" {
		_, err := fmt.Fprintf(out, "\tn%d -> n%d\n", src, dst)
		if err != nil {
			panic(err)
		}
		return
	}
	_, err := fmt.Fprintf(out, "\tn%d -> n%d [label=%s]\n", src, dst, strconv.Quote(label))
	if err != nil {
		panic(err)
	}
}

func exported(n string) bool {
	if len(n) == 0 {
		return true
	}
	r, _ := utf8.DecodeRuneInString(n)
	return unicode.IsUpper(r)
}