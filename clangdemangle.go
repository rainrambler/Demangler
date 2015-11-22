package main

import (
	//"fmt"
	//"strings"
)

const
(
	unknown_error = -4
    invalid_args = -3
    invalid_mangled_name = -2
    memory_alloc_failure = -1
    success = 0
)

type string_pair struct {
	first  string
	second string
}

func (p *string_pair) size() int {
	return len(p.first) + len(p.second)
}

func (p *string_pair) full() string {
	return (p.first) + (p.second)
}

// ???
func (p *string_pair) move_full() string {
	return (p.first) + (p.second)
}

// helper function
func (p *string_pair) second_back() byte {
	size := len(p.second)
	
	if size == 0 {
		return 0
	}
	
	return p.second[size - 1] // the last
}

type sub_type struct {
	content []string_pair
}

func (p *sub_type) push_back(s string_pair) {
	p.content = append(p.content, s)
}

type template_param_type struct {
	content []sub_type
}

type template_param_types struct {
	content []template_param_type
}

type CStyleString struct {
	Content   string
	Pos       int
}

func (p *CStyleString) curChar() byte {
	return p.Content[p.Pos]
}

func (p *CStyleString) prefix(size int) string {
	return p.Content[p.Pos:(p.Pos + size)]
}

func (p *CStyleString) nextChar() byte {
	return p.Content[p.Pos + 1]
}

func (p *CStyleString) isNext(another *CStyleString) bool {
	return (p.Pos + 1) == another.Pos
}

func (p *CStyleString) calcDelta(another *CStyleString) int {
	return p.Pos - another.Pos
}

func (p *CStyleString) equals(another *CStyleString) bool {
	return (p.Pos == another.Pos)
}

func (p *CStyleString) notEach(a1, a2 *CStyleString) bool {
	if p.equals(a1) {
		return false
	}
	if p.equals(a2) {
		return false
	}
	return true
}

type Db struct {
	names                      sub_type
	subs                       template_param_type
	template_param             template_param_types
	cv                         int
	ref                        int
	encoding_depth             int
	parsed_ctor_dtor_cv        bool
	tag_templates              bool
	fix_forward_references     bool
	try_to_parse_template_args bool
}

func (p *Db) names_size() int {
	return len(p.names.content)
}

func (p *Db) names_empty() bool {
	return len(p.names.content) == 0
}

func (p *Db) template_param_empty() bool {
	return len(p.template_param.content) == 0
}

func (p *Db) subs_empty() bool {
	return len(p.subs.content) == 0
}

func (p *Db) subs_pop_back() {
	size := len(p.subs.content)
	if size == 0 {
		return
	}
	p.subs.content = p.subs.content[:size - 1]
}

func (p *Db) subs_push_back(st sub_type) {
	p.subs.content = append(p.subs.content, st)
}

func (p *Db) subs_push_back_pair(sp string_pair) {
	var st sub_type
	st.content = append(st.content, sp)
	p.subs.content = append(p.subs.content, st)
}

func (p *Db) names_pop_back() {
	size := len(p.names.content)
	if size == 0 {
		return
	}
	p.names.content = p.names.content[:size - 1]
}

func (p *Db) names_push_back(s string) {
	var pair string_pair
	pair.first = s
	p.names.content = append(p.names.content, pair)
}

func (p *Db) names_emplace_back() {
	pair := string_pair{}
	p.names.content = append(p.names.content, pair)
}

// the last subs
func (p *Db) subs_back() *sub_type {
	size := len(p.subs.content)
	if size == 0 {
		return nil
	}
	return &p.subs.content[size - 1] // the last
}

// the last template_params
func (p *Db) template_param_back() *template_param_type {
	size := len(p.template_param.content)
	if size == 0 {
		return nil
	}
	return &p.template_param.content[size - 1] // the last
}

// the last sub_type
func (p *Db) template_param_back_back() *sub_type {
	size := len(p.template_param.content)
	if size == 0 {
		return nil
	}
	tpt := &p.template_param.content[size - 1] // the last
	
	size = len(tpt.content)
	return &tpt.content[size - 1]
}

func (p *Db) template_param_pop_back() {
	size := len(p.template_param.content)
	if size == 0 {
		return
	}
	p.template_param.content = p.template_param.content[:size - 1]
}

// the last names
func (p *Db) names_back() *string_pair {
	size := len(p.names.content)
	if size == 0 {
		return nil
	}
	return &p.names.content[size - 1] // the last
}

func cxa_demangle(mangledname string, status *int) string {
	var db Db
	db.cv = 0
	db.ref = 0
    db.encoding_depth = 0
    db.parsed_ctor_dtor_cv = false
    db.tag_templates = true
	db.fix_forward_references = false
    db.try_to_parse_template_args = true
	internal_status := success
	mangledlen := len(mangledname)
	
	demangle2(mangledname, 0, mangledlen, &db,
             &internal_status);
	return mangledname
}

func demangle2(mangled string, first, last int, db *Db, status *int) {
	if (first >= last) {
        *status = invalid_mangled_name
        return
    }
	
	if mangled[first] == '_' {
		if mangled[first + 1] == 'Z' {
			// TODO
		}
	}
}

// <template-param> ::= T_    # first template parameter
//                  ::= T <parameter-2 non-negative number> _
func parse_template_param(first, last *CStyleString, db *Db) CStyleString {
	var cs CStyleString
	cs.Content = first.Content
	cs.Pos = first.Pos
	
	if last.calcDelta(first) < 2 {
		return cs
	}
	
	if cs.curChar() != 'T' {
		return cs
	}
	
	c := cs.nextChar()
	if c == '_' {
		if db.template_param_empty() {
			return cs
		}
		
		lastItem := db.template_param_back()
		if len(lastItem.content) > 0 {
			for _, item := range lastItem.content[0].content {
				db.names_push_back(item.first) // ???
			}
			cs.Pos += 2
		} else {
			db.names_push_back("T_")
			cs.Pos += 2
			db.fix_forward_references = true
		}
	} else if isNumberChar(c) {
		t := &CStyleString{cs.Content, cs.Pos + 1}
		sub := int(t.curChar() - '0')
		t.Pos++
		for !t.equals(last) && isNumberChar(t.curChar()) {
			t.Pos++
			
			sub *= 10
			sub += int(t.curChar() - '0')
		}
		
		if t.equals(last) || (t.curChar() != '_') ||
		    db.template_param_empty() {
			return cs	
		}
		
		sub++
		size := len(db.template_param_back().content)
		if sub < size {
			for _, item := range db.template_param_back().content[sub].content {
				db.names_push_back(item.first) // ???
			}
			
			cs.Pos = t.Pos + 1
		} else {
			db.names_push_back(t.Content[cs.Pos:t.Pos + 1])
			cs.Pos = t.Pos + 1
			db.fix_forward_references = true
		}
	}
	
	return cs
}

//  <ref-qualifier> ::= R                   # & ref-qualifier
//  <ref-qualifier> ::= O                   # && ref-qualifier
//  <function-type> ::= F [Y] <bare-function-type> [<ref-qualifier>] E
func parse_function_type(first, last *CStyleString, db *Db) CStyleString {
	var cs CStyleString
	cs.Content = first.Content
	cs.Pos = first.Pos
	
	if first.Pos == last.Pos {
		return *first
	}
	
	if first.curChar() != 'F' {
		return *first
	}
	
	t := first
	t.Pos++
	
	if t.Pos == last.Pos {
		return cs
	}
	
	if t.curChar() == 'Y' {
		/* extern "C" */
		t.Pos++
		if t.Pos == last.Pos {
			return *first
		}
	}
	
	t1 := parse_type(t, last, db)
	if t1.Pos == t.Pos {
		return cs
	}

	t.Pos = t1.Pos
	sig := "("
	ref_qual := 0
	
	for {
		if t.Pos == last.Pos {
			db.names_pop_back()
			return cs
		}
		
		c := t.curChar()
		if c == 'E' {
			t.Pos++
			break
		}
		if c == 'v' {
			t.Pos++
			continue
		}
		if (c == 'R') && !t.isNext(last) && (t.nextChar() == 'E') {
			t.Pos++
			ref_qual = 1
			continue
		}
		if (c == 'O') && !t.isNext(last) && (t.nextChar() == 'E') {
			t.Pos++
			ref_qual = 2
			continue
		}
		
		k0 := db.names_size()
		t1 := parse_type(t, last, db)
		k1 := db.names_size()
		if (t1.Pos == t.Pos) || (t1.Pos == last.Pos) {
			return cs
		}
		for k := k0; k < k1; k++ {
			if len(sig) > 1 {
				sig += ", "
			}
			sig += db.names.content[k].move_full()
		}
		for k := k0; k < k1; k++ {
			db.names_pop_back()
		}
		
		t.Pos = t1.Pos
	}
	sig += ")"
	
	if ref_qual == 1 {
		sig += " &"
	} else if ref_qual == 2 {
		sig += " &&"
	} else {
		
	}
	
	if db.names_empty() {
		return cs
	}
	
	db.names_back().first += " "
	db.names_back().second = sig + db.names_back().second
	cs.Pos = t.Pos
	
	return cs
}

// <decltype>  ::= Dt <expression> E  # decltype of an id-expression or class member access (C++0x)
//             ::= DT <expression> E  # decltype of an expression (C++0x)
func parse_decltype(first, last *CStyleString, db *Db) CStyleString {
	var cs CStyleString
	cs.Content = first.Content
	cs.Pos = first.Pos
	
	delta := last.Pos - first.Pos
	if delta < 4 {
		return cs
	}
	
	if first.curChar() != 'D' {
		return cs
	}
	
	var noprefix CStyleString
	noprefix.Content = first.Content
	noprefix.Pos = first.Pos + 2
	
	c := first.nextChar()
	if (c == 't') || (c == 'T') {
		t := parse_expression(&noprefix, last, db)
		if (t.Pos != noprefix.Pos) && (t.Pos != last.Pos) && (t.curChar() == 'E') {
			if db.names_empty() {
				return *first
			}
			
			db.names_back().first = "decltype(" + db.names_back().move_full() + ")"
			cs.Pos = t.Pos + 1
		}
	}
	
	return cs
}

// extension:
// <vector-type>           ::= Dv <positive dimension number> _
//                                    <extended element type>
//                         ::= Dv [<dimension expression>] _ <element type>
// <extended element type> ::= <element type>
//                         ::= p # AltiVec vector pixel
func parse_vector_type(first, last *CStyleString, db *Db) CStyleString {
	var cs CStyleString
	cs.Content = first.Content
	cs.Pos = first.Pos
	
	// TODO
	return cs
}

// <encoding> ::= <function name> <bare-function-type>
//            ::= <data name>
//            ::= <special-name>

// clang\cxxabi\src\cxa_demangle.cpp line 4491
func parse_encoding(first, last *CStyleString, db *Db) CStyleString {
	var cs CStyleString
	cs.Content = first.Content
	cs.Pos = first.Pos
	
	if first.equals(last) {
		return cs
	}
	
	su := db.encoding_depth
	db.encoding_depth++
	sb := db.tag_templates
	
	if db.encoding_depth > 1 {
		db.tag_templates = true
	}
	
	c := cs.curChar()
	
	if (c == 'G') || (c == 'T') {
		tmp := parse_special_name(first, last, db)
		return tmp
	}
	
	ends_with_template_args := false
	t := parse_name(first, last, db, &ends_with_template_args)
	cv1 := db.cv
	ref1 := db.ref
	
	if t.equals(first) {
		db.encoding_depth = su // ???
		db.tag_templates = sb
		return cs
	}
	
	c2 := t.curChar()
	if !t.equals(last) && (c2 != 'E') && (c2 != '.') {
		sb2 := db.tag_templates
		db.tag_templates = false
		var t2 CStyleString
		ret2 := ""
		if db.names_empty() {
			db.encoding_depth = su // ???
			db.tag_templates = sb
			return cs
		}
		
		nm := db.names_back().first
		if len(nm) == 0 {
			db.encoding_depth = su // ???
			db.tag_templates = sb
			return cs
		}
		
		if !db.parsed_ctor_dtor_cv && ends_with_template_args {
			t2 = parse_type(&t, last, db)
			
			if t2.equals(&t) {
				db.encoding_depth = su // ???
				db.tag_templates = sb
				return cs
			}
			
			if db.names_size() < 2 {
				db.encoding_depth = su // ???
				db.tag_templates = sb
				return cs
			}
			
			ret1 := db.names_back().first
			ret2 := db.names_back().second
			
			if len(ret2) == 0 {
				ret1 += " "
			}
			
			db.names_pop_back()
			db.names_back().first = ret1 + db.names_back().first
			t.Pos = t2.Pos
		}
		
		db.names_back().first += "("
		if !t.equals(last) && (t.curChar() == 'v') {
			t.Pos++
		} else {
			first_arg := true
			for {
				k0 := db.names_size()
                t2 := parse_type(&t, last, db)
                k1 := db.names_size()
				
				if t2.equals(&t) {
					break
				}
				
				if k1 > k0 {
					tmp := ""
					for k := k0; k < k1; k++ {
						if len(tmp) > 0 {
							tmp += ", "
						}
						
						tmp += db.names.content[k].move_full()
					}
					for k := k0; k < k1; k++ {
						db.names_pop_back()
					}
					
					if len(tmp) > 0 {
						if db.names_empty() {
							return cs
						}
						if !first_arg {
							db.names_back().first += ", "
						} else {
							first_arg = false
						}
						db.names_back().first += tmp
					}
				}
				
				t.Pos = t2.Pos
			}
		}
		
		if db.names_empty() {
			return cs
		}
		
		db.names_back().first += ")"
		if (cv1 & 1) != 0 {
			db.names_back().first += " const"
		}
            
        if (cv1 & 2) != 0 {
			db.names_back().first += " volatile"
		}
         if (cv1 & 4) != 0 {
			db.names_back().first += " restrict"
		}
                        
        if (ref1 == 1) {
			db.names_back().first += " &"
		} else if (ref1 == 2) {
			db.names_back().first += " &&"
		}
                        
        db.names_back().first += ret2
		db.tag_templates = sb2
		cs.Pos = t.Pos
	} else {
		cs.Pos = t.Pos
	}
	
	db.cv = cv1
	db.ref = ref1
	db.encoding_depth = su // ???
	db.tag_templates = sb
	// TODO
	return cs
}

func parse_pointer_to_member_type(first, last *CStyleString, db *Db) CStyleString {
	var cs CStyleString
	cs.Content = first.Content
	cs.Pos = first.Pos
	
	if first.equals(last) {
		return cs
	}
	
	if cs.curChar() != 'M' {
		return cs
	}

	t1 := &CStyleString{cs.Content, cs.Pos + 1}
	t := parse_type(t1, last, db)
	if !t.equals(t1) {
		t2 := parse_type(&t, last, db)
		if !t2.equals(&t) {
			if db.names_size() < 2 {
				return cs
			}
			
			fnc := db.names_back()
			db.names_pop_back()
			class_type := db.names_back()
			
			if (len(fnc.second) > 0) && fnc.second[0] == '(' {
				db.names_back().first = fnc.first + "(" + class_type.move_full() + "::*"
				db.names_back().second = ")" + fnc.second
			} else {
				db.names_back().first = fnc.first + " " + class_type.move_full() + "::*"
				db.names_back().second = fnc.second
			}
			
			cs.Pos = t2.Pos
		}
	}	

	return cs
}

// <discriminator> := _ <non-negative number>      # when number < 10
//                 := __ <non-negative number> _   # when number >= 10
//  extension      := decimal-digit+               # at the end of string
// parse but ignore discriminator
func parse_discriminator(first, last *CStyleString) CStyleString {
	var cs CStyleString
	cs.Content = first.Content
	cs.Pos = first.Pos
	
	if first.equals(last) {
		return cs
	}
	
	if first.curChar() == '_' {
		t1 := &CStyleString{cs.Content, cs.Pos + 1}
		if !t1.equals(last) {
			if isNumberChar(t1.curChar()) {
				cs.Pos = t1.Pos + 1
			} else if t1.curChar() == '_' {
				t1.Pos++
				for !t1.equals(last) && isNumberChar(t1.curChar()) {
					t1.Pos++
				}
				
				if !t1.equals(last) && t1.curChar() == '_' {
					cs.Pos = t1.Pos + 1
				}
			}
		}
	} else if isNumberChar(cs.curChar()) {
		t1 := &CStyleString{cs.Content, cs.Pos + 1}
		for !t1.equals(last) && isNumberChar(t1.curChar()) {
			t1.Pos++
		}
		if t1.equals(last) {
			cs.Pos = t1.Pos
		}
	}
	
	return cs
}

// <template-args> ::= I <template-arg>* E
//     extension, the abi says <template-arg>+
// line 3848
func parse_template_args(first, last *CStyleString, db *Db) CStyleString {
	var cs CStyleString
	cs.Content = first.Content
	cs.Pos = first.Pos
	
	if last.calcDelta(first) < 2 {
		return cs
	}
	
	if first.curChar() != 'I' {
		return cs		
	}
	
	if db.tag_templates {
		db.template_param_back().content = []sub_type{} // clear		
	}
	
	t := &CStyleString{cs.Content, cs.Pos + 1}
	args := "<"
	for t.curChar() != 'E' {
		if db.tag_templates {
			//db.template_param ???
			// db.template_param.emplace_back(db.names.get_allocator());
			
		}
		
		k0 := db.names_size()
		t1 := parse_template_arg(t, last, db)
		k1 := db.names_size()
		
		if db.tag_templates {
			db.template_param_pop_back()
		}
		
		if t1.equals(t) || t1.equals(last) {
			return cs
		}
		
		if db.tag_templates {
			for k := k0; k < k1; k++ {
				db.template_param_back_back().push_back(db.names.content[k])
			}
		}
		
		for k := k0; k < k1; k++ {
			if len(args) > 1 {
				args += ", "
			}
			args += db.names.content[k].move_full()
		}
		
		for k1 != k0 {
			db.names_pop_back()
			k1--
		}
		
		t.Pos = t1.Pos
	}
	
	cs.Pos = t.Pos + 1
	if args[len(args) - 1] != '>' {
		args += ">"
	} else {
		args += " >"
	}
	
	db.names_push_back(args)
	
	return cs
}

// <template-arg> ::= <type>                                             # type or template
//                ::= X <expression> E                                   # expression
//                ::= <expr-primary>                                     # simple expressions
//                ::= J <template-arg>* E                                # argument pack
//                ::= LZ <encoding> E
// line 3795
func parse_template_arg(first, last *CStyleString, db *Db) CStyleString {
	var cs CStyleString
	cs.Content = first.Content
	cs.Pos = first.Pos
	
	if first.equals(last) {
		return cs
	}
	
	switch first.curChar() {
		case 'X':
		tmpPos := &CStyleString{cs.Content, cs.Pos + 1}
		t := parse_expression(tmpPos, last, db)
		if !t.equals(tmpPos) {
			if !t.equals(last) && (t.curChar() == 'E') {
				cs.Pos = t.Pos + 1
			}
		}
		break
		case 'J':
		t := &CStyleString{cs.Content, cs.Pos + 1}
		
		if t.equals(last) {
			return cs
		}
		
		for t.curChar() != 'E' {
			t1 := parse_template_arg(t, last, db)
			if t1.equals(t) {
				return cs
			}
			
			t.Pos = t1.Pos
		}
		cs.Pos = t.Pos + 1
        break
		case 'L':
        // <expr-primary> or LZ <encoding> E
		tmpPos := &CStyleString{cs.Content, cs.Pos + 1}
		if !tmpPos.equals(last) && (tmpPos.curChar() == 'Z') {
			tmpPos.Pos += 1
			t := parse_encoding(tmpPos, last, db)
			if t.notEach(tmpPos, last) && (t.curChar() == 'E') {
				cs.Pos = t.Pos + 1
			}
		} else {
			cs = parse_expr_primary(first, last, db)
		}
		break
		default:
		// <type>
        cs = parse_type(first, last, db)
		break
	}
	
	return cs
}

// <expr-primary> ::= L <type> <value number> E                          # integer literal
//                ::= L <type> <value float> E                           # floating literal
//                ::= L <string type> E                                  # string literal
//                ::= L <nullptr type> E                                 # nullptr literal (i.e., "LDnE")
//                ::= L <type> <real-part float> _ <imag-part float> E   # complex floating point literal (C 2000)
//                ::= L <mangled-name> E     
// line 2664
func parse_expr_primary(first, last *CStyleString, db *Db) CStyleString {
	var cs CStyleString
	cs.Content = first.Content
	cs.Pos = first.Pos
	
	// TODO
	return cs
}

// [gs] nw <expression>* _ <type> E                     # new (expr-list) type
// [gs] nw <expression>* _ <type> <initializer>         # new (expr-list) type (init)
// [gs] na <expression>* _ <type> E                     # new[] (expr-list) type
// [gs] na <expression>* _ <type> <initializer>         # new[] (expr-list) type (init)
// <initializer> ::= pi <expression>* E                 # parenthesized initialization
func parse_new_expr(first, last *CStyleString, db *Db) *CStyleString {
	cs := &CStyleString{first.Content, first.Pos}

	if last.calcDelta(first) < 4 {
		return cs
	}
	
	t := &CStyleString{first.Content, first.Pos}
	parsed_gs := false
	if t.prefix(2) == "gs" {
		t.Pos += 2
		parsed_gs = true
	}
	
	if (t.prefix(2) == "nw") || (t.prefix(2) == "na") {
		is_array := (t.nextChar() == 'a')
		t.Pos += 2
		
		if t.equals(last) {
			return cs
		}
		
		has_expr_list := false
        first_expr := true
		
		for (t.curChar() != '_') {
			t1 := parse_expression(t, last, db)
			if t1.equals(t) || t1.equals(last) {
				return cs
			}
			
			has_expr_list = true
			if !first_expr {
				if db.names_empty() {
					return cs
				}
				
				tmp := db.names_back().move_full()
				db.names_pop_back()
				if len(tmp) > 0 {
					if db.names_empty() {
						return cs
					}
					db.names_back().first += (", ")
                    db.names_back().first += (tmp)
                    first_expr = false
				}
			}
			t.Pos = t1.Pos
		}
		t.Pos++
		t1 := parse_type(t, last, db)
		if t1.equals(t) || t1.equals(last) {
			return cs
		}
		t.Pos = t1.Pos
		has_init := false
		if (last.calcDelta(t) >= 3) && 
			(t.prefix(2) == "pi") {
			t.Pos += 2
			has_init = true
            first_expr = true
			for (t.curChar() != 'E') {
				t1 = parse_expression(t, last, db)
				if t1.equals(t) || t1.equals(last) {
					return cs
				}
				if !first_expr {
					if db.names_empty() {
						return cs
					}
					tmp := db.names_back().move_full()
					db.names_pop_back()
					if len(tmp) > 0 {
						if db.names_empty() {
							return cs
						}
						db.names_back().first += (", ")
                        db.names_back().first += (tmp)
                        first_expr = false
					}
				}
				t.Pos = t1.Pos
			}
		}
		if t.curChar() != 'E' {
			return cs
		}
		init_list := ""
		if (has_init) {
			if db.names_empty() {
				return cs
			}
			init_list = db.names_back().move_full()
            db.names_pop_back()
		}
		if db.names_empty() {
			return cs
		}
		typ := db.names_back().move_full()
        db.names_pop_back()
		expr_list := ""
		if has_expr_list {
			if db.names_empty() {
				return cs
			}
			expr_list = db.names_back().move_full()
            db.names_pop_back()
		}
		r := ""
		if (parsed_gs) {
			r = "::"
		}

        if (is_array) {
			r += "[] "	
		} else {
			r += " "
		}
        if (has_expr_list) {
			r += "(" + expr_list + ") "
		}
        r += typ
        if (has_init) {
			r += " (" + init_list + ")"
		}
        db.names_push_back(r)
        cs.Pos = t.Pos + 1
	}
	
	return cs
}

func parse_noexcept_expression(first, last *CStyleString, db *Db) *CStyleString {
	cs := &CStyleString{first.Content, first.Pos}

	t1 := parse_expression(first, last, db)
	if !t1.equals(first) {
		if db.names_empty() {
			return cs
		}
		db.names_back().first = "noexcept (" + db.names_back().move_full() + ")"
        cs.Pos = t1.Pos
	}
	
	return cs
}

// <call-offset> ::= h <nv-offset> _
//               ::= v <v-offset> _
// 
// <nv-offset> ::= <offset number>
//               # non-virtual base override
// 
// <v-offset>  ::= <offset number> _ <virtual offset number>
//               # virtual base override, with vcall offset
// line 4265
func parse_call_offset(first, last *CStyleString) CStyleString {
	var cs CStyleString
	cs.Content = first.Content
	cs.Pos = first.Pos
	
	if first.equals(last) {
		return cs
	}
	
	switch cs.curChar() {
		case 'h':
		tmpPos := &CStyleString{cs.Content, cs.Pos + 1}
		t := parse_number(tmpPos, last)
		if t.notEach(tmpPos, last) && (t.curChar() == '_') {
			cs.Pos = t.Pos + 1
		}
		break
		case 'v':
		tmpPos := &CStyleString{cs.Content, cs.Pos + 1}
		t := parse_number(tmpPos, last)
		if t.notEach(tmpPos, last) && (t.curChar() == '_') {
			t.Pos++
			t2 := parse_number(&t, last)
			if t2.notEach(tmpPos, last) && (t2.curChar() == '_') {
				cs.Pos = t2.Pos + 1
			}
		}
		break
	}
	
	return cs
}

// <source-name> ::= <positive length number> <identifier>
// line 225
func parse_source_name(first, last *CStyleString, db *Db) CStyleString {
	var cs CStyleString
	cs.Content = first.Content
	cs.Pos = first.Pos
	
	if first.equals(last) {
		return cs
	}
	
	c := cs.curChar()
	if isNumberChar(c) && !first.isNext(last) {
		t := &CStyleString{cs.Content, cs.Pos + 1}
		n := int(c - '0')
		c := t.curChar()
		for isNumberChar(c) {
			n = (n * 10) + int(c - '0')
			t.Pos++
			if t.equals(last) {
				return cs
			}
			c = cs.curChar()
		}
		
		if (last.Pos - t.Pos) >= n {
			r := t.Content[:n]
			
			if r[:10] == "_GLOBAL__N" {
				db.names_push_back("(anonymous namespace)")
			} else {
				db.names_push_back(r)
			}
			
			cs.Pos = t.Pos + n
		}
	}

	return cs
}

// <special-name> ::= TV <type>    # virtual table
//                ::= TT <type>    # VTT structure (construction vtable index)
//                ::= TI <type>    # typeinfo structure
//                ::= TS <type>    # typeinfo name (null-terminated byte string)
//                ::= Tc <call-offset> <call-offset> <base encoding>
//                    # base is the nominal target function of thunk
//                    # first call-offset is 'this' adjustment
//                    # second call-offset is result adjustment
//                ::= T <call-offset> <base encoding>
//                    # base is the nominal target function of thunk
//                ::= GV <object name> # Guard variable for one-time initialization
//                                     # No <type>
//      extension ::= TC <first type> <number> _ <second type> # construction vtable for second-in-first
//      extension ::= GR <object name> # reference temporary for object
func parse_special_name(first, last *CStyleString, db *Db) CStyleString {
	var cs CStyleString
	cs.Content = first.Content
	cs.Pos = first.Pos
	
	if last.Pos - first.Pos <= 2 {
		return cs
	}
	
	switch cs.curChar() {
		case 'T':
		switch cs.nextChar() {
			case 'V':
			// TV <type>    # virtual table
			tmpPos := &CStyleString{cs.Content, cs.Pos + 2}
			t := parse_type(tmpPos, last, db)
			if !t.equals(tmpPos) {
				if (db.names_empty()) {
					return cs
				}  
				db.names_back().first = "typeinfo name for " + db.names_back().first                  
				cs.Pos = t.Pos				
			}
			break
			case 'T':
			// TT <type>    # VTT structure (construction vtable index)
			tmpPos := &CStyleString{cs.Content, cs.Pos + 2}
			t := parse_type(tmpPos, last, db)
			if !t.equals(tmpPos) {
				if (db.names_empty()) {
					return cs
				}  
				db.names_back().first = "VTT for " + db.names_back().first                  
				cs.Pos = t.Pos				
			}
			break
			case 'I':
			// TI <type>    # typeinfo structure
			tmpPos := &CStyleString{cs.Content, cs.Pos + 2}
			t := parse_type(tmpPos, last, db)
			if !t.equals(tmpPos) {
				if (db.names_empty()) {
					return cs
				}  
				db.names_back().first = "typeinfo for " + db.names_back().first                  
				cs.Pos = t.Pos				
			}
			break
			case 'S':
			// TS <type>    # typeinfo name (null-terminated byte string)
			tmpPos := &CStyleString{cs.Content, cs.Pos + 2}
			t := parse_type(tmpPos, last, db)
			if !t.equals(tmpPos) {
				if (db.names_empty()) {
					return cs
				}  
				db.names_back().first = "typeinfo name for " + db.names_back().first                  
				cs.Pos = t.Pos				
			}
			break
			case 'c':
			tmpPos := &CStyleString{cs.Content, cs.Pos + 2}
			t0 := parse_call_offset(tmpPos, last)
			if t0.equals(tmpPos) {
				break
			}
			
			t1 := parse_call_offset(&t0, last)
			if t0.equals(&t1) {
				break
			}
			
			t := parse_encoding(&t1, last, db)
			if t.equals(&t1) {
				if (db.names_empty()) {
					return cs
				}
                
				db.names_back().first = "covariant return thunk to " +
				    db.names_back().first
				cs.Pos = t.Pos
			}
			break
			case 'C':
            // extension ::= TC <first type> <number> _ <second type> # construction vtable for second-in-first
			tmpPos := &CStyleString{cs.Content, cs.Pos + 2}
			t := parse_type(tmpPos, last, db)
			if !t.equals(tmpPos) {
				t0 := parse_number(&t, last)
				if t0.notEach(&t, last) && (t0.curChar() == '_') {
					t0.Pos++
					t1 := parse_type(&t0, last, db)
					
					if !t1.equals(&t0) {
						if (db.names_size() < 2) {
							return cs
						}
						
						left := db.names_back().move_full()
                        db.names_pop_back()
						db.names_back().first = "construction vtable for " +
                            (left) + "-in-" +
                            db.names_back().move_full()
						cs.Pos = t1.Pos
					}	
				}
			}
			break
			default:
			// T <call-offset> <base encoding>
			tmpPos := &CStyleString{cs.Content, cs.Pos + 1}
			t0 := parse_call_offset(tmpPos, last)
			if t0.equals(tmpPos) {
				break
			}
			t := parse_encoding(&t0, last, db)
			if !t.equals(&t0) {
				if db.names_empty() {
					return cs
				}
				
				if cs.nextChar() == 'v' {
					db.names_back().first = "virtual thunk to " +
					    db.names_back().first
					cs.Pos = t.Pos
				} else {
					db.names_back().first = "non-virtual thunk to " +
					    db.names_back().first
					cs.Pos = t.Pos
				}
			}
			break
		}
		break
		case 'G':
        switch cs.nextChar() {
			case 'V':
			// GV <object name> # Guard variable for one-time initialization
			tmpPos := &CStyleString{cs.Content, cs.Pos + 2}
			t := parse_name(tmpPos, last, db, nil)
			if !t.equals(tmpPos) {
				if (db.names_empty()) {
					return cs
				}
				db.names_back().first = "guard variable for " +
					db.names_back().first
			    cs.Pos = t.Pos
			}
			break
			case 'R':
            // extension ::= GR <object name> # reference temporary for object
			tmpPos := &CStyleString{cs.Content, cs.Pos + 2}
			t := parse_name(tmpPos, last, db, nil)
			if !t.equals(tmpPos) {
				if (db.names_empty()) {
					return cs
				}
				db.names_back().first = "reference temporary for " +
					db.names_back().first
			    cs.Pos = t.Pos
			}
			break
		}
		break
	}
	
	return cs
}

func parse_number(first, last *CStyleString) CStyleString {
	if first.Pos == last.Pos {
		return *first
	}
	
	var cs CStyleString
	cs.Content = first.Content
	cs.Pos = first.Pos
	
	t := first
	if t.curChar() == 'n' {
		t.Pos++
	}
	
	if t.Pos != last.Pos {
		if t.curChar() == '0' {
			cs.Pos = t.Pos + 1
		} else if isNonZeroNumberChar(t.curChar()) {
			cs.Pos = t.Pos + 1
			
			for (cs.Pos != last.Pos) && isNumberChar(cs.curChar()) {
				cs.Pos++
			}
		}
	}
	
	return cs
}

// <array-type> ::= A <positive dimension number> _ <element type>
//              ::= A [<dimension expression>] _ <element type>
func parse_array_type(first, last *CStyleString, db *Db) CStyleString {
	var cs CStyleString
	cs.Content = first.Content
	cs.Pos = first.Pos
	if (first.Pos == last.Pos) || ((first.Pos + 1) == last.Pos) {
		return cs
	}
	
	if first.curChar() != 'A' {
		return cs
	}
	
	if first.nextChar() == '_' {
		var noprefix CStyleString
		noprefix.Content = cs.Content
		noprefix.Pos = cs.Pos + 2
		t := parse_type(&noprefix, last, db)
		if t.Pos != noprefix.Pos {
			if db.names_empty() {
				return cs
			}
			
			if db.names_back().second[:2] == " [" {
				// erase(0, 1)
				db.names_back().second = db.names_back().second[1:]
			}
			
			db.names_back().second = " []" + db.names_back().second
			cs.Pos = t.Pos
		}
	} else if isNonZeroNumberChar(first.nextChar()) {
		var noprefix CStyleString
		noprefix.Content = cs.Content
		noprefix.Pos = cs.Pos + 1
		
		t := parse_number(&noprefix, last)
		if (t.Pos != last.Pos) && (t.curChar() == '_') {
			noprefix.Pos = t.Pos + 1
			t2 := parse_type(&noprefix, last, db)
			
			if t2.Pos != noprefix.Pos {
				if db.names_empty() {
					return cs
				}
				
				if db.names_back().second[:2] == " [" {
					db.names_back().second = db.names_back().second[1:]
				}
				
				// " [" + typename C::String(first+1, t) + "]" ???
				tn := " [" + first.Content[first.Pos + 1:t.Pos] + "]"
				db.names_back().second = tn + db.names_back().second
				cs.Pos = t2.Pos
			}
		}
	} else {
		var noprefix CStyleString
		noprefix.Content = cs.Content
		noprefix.Pos = cs.Pos + 1
		
		t := parse_expression(&noprefix, last, db)
		if t.Pos == last.Pos {
			return cs
		}
		
		if t.Pos == noprefix.Pos {
			return cs
		}
		
		if t.curChar() != '_' {
			return cs
		}
		
		noprefix.Pos = t.Pos + 1
		t2 := parse_type(&noprefix, last, db)
		if t2.Pos == noprefix.Pos {
			return cs
		}
		
		if db.names_size() < 2 {
			return *first
		}
		
		ty := db.names_back()
		db.names_pop_back()
		expr := db.names_back()
		db.names_back().first = ty.first
		if ty.second[:2] == " [" {
			ty.second = ty.second[1:]
		}
		db.names_back().second = " [" + expr.move_full() + "]" + ty.second
		cs.Pos = t2.Pos
	}

	return cs
}

// <unqualified-name> ::= <operator-name>
//                    ::= <ctor-dtor-name>
//                    ::= <source-name>   
//                    ::= <unnamed-type-name>
// line 3082
func parse_unqualified_name(first, last *CStyleString, db *Db) CStyleString {
	var cs CStyleString
	cs.Content = first.Content
	cs.Pos = first.Pos
	
	if first.equals(last) {
		return cs
	}
	
	c := cs.curChar()
	if (c == 'C') || (c == 'D') {
		t := parse_ctor_dtor_name(first, last, db)
		if !t.equals(&cs) {
			cs.Pos = t.Pos
		}
	} else if c == 'U' {
		t := parse_unnamed_type_name(first, last, db)
		if !t.equals(&cs) {
			cs.Pos = t.Pos
		}
	} else if isNonZeroNumberChar(c) {
		// 1~9
		t := parse_source_name(first, last, db)
		if !t.equals(&cs) {
			cs.Pos = t.Pos
		}
	} else {
		t := parse_operator_name(first, last, db)
		if !t.equals(&cs) {
			cs.Pos = t.Pos
		}
	}
	
	return cs
}

// <nested-name> ::= N [<CV-qualifiers>] [<ref-qualifier>] <prefix> <unqualified-name> E
//               ::= N [<CV-qualifiers>] [<ref-qualifier>] <template-prefix> <template-args> E
// 
// <prefix> ::= <prefix> <unqualified-name>
//          ::= <template-prefix> <template-args>
//          ::= <template-param>
//          ::= <decltype>
//          ::= # empty
//          ::= <substitution>
//          ::= <prefix> <data-member-prefix>
//  extension ::= L
// 
// <template-prefix> ::= <prefix> <template unqualified-name>
//                   ::= <template-param>
//                   ::= <substitution>
func parse_nested_name(first, last *CStyleString, db *Db, ends_with_template_args *bool) CStyleString {
	var cs CStyleString
	cs.Content = first.Content
	cs.Pos = first.Pos
	
	if first.equals(last) {
		return cs
	}
	
	if first.curChar() != 'N' {
		return cs
	}
	
	cv := 0
	cs.Pos++
	t0 := parse_cv_qualifiers(&cs, last, &cv)
	if t0.equals(last) {
		return *first
	}
	
	db.ref = 0
	if t0.curChar() == 'R' {
		db.ref = 1
		t0.Pos++
	} else if t0.curChar() == 'O' {
		db.ref = 2
		t0.Pos++
	}
	
	db.names_emplace_back()
	
	if (last.calcDelta(&t0) >= 2) && (t0.curChar() == 'S') &&
		(t0.nextChar() == 't') {
		t0.Pos += 2
		db.names_back().first = "std"	
	}
	
	if t0.equals(last) {
		db.names_pop_back()
		return *first
	}
	
	pop_subs := false
	component_ends_with_template_args := false
	for t0.curChar() != 'E' {
		component_ends_with_template_args = false
		var t1 CStyleString
		
		if t0.curChar() == 'S' {
			if !t0.isNext(last) && t0.nextChar() == 't' {
				// do_parse_unqualified_name
				t1 = parse_unqualified_name(&t0, last, db)
				if !t1.equals(&t0) && !t1.equals(last) {
					name := db.names_back().move_full()
					db.names_pop_back()
					if len(db.names_back().first) > 0 {
						db.names_back().first += "::" + name						
					} else {
						db.names_back().first = name
					}
					db.subs_push_back_pair(*db.names_back())
					pop_subs = true
					t0.Pos = t1.Pos
				} else {
					return *first
				}
				// end do_parse_unqualified_name
			} else {
				t1 = parse_substitution(&t0, last, db)
				if !t1.equals(&t0) && !t1.equals(last) {
					name := db.names_back().move_full()
					db.names_pop_back()
					if len(db.names_back().first) > 0 {
						db.names_back().first += "::" + name
						db.subs_push_back_pair(*db.names_back())
					} else {
						db.names_back().first = name
					}
					pop_subs = true
					t0.Pos = t1.Pos
				} else {
					return *first
				}
			}			
		} else if t0.curChar() == 'T' {
			t1 = parse_template_param(&t0, last, db)
			if t1.notEach(&t0, last) {
				name := db.names_back().move_full()
				db.names_pop_back()
				if len(db.names_back().first) > 0 {
					db.names_back().first += "::" + name					
				} else {
					db.names_back().first = name
				}
				db.subs_push_back_pair(*db.names_back())
				pop_subs = true
				t0.Pos = t1.Pos
			} else {
				return *first
			}
		} else if t0.curChar() == 'D' {
			if !t0.isNext(last) && (t0.nextChar() != 't') && (t0.nextChar() != 'T') {
				// do_parse_unqualified_name
				t1 = parse_unqualified_name(&t0, last, db)
				if !t1.equals(&t0) && !t1.equals(last) {
					name := db.names_back().move_full()
					db.names_pop_back()
					if len(db.names_back().first) > 0 {
						db.names_back().first += "::" + name						
					} else {
						db.names_back().first = name
					}
					db.subs_push_back_pair(*db.names_back())
					pop_subs = true
					t0.Pos = t1.Pos
				} else {
					return *first
				}
				// end do_parse_unqualified_name
			} else {
				t1 = parse_decltype(&t0, last, db)
				if t1.notEach(&t0, last) {
					name := db.names_back().move_full()
					db.names_pop_back()
					if len(db.names_back().first) > 0 {
						db.names_back().first += "::" + name						
					} else {
						db.names_back().first = name
					}
					db.subs_push_back_pair(*db.names_back())
					pop_subs = true
					t0.Pos = t1.Pos
				} else {
					return *first
				}
			}
		} else if t0.curChar() == 'I' {
			t1 = parse_template_args(&t0, last, db)
			if t1.notEach(&t0, last) {
				name := db.names_back().move_full()
				db.names_pop_back()
				db.names_back().first += name
				db.subs_push_back_pair(*db.names_back())
				t0.Pos = t1.Pos
				component_ends_with_template_args = true
			} else {
				return *first
			}
		}  else if t0.curChar() == 'L' {
			t0.Pos++
			if t0.equals(last) {
				return *first
			}
		} else {
			// do_parse_unqualified_name
			t1 = parse_unqualified_name(&t0, last, db)
			if !t1.equals(&t0) && !t1.equals(last) {
				name := db.names_back().move_full()
				db.names_pop_back()
				if len(db.names_back().first) > 0 {
					db.names_back().first += "::" + name						
				} else {
					db.names_back().first = name
				}
				db.subs_push_back_pair(*db.names_back())
				pop_subs = true
				t0.Pos = t1.Pos
			} else {
				return *first
			}
			// end do_parse_unqualified_name
		}
	}
	
	cs.Pos = t0.Pos + 1
	db.cv = cv
	if pop_subs && !db.subs_empty() {
		db.subs_pop_back()
	}
	if ends_with_template_args != nil {
		*ends_with_template_args = component_ends_with_template_args
	}
	
	return cs
}

// <local-name> := Z <function encoding> E <entity name> [<discriminator>]
//              := Z <function encoding> E s [<discriminator>]
//              := Z <function encoding> Ed [ <parameter number> ] _ <entity name>
func parse_local_name(first, last *CStyleString, db *Db, ends_with_template_args *bool) CStyleString {
	var cs CStyleString
	cs.Content = first.Content
	cs.Pos = first.Pos
	
	if last.equals(first) {
		return cs
	}
	
	if first.curChar() != 'Z' {
		return cs
	}
	
	tmpPos := &CStyleString{cs.Content, cs.Pos + 1}
	t := parse_encoding(tmpPos, last, db)
	
	if t.curChar() != 'E' {
		return cs
	}
	
	if t.equals(tmpPos) || t.equals(last) || t.isNext(last) {
		return cs
	}
	
	t.Pos++ // ???
	c := t.curChar()
	switch c {
		case 's':
		tmpPos.Pos = t.Pos + 1
		curPos := parse_discriminator(tmpPos, last)
		if db.names_empty() {
			return curPos
		}
		db.names_back().first += "::string literal"
		break
		case 'd':
		tmpPos.Pos = t.Pos + 1
		if !tmpPos.equals(last) {
			t1 := parse_number(tmpPos, last)
			if !t1.equals(last) && (t1.curChar() == '_') {
				t.Pos = t1.Pos + 1
				t1 = parse_name(&t, last, db, ends_with_template_args)
				if !t1.equals(&t) {
					if db.names_size() < 2 {
						return cs
					}
					
					name := db.names_back().move_full()
					db.names_pop_back()
					db.names_back().first += "::" + name
					cs.Pos = t1.Pos
				} else {
					db.names_pop_back()
				}			
			} 
		}
		break
		default:
		t1 := parse_name(&t, last, db, ends_with_template_args)
		if !t1.equals(&t) {
			curPos := parse_discriminator(&t1, last)
			if db.names_size() < 2 {
				return curPos
			}
			name := db.names_back().move_full()
			db.names_pop_back()
			db.names_back().first += "::" + name			
		} else {
			db.names_pop_back()
		}
		break
	}
	
	return cs
}

//   <operator-name>
//                   ::= aa    # &&            
//                   ::= ad    # & (unary)     
//                   ::= an    # &             
//                   ::= aN    # &=            
//                   ::= aS    # =             
//                   ::= cl    # ()            
//                   ::= cm    # ,             
//                   ::= co    # ~             
//                   ::= cv <type>    # (cast)        
//                   ::= da    # delete[]      
//                   ::= de    # * (unary)     
//                   ::= dl    # delete        
//                   ::= dv    # /             
//                   ::= dV    # /=            
//                   ::= eo    # ^             
//                   ::= eO    # ^=            
//                   ::= eq    # ==            
//                   ::= ge    # >=            
//                   ::= gt    # >             
//                   ::= ix    # []            
//                   ::= le    # <=            
//                   ::= li <source-name>  # operator ""
//                   ::= ls    # <<            
//                   ::= lS    # <<=           
//                   ::= lt    # <             
//                   ::= mi    # -             
//                   ::= mI    # -=            
//                   ::= ml    # *             
//                   ::= mL    # *=            
//                   ::= mm    # -- (postfix in <expression> context)           
//                   ::= na    # new[]
//                   ::= ne    # !=            
//                   ::= ng    # - (unary)     
//                   ::= nt    # !             
//                   ::= nw    # new           
//                   ::= oo    # ||            
//                   ::= or    # |             
//                   ::= oR    # |=            
//                   ::= pm    # ->*           
//                   ::= pl    # +             
//                   ::= pL    # +=            
//                   ::= pp    # ++ (postfix in <expression> context)
//                   ::= ps    # + (unary)
//                   ::= pt    # ->            
//                   ::= qu    # ?             
//                   ::= rm    # %             
//                   ::= rM    # %=            
//                   ::= rs    # >>            
//                   ::= rS    # >>=           
//                   ::= v <digit> <source-name>        # vendor extended operator
// line 2333
func parse_operator_name(first, last *CStyleString, db *Db) CStyleString {
	var cs CStyleString
	cs.Content = first.Content
	cs.Pos = first.Pos
	
	if last.calcDelta(first) < 2 {
		return cs
	}
	
	c0 := cs.curChar()
	c1 := cs.nextChar()
	if c0 == 'a' {
		switch (c1)	{
		case 'a':
			db.names_push_back("operator&&");
			cs.Pos += 2
			break;
		case 'd':
		case 'n':
			db.names_push_back("operator&");
			cs.Pos += 2
			break;
		case 'N':
			db.names_push_back("operator&=");
			cs.Pos += 2
			break;
		case 'S':
			db.names_push_back("operator=");
			cs.Pos += 2
			break;
		}
	} else if c0 == 'c' {
		switch (c1) {
		case 'l':
			db.names_push_back("operator()")
			cs.Pos += 2
			break
		case 'm':
			db.names_push_back("operator,")
			cs.Pos += 2
			break
		case 'o':
			db.names_push_back("operator~")
			cs.Pos += 2
			break
		case 'v':
			{
				tmpPos := &CStyleString{cs.Content, cs.Pos + 2}
				try_to_parse_template_args := db.try_to_parse_template_args
				db.try_to_parse_template_args = false
				t := parse_type(tmpPos, last, db)
				db.try_to_parse_template_args = try_to_parse_template_args
				if (!t.equals(tmpPos)) {
					if (db.names_empty()) {
						return cs
					}
					
					s := db.names_back().first
					db.names_back().first = "operator " + s
					db.parsed_ctor_dtor_cv = true
					cs.Pos = t.Pos
				}
			}
			break
		}
	} else if c0 == 'd' {
		switch (c1)	{
		case 'a':
			db.names_push_back("operator delete[]")
			cs.Pos += 2
			break
		case 'e':
			db.names_push_back("operator*")
			cs.Pos += 2
			break
		case 'l':
			db.names_push_back("operator delete")
			cs.Pos += 2
			break
		case 'v':
			db.names_push_back("operator/")
			cs.Pos += 2
			break
		case 'V':
			db.names_push_back("operator/=")
			cs.Pos += 2
			break
		}
	} else if c0 == 'e' {
		switch (c1) {
		case 'o':
			db.names_push_back("operator^")
			cs.Pos += 2
			break
		case 'O':
			db.names_push_back("operator^=")
			cs.Pos += 2
			break
		case 'q':
			db.names_push_back("operator==")
			cs.Pos += 2
			break
		}
	} else if c0 == 'g' {
		switch (c1)	{
		case 'e':
			db.names_push_back("operator>=")
			cs.Pos += 2
			break
		case 't':
			db.names_push_back("operator>")
			cs.Pos += 2
			break
		}
	} else if c0 == 'i' {
		if (c1 == 'x') {
			db.names_push_back("operator[]")
			cs.Pos += 2
		}
	} else if c0 == 'l' {
		switch (c1)	{
		case 'e':
			db.names_push_back("operator<=")
			cs.Pos += 2
			break
		case 'i':
			{
				tmpPos := &CStyleString{cs.Content, cs.Pos + 2}
				t := parse_source_name(tmpPos, last, db)
				if (!t.equals(tmpPos)) {
					if (db.names_empty()) {
						return cs
					}
					s := db.names_back().first
					db.names_back().first = "operator\"\" " + s
					cs.Pos = t.Pos
				}
			}
			break
		case 's':
			db.names_push_back("operator<<")
			cs.Pos += 2
			break
		case 'S':
			db.names_push_back("operator<<=")
			cs.Pos += 2
			break
		case 't':
			db.names_push_back("operator<")
			cs.Pos += 2
			break
		}
	} else if c0 == 'm' {
		switch (c1)	{
		case 'i':
			db.names_push_back("operator-")
			cs.Pos += 2
			break
		case 'I':
			db.names_push_back("operator-=")
			cs.Pos += 2
			break
		case 'l':
			db.names_push_back("operator*")
			cs.Pos += 2
			break
		case 'L':
			db.names_push_back("operator*=")
			cs.Pos += 2
			break
		case 'm':
			db.names_push_back("operator--")
			cs.Pos += 2
			break
		}
	} else if c0 == 'n' {
		switch (c1)	{
		case 'a':
			db.names_push_back("operator new[]")
			cs.Pos += 2
			break
		case 'e':
			db.names_push_back("operator!=")
			cs.Pos += 2
			break
		case 'g':
			db.names_push_back("operator-")
			cs.Pos += 2
			break
		case 't':
			db.names_push_back("operator!")
			cs.Pos += 2
			break
		case 'w':
			db.names_push_back("operator new")
			cs.Pos += 2
			break
		}
	} else if c0 == 'o' {
		switch (c1)	{
		case 'o':
			db.names_push_back("operator||")
			cs.Pos += 2
			break
		case 'r':
			db.names_push_back("operator|")
			cs.Pos += 2
			break
		case 'R':
			db.names_push_back("operator|=")
			cs.Pos += 2
			break
		}
	} else if c0 == 'p' {
		switch (c1)	{
		case 'm':
			db.names_push_back("operator->*")
			cs.Pos += 2
			break
		case 'l':
			db.names_push_back("operator+")
			cs.Pos += 2
			break
		case 'L':
			db.names_push_back("operator+=")
			cs.Pos += 2
			break
		case 'p':
			db.names_push_back("operator++")
			cs.Pos += 2
			break
		case 's':
			db.names_push_back("operator+")
			cs.Pos += 2
			break
		case 't':
			db.names_push_back("operator->")
			cs.Pos += 2
			break
		}
	} else if c0 == 'q' {
		if (c1 == 'u') {
			db.names_push_back("operator?")
			cs.Pos += 2
		}
	} else if c0 == 'r' {
		switch (c1)	{
		case 'm':
			db.names_push_back("operator%")
			cs.Pos += 2
			break
		case 'M':
			db.names_push_back("operator%=")
			cs.Pos += 2
			break
		case 's':
			db.names_push_back("operator>>")
			cs.Pos += 2
			break
		case 'S':
			db.names_push_back("operator>>=")
			cs.Pos += 2
			break
		}
	} else if c0 == 'v' {
		if (isNumberChar(c1)){
			tmpPos := &CStyleString{cs.Content, cs.Pos + 2}
			t := parse_source_name(tmpPos, last, db)
			if (!t.equals(tmpPos)) {
				if (db.names_empty()) {
					return cs
				}
				
				s := db.names_back().first
				db.names_back().first = "operator " + s
				cs.Pos = t.Pos
			}
		}
	}
	
	return cs
}

// <unnamed-type-name> ::= Ut [ <nonnegative number> ] _
//                     ::= <closure-type-name>
// 
// <closure-type-name> ::= Ul <lambda-sig> E [ <nonnegative number> ] _ 
// 
// <lambda-sig> ::= <parameter type>+  # Parameter types or "v" if the lambda has no parameters
func parse_unnamed_type_name(first, last *CStyleString, db *Db) CStyleString {
	var cs CStyleString
	cs.Content = first.Content
	cs.Pos = first.Pos
	
	if last.calcDelta(first) <= 2 {
		return cs
	}
	
	if first.curChar() != 'U' {
		return cs
	}
	
	ty := cs.nextChar()
	if ty == 't' {
		db.names_push_back("'unnamed")
		t0 := cs
		t0.Pos += 2
		if t0.equals(last) {
			db.names_pop_back()
			return cs
		}
		
		if isNumberChar(t0.curChar()) {
			t1 := t0
			t1.Pos = t0.Pos + 1
			
			for !t1.equals(last) && isNumberChar(t1.curChar()) {
				t1.Pos++
			}
			
			// TODO 
			// db.names.back().first.append(t0, t1);
			db.names_back().first += t1.Content[t0.Pos:t1.Pos]
			t0.Pos = t1.Pos
		}
		
		db.names_back().first += "'"
		if t0.equals(last) || (t0.curChar() != '_') {
			db.names_pop_back()
			return cs
		}
		
		cs.Pos = t0.Pos + 1
	} else if ty == 'l' {
		db.names_push_back("'lambda'(")
		t0 := &CStyleString{cs.Content, cs.Pos}
		t0.Pos = cs.Pos + 2
		if t0.curChar() == 'v' {
			db.names_back().first += ")"
			t0.Pos++
		} else {
			t1 := parse_type(t0, last, db)
			if t1.equals(t0) {
				db.names_pop_back()
				return cs
			}
			
			if db.names_size() < 2 {
				return cs
			}
			
			tmp := db.names_back().move_full()
			db.names_pop_back()
			db.names_back().first += tmp
			t0.Pos = t1.Pos
			
			for {
				t1 = parse_type(t0, last, db)
				if t1.equals(t0) {
					break
				}
				if db.names_size() < 2 {
					return cs
				}
				tmp = db.names_back().move_full()
				db.names_pop_back()
				if len(tmp) > 0 {
					db.names_back().first += ", " + tmp					
				}
				t0.Pos = t1.Pos
			}
			db.names_back().first += ")"			
		}
		
		if t0.equals(last) || (t0.curChar() != 'E') {
			db.names_pop_back()
			return cs
		}
		
		t0.Pos++
		if t0.equals(last) {
			db.names_pop_back()
			return cs
		}
		
		if isNumberChar(t0.curChar()) {
			t1 := &CStyleString{t0.Content, t0.Pos + 1}
			for !t1.equals(last) && isNumberChar(t1.curChar()) {
				t1.Pos++
			}
			
			part := t1.Content[t0.Pos:t1.Pos]
			s := db.names_back().first
			db.names_back().first = s[:7] + part + s[7:]
			t0.Pos = t1.Pos
		}
		
		if t0.equals(last) || (t0.curChar() != '_') {
			db.names_pop_back()
			return cs
		}
		
		cs.Pos = t0.Pos + 1
	}
	
	return cs
}

// <ctor-dtor-name> ::= C1    # complete object constructor
//                  ::= C2    # base object constructor
//                  ::= C3    # complete object allocating constructor
//   extension      ::= C5    # ?
//                  ::= D0    # deleting destructor
//                  ::= D1    # complete object destructor
//                  ::= D2    # base object destructor
//   extension      ::= D5    # ?
func parse_ctor_dtor_name(first, last *CStyleString, db *Db) CStyleString {
	var cs CStyleString
	cs.Content = first.Content
	cs.Pos = first.Pos
	
	if last.calcDelta(first) < 2 {
		return *first
	}
	
	if db.names_empty() {
		return *first
	}
	
	c := cs.curChar()
	
	if c == 'C' {
		nc := cs.nextChar()
		if (nc == '1') || (nc == '2') || (nc == '3') || (nc == '5') {
			if db.names_empty() {
				return cs
			}
			
			db.names_push_back(db.names_back().first)
			cs.Pos += 2
			db.parsed_ctor_dtor_cv = true
		}
	} else if c == 'D' {
		nc := cs.nextChar()
		if (nc == '0') || (nc == '1') || (nc == '2') || (nc == '5') {
			if db.names_empty() {
				return cs
			}
			
			db.names_push_back("~" + db.names_back().first)
			cs.Pos += 2
			db.parsed_ctor_dtor_cv = true
		}
	}
	
	return cs
}

// <unscoped-name> ::= <unqualified-name>
//                 ::= St <unqualified-name>   # ::std::
// extension       ::= StL<unqualified-name>
func parse_unscoped_name(first, last *CStyleString, db *Db) CStyleString {
	var cs CStyleString
	cs.Content = first.Content
	cs.Pos = first.Pos
	
	if last.calcDelta(first) < 2 {
		return *first
	}
	
	t0 := *first
	St := false
	if (t0.curChar() == 'S') && (t0.nextChar() == 't') {
		t0.Pos += 2
		St = true
		if !t0.equals(last) && (t0.curChar() == 'L') {
			t0.Pos++
		}
	}
	
	t1 := parse_unqualified_name(&t0, last, db)
	if t1.equals(&t0) {
		if St {
			if db.names_empty() {
				return cs
			}
			
			db.names_back().first = "std::" + db.names_back().first
		}
		
		cs.Pos = t1.Pos
	}
	
	return cs
}

// <name> ::= <nested-name> // N
//        ::= <local-name> # See Scope Encoding below  // Z
//        ::= <unscoped-template-name> <template-args>
//        ::= <unscoped-name>

// <unscoped-template-name> ::= <unscoped-name>
//                          ::= <substitution>
// line 4174
func parse_name(first, last *CStyleString, db *Db, ends_with_template_args *bool) CStyleString {
	var cs CStyleString
	cs.Content = first.Content
	cs.Pos = first.Pos
	
	if last.calcDelta(first) < 2 {
		return cs
	}
	
	t0 := first
	if t0.curChar() == 'L' {
		t0.Pos++
	}
	
	if t0.curChar() == 'N' {
		t1 := parse_nested_name(t0, last, db, ends_with_template_args)
		
		if t1.Pos != t0.Pos {
			cs.Pos = t1.Pos
		}
	} else if t0.curChar() == 'Z' {
		t1 := parse_local_name(t0, last, db, ends_with_template_args)
		
		if t1.Pos != t0.Pos {
			cs.Pos = t1.Pos
		}
	} else {
		t1 := parse_unscoped_name(t0, last, db)
		
		if t1.Pos != t0.Pos {
			if (t1.Pos != last.Pos) && (t1.curChar() == 'I') {
				if db.names_empty() {
					return cs
				}
				
				db.subs_push_back_pair(*db.names_back())
				t0.Pos = t1.Pos
				t1 = parse_template_args(t0, last, db)
				if !t1.equals(t0) {
					if db.names_size() < 2 {
						return cs
					}
					
					tmp := db.names_back().move_full()
					db.names_pop_back()
					db.names_back().first += tmp
					cs.Pos = t1.Pos
					if ends_with_template_args != nil {
						*ends_with_template_args = true
					}
				}
			} else {
				// <unscoped-name>
				cs.Pos = t1.Pos
			}
		} else {
			// try <substitution> <template-args>
			t1 = parse_substitution(t0, last, db)
			if !t1.equals(t0) && !t1.equals(last) && (t1.curChar() == 'I') {
				t0.Pos = t1.Pos
				t1 = parse_template_args(t0, last, db)
				
				if !t1.equals(t0) {
					if db.names_size() < 2 {
						return cs
					}
					
					tmp := db.names_back().move_full()
					db.names_pop_back()
					db.names_back().first += tmp
					cs.Pos = t1.Pos
					if ends_with_template_args != nil {
						*ends_with_template_args = true
					}
				}
			}
		}
	}
	
	return cs
}

func parse_substitution(first, last *CStyleString, db *Db) CStyleString {
	var cs CStyleString
	cs.Content = first.Content
	cs.Pos = first.Pos
	
	// TODO
	return cs
}

func parse_unresolved_type(first, last *CStyleString, db *Db) CStyleString {
	var cs CStyleString
	cs.Content = first.Content
	cs.Pos = first.Pos
	
	// TODO
	return cs
}

// <builtin-type> ::= v    # void
//                ::= w    # wchar_t
//                ::= b    # bool
//                ::= c    # char
//                ::= a    # signed char
//                ::= h    # unsigned char
//                ::= s    # short
//                ::= t    # unsigned short
//                ::= i    # int
//                ::= j    # unsigned int
//                ::= l    # long
//                ::= m    # unsigned long
//                ::= x    # long long, __int64
//                ::= y    # unsigned long long, __int64
//                ::= n    # __int128
//                ::= o    # unsigned __int128
//                ::= f    # float
//                ::= d    # double
//                ::= e    # long double, __float80
//                ::= g    # __float128
//                ::= z    # ellipsis
//                ::= Dd   # IEEE 754r decimal floating point (64 bits)
//                ::= De   # IEEE 754r decimal floating point (128 bits)
//                ::= Df   # IEEE 754r decimal floating point (32 bits)
//                ::= Dh   # IEEE 754r half-precision floating point (16 bits)
//                ::= Di   # char32_t
//                ::= Ds   # char16_t
//                ::= Da   # auto (in dependent new-expressions)
//                ::= Dc   # decltype(auto)
//                ::= Dn   # std::nullptr_t (i.e., decltype(nullptr))
//                ::= u <source-name>    # vendor extended type
func parse_builtin_type(first, last *CStyleString, db *Db) CStyleString {
	var cs CStyleString
	cs.Content = first.Content
	cs.Pos = first.Pos
	
	if first.Pos == last.Pos {
		return cs
	}
	
	cs.Pos++
	
	c := first.curChar()
	if c == 'v' {
		db.names_push_back("void")
	} else if c == 'w' {
		db.names_push_back("wchar_t")
	} else if c == 'b' {
		db.names_push_back("bool")
	} else if c == 'c' {
		db.names_push_back("char")
	} else if c == 'a' {
		db.names_push_back("signed char")
	} else if c == 'h' {
		db.names_push_back("unsigned char")
	} else if c == 's' {
		db.names_push_back("short")
	} else if c == 't' {
		db.names_push_back("unsigned short")
	} else if c == 'i' {
		db.names_push_back("int")
	} else if c == 'j' {
		db.names_push_back("unsigned int")
	} else if c == 'l' {
		db.names_push_back("long")
	} else if c == 'm' {
		db.names_push_back("unsigned long")
	} else if c == 'x' {
		db.names_push_back("long long")
	} else if c == 'y' {
		db.names_push_back("unsigned long long")
	} else if c == 'n' {
		db.names_push_back("__int128")
	} else if c == 'o' {
		db.names_push_back("unsigned __int128")
	} else if c == 'f' {
		db.names_push_back("float")
	} else if c == 'd' {
		db.names_push_back("double")
	} else if c == 'e' {
		db.names_push_back("long double")
	} else if c == 'g' {
		db.names_push_back("__float128")
	} else if c == 'z' {
		db.names_push_back("...")
	} else if c == 'u' {
		t := parse_source_name(&cs, last, db)
		if t.Pos != cs.Pos {
			cs.Pos = t.Pos
		}
	} else if c == 'D' {
		if (first.Pos + 1) == last.Pos {
			return *first
		}
		
		nc := first.nextChar()
		if nc == 'd' {
			db.names_push_back("decimal64")
			cs.Pos = first.Pos + 2
		} else if nc == 'e' {
			db.names_push_back("decimal128")
			cs.Pos = first.Pos + 2
		} else if nc == 'f' {
			db.names_push_back("decimal32")
			cs.Pos = first.Pos + 2
		} else if nc == 'h' {
			db.names_push_back("decimal16")
			cs.Pos = first.Pos + 2
		} else if nc == 'i' {
			db.names_push_back("char32_t")
			cs.Pos = first.Pos + 2
		} else if nc == 's' {
			db.names_push_back("char16_t")
			cs.Pos = first.Pos + 2
		} else if nc == 'a' {
			db.names_push_back("auto")
			cs.Pos = first.Pos + 2
		} else if nc == 'c' {
			db.names_push_back("decltype(auto)")
			cs.Pos = first.Pos + 2
		} else if nc == 'n' {
			db.names_push_back("std::nullptr_t")
			cs.Pos = first.Pos + 2
		} else {
			
		}
	}
	
	return cs
}

// <expression> ::= <unary operator-name> <expression>
//              ::= <binary operator-name> <expression> <expression>
//              ::= <ternary operator-name> <expression> <expression> <expression>
//              ::= cl <expression>+ E                                   # call
//              ::= cv <type> <expression>                               # conversion with one argument
//              ::= cv <type> _ <expression>* E                          # conversion with a different number of arguments
//              ::= [gs] nw <expression>* _ <type> E                     # new (expr-list) type
//              ::= [gs] nw <expression>* _ <type> <initializer>         # new (expr-list) type (init)
//              ::= [gs] na <expression>* _ <type> E                     # new[] (expr-list) type
//              ::= [gs] na <expression>* _ <type> <initializer>         # new[] (expr-list) type (init)
//              ::= [gs] dl <expression>                                 # delete expression
//              ::= [gs] da <expression>                                 # delete[] expression
//              ::= pp_ <expression>                                     # prefix ++
//              ::= mm_ <expression>                                     # prefix --
//              ::= ti <type>                                            # typeid (type)
//              ::= te <expression>                                      # typeid (expression)
//              ::= dc <type> <expression>                               # dynamic_cast<type> (expression)
//              ::= sc <type> <expression>                               # static_cast<type> (expression)
//              ::= cc <type> <expression>                               # const_cast<type> (expression)
//              ::= rc <type> <expression>                               # reinterpret_cast<type> (expression)
//              ::= st <type>                                            # sizeof (a type)
//              ::= sz <expression>                                      # sizeof (an expression)
//              ::= at <type>                                            # alignof (a type)
//              ::= az <expression>                                      # alignof (an expression)
//              ::= nx <expression>                                      # noexcept (expression)
//              ::= <template-param>
//              ::= <function-param>
//              ::= dt <expression> <unresolved-name>                    # expr.name
//              ::= pt <expression> <unresolved-name>                    # expr->name
//              ::= ds <expression> <expression>                         # expr.*expr
//              ::= sZ <template-param>                                  # size of a parameter pack
//              ::= sZ <function-param>                                  # size of a function parameter pack
//              ::= sp <expression>                                      # pack expansion
//              ::= tw <expression>                                      # throw expression
//              ::= tr                                                   # throw with no operand (rethrow)
//              ::= <unresolved-name>                                    # f(p), N::f(p), ::f(p),
//                                                                       # freestanding dependent name (e.g., T::x),
//                                                                       # objectless nonstatic member reference
//              ::= <expr-primary>
// line 3299
func parse_expression(first, last *CStyleString, db *Db) CStyleString {
	var cs CStyleString
	cs.Content = first.Content
	cs.Pos = first.Pos
	
	delta := last.Pos - first.Pos
	if delta < 2 {
		return cs
	}
	
	var t CStyleString
	t.Content = first.Content
	t.Pos = first.Pos
	
	parsed_gs := false
	
	if (last.calcDelta(first) >= 4) &&
	    (t.curChar() == 'g') && 
		(t.nextChar() == 's') {
		t.Pos += 2
		parsed_gs = true	
	}
	
	switch (t.curChar()) {
		case 'L':
        cs = parse_expr_primary(first, last, db)
        break
        case 'T':
        cs = parse_template_param(first, last, db)
        break
        case 'f':
        cs = parse_function_param(first, last, db)
        break
		case 'a': {
			tmpPos := &CStyleString{cs.Content, cs.Pos + 2}
			switch t.nextChar() {
				case 'a':
				t = parse_binary_expression(tmpPos, last, "&&", db)
	            if (!t.equals(tmpPos)) {
					cs.Pos = t.Pos
				}            
	            break
				case 'd':
				t = parse_prefix_expression(tmpPos, last, "&", db)
	            if (!t.equals(tmpPos)) {
					cs.Pos = t.Pos
				}            
	            break
				case 'n':
				t = parse_binary_expression(tmpPos, last, "&", db)
	            if (!t.equals(tmpPos)) {
					cs.Pos = t.Pos
				}            
	            break
				case 'N':
				t = parse_binary_expression(tmpPos, last, "&=", db)
	            if (!t.equals(tmpPos)) {
					cs.Pos = t.Pos
				}            
	            break
				case 'S':
				t = parse_binary_expression(tmpPos, last, "=", db)
	            if (!t.equals(tmpPos)) {
					cs.Pos = t.Pos
				}            
	            break
				case 't':
				cs = parse_alignof_type(first, last, db)
				break
				case 'z':
	            cs = parse_alignof_expr(first, last, db)
	            break
				// TODO
			}
		}
		break
		case 'c': {
			tmpPos := &CStyleString{cs.Content, cs.Pos + 2}
			switch t.nextChar() {
				case 'c':
                	cs = *parse_const_cast_expr(first, last, db)
                	break
	            case 'l':
	                cs = parse_call_expr(first, last, db)
	                break
	            case 'm':
	                t := parse_binary_expression(tmpPos, last, ",", db)
	                if !t.equals(tmpPos) {
						cs.Pos = t.Pos
					}	                    
	                break
	            case 'o':
	                t := parse_prefix_expression(tmpPos, last, "~", db)
	                if !t.equals(tmpPos) {
						cs.Pos = t.Pos
					}	                    
	                break
	            case 'v':
	                cs = parse_conversion_expr(first, last, db)
	                break
			}
			
		}
		break
		case 'd': {
			tmpPos := &CStyleString{cs.Content, cs.Pos + 2}
			
			switch t.nextChar() {
				case 'a': {
                    t1 := parse_expression(tmpPos, last, db)
                    if (!t1.equals(tmpPos)) {
                        if (db.names_empty()) {
							return cs
						}
                        
						tmpstr := ""
						if parsed_gs {
							tmpstr = "::"
						}
                        db.names_back().first = tmpstr + "delete[] " + 
						    db.names_back().move_full()
                        cs.Pos = t1.Pos
                    }
                }
                break
            case 'c':
                cs = parse_dynamic_cast_expr(first, last, db)
                break
            case 'e':
                t = parse_prefix_expression(tmpPos, last, "*", db)
                if (!t.equals(tmpPos)) {
					cs.Pos = t.Pos
				}
                    
                break
            case 'l': {
					t.Pos += 2
                    t1 := parse_expression(&t, last, db)
                    if (!t1.equals(&t)) {
                        if (db.names_empty()) {
							return cs
						}
                        tmpstr := ""
						if parsed_gs {
							tmpstr = "::"
						}
                        db.names_back().first = tmpstr +
                            "delete " + db.names_back().move_full()
                        cs.Pos = t1.Pos
                    }
                }
                break
            case 'n':
                return parse_unresolved_name(first, last, db)
            case 's':
                cs = *parse_dot_star_expr(first, last, db)
                break
            case 't':
                cs = *parse_dot_expr(first, last, db)
                break
            case 'v':
                t = parse_binary_expression(tmpPos, last, "/", db)
                if (!t.equals(tmpPos)) {
					cs.Pos = t.Pos
				}                    
                break
            case 'V':
                t = parse_binary_expression(tmpPos, last, "/=", db)
                if (!t.equals(tmpPos)) {
					cs.Pos = t.Pos
				}                    
                break
			}
		}
		break
		case 'e': {
			tmpPos := &CStyleString{cs.Content, cs.Pos + 2}
			switch t.nextChar() {
				case 'o':
                t = parse_binary_expression(tmpPos, last, "^", db)
                if !t.equals(tmpPos) {
					cs.Pos = t.Pos
				}                    
                break
            case 'O':
                t = parse_binary_expression(tmpPos, last, "^=", db)
                if !t.equals(tmpPos) {
					cs.Pos = t.Pos
				}
                break
            case 'q':
                t = parse_binary_expression(tmpPos, last, "==", db)
                if !t.equals(tmpPos) {
					cs.Pos = t.Pos
				}
                break
			}
		}
		break
		case 'g': {
			tmpPos := &CStyleString{cs.Content, cs.Pos + 2}
			switch t.nextChar() {
			case 'e':
                t = parse_binary_expression(tmpPos, last, ">=", db)
                if !t.equals(tmpPos) {
					cs.Pos = t.Pos
				}                    
                break
            case 't':
                t = parse_binary_expression(tmpPos, last, ">", db)
                if !t.equals(tmpPos) {
					cs.Pos = t.Pos
				}
                break
			}
		}
		break
		case 'i': {
			tmpPos := &CStyleString{cs.Content, cs.Pos + 2}
			if t.nextChar() == 'x' {
                t1 := parse_expression(tmpPos, last, db)
                if !t1.equals(tmpPos) {
					t2 := parse_expression(&t1, last, db)
					if !t2.equals(&t1) {
						if (db.names_size() < 2) {
							return cs
						}
						
						op2 := db.names_back().move_full()
                        db.names_pop_back()
                        op1 := db.names_back().move_full()
                        db.names_back().first = "(" + op1 + ")[" + op2 + "]"
                        cs.Pos = t2.Pos
					}
				} else {
					db.names_pop_back()
				}
			}
		}
		break
		case 'l': {
			tmpPos := &CStyleString{cs.Content, cs.Pos + 2}
			switch t.nextChar() {
			case 'e':
                t = parse_binary_expression(tmpPos, last, "<=", db)
                if !t.equals(tmpPos) {
					cs.Pos = t.Pos
				}                    
                break
            case 's':
                t = parse_binary_expression(tmpPos, last, "<<", db)
                if !t.equals(tmpPos) {
					cs.Pos = t.Pos
				}
                break
            case 'S':
                t = parse_binary_expression(tmpPos, last, "<<=", db)
                if !t.equals(tmpPos) {
					cs.Pos = t.Pos
				}
                break
			case 't':
                t = parse_binary_expression(tmpPos, last, "<", db)
                if !t.equals(tmpPos) {
					cs.Pos = t.Pos
				}
                break
			}
		}
		break
		case 'm': {
			tmpPos := &CStyleString{cs.Content, cs.Pos + 2}
			switch t.nextChar() {
			case 'i':
                t = parse_binary_expression(tmpPos, last, "-", db)
                if !t.equals(tmpPos) {
					cs.Pos = t.Pos
				}                    
                break
            case 'I':
                t = parse_binary_expression(tmpPos, last, "-=", db)
                if !t.equals(tmpPos) {
					cs.Pos = t.Pos
				}
                break
            case 'l':
                t = parse_binary_expression(tmpPos, last, "*", db)
                if !t.equals(tmpPos) {
					cs.Pos = t.Pos
				}
                break
			case 'L':
                t = parse_binary_expression(tmpPos, last, "*=", db)
                if !t.equals(tmpPos) {
					cs.Pos = t.Pos
				}
                break
			case 'm':
				if !tmpPos.equals(last) && (tmpPos.curChar() == '_') {
					tmpPos.Pos++
					t = parse_prefix_expression(tmpPos, last, "--", db)
					if !t.equals(tmpPos) {
						cs.Pos = t.Pos
					}
				} else {
					t1 := parse_expression(tmpPos, last, db)
					if !t1.equals(tmpPos) {
						if (db.names_empty()) {
							return cs
						}
						
						db.names_back().first = "(" + db.names_back().move_full() +
						     ")--"
                        cs.Pos = t1.Pos
					}
				}
				break
			}			
		}
		break
		case 'n': {
			tmpPos := &CStyleString{cs.Content, cs.Pos + 2}
			switch t.nextChar() {
				case 'a', 'w':
				cs = *parse_new_expr(first, last, db)
				break
				case 'e':
				t = parse_binary_expression(tmpPos, last, "!=", db)
                if !t.equals(tmpPos) {
					cs.Pos = t.Pos
				}
                break
				case 'g':
				t = parse_prefix_expression(tmpPos, last, "-", db)
                if !t.equals(tmpPos) {
					cs.Pos = t.Pos
				}
                break
				case 't':
				t = parse_prefix_expression(tmpPos, last, "!", db)
                if !t.equals(tmpPos) {
					cs.Pos = t.Pos
				}
                break
				case 'x':
				t = *parse_noexcept_expression(tmpPos, last, db)
                if !t.equals(tmpPos) {
					cs.Pos = t.Pos
				}
                break
			}
		}
		break
		case 'o':
		tmpPos := &CStyleString{cs.Content, cs.Pos + 2}
		switch t.nextChar() {
			case 'n':
                return parse_unresolved_name(first, last, db)
			case 'o':
				t = parse_binary_expression(tmpPos, last, "||", db)
                if !t.equals(tmpPos) {
					cs.Pos = t.Pos
				}
                break
			case 'r':
				t = parse_binary_expression(tmpPos, last, "|", db)
                if !t.equals(tmpPos) {
					cs.Pos = t.Pos
				}
                break
			case 'R':
				t = parse_binary_expression(tmpPos, last, "|=", db)
                if !t.equals(tmpPos) {
					cs.Pos = t.Pos
				}
                break
		}
		break
		case 'p': {
			tmpPos := &CStyleString{cs.Content, cs.Pos + 2}
			switch t.nextChar() {
				case 'm':
				t = parse_binary_expression(tmpPos, last, "->*", db)
                if !t.equals(tmpPos) {
					cs.Pos = t.Pos
				}
                break
				case 'l':
				t = parse_binary_expression(tmpPos, last, "+", db)
                if !t.equals(tmpPos) {
					cs.Pos = t.Pos
				}
                break
				case 'L':
				t = parse_binary_expression(tmpPos, last, "+=", db)
                if !t.equals(tmpPos) {
					cs.Pos = t.Pos
				}
                break
				case 's':
				t = parse_prefix_expression(tmpPos, last, "+", db)
                if !t.equals(tmpPos) {
					cs.Pos = t.Pos
				}
                break
				case 'p': {
					if !tmpPos.equals(last) && (tmpPos.curChar() == '_') {
						tmpPos.Pos++
						t = parse_prefix_expression(tmpPos, last, "++", db)
						if !t.equals(tmpPos) {
							cs.Pos = t.Pos
						}
					} else {
						t1 := parse_expression(tmpPos, last, db)
						if !t1.equals(tmpPos) {
							if (db.names_empty()) {
								return cs
							}
							
							db.names_back().first = "(" + db.names_back().move_full() +
							     ")++"
	                        cs.Pos = t1.Pos
						}
					}					
				}
				break
				case 't':
                cs = parse_arrow_expr(first, last, db)
                break
			}
		}
		break
		case 'q':
		tmpPos := &CStyleString{cs.Content, cs.Pos + 2}
		if t.nextChar() == 'u' {
			t1 := parse_expression(tmpPos, last, db)
			if !t1.equals(tmpPos) {
				t2 := parse_expression(&t1, last, db)
				if !t2.equals(&t1) {
					t3 := parse_expression(&t2, last, db)
					if !t3.equals(&t2) {
						if (db.names_size() < 3) {
							return cs
						}
						
						op3 := db.names_back().move_full()
                        db.names_pop_back()
                        op2 := db.names_back().move_full()
                        db.names_pop_back()
                        op1 := db.names_back().move_full()
                        db.names_back().first = "(" + op1 + ") ? (" + op2 +
						     ") : (" + op3 + ")"
						cs.Pos = t3.Pos
					} else {
						db.names_pop_back()
						db.names_pop_back()
					}
				} else {
					db.names_pop_back()
				}
			}
		}
		break
		case 'r': {
			tmpPos := &CStyleString{cs.Content, cs.Pos + 2}
			switch t.nextChar() {
				case 'c':
				cs = *parse_reinterpret_cast_expr(first, last, db)
				break
				case 'm':
				t = parse_binary_expression(tmpPos, last, "%", db)
                if !t.equals(tmpPos) {
					cs.Pos = t.Pos
				}
                break
				case 'M':
				t = parse_binary_expression(tmpPos, last, "%=", db)
                if !t.equals(tmpPos) {
					cs.Pos = t.Pos
				}
                break
				case 's':
				t = parse_binary_expression(tmpPos, last, ">>", db)
                if !t.equals(tmpPos) {
					cs.Pos = t.Pos
				}
                break
				case 'S':
				t = parse_binary_expression(tmpPos, last, ">>==", db)
                if !t.equals(tmpPos) {
					cs.Pos = t.Pos
				}
                break
			}
		}
		break
		case 's': {
			switch t.curChar() {
				case 'c':
				cs = *parse_static_cast_expr(first, last, db)
                break
				case 'p':
                cs = parse_pack_expansion(first, last, db)
                break
            	case 'r':
                return parse_unresolved_name(first, last, db)
            	case 't':
                cs = *parse_sizeof_type_expr(first, last, db)
                break
            	case 'z':
                cs = *parse_sizeof_expr_expr(first, last, db)
                break
				case 'Z': {
					if last.calcDelta(&t) >= 3 {
						t.Pos += 2
						switch t.curChar() {
						case 'T':
							cs = *parse_sizeof_param_pack_expr(first, last, db)
                        	break
                    	case 'f':
                        	cs = *parse_sizeof_function_param_pack_expr(first, last, db)
                        	break
						}
					}
				}
				break
			}
		}
		break
		case 't': {
			switch t.nextChar() {
				case 'e', 'i':
				cs = *parse_typeid_expr(first, last, db)
                break
				case 'r':
				db.names_push_back("throw")
				cs.Pos += 2
				break
				case 'w':
				cs = *parse_throw_expr(first, last, db)
				break
			}
		}
		break
		case '1','2','3','4','5','6','7','8','9':
		return parse_unresolved_name(first, last, db)
	}
	
	return cs
}

func parse_function_param(first, last *CStyleString, db *Db) CStyleString {
	var cs CStyleString
	cs.Content = first.Content
	cs.Pos = first.Pos
	
	// TODO
	return cs
}

func parse_prefix_expression(first, last *CStyleString, op string, db *Db) CStyleString {
	var cs CStyleString
	cs.Content = first.Content
	cs.Pos = first.Pos
	
	t1 := parse_expression(first, last, db)
	if !t1.equals(first) {
		if (db.names_empty()) {
			return cs
		}
		
		db.names_back().first =  op + "(" + db.names_back().move_full() + ")"
        cs.Pos = t1.Pos
	}
	return cs
}

func parse_binary_expression(first, last *CStyleString, op string, db *Db) CStyleString {
	var cs CStyleString
	cs.Content = first.Content
	cs.Pos = first.Pos
	
	t1 := parse_expression(first, last, db)
	if !t1.equals(first) {
		t2 := parse_expression(&t1, last, db)
		if !t2.equals(&t1) {
			if (db.names_size() < 2) {
				return cs
			}
			
			op2 := db.names_back().move_full()
	        db.names_pop_back()
	        op1 := db.names_back().move_full()
			db.names_back().first = ""
			if (op == ">") {
				db.names_back().first += "("
			}
			db.names_back().first += "(" + op1 + ") " + op + " (" + op2 + ")"
	        if (op == ">") {
				db.names_back().first += ")"
			}
			
			cs.Pos = t2.Pos
		}	
	} else {
		db.names_pop_back()
	}
	return cs
}

// pt <expression> <expression>                    # expr->name
// line 1549
func parse_arrow_expr(first, last *CStyleString, db *Db) CStyleString {
	var cs CStyleString
	cs.Content = first.Content
	cs.Pos = first.Pos
	
	if last.calcDelta(first) < 3 {
		return cs
	}
	
	if (first.curChar() != 'p') || (first.nextChar() != 't') {
		return cs
	}
	
	tmpPos := &CStyleString{cs.Content, cs.Pos + 2}
	t := parse_expression(tmpPos, last, db)
	if !t.equals(tmpPos) {
		t1 := parse_expression(&t, last, db)
		
		if !t1.equals(&t) {
			if db.names_size() < 2 {
				return cs
			}
			
			tmp := db.names_back().move_full()
			db.names_pop_back()
            db.names_back().first += "->"
            db.names_back().first += tmp
            cs.Pos = t1.Pos
		}
	}
	
	return cs
}

// rc <type> <expression>       # reinterpret_cast<type> (expression)
func parse_reinterpret_cast_expr(first, last *CStyleString, db *Db) *CStyleString {
	cs := &CStyleString{first.Content, first.Pos}

	if last.calcDelta(first) < 3 {
		return cs
	}
	
	if first.prefix(2) != "rc" {
		return cs
	}

	tmpPos := &CStyleString{cs.Content, cs.Pos + 2}
	t := parse_type(tmpPos, last, db)
	
	if !t.equals(tmpPos) {
		t1 := parse_expression(&t, last, db)
		
		if !t1.equals(&t) {
			if db.names_size() < 2 {
				return cs
			}
			
			expr := db.names_back().move_full()
            db.names_pop_back()
            db.names_back().first = "reinterpret_cast<" + 
			    db.names_back().move_full() + ">(" + expr + ")"
			cs.Pos = t1.Pos
		}
	}
	
	return cs
}

// sc <type> <expression>         # static_cast<type> (expression)
func parse_static_cast_expr(first, last *CStyleString, db *Db) *CStyleString {
	cs := &CStyleString{first.Content, first.Pos}

	if last.calcDelta(first) < 3 {
		return cs
	}
	
	if first.prefix(2) != "sc" {
		return cs
	}

	tmpPos := &CStyleString{cs.Content, cs.Pos + 2}
	t := parse_type(tmpPos, last, db)
	
	if !t.equals(tmpPos) {
		t1 := parse_expression(&t, last, db)
		
		if !t1.equals(&t) {
			if db.names_size() < 2 {
				return cs
			}
			
			expr := db.names_back().move_full()
            db.names_pop_back()
            db.names_back().first = "static_cast<" + 
				db.names_back().move_full() + ">(" + expr + ")"
			cs.Pos = t1.Pos
		}
	}
	
	return cs
}

func parse_pack_expansion(first, last *CStyleString, db *Db) CStyleString {
	var cs CStyleString
	cs.Content = first.Content
	cs.Pos = first.Pos
	
	// TODO
	return cs
}

func parse_sizeof_type_expr(first, last *CStyleString, db *Db) *CStyleString {
	cs := &CStyleString{first.Content, first.Pos}

	if last.calcDelta(first) < 3 {
		return cs
	}
	
	if first.prefix(2) != "st"{
		return cs
	}
	
	tmpPos := &CStyleString{cs.Content, cs.Pos + 2}
	t := parse_type(tmpPos, last, db)
	
	if !t.equals(tmpPos) {
		if db.names_empty() {
			return cs
		}

        db.names_back().first = "sizeof (" + db.names_back().move_full() + ")"
		cs.Pos = t.Pos
	}
	
	return cs
}

// sz <expr>       # sizeof (a expression)
func parse_sizeof_expr_expr(first, last *CStyleString, db *Db) *CStyleString {
	cs := &CStyleString{first.Content, first.Pos}

	if last.calcDelta(first) < 3 {
		return cs
	}
	
	if first.prefix(2) != "sz"{
		return cs
	}
	
	tmpPos := &CStyleString{cs.Content, cs.Pos + 2}
	t := parse_expression(tmpPos, last, db)
	
	if !t.equals(tmpPos) {
		if db.names_empty() {
			return cs
		}

        db.names_back().first = "sizeof (" + db.names_back().move_full() + ")"
		cs.Pos = t.Pos
	}
	
	return cs
}

// sZ <template-param>                                  # size of a parameter pack
func parse_sizeof_param_pack_expr(first, last *CStyleString, db *Db) *CStyleString {
	cs := &CStyleString{first.Content, first.Pos}

	if last.calcDelta(first) < 3 {
		return cs
	}
	
	if first.prefix(3) != "sZT"{
		return cs
	}
	
	k0 := db.names_size()
	tmpPos := &CStyleString{cs.Content, cs.Pos + 2}
	t := parse_template_param(tmpPos, last, db)
	k1 := db.names_size()

	if !t.equals(tmpPos) {
		tmp := "sizeof...("
		k := k0
		if k != k1 {
			tmp += db.names.content[k].move_full()
			k++
			for k != k1 {
				tmp += ", " + db.names.content[k].move_full()
				k++
			}
		}
		tmp += ")"
		for k1 != k0 {
			db.names_pop_back()
			k1--
		}
		db.names_push_back(tmp)
		cs.Pos = t.Pos
	}
	
	return cs
}

// sZ <function-param>                                  # size of a function parameter pack
func parse_sizeof_function_param_pack_expr(first, last *CStyleString, db *Db) *CStyleString {
	cs := &CStyleString{first.Content, first.Pos}

	if last.calcDelta(first) < 3 {
		return cs
	}
	
	if first.prefix(3) != "sZf"{
		return cs
	}
	
	tmpPos := &CStyleString{cs.Content, cs.Pos + 2}
	t := parse_function_param(tmpPos, last, db)
	
	if !t.equals(tmpPos) {
		if db.names_empty() {
			return cs
		}

        db.names_back().first = "sizeof...(" + db.names_back().move_full() + ")"
		cs.Pos = t.Pos
	}
	
	return cs
}

// te <expression>                                      # typeid (expression)
// ti <type>                                            # typeid (type)
func parse_typeid_expr(first, last *CStyleString, db *Db) *CStyleString {
	cs := &CStyleString{first.Content, first.Pos}

	if last.calcDelta(first) < 3 {
		return cs
	}
	
	if (first.prefix(2) != "te") &&
	    (first.prefix(2) != "ti") {
		return cs
	}
	
	tmpPos := &CStyleString{cs.Content, cs.Pos + 2}
	var t CStyleString
	if cs.nextChar() == 'e' {
		t = parse_expression(tmpPos, last, db)
	} else {
		t = parse_type(tmpPos, last, db)
	}
	
	if !t.equals(tmpPos) {
		if db.names_empty() {
			return cs
		}

        db.names_back().first = "typeid(" + db.names_back().move_full() + ")"
		cs.Pos = t.Pos
	}
	
	return cs
}

func parse_throw_expr(first, last *CStyleString, db *Db) *CStyleString {
	cs := &CStyleString{first.Content, first.Pos}

	if last.calcDelta(first) < 3 {
		return cs
	}
	
	if first.prefix(2) != "tw" {
		return cs
	}

	tmpPos := &CStyleString{cs.Content, cs.Pos + 2}
	t := parse_expression(tmpPos, last, db)
	
	if !t.equals(tmpPos) {
		if db.names_empty() {
			return cs
		}
        db.names_back().first = "throw " + db.names_back().move_full()
		cs.Pos = t.Pos
	}
	
	return cs
}

// cc <type> <expression>                               # const_cast<type> (expression)
func parse_const_cast_expr(first, last *CStyleString, db *Db) *CStyleString {
	cs := &CStyleString{first.Content, first.Pos}

	if last.calcDelta(first) < 3 {
		return cs
	}
	
	if first.prefix(2) != "cc" {
		return cs
	}

	tmpPos := &CStyleString{cs.Content, cs.Pos + 2}
	t := parse_type(tmpPos, last, db)
	
	if !t.equals(tmpPos) {
		t1 := parse_expression(&t, last, db)
		
		if !t1.equals(&t) {
			if db.names_size() < 2 {
				return cs
			}
			
			expr := db.names_back().move_full()
            db.names_pop_back()
            db.names_back().first = "const_cast<" + db.names_back().move_full() +
			     ">(" + expr + ")"
			cs.Pos = t1.Pos
		}
	}
	
	return cs
}

// cl <expression>+ E      # call

func parse_call_expr(first, last *CStyleString, db *Db) CStyleString {
	var cs CStyleString
	cs.Content = first.Content
	cs.Pos = first.Pos
	
	if last.calcDelta(first) < 4 {
		return cs
	}
	
	if (first.curChar() != 'c') || (first.nextChar() != 'l') {
		return cs
	}
	
	tmpPos := &CStyleString{cs.Content, cs.Pos + 2}
	t := parse_expression(tmpPos, last, db)
	if !t.equals(tmpPos) {
		if t.equals(last) {
			return cs
		}
		if (db.names_empty()) {
			return cs
		}
		
		db.names_back().first += db.names_back().second
		db.names_back().second = ""
		db.names_back().first += "("
		first_expr := true
		for t.curChar() != 'E' {
			t1 := parse_expression(&t, last, db)
			if !t1.notEach(&t, last) {
				return cs
			}
			if (db.names_empty()) {
				return cs
			}
			tmp := db.names_back().move_full()
            db.names_pop_back()
			if len(tmp) > 0 {
				if db.names_empty() {
					return cs
				}
				if !first_expr {
					db.names_back().first += ", "
					first_expr = false
				}
				db.names_back().first += tmp
			}
			t.Pos = t1.Pos
		}
		t.Pos++
		if db.names_empty() {
			return cs
		}
		db.names_back().first += ")"
		cs.Pos = t.Pos
	}
	
	return cs
}

// <unresolved-name>
//  extension        ::= srN <unresolved-type> [<template-args>] <unresolved-qualifier-level>* E <base-unresolved-name>
//                   ::= [gs] <base-unresolved-name>                     # x or (with "gs") ::x
//                   ::= [gs] sr <unresolved-qualifier-level>+ E <base-unresolved-name>  
//                                                                       # A::x, N::y, A<T>::z; "gs" means leading "::"
//                   ::= sr <unresolved-type> <base-unresolved-name>     # T::x / decltype(p)::x
//  extension        ::= sr <unresolved-type> <template-args> <base-unresolved-name>
//                                                                       # T::N::x /decltype(p)::N::x
//  (ignored)        ::= srN <unresolved-type>  <unresolved-qualifier-level>+ E <base-unresolved-name>

// line 1131
func parse_unresolved_name(first, last *CStyleString, db *Db) CStyleString {
	var cs CStyleString
	cs.Content = first.Content
	cs.Pos = first.Pos
	
	if last.calcDelta(first) <= 2 {
		return cs
	}
	
	t := &CStyleString{cs.Content, cs.Pos}
	global := false
	if (t.curChar() == 'g') && (t.nextChar() == 's') {
		global = true
		t.Pos += 2
	}
	
	t2 := parse_base_unresolved_name(t, last, db)
	if !t2.equals(t) {
		if global {
			if db.names_empty() {
				return cs
			}
			
			db.names_back().first = "::" + db.names_back().first
		}
		
		cs.Pos = t2.Pos
	} else if (last.calcDelta(t) > 2) && (t.Content[t.Pos:t.Pos + 2] == "sr") {
		t.Pos += 2
		if t.curChar() == 'N' {
			t.Pos++
			t1 := parse_unresolved_type(t, last, db)
			if !t1.notEach(t, last) {
				return cs
			}
			
			t.Pos = t1.Pos
            t1 = parse_template_args(t, last, db)
			if !t1.equals(t) {
				if (db.names_size() < 2) {
					return cs
				}
				args := db.names_back().move_full()
                db.names_pop_back()
                db.names_back().first += (args)
                t.Pos = t1.Pos
                if t.equals(last) {
					db.names_pop_back()
                    return cs
				}
			}
			
			for t.curChar() != 'E' {
				t1 := parse_unresolved_qualifier_level(t, last, db)
				if t1.equals(t) || t1.equals(last) ||
				    (db.names_size() < 2) {
					return cs	
				}
				
				s := db.names_back().move_full()
				db.names_pop_back()
				db.names_back().first += "::" + (s)
				t.Pos = t1.Pos
			}
			
			t.Pos++
			t1 = parse_base_unresolved_name(t, last, db)
			if t1.equals(t) {
				if !db.names_empty() {
					db.names_pop_back()
				}
				
				return cs
			}
			if db.names_size() < 2 {
				return cs
			}
			s := db.names_back().move_full()
            db.names_pop_back()
            db.names_back().first += "::" + (s)
            cs.Pos = t1.Pos
		} else {
			t1 := parse_unresolved_type(t, last, db)
			if !t1.equals(t) {
				t.Pos = t1.Pos
                t1 = parse_template_args(t, last, db)
				if !t1.equals(t) {
					if db.names_size() < 2 {
						return cs
					}
					args := db.names_back().move_full()
            		db.names_pop_back()
            		db.names_back().first += "::" + (args)
            		t.Pos = t1.Pos
				}
				t1 = parse_base_unresolved_name(t, last, db)
				if t1.equals(t) {
					if !db.names_empty() {
						db.names_pop_back()
					}
					return cs
				}
				if db.names_size() < 2 {
					return cs
				}
				s := db.names_back().move_full()
	            db.names_pop_back()
	            db.names_back().first += "::" + (s)
	            cs.Pos = t1.Pos
			} else {
				t1 = parse_unresolved_qualifier_level(t, last, db)
				if t1.equals(t) || t1.equals(last) {
					return cs
				}
				
				t.Pos = t1.Pos
                if (global) {
					if db.names_empty() {
						return cs
					}
					db.names_back().first = "::" + db.names_back().first
				}
				
				for t.curChar() != 'E' {
					t1 = parse_unresolved_qualifier_level(t, last, db)
					if t1.equals(t) || t1.equals(last) || (db.names_size() < 2) {
						return cs
					}
					
					s := db.names_back().move_full()
	            	db.names_pop_back()
	            	db.names_back().first += "::" + (s)
	            	cs.Pos = t1.Pos
				}
				t.Pos++
				t1 = parse_base_unresolved_name(t, last, db)
				if t1.equals(t) {
					if !db.names_empty() {
						db.names_pop_back()
					}
					return cs
				}
				if (db.names_size() < 2) {
					return cs
				}
				s := db.names_back().move_full()
	            db.names_pop_back()
	            db.names_back().first += "::" + (s)
	            cs.Pos = t1.Pos
			}
		}
	}
	return cs
}

// <unresolved-qualifier-level> ::= <simple-id>
func parse_unresolved_qualifier_level(first, last *CStyleString, db *Db) CStyleString {
	return parse_simple_id(first, last, db)
}

func parse_simple_id(first, last *CStyleString, db *Db) CStyleString {
	var cs CStyleString
	cs.Content = first.Content
	cs.Pos = first.Pos
	
	if first.equals(last) {
		return cs
	}
	t := parse_source_name(first, last, db)
	if !t.equals(first) {
		t1 := parse_template_args(&t, last, db)
		if !t1.equals(&t) {
			if (db.names_size() < 2) {
				return cs
			}
			args := db.names_back().move_full()
            db.names_pop_back()
            db.names_back().first += (args)
		}
		cs.Pos = t1.Pos
	} else {
		cs.Pos = t.Pos
	}
	
	return cs
}

// <base-unresolved-name> ::= <simple-id>                                # unresolved name
//          extension     ::= <operator-name>                            # unresolved operator-function-id
//          extension     ::= <operator-name> <template-args>            # unresolved operator template-id
//                        ::= on <operator-name>                         # unresolved operator-function-id
//                        ::= on <operator-name> <template-args>         # unresolved operator template-id
//                        ::= dn <destructor-name>                       # destructor or pseudo-destructor;
//     

// line 1055
func parse_base_unresolved_name(first, last *CStyleString, db *Db) CStyleString {
	var cs CStyleString
	cs.Content = first.Content
	cs.Pos = first.Pos
	
	if last.calcDelta(first) < 2 {
		return cs
	}
	
	s := cs.prefix(2)
	if (s == "on") || (s == "dn") {
		tmpPos := &CStyleString{cs.Content, cs.Pos + 2}
		if s[0] == 'o' {			
			t := parse_operator_name(tmpPos, last, db)
			if !t.equals(tmpPos) {
				cs = parse_template_args(&t, last, db)
				if !cs.equals(&t) {
					if db.names_size() < 2 {
						return cs
					}
					args := db.names_back().move_full()
                    db.names_pop_back()
                    db.names_back().first += (args)
				}
			}
		} else {
			t := parse_destructor_name(tmpPos, last, db)
			if !t.equals(tmpPos) {
				cs.Pos = t.Pos
			}
		}
	} else {
		t := parse_simple_id(first, last, db)
		if t.equals(first) {
			t = parse_operator_name(first, last, db)
			if !t.equals(first) {
				cs = parse_template_args(&t, last, db)
				if !cs.equals(&t) {
					if db.names_size() < 2 {
						return cs
					}
					args := db.names_back().move_full()
                    db.names_pop_back()
                    db.names_back().first += (args)
				}
			}
		} else {
			cs.Pos = t.Pos
		}
	}
	
	return cs
}

// <destructor-name> ::= <unresolved-type>                               # e.g., ~T or ~decltype(f())
//                   ::= <simple-id>  
func parse_destructor_name(first, last *CStyleString, db *Db) CStyleString {
	var cs CStyleString
	cs.Content = first.Content
	cs.Pos = first.Pos
	
	if first.equals(last) {
		return cs
	}
	
	t := parse_unresolved_type(first, last, db)
	if t.equals(first) {
		t = parse_simple_id(first, last, db)
	}
	if !t.equals(first) {
		if db.names_empty() {
			return cs
		}
		db.names_back().first = "~" + db.names_back().first
		cs.Pos = t.Pos
	}
	
	return cs
}

// ds <expression> <expression>                         # expr.*expr
func parse_dot_star_expr(first, last *CStyleString, db *Db) *CStyleString {
	cs := &CStyleString{first.Content, first.Pos}

	if last.calcDelta(first) < 3 {
		return cs
	}
	
	if first.prefix(2) != "ds" {
		return cs
	}

	tmpPos := &CStyleString{cs.Content, cs.Pos + 2}
	t := parse_expression(tmpPos, last, db)
	
	if !t.equals(tmpPos) {
		t1 := parse_expression(&t, last, db)
		
		if !t1.equals(&t) {
			if db.names_size() < 2 {
				return cs
			}
			
			expr := db.names_back().move_full()
            db.names_pop_back()
            db.names_back().first += ".*" + (expr)
			cs.Pos = t1.Pos
		}
	}
	
	return cs
}

// dt <expression> <unresolved-name>                    # expr.name
func parse_dot_expr(first, last *CStyleString, db *Db) *CStyleString {
	cs := &CStyleString{first.Content, first.Pos}

	if last.calcDelta(first) < 3 {
		return cs
	}
	
	if first.prefix(2) != "dt" {
		return cs
	}

	tmpPos := &CStyleString{cs.Content, cs.Pos + 2}
	t := parse_expression(tmpPos, last, db)
	
	if !t.equals(tmpPos) {
		t1 := parse_unresolved_name(&t, last, db)
		
		if !t1.equals(&t) {
			if db.names_size() < 2 {
				return cs
			}
			
			nm := db.names_back().move_full()
            db.names_pop_back()
            db.names_back().first += "." + (nm)
			cs.Pos = t1.Pos
		}
	}
	
	return cs
}

// cv <type> <expression>                               # conversion with one argument
// cv <type> _ <expression>* E                          # conversion with a different number of arguments
func parse_conversion_expr(first, last *CStyleString, db *Db) CStyleString {
	var cs CStyleString
	cs.Content = first.Content
	cs.Pos = first.Pos
	
	if last.calcDelta(first) < 3 {
		return cs
	}
	
	if (first.curChar() != 'c') || (first.nextChar() != 'v') {
		return cs
	}
	
	try_to_parse_template_args := db.try_to_parse_template_args
    db.try_to_parse_template_args = false
	
	tmpPos := &CStyleString{cs.Content, cs.Pos + 2}
	t := parse_type(tmpPos, last, db)
	db.try_to_parse_template_args = try_to_parse_template_args
	
	if t.notEach(tmpPos, last) {
		if t.curChar() != '_' {
			t1 := parse_expression(&t, last, db)
			if t1.equals(&t) {
				return cs
			}
			t.Pos = t1.Pos
		} else {
			t.Pos++
			if t.equals(last) {
				return cs
			}
			
			if (t.curChar() == 'E') {
				// db.names.emplace_back();
				
			} else {
				first_expr := true
				for t.curChar() != 'E' {
					t1 := parse_expression(&t, last, db)
					if !t1.notEach(&t, last) {
						return cs
					}
					if !first_expr {
						if (db.names_empty()) {
							return cs
						}
						
						tmp := db.names_back().move_full()
						db.names_pop_back()
						if len(tmp) > 0 {
							if (db.names_empty()) {
								return cs
							}
							db.names_back().first += ", "
							db.names_back().first += tmp
							first_expr = false
						}
					}
					
					t.Pos = t1.Pos
				}
			}
			t.Pos++
		}
		
		if (db.names_size() < 2) {
			return cs
		}
		
		tmp := db.names_back().move_full()
        db.names_pop_back()
        db.names_back().first = "(" + db.names_back().move_full() + ")(" + tmp + ")"
        cs.Pos = t.Pos
	}
	return cs
}

// dc <type> <expression>                               # dynamic_cast<type> (expression)
func parse_dynamic_cast_expr(first, last *CStyleString, db *Db) CStyleString {
	var cs CStyleString
	cs.Content = first.Content
	cs.Pos = first.Pos
	
	if last.calcDelta(first) < 3 {
		return cs
	}
	
	if (first.curChar() != 'd') || (first.nextChar() != 'c') {
		return cs
	}
	
	tmpPos := &CStyleString{cs.Content, cs.Pos + 2}
	t := parse_type(tmpPos, last, db)
	
	if !t.equals(tmpPos) {
		t1 := parse_expression(&t, last, db)
		if !t1.equals(&t) {
			if (db.names_size() < 2) {
				return cs
			}
			
			expr := db.names_back().move_full()
            db.names_pop_back()
            db.names_back().first = "dynamic_cast<" +
			    db.names_back().move_full() + ">(" + expr + ")"
            cs.Pos = t1.Pos
		}
	}
	
	return cs
}

// at <type>      # alignof (a type)
// line 3161
func parse_alignof_type(first, last *CStyleString, db *Db) CStyleString {
	var cs CStyleString
	cs.Content = first.Content
	cs.Pos = first.Pos
	
	if last.calcDelta(first) < 3 {
		return cs
	}
	
	if (first.curChar() != 'a') || (first.nextChar() != 't') {
		return cs
	}
	
	tmpPos := &CStyleString{cs.Content, cs.Pos + 2}
	t := parse_type(tmpPos, last, db)
	if !t.equals(tmpPos) {
		if (db.names_empty()) {
			return cs
		}
		
		db.names_back().first = "alignof (" + db.names_back().move_full() + ")"
        cs.Pos = t.Pos
	}
	
	return cs
}

// az <expression>                                            # alignof (a expression)
// line 3181
func parse_alignof_expr(first, last *CStyleString, db *Db) CStyleString {
	var cs CStyleString
	cs.Content = first.Content
	cs.Pos = first.Pos
	
	if last.calcDelta(first) < 3 {
		return cs
	}
	
	if (first.curChar() != 'a') || (first.nextChar() != 'z') {
		return cs
	}
	
	tmpPos := &CStyleString{cs.Content, cs.Pos + 2}
	t := parse_expression(tmpPos, last, db)
	if !t.equals(tmpPos) {
		if (db.names_empty()) {
			return cs
		}
		
		db.names_back().first = "alignof (" + db.names_back().move_full() + ")"
        cs.Pos = t.Pos
	}
	
	return cs
}

// line 1891
func parse_type(first, last *CStyleString, db *Db) CStyleString {
	var cs CStyleString
	cs.Content = first.Content
	cs.Pos = first.Pos
	
	if first.Pos == last.Pos {
		return cs
	}
	
	if (first.curChar() == 'r') ||
		(first.curChar() == 'V') ||
		(first.curChar() == 'K') {
		cv := 0
		t := parse_cv_qualifiers(first, last, &cv)
		if first.Pos != t.Pos {
			is_function := (t.curChar() == 'F')
			k0 := db.names_size()
			t1 := parse_type(&t, last, db)
			k1 := db.names_size()
			
			if t1.Pos != t.Pos {
				if is_function {
					db.subs_pop_back()
				}
				
				// ??? db.subs.emplace_back(db.names.get_allocator());
				for k := k0; k < k1; k++ {
					if is_function {
						p := len(db.names.content[k].second)
						if db.names.content[k].second[p - 2] == '&' {
							p -= 3
						} else if db.names.content[k].second_back() == '&' {
							p -= 2
						}
						
						if (cv & 1) != 0 {
							s := db.names.content[k].second
							db.names.content[k].second = s[:p] + " const" + s[p:]
							p += 6
						}
						
						if (cv & 2) != 0 {
							s := db.names.content[k].second
							db.names.content[k].second = s[:p] + " volatile" + s[p:]
							p += 9
						}
						
						if (cv & 4) != 0 {
							s := db.names.content[k].second
							db.names.content[k].second = s[:p] + " restrict" + s[p:]
						}						
					} else {
						if (cv & 1) != 0 {
							db.names.content[k].first += " const"
						}
						
						if (cv & 2) != 0 {
							db.names.content[k].first += " volatile"
						}
						
						if (cv & 4) != 0 {
							db.names.content[k].first += " restrict"
						}
					}
					
					arr := db.subs_back().content
					arr = append(arr, db.names.content[k])
				}
				
				cs.Pos = t1.Pos
			}
		}
		
		return cs
	} else {
		t := parse_builtin_type(first, last, db)
		if first.Pos != t.Pos {
			cs.Pos = t.Pos
		} else {
			c := first.Content[0]
			
			if c == 'A' {
				t = parse_array_type(first, last, db)
				
				if t.Pos != cs.Pos {
					if db.names_size() == 0 {
						return cs
					}
					
					cs.Pos = t.Pos
					// TODO
					// ??? db.subs.push_back(typename C::sub_type(1, db.names.back(), db.names.get_allocator()));
				}
			} else if c == 'C' {
				cs.Pos++
				t = parse_type(&cs, last, db)
				
				if t.Pos != cs.Pos {
					if db.names_size() == 0 {
						return *first // or cs???						
					}
					
					db.names_back().first += " complex"
					cs.Pos = t.Pos
					// TODO
				}
			} else if c == 'F' {
				t = parse_function_type(&cs, last, db)
				
				if t.Pos != cs.Pos {
					if db.names_size() == 0 {
						return *first // or cs???						
					}
					
					cs.Pos = t.Pos
					// TODO
				}
			} else if c == 'G' {
				cs.Pos++
				t = parse_type(&cs, last, db)
				
				if t.Pos != cs.Pos {
					if db.names_size() == 0 {
						return *first // or cs???						
					}
					
					db.names_back().first += " imaginary"
					cs.Pos = t.Pos
					// TODO
				}
			} else if c == 'M' {
				t = parse_pointer_to_member_type(&cs, last, db)
				
				if t.Pos != cs.Pos {
					if db.names_size() == 0 {
						return *first // or cs???						
					}
					
					cs.Pos = t.Pos
					// TODO
				}
			} else if c == 'O' {
				cs.Pos++
				k0 := db.names_size()
				t := parse_type(&cs, last, db)
				k1 := db.names_size()
				
				if t.Pos != cs.Pos {
					// ??? db.subs.emplace_back(db.names.get_allocator());
					for k := k0; k < k1; k++ {
						s := db.names.content[k].second
						
						if s[:2] == " [" {
							db.names.content[k].first += " ("
							db.names.content[k].second = ")" + s
						} else if (len(s) > 0) && (s[0] == '(') {
							db.names.content[k].first += "("
							db.names.content[k].second = ")" + s
						}
						
						db.names.content[k].first += "&&"
						db.subs_back().content = append(db.subs_back().content, db.names.content[k])
					}
					
					cs.Pos = t.Pos
				}
			} else if c == 'P' {
				cs.Pos++
				k0 := db.names_size()
				t := parse_type(&cs, last, db)
				k1 := db.names_size()
				
				if t.Pos != cs.Pos {
					// ??? db.subs.emplace_back(db.names.get_allocator());
					for k := k0; k < k1; k++ {
						s := db.names.content[k].second
						
						if s[:2] == " [" {
							db.names.content[k].first += " ("
							db.names.content[k].second = ")" + s
						} else if (len(s) > 0) && (s[0] == '(') {
							db.names.content[k].first += "("
							db.names.content[k].second = ")" + s
						}
						
						if (first.Content[1] != 'U') || (db.names.content[k].first[:12] != "objc_object<") {
							db.names.content[k].first += "*"
						} else {
							db.names.content[k].first = "id" + db.names.content[k].first[11:]
						}
						db.subs_back().content = append(db.subs_back().content, db.names.content[k])
					}
					
					cs.Pos = t.Pos
				}
			} else if c == 'R' {
				cs.Pos++
				k0 := db.names_size()
				t := parse_type(&cs, last, db)
				k1 := db.names_size()
				
				if t.Pos != cs.Pos {
					// ??? db.subs.emplace_back(db.names.get_allocator());
					for k := k0; k < k1; k++ {
						s := db.names.content[k].second
						
						if s[:2] == " [" {
							db.names.content[k].first += " ("
							db.names.content[k].second = ")" + s
						} else if (len(s) > 0) && (s[0] == '(') {
							db.names.content[k].first += "("
							db.names.content[k].second = ")" + s
						}
						
						db.names.content[k].first += "&"
						db.subs_back().content = append(db.subs_back().content, db.names.content[k])
					}
					
					cs.Pos = t.Pos
				}
			} else if c == 'T' {
				k0 := db.names_size()
				t := parse_template_param(&cs, last, db)
				k1 := db.names_size()
				
				if t.Pos != first.Pos {
					// ??? db.subs.emplace_back(db.names.get_allocator());
					for k := k0; k < k1; k++ {
						db.subs_back().content = append(db.subs_back().content, db.names.content[k])
					}
					
					if db.try_to_parse_template_args && (k1 == (k0 + 1)) {
						t1 := parse_template_args(&t, last, db)
						
						if t1.Pos != t.Pos {
							args := db.names_back().move_full()
							db.names_pop_back()
							db.names_back().first += args
							t.Pos = t1.Pos
						}
					}
					
					cs.Pos = t.Pos
				}
			} else if c == 'U' {
				cs.Pos++
				if cs.Pos != last.Pos {
					t := parse_source_name(&cs, last, db)
					if t.Pos != cs.Pos {						
						t2 := parse_type(&t, last, db)
						if t2.Pos != t.Pos {
							if db.names_size() < 2 {
								return *first
							}
							
							ty := db.names_back().move_full()
							db.names_pop_back()
							if db.names_back().first[:9] != "objcproto" {
								db.names_back().first = ty + " " + 
									db.names_back().move_full()
							} else {
								// TODO
								proto := db.names_back().move_full()
								db.names_pop_back()
								
								var start CStyleString
								start.Content = proto
								start.Pos = 9
								
								var protoLast CStyleString
								protoLast.Content = proto
								protoLast.Pos = len(proto)
								t = parse_source_name(&start, &protoLast, db)
								if t.Pos != start.Pos {
									db.names_back().first = ty + "<" + db.names_back().move_full() + ">"
								} else {
									//db.names.content = append(db.names.content, ty + " " + proto)
									db.names_push_back(ty + " " + proto)
								}
							}
							
							cs.Pos = t2.Pos
						}
					}
				}
				return cs
			} else if c == 'S' {
				if ((first.Pos + 1) != last.Pos) && (first.Content[first.Pos + 1] == 't') {
					t := parse_name(first, last, db, nil)
					if t.Pos != first.Pos {
						if db.names_size() == 0 {
							return cs
						}
						
						cs.Pos = t.Pos
					} 
				} else {
					t = parse_substitution(first, last, db)
					if t.Pos != first.Pos {
						cs.Pos = t.Pos
						// Parsed a substitution.  If the substitution is a
                           //  <template-param> it might be followed by <template-args>.
						t = parse_template_args(first, last, db)
						if t.Pos != first.Pos {
							if db.names_size() < 2 {
								return cs
							}
							
							template_args := db.names_back().move_full()
							db.names_pop_back()
							db.names_back().first += template_args
							// Need to create substitution for <template-template-param> <template-args>
							// TODO
							cs.Pos = t.Pos
						}
					}
				}
				
				return cs
			} else if c == 'D' {
				if (first.Pos + 1) == last.Pos {
					return cs
				}
				
				c := first.Content[first.Pos + 1]
				if c == 'p' {
					k0 := db.names_size()
					
					var startStr CStyleString
					startStr.Content = first.Content
					startStr.Pos = first.Pos + 2
					t := parse_type(&startStr, last, db)
					k1 := db.names_size()
					if t.Pos != startStr.Pos {
						// db.subs.emplace_back(db.names.get_allocator());
						for k := k0; k < k1; k++ {
							db.subs_back().push_back(db.names.content[k])
						}
						cs.Pos = t.Pos						
					}
					
					return cs
				} else if (c == 't') || (c == 'T') {
					t = parse_decltype(first, last, db)
					if t.Pos != first.Pos {
						if db.names_empty() {
							return cs
						}
						
						cs.Pos = t.Pos						
					}
					return cs
				} else if (c == 'v') {
					t = parse_vector_type(first, last, db)
					if t.Pos != first.Pos {
						if db.names_empty() {
							return cs
						}
						
						cs.Pos = t.Pos	
					}
					return cs
				}
			} else {
				// default
				// must check for builtin-types before class-enum-types to avoid
                // ambiguities with operator-names
				t = parse_builtin_type(first, last, db)
				if t.Pos != first.Pos {
					cs.Pos = t.Pos
				} else {
					t = parse_name(first, last, db, nil)
					if t.Pos != first.Pos {
						if db.names_empty() {
							return cs
						}
						
						cs.Pos = t.Pos
					}
				}
				
			}
		}
	}
	
	return cs
}

func parse_cv_qualifiers(first, last *CStyleString, cv *int) CStyleString {
	*cv = 0
	var cs CStyleString
	cs.Content = first.Content
	cs.Pos = first.Pos
	
	if first.Pos == last.Pos {
		return cs
	}
	
	if first.curChar() == 'r' {
		*cv += 4
		cs.Pos++
	} 
	if first.curChar() == 'V' {
		*cv += 2
		cs.Pos++
	} 
	if first.curChar() == 'K' {
		*cv += 1
		cs.Pos++
	}  else {
		
	}
	return cs
}
