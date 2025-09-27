package main

import (
	"fmt"
	"io"
	"log"
	"testing"
)

func TestStd(t *testing.T) {
	var a = []any{1, 2, 3}

	fmt.Println(a)      // want `pass \[\]any as any to func fmt.Println func\(a \.\.\.any\) \(n int, err error\)`
	fmt.Printf("%v", a) // want `pass \[\]any as any to func fmt.Printf func\(format string, a \.\.\.any\) \(n int, err error\)`
	fmt.Print(a)        // want `pass \[\]any as any to func fmt.Print func\(a \.\.\.any\) \(n int, err error\)`

	fmt.Sprintln(a)      // want `pass \[\]any as any to func fmt.Sprintln func\(a \.\.\.any\) string`
	fmt.Sprintf("%v", a) // want `pass \[\]any as any to func fmt.Sprintf func\(format string, a \.\.\.any\) string`
	fmt.Sprint(a)        // want `pass \[\]any as any to func fmt.Sprint func\(a \.\.\.any\) string`

	var w io.Writer
	fmt.Fprintln(w, a)      // want `pass \[\]any as any to func fmt.Fprintln func\(w io.Writer, a \.\.\.any\) \(n int, err error\)`
	fmt.Fprintf(w, "%v", a) // want `pass \[\]any as any to func fmt.Fprintf func\(w io.Writer, format string, a \.\.\.any\) \(n int, err error\)`
	fmt.Fprint(w, a)        // want `pass \[\]any as any to func fmt.Fprint func\(w io.Writer, a \.\.\.any\) \(n int, err error\)`

	log.Println(a)      // want `pass \[\]any as any to func log.Println func\(v \.\.\.any\)`
	log.Printf("%v", a) // want `pass \[\]any as any to func log.Printf func\(format string, v \.\.\.any\)`
	log.Print(a)        // want `pass \[\]any as any to func log.Print func\(v \.\.\.any\)`

	log.Fatalln(a)      // want `pass \[\]any as any to func log.Fatalln func\(v \.\.\.any\)`
	log.Fatalf("%v", a) // want `pass \[\]any as any to func log.Fatalf func\(format string, v \.\.\.any\)`
	log.Fatal(a)        // want `pass \[\]any as any to func log.Fatal func\(v \.\.\.any\)`

	log.Panicln(a)      // want `pass \[\]any as any to func log.Panicln func\(v \.\.\.any\)`
	log.Panicf("%v", a) // want `pass \[\]any as any to func log.Panicf func\(format string, v \.\.\.any\)`
	log.Panic(a)        // want `pass \[\]any as any to func log.Panic func\(v \.\.\.any\)`
}
