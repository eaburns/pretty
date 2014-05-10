//Package pp provides a pretty-printer for Go types. It produces a
// lightweight, Go-syntax-like output.  It elides some type information
// and syntactic details.  The intent is to show a data structure, such
// as an abstract syntax tree, without much clutter.  It also supports
// printing to the dot language of graphviz!
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
// Print prunes cycles.  Recall that if you pass a cyclic object as a
// value, a copy is made.  The copy is not part of the cycle.
func Print(out io.Writer, v interface{}) (err error) {
	defer recoverErr(&err)
	print(out, make(map[reflect.Value]bool), "\n", reflect.ValueOf(v))
	return err
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
	if strer, ok := v.Interface().(fmt.Stringer); ok {
		pr(out, "%s", strer)
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
			print(out, path, indent2, v.Field(i))
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

// Dot writes value to the writer using the dot language of graphviz.
func Dot(out io.Writer, v interface{}) (err error) {
	defer recoverErr(&err)
	if _, err = io.WriteString(out, "digraph {\n"); err != nil {
		return err
	}
	dot(out, make(map[reflect.Value]int), 0, reflect.ValueOf(v))
	_, err = io.WriteString(out, "}")
	return err
}

func dot(out io.Writer, seen map[reflect.Value]int, n int, v reflect.Value) (nd, next int) {
	if m, ok := seen[v]; ok {
		return m, n
	}
	defer func() { seen[v] = nd }()
	if strer, ok := v.Interface().(fmt.Stringer); ok {
		return node(out, n, "%s", strer)
	}
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
		_, next := node(out, n, "%s[]", v.Type().Elem().Name())
		// Must see 'v' before recurring on dot.
		seen[v] = n
		for i := 0; i < v.Len(); i++ {
			var m int
			m, next = dot(out, seen, next, v.Index(i))
			arc(out, n, m, "")
		}
		return n, next

	case reflect.Interface, reflect.Ptr:
		if v.IsNil() {
			return node(out, n, "nil")
		}
		return dot(out, seen, n, v.Elem())

	case reflect.String:
		return node(out, n, "%s", strconv.Quote(v.String()))

	case reflect.Struct:
		t := v.Type()
		_, next := node(out, n, t.Name())
		// Must see 'v' before recurring on dot.
		seen[v] = n
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			if !exported(f.Name) {
				// Don't output unexported fields.
				continue
			}
			var m int
			m, next = dot(out, seen, next, v.Field(i))
			arc(out, n, m, f.Name)
		}
		return n, next

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

func node(out io.Writer, n int, f string, args ...interface{}) (nd, next int) {
	s := fmt.Sprintf(f, args...)
	_, err := fmt.Fprintf(out, "\tn%d [label=%s]\n", n, strconv.Quote(s))
	if err != nil {
		panic(err)
	}
	return n, n + 1
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

func recoverErr(err *error) {
	if r := recover(); r == nil {
		return
	} else if e, ok := r.(error); ok {
		*err = e
	} else {
		panic(*err)
	}
}
