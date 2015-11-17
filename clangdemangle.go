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

/*
func (p *CStyleString) firstChar() byte {
	return p.Content[0]
}
*/

func (p *CStyleString) currentChar() byte {
	return p.Content[p.Pos]
}

func (p *CStyleString) curChar() byte {
	return p.Content[p.Pos]
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
	
	if first.currentChar() != 'F' {
		return *first
	}
	
	t := first
	t.Pos++
	
	if t.Pos == last.Pos {
		return cs
	}
	
	if t.currentChar() == 'Y' {
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
		
		c := t.currentChar()
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
		if (t.Pos != noprefix.Pos) && (t.Pos != last.Pos) && (t.currentChar() == 'E') {
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
	//cv1 := db.cv
	//ref1 := db.ref
	
	if t.equals(first) {
		db.encoding_depth = su // ???
		db.tag_templates = sb
		return cs
	}
	
	c2 := t.curChar()
	if !t.equals(last) && (c2 != 'E') && (c2 != '.') {
		sb2 := db.tag_templates
		db.tag_templates = false
		
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
			t2 := parse_type(&t, last, db)
			
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
			//first_arg := true
			for {
				k0 := db.names_size()
                t2 := parse_type(&t, last, db)
                k1 := db.names_size()
				
				if t2.equals(&t) {
					break
				}
				
				if k1 > k0 {
					// TODO
				}
			}
		}
		
		db.tag_templates = sb2
	}
	
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

func parse_template_args(first, last *CStyleString, db *Db) CStyleString {
	var cs CStyleString
	cs.Content = first.Content
	cs.Pos = first.Pos
	
	// TODO
	return cs
}

func parse_source_name(first, last *CStyleString, db *Db) CStyleString {
	var cs CStyleString
	cs.Content = first.Content
	cs.Pos = first.Pos
	
	// TODO
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
	
	// TODO
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
	if t.currentChar() == 'n' {
		t.Pos++
	}
	
	if t.Pos != last.Pos {
		if t.currentChar() == '0' {
			cs.Pos = t.Pos + 1
		} else if isNonZeroNumberChar(t.currentChar()) {
			cs.Pos = t.Pos + 1
			
			for (cs.Pos != last.Pos) && isNumberChar(cs.currentChar()) {
				cs.Pos++
			}
		}
	}
	
	return cs
}

func parse_array_type(first, last *CStyleString, db *Db) CStyleString {
	var cs CStyleString
	cs.Content = first.Content
	cs.Pos = first.Pos
	if (first.Pos == last.Pos) || ((first.Pos + 1) == last.Pos) {
		return cs
	}
	
	if first.currentChar() != 'A' {
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
		if (t.Pos != last.Pos) && (t.currentChar() == '_') {
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
		
		if t.currentChar() != '_' {
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
	if t0.currentChar() == 'L' {
		t0.Pos++
	}
	
	if t0.currentChar() == 'N' {
		t1 := parse_nested_name(t0, last, db, ends_with_template_args)
		
		if t1.Pos != t0.Pos {
			cs.Pos = t1.Pos
		}
	} else if t0.currentChar() == 'Z' {
		t1 := parse_local_name(t0, last, db, ends_with_template_args)
		
		if t1.Pos != t0.Pos {
			cs.Pos = t1.Pos
		}
	} else {
		t1 := parse_unscoped_name(t0, last, db)
		
		if t1.Pos != t0.Pos {
			if (t1.Pos != last.Pos) && (t1.currentChar() == 'I') {
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
			if !t1.equals(t0) && !t1.equals(last) && (t1.currentChar() == 'I') {
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

func parse_builtin_type(first, last *CStyleString, db *Db) CStyleString {
	var cs CStyleString
	cs.Content = first.Content
	cs.Pos = first.Pos
	
	if first.Pos == last.Pos {
		return cs
	}
	
	cs.Pos++
	
	c := first.currentChar()
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

func parse_expression(first, last *CStyleString, db *Db) CStyleString {
	var cs CStyleString
	cs.Content = first.Content
	cs.Pos = first.Pos
	
	delta := last.Pos - first.Pos
	if delta < 2 {
		return cs
	}
	
	//t := first
	//parsed_gs := false
	
	
	// TODO
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
	
	if (first.currentChar() == 'r') ||
		(first.currentChar() == 'V') ||
		(first.currentChar() == 'K') {
		cv := 0
		t := parse_cv_qualifiers(first, last, &cv)
		if first.Pos != t.Pos {
			is_function := (t.currentChar() == 'F')
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
	
	if first.currentChar() == 'r' {
		*cv += 4
		cs.Pos++
	} 
	if first.currentChar() == 'V' {
		*cv += 2
		cs.Pos++
	} 
	if first.currentChar() == 'K' {
		*cv += 1
		cs.Pos++
	}  else {
		
	}
	return cs
}
