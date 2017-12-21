package xmllib

/*
TODO in main
*/

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"unicode/utf8"

	"github.com/pschlump/MiscLib"
	"github.com/pschlump/godebug"
	"github.com/pschlump/xml-diff/cfgLib"
)

// An Encoder writes JSON objects to an output stream.
type Encoder struct {
	w                 io.Writer
	err               error
	contentPrefix     string
	attributePrefix   string
	indent            bool
	indentText        string
	outputFmt         string
	combineAllAttrs   bool
	combineAttrs      map[string]map[string]bool // tagName . Attr or tagName . * - for all on tag
	combineAllContent bool
	combineContent    map[string]string
	noSortTag         bool
	noSortTagName     map[string]bool
	Config            CfgType
}

// NewEncoder returns a new encoder that writes to w.
func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{
		w:               w,
		outputFmt:       "json",
		combineAttrs:    make(map[string]map[string]bool),
		combineContent:  make(map[string]string),
		noSortTagName:   make(map[string]bool),
		contentPrefix:   contentPrefix,
		attributePrefix: attrPrefix,
	}
}

func (enc *Encoder) IndentOption(s string) *Encoder {
	enc.indent = true
	enc.indentText = s
	return enc
}

func (enc *Encoder) OutputFormatOption(s string) *Encoder {
	// Xyzzy - should check 's' is valid format, "json", "xml"
	enc.outputFmt = s
	return enc
}

func (enc *Encoder) CustomPrefixesOption(contentPrefix string, attributePrefix string) *Encoder {
	enc.contentPrefix = contentPrefix
	enc.attributePrefix = attributePrefix
	return enc
}

func (enc *Encoder) SetAttributePrefix(prefix string) {
	enc.attributePrefix = prefix
}

func (enc *Encoder) SetContentPrefix(prefix string) {
	enc.contentPrefix = prefix
}

//func (enc *Encoder) EncodeWithCustomPrefixes(root *Node, contentPrefix string, attributePrefix string) error {
//	enc.contentPrefix = contentPrefix
//	enc.attributePrefix = attributePrefix
//	return enc.Encode(root)
//}

// Encode writes the JSON encoding of v to the stream
func (enc *Encoder) Encode(root *Node) error {
	if enc.err != nil {
		return enc.err
	}
	if root == nil {
		return nil
	}

	// csv, insert - other formats.

	// pick format JSON, XML, Text
	switch enc.outputFmt {
	case "json", "JSON":
		enc.err = enc.formatJson(root, 0)
		enc.write("\n")
	case "xml", "XML":
		// fmt.Printf("root=%s\n", godebug.SVarI(root))
		godebug.Printf(db2, "Encode Original: %s\n", godebug.SVarI(root))
		enc.err = enc.processXml(root, "", 0)
		if enc.err != nil {
			return enc.err
		}
		godebug.Printf(db2, "Encode After processXml:  %s \n", godebug.SVarI(root))
		enc.err = enc.formatXml(root, "", 0)
	default:
		fmt.Fprintf(os.Stderr, "Invalid format %s AT: %s\n", enc.outputFmt, godebug.LF())
		return fmt.Errorf("Invalid format %s AT: %s\n", enc.outputFmt, godebug.LF())
	}

	return enc.err
}

func (enc *Encoder) formatXml(curNode *Node, tag string, lvl int) (err error) {

	if db0 {
		fmt.Printf("curNode -- at %s -- =%s, depth=%d\n", tag, godebug.SVarI(curNode), lvl)
	}

	var getAttrs = func(curNode *Node) (rv map[string]string) {
		rv = make(map[string]string)
		for name, it := range curNode.Children {
			if db0 {
				fmt.Printf("name=%s it=%s AT: %s\n", name, godebug.SVarI(it), godebug.LF())
			}
			for _, it2 := range it {
				if it2.NType == AttrNode {
					rv[name[len(enc.attributePrefix):]] = it2.Data
				}
			}
		}
		return
	}

	var renderAttrs = func(attrs map[string]string) {
		// for key, val := range attrs { // sort the keys at this point!
		keys := KeysFromMap(attrs)
		sort.Strings(keys)
		for _, key := range keys {
			val := attrs[key]
			enc.write(" ", key, `="`, val, `"`) // xyzzy2 - escape quotes in val
		}
	}

	var HasVal = func() (has bool) {
		for _, it := range curNode.Children {
			for _, it2 := range it {
				if it2.NType == ValNode {
					return true
				}
			}
		}
		return
	}

	in := strings.Repeat("\t", lvl)

	// ------------------------------------------------------------------------------------------------------

	if curNode.NType == RootNode {
		enc.write(`<?xml version="1.0" encoding="UTF-8"?>`, "\n")
		for tag, cur := range curNode.Children {
			for _, aNode := range cur {
				enc.formatXml(aNode, tag, lvl) // no add - at top.
			}
		}
		return
	}

	// Open Tag -----------------------------------------
	enc.write(in, "<", tag)
	attrs := getAttrs(curNode)
	if db0 {
		fmt.Printf("X7 - has children, not RootNode \nattrs = %s\n\n", godebug.SVarI(attrs))
	}
	renderAttrs(attrs)

	if !HasVal() && len(curNode.Data) == 0 {
		enc.write("/>\n")
		return
	} else {
		enc.write(">") // if no body, then short cut end!
	}

	if len(curNode.Data) > 0 {
		enc.write(curNode.Data)
		in = ""
	}

	if HasVal() {
		enc.write("\n") // xyzzy - if no body, then short cut end!
		// The Body -----------------------------------------
		// for name, it := range curNode.Children {
		keys := KeysFromMap(curNode.Children)
		if !enc.noSortTag || !enc.noSortTagName[tag] {
			sort.Strings(keys)
		}
		for _, name := range keys {
			it := curNode.Children[name]
			for _, it2 := range it {
				if it2.NType == ValNode {
					if db0 {
						fmt.Printf("X8 node[%s] value = %s\n\n", name, godebug.SVarI(it2))
					}
					e0 := enc.formatXml(it2, name, lvl+1)
					if e0 != nil {
						return e0
					}
				}
			}
		}
	}

	// Close Tag -----------------------------------------
	enc.write(in, "</", tag, ">\n")

	return nil
}

// xyzzy - iterate over tree
func (enc *Encoder) formatJson(curNode *Node, lvl int) error {
	var indentN = func(n int) {
		if enc.indent {
			for ii := 0; ii < n; ii++ {
				enc.write(enc.indentText)
			}
		}
	}
	if curNode.HasChildren() {
		enc.write("{")
		if enc.indent {
			enc.write("\n")
		}

		// xyzzy - must sort names before print?  Attributes must be in order for compare.

		// Add data as an additional attibute (if any)
		if len(curNode.Data) > 0 {
			indentN(lvl + 1)
			enc.write(`"`, enc.contentPrefix, "content", `": `, sanitiseString(curNode.Data), ", ")
			if enc.indent {
				enc.write("\n")
			}
		}

		sl := make([]string, 0, len(curNode.Children))
		for label := range curNode.Children {
			sl = append(sl, label)
		}
		// fmt.Printf("sl->%s<-\n", sl)
		if len(sl) > 1 {
			// fmt.Printf("Must sort")
			sort.Strings(sl)
		}
		// fmt.Printf("sorted: sl->%s<-\n", sl)

		com := ""
		// for label, children := range curNode.Children {
		for ii := range sl {
			label, children := sl[ii], curNode.Children[sl[ii]]
			enc.write(com)
			indentN(lvl + 1)
			enc.write(`"`, label, `": `)

			if len(children) > 1 {
				// Array
				// xyzzy - may need to sort?
				enc.write("[") // xyzzy - need to estimate if length is less than X- then one line - else - multi-line
				com1 := ""
				for _, ch := range children {
					enc.write(com1)
					enc.formatJson(ch, lvl+2)
					com1 = ", "
				}
				enc.write("]")
			} else {
				// Map
				enc.formatJson(children[0], lvl+1)
			}

			if enc.indent {
				com = ",\n"
			} else {
				com = ", "
			}
		}

		enc.write("\n")
		indentN(lvl)
		enc.write("}")
	} else {
		// TODO : Extract data type
		enc.write(sanitiseString(curNode.Data))
	}

	return nil
}

func (enc *Encoder) write(s ...string) {
	for _, ss := range s {
		enc.w.Write([]byte(ss))
	}
}

// https://golang.org/src/encoding/json/encode.go?s=5584:5627#L788
var hex = "0123456789abcdef"

func sanitiseString(s string) string {
	var buf bytes.Buffer

	buf.WriteByte('"')
	start := 0
	for i := 0; i < len(s); {
		if b := s[i]; b < utf8.RuneSelf {
			if 0x20 <= b && b != '\\' && b != '"' && b != '<' && b != '>' && b != '&' {
				i++
				continue
			}
			if start < i {
				buf.WriteString(s[start:i])
			}
			switch b {
			case '\\', '"':
				buf.WriteByte('\\')
				buf.WriteByte(b)
			case '\n':
				buf.WriteByte('\\')
				buf.WriteByte('n')
			case '\r':
				buf.WriteByte('\\')
				buf.WriteByte('r')
			case '\t':
				buf.WriteByte('\\')
				buf.WriteByte('t')
			default:
				// This encodes bytes < 0x20 except for \n and \r,
				// as well as <, > and &. The latter are escaped because they
				// can lead to security holes when user-controlled strings
				// are rendered into JSON and served to some browsers.
				buf.WriteString(`\u00`)
				buf.WriteByte(hex[b>>4])
				buf.WriteByte(hex[b&0xF])
			}
			i++
			start = i
			continue
		}
		c, size := utf8.DecodeRuneInString(s[i:])
		if c == utf8.RuneError && size == 1 {
			if start < i {
				buf.WriteString(s[start:i])
			}
			buf.WriteString(`\ufffd`)
			i += size
			start = i
			continue
		}
		// U+2028 is LINE SEPARATOR.
		// U+2029 is PARAGRAPH SEPARATOR.
		// They are both technically valid characters in JSON strings,
		// but don't work in JSONP, which has to be evaluated as JavaScript,
		// and can lead to security holes there. It is valid JSON to
		// escape them, so we do so unconditionally.
		// See http://timelessrepo.com/json-isnt-a-javascript-subset for discussion.
		if c == '\u2028' || c == '\u2029' {
			if start < i {
				buf.WriteString(s[start:i])
			}
			buf.WriteString(`\u202`)
			buf.WriteByte(hex[c&0xF])
			i += size
			start = i
			continue
		}
		i += size
	}
	if start < len(s) {
		buf.WriteString(s[start:])
	}
	buf.WriteByte('"')
	return buf.String()
}

func (enc *Encoder) processXml(curNode *Node, tag string, lvl int) (err error) {

	godebug.Printf(db1, "procesXml: curNode -- at %s -- =%s, depth=%d\n", tag, godebug.SVarI(curNode), lvl)

	cfg := enc.Config

	var getAttrs = func(curNode *Node) {
		match := false
		for name, it := range curNode.Children {
			// name = name[len(enc.attributePrefix):]
			match = false
			godebug.Printf(db1, "getAttrs: name=%s it=%s AT: %s\n", name, godebug.SVarI(it), godebug.LF())
			for kk, it2 := range it {
				godebug.Printf(db3, "kk=%v it2=%v AT: %s\n", kk, SVar(it2), godebug.LF())
				if it2.NType == AttrNode {
					godebug.Printf(db3, "fond AttrNode AT: %s\n", godebug.LF())
					// xyzzy - wild card
					if n, ok := cfg.attrsLookup[tag]; ok {
						godebug.Printf(db3, "found attribute in attrsLookup for tag=[%s], n=%d AT: %s\n", tag, n, godebug.LF())
						// xyzzy - wild card
						if cfg.AttrsToValue[n].AttrName == name {
							godebug.Printf(db3, "%sfound AttrsToValue[%d] name=[%s] AT: %s%s\n", MiscLib.ColorGreen, n, name, godebug.LF(), MiscLib.ColorReset)
							it2.NType = ValNode
							it[kk] = it2
							match = true
						}
					}
				}
			}
			if match {
				curNode.Children[name] = it
				godebug.Printf(db3, "%smatch=true it=%s AT: %s%s\n", MiscLib.ColorGreen, SVar(it), godebug.LF(), MiscLib.ColorReset)
			}
		}
	}

	var HasVal = func() (has bool) {
		for _, it := range curNode.Children {
			for _, it2 := range it {
				if it2.NType == ValNode {
					return true
				}
			}
		}
		return
	}

	// ------------------------------------------------------------------------------------------------------

	if curNode.NType == RootNode {
		for tag, cur := range curNode.Children {
			for _, aNode := range cur {
				enc.processXml(aNode, tag, lvl) // no add - at top.
			}
		}
		return
	}

	getAttrs(curNode)

	if !HasVal() && len(curNode.Data) == 0 {
		return
	}

	if HasVal() {
		// for name, it := range curNode.Children {
		keys := KeysFromMap(curNode.Children)
		if !enc.noSortTag || !enc.noSortTagName[tag] {
			sort.Strings(keys)
		}
		match := false
		for _, name := range keys {
			match = false
			it := curNode.Children[name]
			godebug.Printf(db5, "%sgetAttrs: tag=%s name=%s it=%s AT: %s%s\n", MiscLib.ColorYellow, tag, name, godebug.SVarI(it), godebug.LF(), MiscLib.ColorReset)
			for kk, it2 := range it {
				godebug.Printf(db5, "kk=%v it2=%v AT: %s\n", kk, SVar(it2), godebug.LF())
				if it2.NType == ValNode {
					godebug.Printf(db5, "fond ValNode tag=[%s] valsLookup=%s AT: %s\n", tag, SVar(cfg.valsLookup), godebug.LF())
					godebug.Printf(db1, "X8 node[%s] value = %s\n\n", name, godebug.SVarI(it2))
					// xyzzy - wild card
					if n, ok := cfg.valsLookup[tag]; ok {
						godebug.Printf(db5, "found attribute in attrsLookup for tag=[%s], n=%d AT: %s\n", tag, n, godebug.LF())
						// xyzzy - wild card
						if cfg.ValueToAttr[n].AttrName == name {
							godebug.Printf(db5, "%sfound AttrsToValue[%d] name=[%s] AT: %s%s\n", MiscLib.ColorGreen, n, name, godebug.LF(), MiscLib.ColorReset)
							it2.NType = AttrNode
							match = true
							it[kk] = it2
						}
					}
					if match {
						curNode.Children[name] = it
						godebug.Printf(db5, "%smatch=true it=%s AT: %s%s\n", MiscLib.ColorGreen, SVar(it), godebug.LF(), MiscLib.ColorReset)
					}
					e0 := enc.processXml(it2, name, lvl+1)
					if e0 != nil {
						return e0
					}
				}
			}
		}
	}

	return nil
}

/*
	attrs:
		tag-name attr-name -> nested-tag-name
		tag-name * -> nested-tag-name
		* attr-name -> nested-tag-name
		* * -> *

	content:
		tag-name -> attr-name
		* -> attr-name
		* -> *

	nosort:
		tag
		tag
*/

type AttrToValueType struct {
	TagName  string `json:"TagName"`
	AttrName string `json:"AttrName"`
	ValName  string `json:"ValName"`
}

type ValueToAttrType struct {
	TagName  string `json:"TagName"`
	AttrName string `json:"AttrName"`
}

type NoSortType struct {
	TagName string `json:"TagName"`
}

type CfgType struct {
	AttrsToValue []AttrToValueType `json:"AttrsToValue"`
	ValueToAttr  []ValueToAttrType `json:"ValueToAttr"`
	NoSort       []NoSortType      `json:"NoSort"`
	attrsLookup  map[string]int
	valsLookup   map[string]int
	noSortLookup map[string]bool
}

func ReadCfg(fn string) (cfg CfgType) {

	if fn == "" {
		return
	}

	if !Exists(fn) {
		fmt.Fprintf(os.Stderr, "Missing configuration file %s\n", fn)
		os.Exit(1)
	}

	cfgLib.ReadConfigFile(fn, &cfg)

	cfg.attrsLookup = make(map[string]int)
	cfg.valsLookup = make(map[string]int)
	cfg.noSortLookup = make(map[string]bool)
	for ii, vv := range cfg.AttrsToValue {
		cfg.attrsLookup[vv.TagName] = ii
	}
	for ii, vv := range cfg.ValueToAttr {
		cfg.valsLookup[vv.TagName] = ii
	}
	for _, vv := range cfg.NoSort {
		cfg.noSortLookup[vv.TagName] = true
	}

	godebug.Printf(db4, "Cfg=%s\n", SVarI(cfg))
	godebug.Printf(db4, "valsLookup=%s\n", SVarI(cfg.valsLookup))

	return
}

const db0 = false
const db1 = false
const db2 = false
const db3 = false
const db4 = false
const db5 = false

/* vim: set noai ts=4 sw=4: */
