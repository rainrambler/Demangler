package main

import (
	"fmt"
	//"fmt"
)

func demanglefunc(s, want string) {
	stat := 0
	unmangled := cxa_demangle(s, &stat)
	fmt.Printf("Should: %v, Actual: %v\n", want, unmangled)
}

func main() {
	//demanglefunc("_ZNSt8ios_base4InitC1Ev", "std::ios_base::Init::Init()")
	demanglefunc("_ZplR1XS0_", "operator+(X&, X&)")
}
