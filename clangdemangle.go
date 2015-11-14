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

func (p *CStyleString) firstChar() byte {
	return p.Content[0]
}

func (p *CStyleString) currentChar() byte {
	return p.Content[p.Pos]
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

func (p *Db) subs_pop_back() {
	size := len(p.subs.content)
	if size == 0 {
		return
	}
	p.subs.content = p.subs.content[:size - 1]
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
							
						}
					}
				}
			}
		}
		
		return cs
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