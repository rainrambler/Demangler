package main

import (
	"fmt"
	"strings"
)

// http://mentorembedded.github.io/cxx-abi/abi.html#mangling
func demangle(mangledname string) string {
	return ""
}

type Demangler struct {
	Mangled      string
	Remain       string
	CVqualifiers []byte
	RefQualifer  byte
	FuncName     string
	isCtor       bool
	isDtor       bool
}

func isNestedName(mangled string) bool {
	if len(mangled) < 1 {
		return false
	}
	// the first char
	return mangled[0] == 'N'
}

func isCVqualifier(mangled string) bool {
	if len(mangled) < 1 {
		return false
	}
	// the first char
	return (mangled[0] == 'r') ||
		(mangled[0] == 'V') || 
		(mangled[0] == 'K')
}

func isRefqualifier(mangled string) bool {
	if len(mangled) < 1 {
		return false
	}
	// the first char
	return (mangled[0] == 'R') ||
		(mangled[0] == 'O')
}

// 5.1.4.1 Virtual Tables and RTTI
// 5.1.4.2 Virtual Override Thunks
func isSpecialName(mangled string) bool {
	if len(mangled) < 2 {
		return false
	}
	
	c0 := mangled[0]
	c1 := mangled[1]
	
	if c0 != 'T' {
		return false
	}
	
	return (c1 == 'V') || (c1 == 'T') || (c1 == 'I') || (c1 == 'S') ||
		(c1 == 'h') || (c1 == 'v')
}

func isNumberChar(c byte) bool {
	return (c >= '0') && (c <= '9')
}

func (p *Demangler) unmangle(mangled string) {
	p.Mangled = mangled
	p.Remain = mangled
	
	if strings.HasPrefix(p.Remain, "_Z") {
		p.Remain = p.Remain[2:]
	} else {
		fmt.Printf("WARN: Not valid mangled name: %v\n", p.Remain)
		return
	}
	
	p.parseEncoding()
}

func (p *Demangler) parseEncoding() {
	if isSpecialName(p.Remain) {
		p.parseSpecialName()
	}
}

func (p *Demangler) parseSpecialName() {
	
}

func (p *Demangler) parseNestedName() {
	if !isNestedName(p.Remain) {
		fmt.Printf("Mangled remain not a nested name: %v\n", p.Remain)
		return
	}
	
	p.parseCvQualifiers()
	
	p.parseRefQualifier()
	
	p.parsePrefix()
}

func (p *Demangler) parseCvQualifiers() {
	for isCVqualifier(p.Remain) {
		p.CVqualifiers = append(p.CVqualifiers, p.Remain[0])
		p.Remain = p.Remain[1:]
	}
}

func (p *Demangler) parseRefQualifier() {
	if isRefqualifier(p.Remain) {
		p.RefQualifer = p.Remain[0]
		p.Remain = p.Remain[1:]
	}
}

// TODO
func (p *Demangler) parsePrefix() {
	
}

/*
<unqualified-name> ::= <operator-name>
                   ::= <ctor-dtor-name>  
                   ::= <source-name>   
                   ::= <unnamed-type-name>   
*/
func (p *Demangler) parseUnqualifiedName() {
	var res bool
	res = p.parseCtorDtorName()
	
	if !res {
		res = p.parseOperatorName()
	}
	
	if !res {
		res = p.parseUnnamedTypeName()
	}
	
	if !res {
		//s := p.parseSourceName()
		
	}
}

func (p *Demangler) parseOperatorName() bool {
	if len(p.Remain) < 2 {
		return false
	}
	
	op := p.Remain[0:2]
	if (op == "nw") {
		p.FuncName = "new"
	} else if (op == "na") {
		p.FuncName = "new[]"
	} else if (op == "dl") {
		p.FuncName = "delete"
	} else if (op == "da") {
		p.FuncName = "delete[]"
	} else if (op == "ps") {
		p.FuncName = "+ (unary)"
	} else if (op == "ng") {
		p.FuncName = "- (unary)"
	} else if (op == "ad") {
		p.FuncName = "& (unary)"
	} else if (op == "de") {
		p.FuncName = "* (unary)"
	} else if (op == "co") {
		p.FuncName = "~"
	} else if (op == "pl") {
		p.FuncName = "+"
	} else if (op == "mi") {
		p.FuncName = "-"
	} else if (op == "ml") {
		p.FuncName = "*"
	} else if (op == "dv") {
		p.FuncName = "/"
	} else if (op == "rm") {
		p.FuncName = "%"
	} else if (op == "an") {
		p.FuncName = "&"
	} else if (op == "or") {
		p.FuncName = "|"
	} else if (op == "eo") {
		p.FuncName = "^"
	} else if (op == "aS") {
		p.FuncName = "="
	} else if (op == "pL") {
		p.FuncName = "+="
	} else if (op == "mI") {
		p.FuncName = "-="
	} else if (op == "mL") {
		p.FuncName = "*="
	} else if (op == "dV") {
		p.FuncName = "/="
	} else if (op == "rM") {
		p.FuncName = "%="
	} else if (op == "aN") {
		p.FuncName = "&="
	} else if (op == "oR") {
		p.FuncName = "|="
	} else if (op == "eO") {
		p.FuncName = "^="
	} else if (op == "ls") {
		p.FuncName = "<<"
	} else if (op == "rs") {
		p.FuncName = ">>"
	} else if (op == "lS") {
		p.FuncName = "<<="
	} else if (op == "rS") {
		p.FuncName = ">>="
	} else if (op == "eq") {
		p.FuncName = "=="
	} else if (op == "ne") {
		p.FuncName = "!="
	} else if (op == "lt") {
		p.FuncName = "<"
	} else if (op == "gt") {
		p.FuncName = ">"
	} else if (op == "le") {
		p.FuncName = "<="
	} else if (op == "ge") {
		p.FuncName = ">="
	} else if (op == "nt") {
		p.FuncName = "!"
	} else if (op == "aa") {
		p.FuncName = "&&"
	} else if (op == "oo") {
		p.FuncName = "||"
	} else if (op == "pp") {
		p.FuncName = "++"
	} else if (op == "mm") {
		p.FuncName = "--"
	} else if (op == "cm") {
		p.FuncName = ","
	} else if (op == "pm") {
		p.FuncName = "->*"
	} else if (op == "pt") {
		p.FuncName = "->"
	} else if (op == "cl") {
		p.FuncName = "()"
	} else if (op == "ix") {
		p.FuncName = "[]"
	} else if (op == "qu") {
		p.FuncName = "?"
	} else if (op == "cv") {
		p.FuncName = "(cast)" // ?
	} else if (op == "li") {
		p.FuncName = ""
	} else if (op[0] == 'v') {
		p.FuncName = "" // ?
	} else {
		p.FuncName = ""
		return false
	}
	
	p.Remain = p.Remain[2:]
	return true
}

func (p *Demangler) parseCtorDtorName() bool {
	if len(p.Remain) < 2 {
		return false
	}
	
	nm := p.Remain[0:2]
	if nm == "C1" {
		p.isCtor = true
	} else if nm == "C2" {
		p.isCtor = true
	} else if nm == "C3" {
		p.isCtor = true
	} else if nm == "D0" {
		p.isDtor = true
	} else if nm == "D1" {
		p.isDtor = true
	} else if nm == "D2" {
		p.isDtor = true
	} else {
		return false
	}
	
	p.Remain = p.Remain[2:]
	return true
}

func (p *Demangler) parseUnnamedTypeName() bool {
	if len(p.Remain) < 2 {
		return false
	}
	
	nm := p.Remain[0:2]
	if nm == "Ut" {
		// TODO
		fmt.Printf("WARN: parseUnnamedTypeName not implemented! Value: %v\n", p.Remain)
	} else {
		return false
	}
	
	p.Remain = p.Remain[2:]
	return true
}

func (p *Demangler) parseTemplatePrefix() {
	
}

// <source-name> ::= <positive length number> <identifier>
func (p *Demangler) parseSourceName() string {
	if len(p.Remain) < 2 {
		return ""
	}
	
	c0 := p.Remain[0]
	c1 := p.Remain[1]
	charLen := 0
	if isNumberChar(c0) && isNumberChar(c1) {
		charLen = (int(c0 - '0') * 10) + int(c1 - '0')
		p.Remain = p.Remain[2:]
		
		return p.parseIdentifier(charLen)
	} else if isNumberChar(c0) {
		charLen = int(c0 - '0')
		p.Remain = p.Remain[1:]
		
		return p.parseIdentifier(charLen)
	} else {
		// some error
		fmt.Printf("WARN: Invalid source name: %v\n", p.Remain)
		return ""
	}
}

func (p *Demangler) parseIdentifier(size int) string {
	s := p.Remain[0:size]
	return s
}

func (p *Demangler) parseType() {
	
}

func (p *Demangler) parseBuildinType() {
	
}
