package ast

import "monkey/internal/token"

type LetStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (s *LetStatement) statementNode()       {}
func (s *LetStatement) TokenLiteral() string { return s.Token.Literal }

type ReturnStatement struct {
	Token token.Token
	Value Expression
}

func (s *ReturnStatement) statementNode() {}

func (s *ReturnStatement) TokenLiteral() string { return s.Token.Literal }
