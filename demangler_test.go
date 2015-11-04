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

// TODO substitution
func testDemangle4(t *testing.T) {
	s := "_Zrm1XS_"
	res := "operator%(X, X)"
	
	unmangled := demangle(s)
	
	if res != unmangled {
		t.Errorf("TestDemangle4 result = [%v], want [%v]",
			unmangled, res)
	}
}

func TestDemangle5(t *testing.T) {
	s := "_ZN6System5Sound4beepEv"
	res := "System::Sound::beep(void)"
	
	unmangled := demangle(s)
	
	if res != unmangled {
		t.Errorf("TestDemangle5 result = [%v], want [%v]",
			unmangled, res)
	}
}

func testDemangle6(t *testing.T) {
	s := "_ZNSt8ios_base4InitC1Ev"
	res := "std::ios_base::Init::Init()"
	
	unmangled := demangle(s)
	
	if res != unmangled {
		t.Errorf("TestDemangle6 result = [%v], want [%v]",
			unmangled, res)
	}
}
