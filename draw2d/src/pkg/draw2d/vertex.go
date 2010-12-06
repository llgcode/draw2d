package draw2d

type VertexCommand byte

const (
	VertexNoCommand VertexCommand = iota
	VertexStartCommand
	VertexJoinCommand
	VertexCloseCommand
	VertexStopCommand
)

type VertexConverter interface {
	NextCommand(cmd VertexCommand)
	Vertex(x, y float)
}




