package main

import (
	"testing"
)

func TestClangDemangle1(t *testing.T) {
	s := "_Z1fv"
	res := "f()"
	
	stat := 0
	unmangled := cxa_demangle(s, &stat)
	
	if res != unmangled {
		t.Errorf("TestDemangle1 result = [%v], want [%v]",
			unmangled, res)
	}
}

func TestClangDemangle2(t *testing.T) {
	s := "_Z1fi"
	res := "f(int)"
	
	stat := 0
	unmangled := cxa_demangle(s, &stat)
	
	if res != unmangled {
		t.Errorf("TestDemangle2 result = [%v], want [%v]",
			unmangled, res)
	}
}

func TestClangDemangle3(t *testing.T) {
	s := "_Z3foo3bar"
	res := "foo(bar)"
	
	stat := 0
	unmangled := cxa_demangle(s, &stat)
	
	if res != unmangled {
		t.Errorf("TestDemangle3 result = [%v], want [%v]",
			unmangled, res)
	}
}

// substitution
func TestClangDemangle4(t *testing.T) {
	s := "_Zrm1XS_"
	res := "operator%(X, X)"
	
	stat := 0
	unmangled := cxa_demangle(s, &stat)
	
	if res != unmangled {
		t.Errorf("TestDemangle4 result = [%v], want [%v]",
			unmangled, res)
	}
}

func TestClangDemangle5(t *testing.T) {
	s := "_ZN6System5Sound4beepEv"
	res := "System::Sound::beep()"
	
	stat := 0
	unmangled := cxa_demangle(s, &stat)
	
	if res != unmangled {
		t.Errorf("TestDemangle5 result = [%v], want [%v]",
			unmangled, res)
	}
}

func TestClangDemangle6(t *testing.T) {
	s := "_ZNSt8ios_base4InitC1Ev"
	res := "std::ios_base::Init::Init()"
	
	stat := 0
	unmangled := cxa_demangle(s, &stat)
	
	if res != unmangled {
		t.Errorf("TestDemangle6 result = [%v], want [%v]",
			unmangled, res)
	}
}

// http://mentorembedded.github.io/cxx-abi/abi-examples.html#mangling
func TestClangDemangle7(t *testing.T) {
	s := "_ZplR1XS0_"
	res := "operator+(X&, X&)"
	
	stat := 0
	unmangled := cxa_demangle(s, &stat)
	
	if res != unmangled {
		t.Errorf("TestDemangle7 result = [%v], want [%v]",
			unmangled, res)
	}
}

// http://mentorembedded.github.io/cxx-abi/abi-examples.html#mangling
func TestClangDemangle8(t *testing.T) {
	s := "_ZN3FooIA4_iE3barE"
	res := "Foo<int [4]>::bar"
	
	stat := 0
	unmangled := cxa_demangle(s, &stat)
	
	if res != unmangled {
		t.Errorf("TestDemangle8 result = [%v], want [%v]",
			unmangled, res)
	}
}

// http://mentorembedded.github.io/cxx-abi/abi-examples.html#mangling
func TestClangDemangle9(t *testing.T) {
	s := "_ZlsRK1XS1_"
	res := "operator<<(X const&, X const&)"
	
	stat := 0
	unmangled := cxa_demangle(s, &stat)
	
	if res != unmangled {
		t.Errorf("TestDemangle9 result = [%v], want [%v]",
			unmangled, res)
	}
}

// http://mentorembedded.github.io/cxx-abi/abi-examples.html#mangling
func TestClangDemangle10(t *testing.T) {
	s := "_Z1fIiEvi"
	res := "void f<int>(int)"
	
	stat := 0
	unmangled := cxa_demangle(s, &stat)
	
	if res != unmangled {
		t.Errorf("TestDemangle10 result = [%v], want [%v]",
			unmangled, res)
	}
}

// http://mentorembedded.github.io/cxx-abi/abi-examples.html#mangling
func TestClangDemangle11(t *testing.T) {
	s := "_ZN5StackIiiE5levelE"
	res := "Stack<int, int>::level"
	
	stat := 0
	unmangled := cxa_demangle(s, &stat)
	
	if res != unmangled {
		t.Errorf("TestDemangle11 result = [%v], want [%v]",
			unmangled, res)
	}
}

// http://mentorembedded.github.io/cxx-abi/abi-examples.html#mangling
func TestClangDemangle12(t *testing.T) {
	s := "_Z3fooIiPFidEiEvv"
	res := "void foo<int, int (*)(double), int>()"
	
	stat := 0
	unmangled := cxa_demangle(s, &stat)
	
	if res != unmangled {
		t.Errorf("TestDemangle12 result = [%v], want [%v]",
			unmangled, res)
	}
}

// http://mentorembedded.github.io/cxx-abi/abi-examples.html#mangling
func TestClangDemangle13(t *testing.T) {
	s := "_ZlsRSoRKSs"
	res := "operator<<(std::ostream&, std::string const&)"
	
	stat := 0
	unmangled := cxa_demangle(s, &stat)
	
	if res != unmangled {
		t.Errorf("TestDemangle13 result = [%v], want [%v]",
			unmangled, res)
	}
}

// http://mentorembedded.github.io/cxx-abi/abi-examples.html#mangling
func TestClangDemangle14(t *testing.T) {
	s := "_ZN5StackIiiE5levelE"
	res := "Stack<int, int>::level"
	
	stat := 0
	unmangled := cxa_demangle(s, &stat)
	
	if res != unmangled {
		t.Errorf("TestDemangle14 result = [%v], want [%v]",
			unmangled, res)
	}
}

// http://mentorembedded.github.io/cxx-abi/abi-examples.html#mangling
func TestClangDemangle15(t *testing.T) {
	s := "_Z5firstI3DuoEvT_"
	res := "void first<Duo>(Duo)"
	
	stat := 0
	unmangled := cxa_demangle(s, &stat)
	
	if res != unmangled {
		t.Errorf("TestDemangle15 result = [%v], want [%v]",
			unmangled, res)
	}
}

func Test_align_up(t *testing.T) {
	var a arena
	res := 16
	val := a.align_up(15)
	if res != val {
		t.Errorf("Test_align_up result = [%v], want [%v]",
			val, res)
	}
	val = a.align_up(1)
	if res != val {
		t.Errorf("Test_align_up result = [%v], want [%v]",
			val, res)
	}
	res = 32
	val = a.align_up(17)
	if res != val {
		t.Errorf("Test_align_up result = [%v], want [%v]",
			val, res)
	}
}
