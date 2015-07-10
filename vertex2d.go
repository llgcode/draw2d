// Copyright 2010 The draw2d Authors. All rights reserved.
// created: 21/11/2010 by Laurent Le Goff

package draw2d

// VertexCommand defines different commands to describe the vertex of a path.
type VertexCommand byte

const (
	// VertexNoCommand does nothing
	VertexNoCommand VertexCommand = iota
	// VertexStartCommand starts a (sub)path
	VertexStartCommand
	// VertexJoinCommand joins the two edges at the vertex
	VertexJoinCommand
	// VertexCloseCommand closes the subpath
	VertexCloseCommand
	// VertexStopCommand is the endpoint of the path.
	VertexStopCommand
)

// VertexConverter allows to convert vertices.
type VertexConverter interface {
	NextCommand(cmd VertexCommand)
	Vertex(x, y float64)
}
