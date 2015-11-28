package main

import (
	"fmt"
	//"fmt"
)

func main() {
	s := "_Z3foo3bar"
	res := "foo(bar)"
	
	stat := 0
	unmangled := cxa_demangle(s, &stat)
	
	fmt.Printf("Should: %v, Actual: %v\n", res, unmangled)
}
