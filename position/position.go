package position

type Position struct {
	FileName string
	Index    int
	Line     int
	Col      int
}

func (p *Position) Advance(currentChar byte) *Position {
	p.Index++
	p.Col++

	if currentChar == '\n' {
		p.Line++
		p.Col = 0
	}

	return p
}

func (p *Position) Copy() *Position {
	return &Position{
		FileName: p.FileName,
		Index:    p.Index,
		Line:     p.Line,
		Col:      p.Col,
	}
}
