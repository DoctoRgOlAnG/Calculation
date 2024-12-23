package calculation

import "errors"

var (
	ErrInvalidExpression     = errors.New("invalid expression")
	ErrDivisionByZero        = errors.New("division by zero")
	ErrEmptyInput            = errors.New("empty expression")
	ErrMismatchedParentheses = errors.New("mismatched parentheses")
	ErrUnmatchedParenthesess = errors.New("unmatched parenthesess")
	ErrInvalidCharacter      = errors.New("invalid character")
)
