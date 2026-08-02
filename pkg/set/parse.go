package set

import "fmt"

type Parser struct {
	lexer   *Lexer
	current Token
}

func Parse(name string, data []byte) (*Set, error) {
	p := &Parser{
		lexer: NewLexer(data),
	}

	if err := p.next(); err != nil {
		return nil, err
	}

	root := &Set{
		Name: name,
	}

	for p.current.Type != TokenEOF {
		node, err := p.parseSet()
		if err != nil {
			return nil, err
		}

		root.Body = append(root.Body, node)
	}

	return root, nil
}

func (p *Parser) next() error {
	token, err := p.lexer.Next()
	if err != nil {
		return err
	}
	p.current = token
	return nil
}

func (p *Parser) parseSet() (*Set, error) {
	var end TokenType
	switch p.current.Type {
	case TokenLBrace:
		end = TokenRBrace
	case TokenLParen:
		end = TokenRParen
	default:
		return nil, fmt.Errorf(
			"expected '{' or '(', got %s at %d:%d",
			p.current.Type,
			p.current.Line,
			p.current.Column,
		)
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
		case TokenWord, TokenString:
			if node.Name == "" {
				node.Name = p.current.Value
			} else {
				node.Args = append(node.Args, p.current.Value)
			}
			if err := p.next(); err != nil {
				return nil, err
			}
		case TokenLBrace, TokenLParen:
			child, err := p.parseSet()
			if err != nil {
				return nil, err
			}
			node.Body = append(node.Body, child)
		case end:
			if err := p.next(); err != nil {
				return nil, err
			}
			return node, nil
		case TokenEOF:
			return nil, fmt.Errorf(
				"unexpected EOF while parsing %q",
				node.Name,
			)
		}
	}
}
