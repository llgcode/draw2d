// Copyright 2015 The draw2d Authors. All rights reserved.
// created: 16/12/2017 by Drahoslav Bednář

package draw2dsvg

import (
	"encoding/xml"
)

/* svg elements */

type Svg struct {
	XMLName xml.Name `xml:"svg"`
	Xmlns   string   `xml:"xmlns,attr"`
	Groups  []Group  `xml:"g"`
	FillStroke
}

type Group struct {
	FillStroke
	Transform string  `xml:"transform,attr,omitempty"`
	Groups    []Group `xml:"g"`
	Paths     []Path  `xml:"path"`
	Texts     []Text  `xml:"text"`
}

type Path struct {
	FillStroke
	Desc string `xml:"d,attr"`
}

type Text struct {
	FillStroke
	Text  string `xml:",innerxml"`
	Style string `xml:"style,attr,omitempty"`
}

/* shared attrs */

type FillStroke struct {
	Fill     string `xml:"fill,attr,omitempty"`
	FillRule string `xml:"fill-rule,attr,omitempty"`

	Stroke           string `xml:"stroke,attr,omitempty"`
	StrokeWidth      string `xml:"stroke-width,attr,omitempty"`
	StrokeLinecap    string `xml:"stroke-linecap,attr,omitempty"`
	StrokeLinejoin   string `xml:"stroke-linejoin,attr,omitempty"`
	StrokeDasharray  string `xml:"stroke-dasharray,attr,omitempty"`
	StrokeDashoffset string `xml:"stroke-dashoffset,attr,omitempty"`
}
