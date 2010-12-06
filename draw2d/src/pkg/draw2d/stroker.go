package draw2d

type Cap int

const (
	RoundCap Cap = iota
	ButtCap
	SquareCap
)

type LineStroker struct {
	Next VertexConverter
	HalfLineWidth float
	Cap Cap
	Join Join
	vertices []float
	rewind []float
	x, y, nx, ny float
	command VertexCommand
}

func NewLineStroker(converter VertexConverter) (*LineStroker){
	l := new(LineStroker)
	l.Next = converter
	l.HalfLineWidth = 0.5
	l.vertices = make([]float, 0)
	l.rewind = make([]float, 0)
	l.Cap = ButtCap 
	l.Join = MiterJoin
	l.command = VertexNoCommand
	return l
}


func (l *LineStroker) NextCommand(command VertexCommand) {
	l.command = command
	if(command == VertexStopCommand) {
		l.Next.NextCommand(VertexStartCommand)
		for i,j:=0,1; j < len(l.vertices); i,j=i+2,j+2 {
			l.Next.Vertex(l.vertices[i], l.vertices[j])
			l.Next.NextCommand(VertexNoCommand)
		}
		for i,j:=len(l.rewind) - 2 ,len(l.rewind) - 1; j > 0; i,j=i-2,j-2 {
			l.Next.NextCommand(VertexNoCommand)
			l.Next.Vertex(l.rewind[i], l.rewind[j])
		}
		if len(l.vertices) > 1 {
			l.Next.NextCommand(VertexNoCommand)	
			l.Next.Vertex(l.vertices[0] , l.vertices[1])
		}
		l.Next.NextCommand(VertexStopCommand)
		// reinit vertices	
		l.vertices = make([]float, 0)
		l.rewind = make([]float, 0)
		l.x, l.y, l.nx, l.ny = 0, 0, 0, 0
	}
}

func (l *LineStroker) Vertex(x, y float) {
	switch l.command {
		case VertexNoCommand:
			l.line(l.x, l.y, x, y)
		case VertexStartCommand:
			l.x, l.y = x, y
		case VertexJoinCommand:
			l.joinLine(l.x, l.y, l.nx, l.ny, x, y)
		case VertexCloseCommand:
			l.line(l.x, l.y, x, y)
			l.joinLine(l.x, l.y, l.nx, l.ny, x, y)
			l.closePolygon()
	}
	l.command = VertexNoCommand
}

func (l *LineStroker) closePolygon() {
	if len(l.vertices) > 1 {
		l.vertices = append(l.vertices, l.vertices[0] , l.vertices[1])
		l.rewind = append(l.rewind, l.rewind[0] , l.rewind[1])
	}
}


func (l *LineStroker) line(x1, y1, x2, y2 float) {
	dx := (x2 - x1)
	dy := (y2 - y1)
	d := vectorDistance(dx, dy)
	if d != 0 {
		nx := dy * l.HalfLineWidth / d
		ny := -(dx * l.HalfLineWidth / d)
		l.vertices = append(l.vertices, x1 + nx, y1 + ny, x2 + nx , y2 + ny)
		l.rewind = append(l.rewind, x1 - nx, y1 - ny, x2 - nx, y2 - ny)
		l.x, l.y, l.nx, l.ny = x2 , y2 , nx, ny
	}
}

func (l *LineStroker) joinLine(x1, y1, nx1, ny1, x2, y2 float) {
	dx := (x2 - x1)
	dy := (y2 - y1)
	d := vectorDistance(dx, dy)
	
	if(d != 0) {
		nx := dy * l.HalfLineWidth / d
		ny := -(dx * l.HalfLineWidth / d)
	/*	l.join(x1, y1, x1 + nx, y1 - ny, nx, ny, x1 + ny2, y1 + nx2, nx2, ny2)
		l.join(x1, y1, x1 - ny1, y1 - nx1, nx1, ny1, x1 - ny2, y1 - nx2, nx2, ny2)*/
		
		l.vertices = append(l.vertices, x1 + nx, y1 + ny, x2 + nx , y2 + ny)
		l.rewind = append(l.rewind, x1 - nx, y1 - ny, x2 - nx, y2 - ny)
		l.x, l.y, l.nx, l.ny = x2 , y2 ,nx, ny
	}
}

/*
void math_stroke<VC>::calc_arc(VC& vc,
                                   double x,   double y, 
                                   double dx1, double dy1, 
                                   double dx2, double dy2)
    {
        double a1 = atan2(dy1 * m_width_sign, dx1 * m_width_sign);
        double a2 = atan2(dy2 * m_width_sign, dx2 * m_width_sign);
        double da = a1 - a2;
        int i, n;

        da = acos(m_width_abs / (m_width_abs + 0.125 / m_approx_scale)) * 2;

        add_vertex(vc, x + dx1, y + dy1);
        if(m_width_sign > 0)
        {
            if(a1 > a2) a2 += 2 * pi;
            n = int((a2 - a1) / da);
            da = (a2 - a1) / (n + 1);
            a1 += da;
            for(i = 0; i < n; i++)
            {
                add_vertex(vc, x + cos(a1) * m_width, y + sin(a1) * m_width);
                a1 += da;
            }
        }
        else
        {
            if(a1 < a2) a2 -= 2 * pi;
            n = int((a1 - a2) / da);
            da = (a1 - a2) / (n + 1);
            a1 -= da;
            for(i = 0; i < n; i++)
            {
                add_vertex(vc, x + cos(a1) * m_width, y + sin(a1) * m_width);
                a1 -= da;
            }
        }
        add_vertex(vc, x + dx2, y + dy2);
    }
*/
