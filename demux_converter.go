// Copyright 2010 The draw2d Authors. All rights reserved.
// created: 13/12/2010 by Laurent Le Goff

package draw2d

type DemuxConverter struct {
	converters []LineBuilder
}

func NewDemuxConverter(converters ...LineBuilder) *DemuxConverter {
	return &DemuxConverter{converters}
}

func (dc *DemuxConverter) NextCommand(cmd LineMarker) {
	for _, converter := range dc.converters {
		converter.NextCommand(cmd)
	}
}

func (dc *DemuxConverter) MoveTo(x, y float64) {
	for _, converter := range dc.converters {
		converter.MoveTo(x, y)
	}
}

func (dc *DemuxConverter) LineTo(x, y float64) {
	for _, converter := range dc.converters {
		converter.LineTo(x, y)
	}
}
