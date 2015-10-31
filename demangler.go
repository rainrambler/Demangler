package main

import (
	"fmt"
	"strings"
)

// http://mentorembedded.github.io/cxx-abi/abi.html#mangling
func demangle(mangledname string) string {
	var d Demangler
	d.unmangle(mangledname)
	return ""
}

type ParamType struct {
	isPointer     bool
	isRef         bool
	isRValue      bool
	isComplex     bool
	isImaginary   bool
	isArray       bool
	isPtrToMember bool
	isTemplate    bool
	isDecltype    bool
	CvQualifiers  []string
	TypeName      string
}

type Demangler struct {
	Mangled      string
	Remain       string
	CVqualifiers []byte
	RefQualifer  byte
	FuncName     string
	isCtor       bool
	isDtor       bool
	isStd        bool // from St <unqualified-name>   # ::std::
	AllParams    []ParamType
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

func (p *Demangler) generate() string {
	return ""
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

/*
<mangled-name> ::= _Z <encoding>
    <encoding> ::= <function name> <bare-function-type>
	      	   ::= <data name>
	           ::= <special-name>
*/
func (p *Demangler) parseEncoding() {
	if isSpecialName(p.Remain) {
		p.parseSpecialName()
	} else {
		p.parseName()
		p.parseBareFunctionTypes()
	}
}

func (p *Demangler) parseBareFunctionTypes() {
	for len(p.Remain) > 0 {
		p.parseBareFunctionType()
	}
}

func (p *Demangler) parseBareFunctionType() {
	var para ParamType
	res, remain := para.parseType(p.Remain)
	
	if !res {
		return
	}
	
	p.AllParams = append(p.AllParams, para)
	p.Remain = remain
}

// Virtual Tables and RTTI
func isVirtualTableAndRTTI(mangledname string) bool {
	if len(mangledname) < 2 {
		return false
	}
	
	nm := mangledname[0:2]
	return (nm == "TV") || (nm == "TT") || (nm == "TI") || (nm == "TS")
}

func (p *Demangler) parseVirtualTableAndRTTI() {
	nm := p.Remain[0:2]
	// Virtual Tables and RTTI
	if nm == "TV" {
		p.Remain = p.Remain[2:]
	} else if nm == "TT" {
		p.Remain = p.Remain[2:]
	} else if nm == "TI" {
		p.Remain = p.Remain[2:]
	} else if nm == "TS" {
		p.Remain = p.Remain[2:]
	} else {
		fmt.Printf("WARN: Not valid VirtualTableAndRTTI: %v\n", p.Remain)
		return
	}
	
	var para ParamType
	res, remain := para.parseType(p.Remain)
	
	if !res {
		return
	}
	
	p.AllParams = append(p.AllParams, para)
	p.Remain = remain
}

// Virtual Tables and RTTI
func isVirtualOverrideTrunks(mangledname string) bool {
	if len(mangledname) < 2 {
		return false
	}
	
	nm := mangledname[0:2]
	return (nm == "Th") || (nm == "Tv") || (nm == "Tc")
}

func (p *Demangler) parseVirtualOverrideTrunks() {
	// Virtual Tables and RTTI
	nm := p.Remain[0:2]
	if nm == "Th" {
		p.Remain = p.Remain[2:]
		
		// TODO parse offset number		
		panic("Not implemented")
	} else if nm == "Tv" {
		p.Remain = p.Remain[2:]
		// TODO parse offset number
		panic("Not implemented")
	} else if nm == "Tc" {
		p.Remain = p.Remain[2:]
		// TODO parse offset number
		panic("Not implemented")
	} else {
		fmt.Printf("WARN: Not valid VirtualOverrideTrunks: %v\n", p.Remain)
	}
}

func (p *Demangler) parseSpecialName() {
	if isVirtualTableAndRTTI(p.Remain) {
		p.parseVirtualTableAndRTTI()
	} else if isVirtualOverrideTrunks(p.Remain) {
		p.parseVirtualOverrideTrunks()
	} else {
		fmt.Printf("WARN: Unknown SpecialName: %v\n", p.Remain)
	}
}

/*
<substitution> ::= S <seq-id> _
		 ::= S_
*/
func isSubstitution(mangledname string) bool {
	if len(mangledname) < 1 {
		return false
	}
	
	return mangledname[0] == 'S'
}

// 5.1.4.1 Virtual Tables and RTTI
// 5.1.4.2 Virtual Override Thunks
func isLocalName(mangled string) bool {
	if len(mangled) < 1 {
		return false
	}
	
	c0 := mangled[0]
		
	return (c0 == 'Z')
}

/*
<name> ::= <nested-name>
	   ::= <unscoped-name>
	   ::= <unscoped-template-name> <template-args>
	   ::= <local-name>	# See Scope Encoding below

    <unscoped-name> ::= <unqualified-name>
		    ::= St <unqualified-name>   # ::std::

    <unscoped-template-name> ::= <unscoped-name>
			     ::= <substitution>
*/
func (p *Demangler) parseName() {
	if isNestedName(p.Remain) {
		p.parseNestedName()
	} else if isSubstitution(p.Remain) {
		p.parseSubstitution()
	} else if isLocalName(p.Remain) {
		p.parseLocalName()
	} else {
		p.parseUnscopedName()
	}
}

/*
<unscoped-name> ::= <unqualified-name>
		    ::= St <unqualified-name>   # ::std::
*/
func (p *Demangler) parseUnscopedName() {
	mangledname := p.Remain
	if len(mangledname) < 2 {
		return
	}
	
	op := mangledname[0:2]
	
	if op == "St" {
		p.isStd = true
		p.Remain = p.Remain[2:]
	}
	
	p.parseUnqualifiedName()
}

func (p *Demangler) parseSubstitution() {
	panic("Not Implemented")
}

func (p *Demangler) parseLocalName() {
	panic("Not Implemented")
}

func (p *Demangler) parseNestedName() {
	if !isNestedName(p.Remain) {
		fmt.Printf("WARN:Mangled remain not a nested name: %v\n", p.Remain)
		return
	}
	
	p.parseCvQualifiers()
	
	p.parseRefQualifier()
	
	p.parsePrefix()
	
	p.parseUnqualifiedName()
}

func (p *Demangler) parseCvQualifiers() {
	for isCVqualifier(p.Remain) {
		p.CVqualifiers = append(p.CVqualifiers, p.Remain[0])
		p.Remain = p.Remain[1:]
	}
}

func parseCvQualifiers(mangled string) (cvs []string, remain string) {
	part := mangled
	qualifiers := []string{}
	for isCVqualifier(part) {
		c := part[0]
		if c == 'r' {
			qualifiers = append(qualifiers, "restrict")
		} else if c == 'V' {
			qualifiers = append(qualifiers, "volatile")
		} else if c == 'K' {
			qualifiers = append(qualifiers, "const")
		} else {
			fmt.Printf("WARN: Invalid CV Qualifiers: %v\n", mangled)
		}
		
		part = part[1:]
	}
	
	return qualifiers, part
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
func (p *Demangler) parseUnqualifiedName() string {
	var res bool
	var nm string
	res = p.parseCtorDtorName()
	if res {
		//  ?
		return ""
	}
	
	res, nm = parseOperatorName(p.Remain)	
	if res {
		p.Remain = p.Remain[2:]
		return nm
	}
	
	res = p.parseUnnamedTypeName()
	if res {
		p.Remain = p.Remain[2:]
		return nm		
	}
	
	var keyword string
	keyword, p.Remain = parseSourceName(p.Remain)
	if keyword != "" {
		return keyword
	}
	
	res, p.Remain = parseUnnamedTypeName(p.Remain)
	if res {
		
	}
	
	return ""
}

func parseOperatorName(mangledname string) (bool, string) {
	var nm string
	if len(mangledname) < 2 {
		return false, ""
	}
	
	op := mangledname[0:2]
	if (op == "nw") {
		nm = "new"
	} else if (op == "na") {
		nm = "new[]"
	} else if (op == "dl") {
		nm = "delete"
	} else if (op == "da") {
		nm = "delete[]"
	} else if (op == "ps") {
		nm = "+ (unary)"
	} else if (op == "ng") {
		nm = "- (unary)"
	} else if (op == "ad") {
		nm = "& (unary)"
	} else if (op == "de") {
		nm = "* (unary)"
	} else if (op == "co") {
		nm = "~"
	} else if (op == "pl") {
		nm = "+"
	} else if (op == "mi") {
		nm = "-"
	} else if (op == "ml") {
		nm = "*"
	} else if (op == "dv") {
		nm = "/"
	} else if (op == "rm") {
		nm = "%"
	} else if (op == "an") {
		nm = "&"
	} else if (op == "or") {
		nm = "|"
	} else if (op == "eo") {
		nm = "^"
	} else if (op == "aS") {
		nm = "="
	} else if (op == "pL") {
		nm = "+="
	} else if (op == "mI") {
		nm = "-="
	} else if (op == "mL") {
		nm = "*="
	} else if (op == "dV") {
		nm = "/="
	} else if (op == "rM") {
		nm = "%="
	} else if (op == "aN") {
		nm = "&="
	} else if (op == "oR") {
		nm = "|="
	} else if (op == "eO") {
		nm = "^="
	} else if (op == "ls") {
		nm = "<<"
	} else if (op == "rs") {
		nm = ">>"
	} else if (op == "lS") {
		nm = "<<="
	} else if (op == "rS") {
		nm = ">>="
	} else if (op == "eq") {
		nm = "=="
	} else if (op == "ne") {
		nm = "!="
	} else if (op == "lt") {
		nm = "<"
	} else if (op == "gt") {
		nm = ">"
	} else if (op == "le") {
		nm = "<="
	} else if (op == "ge") {
		nm = ">="
	} else if (op == "nt") {
		nm = "!"
	} else if (op == "aa") {
		nm = "&&"
	} else if (op == "oo") {
		nm = "||"
	} else if (op == "pp") {
		nm = "++"
	} else if (op == "mm") {
		nm = "--"
	} else if (op == "cm") {
		nm = ","
	} else if (op == "pm") {
		nm = "->*"
	} else if (op == "pt") {
		nm = "->"
	} else if (op == "cl") {
		nm = "()"
	} else if (op == "ix") {
		nm = "[]"
	} else if (op == "qu") {
		nm = "?"
	} else if (op == "cv") {
		nm = "(cast)" // ?
	} else if (op == "li") {
		nm = ""
	} else if (op[0] == 'v') {
		nm = "" // ?
	} else {
		nm = ""
		return false, nm
	}
	
	return true, nm
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

func parseUnnamedTypeName(mangled string) (isSuccess bool, remain string) {
	if len(mangled) < 2 {
		return false, mangled
	}
	
	nm := mangled[0:2]
	if nm == "Ut" {
		// TODO
		fmt.Printf("WARN: parseUnnamedTypeName not implemented! Value: %v\n", mangled)
		return true, mangled[2:]
	} else {
		return false, mangled
	}
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

func parseNumber(mangled string) (num int, remain string) {
	if len(mangled) == 0 {
		return -1, mangled
	}
	
	c0 := mangled[0]
	if len(mangled) == 1 {
		if isNumberChar(c0) {
			return int(c0 - '0'), ""
		} else {
			return -1, mangled
		}		
	}
	
	c1 := mangled[1]
	if isNumberChar(c0) && isNumberChar(c1) {
		charLen := (int(c0 - '0') * 10) + int(c1 - '0')
		part := mangled[2:]
		
		return charLen, part
	} else if isNumberChar(c0) {
		charLen := int(c0 - '0')
		part := mangled[1:]
		
		return charLen, part
	} else {
		// some error
		fmt.Printf("WARN: Invalid number: %v\n", mangled)
		return -1, mangled
	}
}

// <source-name> ::= <positive length number> <identifier>
func parseSourceName(mangled string) (identifier, remain string) {
	if len(mangled) < 2 {
		return "", mangled
	}
	
	size, part := parseNumber(mangled)
	if size < 0 {
		return "", mangled
	}
	
	return parseIdentifier(part, size)
}

func parseIdentifier(mangled string, size int) (keyword, remain string) {
	s := mangled[0:size]
	r := mangled[size:]
	return s, r
}

func (p *Demangler) parseTemplatePrefix() {
	
}

func (p *ParamType) parseType(mangled string) (result bool, remains string) {
	qualifiers, remain := parseCvQualifiers(mangled)
	
	if len(qualifiers) > 0 {
		p.CvQualifiers = append(p.CvQualifiers, qualifiers...)
	}
	
	c0 := remain[0]
	if c0 == 'P' {
		p.isPointer = true
		remain = remain[1:]
	} else if c0 == 'R' {
		p.isRef = true
		remain = remain[1:]
	} else if c0 == 'O' {
		p.isRValue = true
		remain = remain[1:]
	} else if c0 == 'C' {
		p.isComplex = true
		remain = remain[1:]
	} else if c0 == 'G' {
		p.isImaginary = true
		remain = remain[1:]
	} else if c0 == 'U' {
		panic("Not implemented")
	}
	
	c1 := remain[0]
	if c1 == 'A' {
		// array type
		p.isArray = true
		
		remain = p.parseArrayType(remain)
	} else if c1 == 'M' {
		panic("Not Implemented")
	} else if c1 == 'T' {
		// Template parameters
		panic("Not Implemented")
	}
	
	return true, remain
}

/*
<array-type> ::= A <positive dimension number> _ <element type>
	         ::= A [<dimension expression>] _ <element type>
*/
func (p *ParamType) parseArrayType(mangledname string) string {
	if len(mangledname) < 4 {
		return ""
	}
	
	if mangledname[0] != 'A' {
		fmt.Printf("WARN: Invalid array type: [%v]\n", mangledname)
		return ""
	}
	
	remain := mangledname[1:]
	size, part := parseNumber(remain)
	if size < 0 {
		return remain
	}
	
	if part[0] != '_' {
		fmt.Printf("WARN: Invalid array size type: [%v]\n", mangledname)
		return remain
	}
	
	_, remain = p.parseType(part[1:]) // element type
	return remain
}

func isBuildinType(mangled string) bool {
	alltype := "vwbcahstijlmxynofdegzu"
	firstChar := mangled[:1]
	if strings.Contains(alltype, firstChar) {
		return true
	}
	
	if firstChar != "D" {
		return false
	}
	
	if len(mangled) < 2 {
		return false
	}
	
	secondChar := mangled[1]
	return (secondChar == 'd') || (secondChar == 'e') ||
		(secondChar == 'f') || (secondChar == 'h') ||
		(secondChar == 'i') || (secondChar == 's') ||
		(secondChar == 'a') || (secondChar == 'c') ||
		(secondChar == 'n')
}

// Builtin types are represented by single-letter codes
func parseBuildinType(c0, c1 byte) string {
	if c0 == 'v' {
		return "void"
	} else if c0 == 'w' {
		return "wchar_t"
	} else if c0 == 'b' {
		return "bool"
	} else if c0 == 'c' {
		return "char"
	} else if c0 == 'a' {
		return "signed char"
	} else if c0 == 'h' {
		return "unsigned char"
	} else if c0 == 's' {
		return "short"
	} else if c0 == 't' {
		return "unsigned short"
	} else if c0 == 'i' {
		return "int"
	} else if c0 == 'j' {
		return "unsigned int"
	} else if c0 == 'l' {
		return "long"
	} else if c0 == 'm' {
		return "unsigned long"
	} else if c0 == 'x' {
		return "long long"
	} else if c0 == 'y' {
		return "unsigned long long"
	} else if c0 == 'n' {
		return "__int128"
	} else if c0 == 'o' {
		return "unsigned __int128"
	} else if c0 == 'f' {
		return "float"
	} else if c0 == 'd' {
		return "double"
	} else if c0 == 'e' {
		return "long double"
	} else if c0 == 'g' {
		return "__float128"
	} else if c0 == 'z' {
		return "ellipsis"
	} else if c0 == 'u' {
		panic("Not implemented")
	} else if c0 == 'D' {
		if (c1 == 'd') || (c1 == 'e') ||
			(c1 == 'f') || (c1 == 'h') {
				panic("Not implemented")
		} else if c1 == 'i' {
			return "char32_t"
		} else if c1 == 's' {
			return "char16_t"
		} else if c1 == 'a' {
			return "auto"
		} else if c1 == 'c' {
			return "decltype(auto)"
		} else {
			fmt.Printf("WARN: Invalid buildin type: [%c, %c]\n", c0, c1)
			return ""
		}
	} else {
		fmt.Printf("WARN: Invalid buildin type 2: [%c, %c]\n", c0, c1)
		return ""
	}
}
