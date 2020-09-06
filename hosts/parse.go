package hosts

import (
	"bufio"
	"io"
	"net"
)

type TokenType int

const (
	Unknown TokenType = iota
	Text
	Separator
	Comment
)

type Token struct {
	Type TokenType
	Text string
}

type Line struct {
	Tokens []Token
}

func (l *Line) Append(token Token) {
	l.Tokens = append(l.Tokens, token)
}

type ParseResult struct {
	Lines []Line
}

func (p *ParseResult) CheckSyntax() bool {

	for _, l := range p.Lines {
		textCount := 0
		for _, t := range l.Tokens {
			switch t.Type {
			case Text:
				if textCount == 0 {
					ip := net.ParseIP(t.Text)
					if ip == nil {
						return false
					}
				}
				textCount++
			}
		}
		if textCount == 1 {
			return false
		}
	}

	return true
}

func Parse(reader io.Reader) (*ParseResult, error) {
	scanner := bufio.NewScanner(reader)
	result := ParseResult{Lines: []Line{}}
	for scanner.Scan() {
		text := scanner.Text()
		line := Line{Tokens: []Token{}}
		for i := 0; i < len(text); {
			if text[i] == ' ' || text[i] == '\t' {
				j := i + 1
				for ; j < len(text); j++ {
					if text[j] != ' ' && text[j] != '\t' {
						break
					}
				}
				line.Append(Token{Type: Separator, Text: text[i:j]})
				i = j
			} else if text[i] == '#' {
				line.Append(Token{Type: Comment, Text: text[i:]})
				break
			} else {
				j := i + 1
				for ; j < len(text); j++ {
					if text[j] == ' ' || text[j] == '\t' || text[j] == '#' {
						break
					}
				}
				line.Append(Token{Type: Text, Text: text[i:j]})
				i = j
			}
		}
		result.Lines = append(result.Lines, line)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return &result, nil
}
