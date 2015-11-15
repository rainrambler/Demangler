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
			
		}
	}
}

// clang\cxxabi\src\cxa_demangle.cpp line 4491
func parse_encoding(mangled string, first, last int, db *Db) string {
	if first == last {
		return mangled
	}
	
	//prevDepth := db.encoding_depth
	db.encoding_depth++
	
	//prevTagTempl := db.tag_templates
	if db.encoding_depth > 1 {
		db.tag_templates = true
	}
	
	c := mangled[first]
	if (c == 'G') || (c == 'T') {
		
	} else {
		
	}
	
	return mangled
}

func parse_special_name(mangled string, first, last int, db *Db, result *CStyleString) {
	if (last - first) <= 2 {
		// invalid, not parsed
		return
	}
	
	c0 := mangled[first]
	c1 := mangled[first + 1]
	if c0 == 'T' {
		if c1 == 'V' {
			// TV <type>    # virtual table
			
		} else if c1 == 'T' {
			
		}
	} else {
		
	}
}

func parse_template_param(first, last *CStyleString, db *Db) CStyleString {
	var cs CStyleString
	cs.Content = first.Content
	cs.Pos = first.Pos
	
	// TODO
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
	
	if first.currentChar() != 'D' {
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

func parse_vector_type(first, last *CStyleString, db *Db) CStyleString {
	var cs CStyleString
	cs.Content = first.Content
	cs.Pos = first.Pos
	
	// TODO
	return cs
}

func parse_pointer_to_member_type(first, last *CStyleString, db *Db) CStyleString {
	var cs CStyleString
	cs.Content = first.Content
	cs.Pos = first.Pos
	
	// TODO
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

func parse_unqualified_name(first, last *CStyleString, db *Db) CStyleString {
	var cs CStyleString
	cs.Content = first.Content
	cs.Pos = first.Pos
	
	// TODO
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

func parse_local_name(first, last *CStyleString, db *Db, ends_with_template_args *bool) CStyleString {
	var cs CStyleString
	cs.Content = first.Content
	cs.Pos = first.Pos
	
	// TODO
	return cs
}

func parse_unscoped_name(first, last *CStyleString, db *Db) CStyleString {
	var cs CStyleString
	cs.Content = first.Content
	cs.Pos = first.Pos
	
	// TODO
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
