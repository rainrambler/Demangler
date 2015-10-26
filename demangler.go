package main

import (
	"fmt"
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

func (p *Demangler) unmangle(mangled string) {
	
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

func (p *Demangler) parsePrefix() {
	
}

func (p *Demangler) parseUnqualifiedName() {
	
}

func (p *Demangler) parseOperatorName() {
	if len(p.Remain < 2) {
		return
	}
	
	op := p.Remain[0:2]
	if (op == "nw") {
		p.FuncName = "new"
	} else if (op == "na") {
		p.FuncName = "new[]"
	} else if (op == "dl") {
		
	} else if (op == "da") {
		
	} else if (op == "ps") {
		
	} else if (op == "ng") {
		
	} else if (op == "ad") {
		
	} else if (op == "de") {
		
	} else if (op == "co") {
		
	} else if (op == "pl") {
		
	} else if (op == "mi") {
		
	} else if (op == "ml") {
		
	} else if (op == "dv") {
		
	} else if (op == "rm") {
		
	} else if (op == "an") {
		
	} else if (op == "or") {
		
	} else if (op == "eo") {
		
	} else if (op == "aS") {
		
	} else if (op == "pL") {
		
	} else if (op == "mI") {
		
	} else if (op == "mL") {
		
	} else if (op == "dV") {
		
	} else if (op == "rM") {
		
	} else if (op == "aN") {
		
	} else if (op == "oR") {
		
	} else if (op == "eO") {
		
	} 
}

func (p *Demangler) parseTemplatePrefix() {
	
}

func (p *Demangler) parseSourceName() {
	
}

func (p *Demangler) parseIdentifier() {
	
}
