package draw2d


type LineDasher struct {

}

func (d *LineDasher) Start() {
	
}

func (d *LineDasher) Stop() {
	
}

func (d *LineDasher) Vertex(x, y float) {

}




/*
type PathAdapter struct {
	path           *raster.Path
	x, y, distance float
	dash           []float
	currentDash    int
	dashOffset     float
}

func tracePath(dash []float, dashOffset float, paths ...*Path) *raster.Path {
	var adapter PathAdapter
	if dash != nil && len(dash) > 0 {
		adapter.dash = dash
	} else {
		adapter.dash = nil
	}
	adapter.currentDash = 0
	adapter.dashOffset = dashOffset
	adapter.path = new(raster.Path)
	for _, path := range paths {
		path.TraceLine(&adapter)
	}
	return adapter.path
}

func floatToPoint(x, y float) raster.Point {
	return raster.Point{raster.Fix32(x * 256), raster.Fix32(y * 256)}
}

func (p *PathAdapter) MoveTo(x, y float) {
	p.path.Start(floatToPoint(x, y))
	p.x, p.y = x, y
	p.distance = p.dashOffset
	p.currentDash = 0
}

func (p *PathAdapter) LineTo(x, y float) {
	if p.dash != nil {
		rest := p.dash[p.currentDash] - p.distance
		for rest < 0 {
			p.distance = p.distance - p.dash[p.currentDash]
			p.currentDash = (p.currentDash + 1) % len(p.dash)
			rest = p.dash[p.currentDash] - p.distance
		}
		d := distance(p.x, p.y, x, y)
		for d >= rest {
			k := rest / d
			lx := p.x + k*(x-p.x)
			ly := p.y + k*(y-p.y)
			if p.currentDash%2 == 0 {
				// line
				p.path.Add1(floatToPoint(lx, ly))
			} else {
				// gap
				p.path.Start(floatToPoint(lx, ly))
			}
			d = d - rest
			p.x, p.y = lx, ly
			p.currentDash = (p.currentDash + 1) % len(p.dash)
			rest = p.dash[p.currentDash]
		}
		p.distance = d
		if p.currentDash%2 == 0 {
			p.path.Add1(floatToPoint(x, y))
		} else {
			p.path.Start(floatToPoint(x, y))
		}
		if p.distance >= p.dash[p.currentDash] {
			p.distance = p.distance - p.dash[p.currentDash]
			p.currentDash = (p.currentDash + 1) % len(p.dash)
		}
	} else {
		p.path.Add1(floatToPoint(x, y))
	}
	p.x, p.y = x, y
}
*/