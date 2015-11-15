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

func (p *CStyleString) nextChar() byte {
	return p.Content[p.Pos + 1]
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

func (p *Db) subs_pop_back() {
	size := len(p.subs.content)
	if size == 0 {
		return
	}
	p.subs.content = p.subs.content[:size - 1]
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

func parse_function_type(first, last *CStyleString, db *Db) CStyleString {
	var cs CStyleString
	cs.Content = first.Content
	cs.Pos = first.Pos
	
	// TODO
	return cs
}

func parse_decltype(first, last *CStyleString, db *Db) CStyleString {
	var cs CStyleString
	cs.Content = first.Content
	cs.Pos = first.Pos
	
	// TODO
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
		t2 := parse_type(noprefix, last, db)
		if t2.Pos == noprefix.Pos {
			return cs
		}
		
		if db.names_size() < 2 {
			return first
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

func parse_name(first, last *CStyleString, db *Db) CStyleString {
	var cs CStyleString
	cs.Content = first.Content
	cs.Pos = first.Pos
	
	// TODO
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
		t := parse_source_name(cs, last, db)
		if t.Pos != cs.Pos {
			cs.Pos = t.Pos
		}
	} else if c == 'D' {
		if (first.Pos + 1) == last.Pos {
			return first
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
					t := parse_name(first, last, db)
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
					t = parse_name(first, last, db)
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