package xmllib

import (
	"bytes"
	"io"
)

// Convert converts the given XML document to JSON
func Convert(r io.Reader) (*bytes.Buffer, error) {
	// Decode XML document
	root := &Node{}
	err := NewDecoder(r).Decode(root)
	if err != nil {
		return nil, err
	}

	// Then encode it in JSON
	buf := new(bytes.Buffer)
	// PJS - func (enc *Encoder) IndentOption(s string) *Encoder {
	enc := NewEncoder(buf)
	//	enc.IndentOption("\t")
	//	enc.OutputFormatOption("xml")

	err = enc.Encode(root)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

// Convert converts the given XML document to JSON
func ConvertXML(r io.Reader) (*bytes.Buffer, error) {
	// Decode XML document
	root := &Node{}
	err := NewDecoder(r).Decode(root)
	if err != nil {
		return nil, err
	}

	// Then encode it in JSON
	buf := new(bytes.Buffer)
	// PJS - func (enc *Encoder) IndentOption(s string) *Encoder {
	enc := NewEncoder(buf)
	enc.IndentOption("\t")
	enc.OutputFormatOption("xml")

	err = enc.Encode(root)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

// PJS - ReadXMLAsNode reads a source of XML and returns the Node tree
func ReadXMLAsNode(r io.Reader) (root *Node, err error) {
	// Decode XML document
	root = &Node{}
	err = NewDecoder(r).Decode(root)
	if err != nil {
		return nil, err
	}
	return
}
