package draw2d

type DemuxConverter struct {
	converters []VertexConverter
}

func NewDemuxConverter(converters ...VertexConverter) *DemuxConverter {
	return &DemuxConverter{converters}
}

func (dc *DemuxConverter) NextCommand(cmd VertexCommand) {
	for _, converter := range dc.converters {
		converter.NextCommand(cmd)
	}
}
func (dc *DemuxConverter) Vertex(x, y float) {
	for _, converter := range dc.converters {
		converter.Vertex(x, y)
	}
}
