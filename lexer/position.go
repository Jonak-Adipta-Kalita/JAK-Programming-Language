package lexer

type Position struct {
	Idx  int
	Ln   int
	Col  int
	Fn   string
	Ftxt string
}

func (p *Position) Advance(currentChar byte) *Position {
	p.Idx += 1
	p.Col += 1

	if currentChar == '\n' {
		p.Ln += 1
		p.Col = 0
	}

	return p
}

func (p *Position) Copy() *Position {
	return &Position{p.Idx, p.Ln, p.Col, p.Fn, p.Ftxt}
}
