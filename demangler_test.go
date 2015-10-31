package main

import (
	"testing"
)

func TestDemangle1(t *testing.T) {
	s := "_Z1fv"
	res := "f(void)"
	
	unmangled := demangle(s)
	
	if res != unmangled {
		t.Errorf("TestDemangle1 result = [%v], want [%v]",
			unmangled, res)
	}
}

func TestDemangle2(t *testing.T) {
	s := "_Z1fi"
	res := "f(int)"
	
	unmangled := demangle(s)
	
	if res != unmangled {
		t.Errorf("TestDemangle2 result = [%v], want [%v]",
			unmangled, res)
	}
}

func TestDemangle3(t *testing.T) {
	s := "_Z3foo3bar"
	res := "foo(bar)"
	
	unmangled := demangle(s)
	
	if res != unmangled {
		t.Errorf("TestDemangle3 result = [%v], want [%v]",
			unmangled, res)
	}
}
