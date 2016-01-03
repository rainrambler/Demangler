package main

import (
	"fmt"
	"os"
)

func demanglefunc(s, want string) {
	stat := 0
	unmangled := cxa_demangle(s, &stat)
	fmt.Printf("Should: %v, Actual: %v\n", want, unmangled)
}

func main() {
	argsWithProg := os.Args
	if len(argsWithProg) < 2 {
		fmt.Println("Usage: Demangler [MangledName]")
		fmt.Println("eg: demangler _ZNSt8ios_base4InitC1Ev")
		return
	}
	
	mangledname := argsWithProg[1]
	stat := 0
	unmangled := cxa_demangle(mangledname, &stat)
	fmt.Println(unmangled)
}
