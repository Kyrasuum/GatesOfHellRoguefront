package set

import "fmt"

type Parser struct {
	lexer   *Lexer
	current Token
}

func Parse(data []byte) ([]*Set, error) {
	p := &Parser{
		lexer: NewLexer(data),
	}

	if err := p.next(); err != nil {
		return nil, err
	}

	var roots []*Set

	for p.current.Type != TokenEOF {
		root, err := p.parseSet()
		if err != nil {
			return nil, err
		}

		roots = append(roots, root)
	}

	return roots, nil
}

func (p *Parser) next() error {
	token, err := p.lexer.Next()
	if err != nil {
		return err
	}
	p.current = token
	return nil
}

func (p *Parser) expect(expected TokenType) error {
	if p.current.Type != expected {
		return fmt.Errorf(
			"expected %s, got %s at %d:%d",
			expected,
			p.current.Type,
			p.current.Line,
			p.current.Column,
		)
	}
	return nil
}

func (p *Parser) parseSet() (*Set, error) {
	if err := p.expect(TokenLBrace); err != nil {
		return nil, err
	}
	if err := p.next(); err != nil {
		return nil, err
	}
	if p.current.Type != TokenWord && p.current.Type != TokenString {
		return nil, fmt.Errorf(
			"expected node name, got %s at %d:%d",
			p.current.Type,
			p.current.Line,
			p.current.Column,
		)
	}

	node := &Set{
		Name: p.current.Value,
	}

	if err := p.next(); err != nil {
		return nil, err
	}

	for {
		switch p.current.Type {
		case TokenEOF:
			return nil, fmt.Errorf(
				"unexpected EOF while parsing '%s'",
				node.Name,
			)
		case TokenWord, TokenString:
			node.Args = append(node.Args, p.current.Value)
			if err := p.next(); err != nil {
				return nil, err
			}
		case TokenLBrace:
			child, err := p.parseSet()
			if err != nil {
				return nil, err
			}
			node.Body = append(node.Body, child)
		case TokenRBrace:
			if err := p.next(); err != nil {
				return nil, err
			}
			return node, nil
		}
	}
}
