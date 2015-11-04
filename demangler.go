package main

import (
	"fmt"
	"strings"
)

const
(
	CTOR_FUNC   = 1
	DTOR_FUNC   = 2
	NORMAL_FUNC = 3
)

// http://mentorembedded.github.io/cxx-abi/abi.html#mangling
func demangle(mangledname string) string {
	var d Demangler
	d.unmangle(mangledname)
	return d.generate()
}

type EntityName struct {
	Remain       string
	isStd        bool // from St <unqualified-name>   # ::std::
	Result       string
	CVqualifiers []string // restrict (C99), volatile, const
	RefQualifer  byte // R: &, O: &&
	NestedNames  []string 
	FunctionType int
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
	CVqualifiers []string // restrict (C99), volatile, const
	RefQualifer  byte // R: &, O: &&
	NestedNames  []string 
	FuncName     string
	isCtor       bool
	isDtor       bool
	AllParams    []ParamType
}

func isNestedName(mangled string) bool {
	if len(mangled) < 1 {
		return false
	}
	// the first char
	return mangled[0] == 'N'
}

// <CV-qualifiers> ::= [r] [V] [K] 	# restrict (C99), volatile, const
func isCVqualifier(mangled string) bool {
	if len(mangled) < 1 {
		return false
	}
	// the first char
	return (mangled[0] == 'r') ||
		(mangled[0] == 'V') || 
		(mangled[0] == 'K')
}

/*
<ref-qualifier> ::= R                   # & ref-qualifier
<ref-qualifier> ::= O                   # && ref-qualifier
*/
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

func (p *Demangler) fillFuncName() {
	if p.FuncName == "" {
		size := len(p.NestedNames)
		if size > 0 {
			p.FuncName = p.NestedNames[size - 1] // the last
			p.NestedNames = p.NestedNames[:size - 1] // remove the last
		}
	}
}

func (p *Demangler) generate() string {
	s := ""
		
	p.fillFuncName()
	
	for _, item := range p.NestedNames {
		s += item + "::"
	}
	
	if p.isCtor {
		s += p.FuncName + "::" + p.FuncName + "()"
		return s
	}
	
	if p.isDtor {
		s += p.FuncName + "::~" + p.FuncName + "()"
		return s
	}
	
	s += p.FuncName + "("
	
	// params
	for _, item := range p.AllParams {
		s += item.toString() + ","
	}
	
	s = strings.TrimRight(s, ",") // trim the last comma
	
	s += ")"
	
	for _, item := range p.CVqualifiers {
		s += item
	}
	
	return s
}

func (p *Demangler) unmangle(mangled string) {
	p.Mangled = mangled
	p.Remain = mangled
	p.CVqualifiers = []string{}
	p.NestedNames = []string{}
	p.AllParams = []ParamType{}
	
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
	fmt.Printf("DBG: Demangler.parseEncoding mangledname: %v\n", p.Remain)
	if isSpecialName(p.Remain) {
		p.parseSpecialName()
	} else {
		var en EntityName
		en.Remain = p.Remain
		en.parseName()
		
		p.AppendResult(en)
		p.parseBareFunctionTypes()
	}
}

func (p *Demangler) AppendResult(en EntityName) {
	if en.FunctionType == CTOR_FUNC {
		p.isCtor = true
	} else if en.FunctionType == DTOR_FUNC {
		p.isDtor = true
	}
		
	if en.isStd {
			
	}
	
	p.RefQualifer = en.RefQualifer
	if len(en.CVqualifiers) > 0 {
		p.CVqualifiers = append(p.CVqualifiers, en.CVqualifiers...)
	}
	
	if len(en.NestedNames) > 0 {
		p.NestedNames = append(p.NestedNames, en.NestedNames...)
	}
	
	p.FuncName = en.Result
	p.Remain = en.Remain
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

/*
<special-name> ::= T <call-offset> <base encoding>
		      # base is the nominal target function of thunk
  <call-offset> ::= h <nv-offset> _
		::= v <v-offset> _
  <nv-offset> ::= <offset number>
		      # non-virtual base override
  <v-offset>  ::= <offset number> _ <virtual offset number>
		      # virtual base override, with vcall offset
*/
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

/*
<special-name> ::= T <call-offset> <base encoding>
		      # base is the nominal target function of thunk
*/
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
func (p *EntityName) parseName() {
	fmt.Printf("DBG: EntityName.parseName Remain: %v\n", p.Remain)
	if isNestedName(p.Remain) {
		p.parseNestedName()
	} else if isSubstitution(p.Remain) {
		p.parseSubstitution()
	} else if isLocalName(p.Remain) {
		p.parseLocalName()
	} else {
		p.Result = p.parseUnscopedName()
	}
}

/*
<unscoped-name> ::= <unqualified-name>
		    ::= St <unqualified-name>   # ::std::
*/
// return value: parsed name
func (p *EntityName) parseUnscopedName() string {
	fmt.Printf("DBG: EntityName.parseUnscopedName Remain: %v\n", p.Remain)
	mangledname := p.Remain
	if len(mangledname) < 2 {
		fmt.Printf("WARN: EntityName.parseUnscopedName Unknown name: %v\n", mangledname)
		return mangledname
	}
	
	op := mangledname[0:2]
	
	if op == "St" {
		p.isStd = true
		p.Remain = p.Remain[2:]
	}
	
	return p.parseUnqualifiedName()
}

/*
<substitution> ::= S <seq-id> _
		 ::= S_
*/
func (p *EntityName) parseSubstitution() {
	// TODO
	panic("Not Implemented")
}

func (p *EntityName) parseLocalName() {
	// TODO
	panic("Not Implemented")
}

/*
<nested-name> ::= N [<CV-qualifiers>] [<ref-qualifier>] <prefix> <unqualified-name> E
		      ::= N [<CV-qualifiers>] [<ref-qualifier>] <template-prefix> <template-args> E
*/
func (p *EntityName) parseNestedName() {
	fmt.Printf("DBG: EntityName.parseNestedName: %v\n", p.Remain)
	if !isNestedName(p.Remain) {
		fmt.Printf("WARN:Mangled remain not a nested name: %v\n", p.Remain)
		return
	}
	
	p.Remain = p.Remain[1:] // remove the first 'N' 
	
	p.parseCvQualifiers()
	
	p.parseRefQualifier()
	
	p.parsePrefix()
	
	for (p.Remain[0] != 'E') && (len(p.Remain) > 0) {
		partName := p.parseUnqualifiedName()
		//partName := p.parseUnscopedName()
		p.NestedNames = append(p.NestedNames, partName)
	}
		
	if p.Remain[0] == 'E' {
		p.Remain = p.Remain[1:]
	} else {
		fmt.Printf("WARN:EntityName.parseNestedName: invalid remain: %v\n", p.Remain)
	}
}

// <CV-qualifiers> ::= [r] [V] [K] 	# restrict (C99), volatile, const
func (p *EntityName) parseCvQualifiers() {
	qualifiers, remain := parseCvQualifiers(p.Remain)
	p.Remain = remain
	
	if len(qualifiers) > 0 {
		p.CVqualifiers = append(p.CVqualifiers, qualifiers...)
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

func parseRefQualifier(mangled string) (refq byte, remain string) {
	if isRefqualifier(mangled) {
		r := mangled[0]
		rem := mangled[1:]
		
		return r, rem
	} else {
		return 0, mangled
	}
}

func (p *EntityName) parseRefQualifier() {
	if isRefqualifier(p.Remain) {
		p.RefQualifer = p.Remain[0]
		p.Remain = p.Remain[1:]
	}
}

/*
<prefix> ::= <unqualified-name>                 # global class or namespace
         ::= <prefix> <unqualified-name>        # nested class or namespace
	     ::= <template-prefix> <template-args>  # class template specialization
         ::= <template-param>                   # template type parameter
         ::= <decltype>                         # decltype qualifier
         ::= <prefix> <data-member-prefix>      # initializer of a data member
	     ::= <substitution>
*/
func (p *EntityName) parsePrefix() {
	fmt.Printf("DBG: EntityName.parsePrefix: %v\n", p.Remain)
	if len(p.Remain) < 2 {
		return 
	}
	
	c0 := p.Remain[0]
	if c0 == 'S' {
		p.parseSubstitution()
		return
	} else if c0 == 'T' {
		p.parseTemplateParam()
		return
	}
	
	twochar := p.Remain[:2]
	if (twochar == "Dt") || (twochar == "DT") {
		parseDeclType()
		return
	}
	
	s := p.parseUnqualifiedName()
	if s == "" {
		p.parsePrefix()
		s = p.parseUnqualifiedName()
		p.NestedNames = append(p.NestedNames, s)
	} else {
		p.NestedNames = append(p.NestedNames, s)
	}
}

/*
<decltype>  ::= Dt <expression> E  # decltype of an id-expression 
                                      or class member access (C++0x)
            ::= DT <expression> E  # decltype of an expression (C++0x)
*/
func parseDeclType() {
	panic("Not Implemented")
}

/*
<template-param> ::= T_	# first template parameter
		         ::= T <parameter-2 non-negative number> _
*/
func (p *EntityName) parseTemplateParam() {
	panic("Not Implemented")
}

/*
<unqualified-name> ::= <operator-name>
                   ::= <ctor-dtor-name>  
                   ::= <source-name>   
                   ::= <unnamed-type-name>   
*/
func (p *EntityName) parseUnqualifiedName() string {
	fmt.Printf("DBG: parseUnqualifiedName: %v\n", p.Remain)
	
	var res bool
	var nm string
	p.FunctionType, p.Remain = parseCtorDtorName(p.Remain)
	if p.FunctionType != NORMAL_FUNC {
		// Ctor or Dtor confirmed
		return p.Remain
	}
	
	res, nm = parseOperatorName(p.Remain)	
	if res {
		p.Remain = p.Remain[2:]
		return nm
	}
	
	res, nm = parseUnnamedTypeName(p.Remain)
	if res {
		p.Remain = nm
		return nm		
	}
	
	var keyword string
	keyword, p.Remain = parseSourceName(p.Remain)
	return keyword
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

func parseCtorDtorName(mangled string) (funcType int, remain string) {
	if len(mangled) < 2 {
		return NORMAL_FUNC, mangled
	}
	
	nm := mangled[0:2]
	if nm == "C1" {
		return CTOR_FUNC, mangled[2:]
	} else if nm == "C2" {
		return CTOR_FUNC, mangled[2:]
	} else if nm == "C3" {
		return CTOR_FUNC, mangled[2:]
	} else if nm == "D0" {
		return DTOR_FUNC, mangled[2:]
	} else if nm == "D1" {
		return DTOR_FUNC, mangled[2:]
	} else if nm == "D2" {
		return DTOR_FUNC, mangled[2:]
	} else {
		return NORMAL_FUNC, mangled
	}
}

/*
<unnamed-type-name> ::= Ut [ <nonnegative number> ] _ 
*/
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
		return -1, mangled[1:]
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

func (p *EntityName) parseTemplatePrefix() {
	
}

/*
<type> ::= <builtin-type>
	 ::= <function-type>
	 ::= <class-enum-type>
	 ::= <array-type>
	 ::= <pointer-to-member-type>
	 ::= <template-param>
	 ::= <template-template-param> <template-args>
	 ::= <decltype>
	 ::= <substitution> # See Compression below
	
<type> ::= <CV-qualifiers> <type>
	 ::= P <type>	# pointer-to
	 ::= R <type>	# reference-to
	 ::= O <type>	# rvalue reference-to (C++0x)
	 ::= C <type>	# complex pair (C 2000)
	 ::= G <type>	# imaginary (C 2000)
	 ::= U <source-name> [<template-args>] <type>	# vendor extended type qualifier
*/
func (p *ParamType) parseType(mangled string) (result bool, remains string) {
	fmt.Printf("DBG: ParamType.parseType mangled: %v\n", mangled)
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
	
	if isBuildinType(remain) {
		if len(remain) > 1 {
			s := parseBuildinType(c1, remain[1])
			p.TypeName = s
			return (s != ""), remain[1:]
		} else {
			s := parseBuildinType(c1, 0)
			p.TypeName = s
			return (s != ""), remain[1:]
		}
	} else {
		return p.parseClassEnumType(remain)
	}
}

/*
<class-enum-type> ::= <name>     # non-dependent type name, dependent type name, or dependent typename-specifier
                  ::= Ts <name>  # dependent elaborated type specifier using 'struct' or 'class'
                  ::= Tu <name>  # dependent elaborated type specifier using 'union'
                  ::= Te <name>  # dependent elaborated type specifier using 'enum'
*/
// Return value: remained string
func (p *ParamType) parseClassEnumType(mangledname string) (result bool, remains string) {
	fmt.Printf("DBG: ParamType.parseClassEnumType mangledname: %v\n", mangledname)
	if len(mangledname) < 2 {
		return false, mangledname
	}
	
	var mangled string
	prefix := mangledname[:2]
	if prefix == "Ts" {
		// TODO
		// parse each different class enum type
		mangled = mangledname[2:]			
	} else if prefix == "Tu" {
		mangled = mangledname[2:]
	} else if prefix == "Te" {
		mangled = mangledname[2:]
	} else {
		mangled = mangledname
	}
	
	var en EntityName
	en.Remain = mangled
	en.parseName()

	p.TypeName = en.Result
	return true, en.Remain
}

/*
<array-type> ::= A <positive dimension number> _ <element type>
	         ::= A [<dimension expression>] _ <element type>
*/
func (p *ParamType) parseArrayType(mangledname string) string {
	fmt.Printf("DBG: ParamType.parseArrayType: %v\n", mangledname)
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

func (p *ParamType) toString() string {
	s := ""
	for _, item := range p.CvQualifiers {
		s += item + " "
	}
	
	s += p.TypeName
	
	if p.isPointer {
		s += "*"
	} else if p.isRef {
		s += "&"
	} else if p.isRValue {
		s += "&&"
	}
	
	return s
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
