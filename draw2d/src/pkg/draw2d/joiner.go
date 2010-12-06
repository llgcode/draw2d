package draw2d


type Join int

const (
	BevelJoin Join = iota
	RoundJoin
	MiterJoin
)

type JoinerFunc func(x1, y1, nx1, ny1, x2, y2, nx2, ny2  float)

func emptyJoiner(x1, y1, nx1, ny1, x2, y2, nx2, ny2  float) {

}