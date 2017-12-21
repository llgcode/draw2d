// Copyright 2015 The draw2d Authors. All rights reserved.
// created: 16/12/2017 by Drahoslav Bednář

// Package draw2dsvg_test gives test coverage with the command:
// go test -cover ./... | grep -v "no test"
// (It should be run from its parent draw2d directory.)
package draw2dsvg

import (
	"testing"
	"encoding/xml"
)

// Test basic encoding of svg/xml elements
func TestXml(t *testing.T) {
	
	svg := NewSvg()
	svg.Groups = []Group{Group{
		Groups: []Group{
			Group{},  // nested groups
			Group{},
		},
		Texts: []Text{
			Text{Text: "Hello"}, // text
			Text{Text: "world", Style: "opacity: 0.5"}, // text with style
		},
		Paths: []Path{
			Path{Data: "M100,200 C100,100 250,100 250,200 S400,300 400,200"}, // simple path
			Path{}, // empty path
		},
	}}

	expectedOut := `<svg xmlns="http://www.w3.org/2000/svg">
  <g>
    <g></g>
    <g></g>
    <path d="M100,200 C100,100 250,100 250,200 S400,300 400,200"></path>
    <path d=""></path>
    <text>Hello</text>
    <text Style="opacity: 0.5">world</text>
  </g>
</svg>`

	out, err := xml.MarshalIndent(svg, "", "  ")

	if err != nil {
		t.Error(err)
	}
	if string(out) != expectedOut {
		t.Errorf("svg output is not as expected\n"+
			"got:\n%s\n\n"+
			"want:\n%s\n",
			string(out),
			expectedOut,
		)
	}
}
