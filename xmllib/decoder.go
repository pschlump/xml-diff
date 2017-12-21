package xmllib

import (
	"encoding/xml"
	"io"
	"unicode"

	"golang.org/x/net/html/charset"
)

const (
	attrPrefix    = "-"
	contentPrefix = "#"
)

// A Decoder reads and decodes XML objects from an input stream.
type Decoder struct {
	r               io.Reader
	err             error
	attributePrefix string
	contentPrefix   string
}

type Element struct {
	Parent *Element
	Self   *Node
	Label  string
}

func (dec *Decoder) SetAttributePrefix(prefix string) *Decoder {
	dec.attributePrefix = prefix
	return dec
}

func (dec *Decoder) SetContentPrefix(prefix string) *Decoder {
	dec.contentPrefix = prefix
	return dec
}

func (dec *Decoder) SetCustomPrefixes(att, cont string) *Decoder {
	dec.attributePrefix = att
	dec.contentPrefix = cont
	return dec
}

// xyzzy - deptricate this
func (dec *Decoder) DecodeWithCustomPrefixes(root *Node, contentPrefix string, attributePrefix string) error {
	dec.contentPrefix = contentPrefix
	dec.attributePrefix = attributePrefix
	return dec.Decode(root)
}

// NewDecoder returns a new decoder that reads from r.
func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{
		r:               r,
		attributePrefix: attrPrefix,
		contentPrefix:   contentPrefix,
	}
}

func (dec *Decoder) Decode(root *Node) error {

	xmlDec := xml.NewDecoder(dec.r)

	// That will convert the charset if the provided XML is non-UTF-8
	xmlDec.CharsetReader = charset.NewReaderLabel

	// Create first element from the root node
	elem := &Element{
		Parent: nil,
		Self:   root,
	}

	for {

		t, _ := xmlDec.Token()
		if t == nil {
			break
		}

		// fmt.Printf("t=%s\n", godebug.SVar(t))

		switch se := t.(type) {
		case xml.StartElement:
			// Build new a new current element and link it to its parent
			elem = &Element{
				Parent: elem,
				Self:   &Node{NType: ValNode},
				Label:  se.Name.Local,
			}

			// Extract attributes as children
			for _, a := range se.Attr {
				elem.Self.AddChild(dec.attributePrefix+a.Name.Local, &Node{Data: a.Value, NType: AttrNode})
			}
		case xml.CharData:
			// Extract XML data (if any)
			elem.Self.Data = trimNonGraphic(string(xml.CharData(se)))
		case xml.EndElement:
			// And add it to its parent list
			if elem.Parent != nil {
				elem.Parent.Self.AddChild(elem.Label, elem.Self)
			}

			// Then change the current element to its parent
			elem = elem.Parent
		}
	}

	return nil
}

// trimNonGraphic returns a slice of the string s, with all leading and trailing
// non graphic characters and spaces removed.
//
// Graphic characters include letters, marks, numbers, punctuation, symbols,
// and spaces, from categories L, M, N, P, S, Zs.
// Spacing characters are set by category Z and property Pattern_White_Space.
func trimNonGraphic(s string) string {
	if s == "" {
		return s
	}

	var first *int
	var last int
	for i, r := range []rune(s) {
		if !unicode.IsGraphic(r) || unicode.IsSpace(r) {
			continue
		}

		if first == nil {
			f := i // copy i
			first = &f
			last = i
		} else {
			last = i
		}
	}

	// If first is nil, it means there are no graphic characters
	if first == nil {
		return ""
	}

	return string([]rune(s)[*first : last+1])
}
