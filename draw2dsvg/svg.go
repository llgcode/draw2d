// Copyright 2015 The draw2d Authors. All rights reserved.
// created: 16/12/2017 by Drahoslav Bednář

package draw2dsvg

import (
	"encoding/xml"
)

type Svg struct {
	XMLName xml.Name `xml:"svg"`
	Xmlns string `xml:"xmlns,attr"`
	Groups []Group `xml:"g"`
}

type Group struct {
	Groups []Group `xml:"g"`
	Paths []Path `xml:"path"`
	Texts []Text `xml:"text"`
}

type Path struct {
	Data string `xml:"d,attr"`
}

type Text struct {
	Text string `xml:",innerxml"`
	Style string `xml:",attr,omitempty"`
}