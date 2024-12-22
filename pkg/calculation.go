package calculation

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func Calc(expression string) (float64, error) {
	tokens := tokenize(expression)
	if len(tokens) == 0 {
		return 0, errors.New("empty expression")
	}

	output := []float64{}
	operators := []string{}

	for _, token := range tokens {
		if isNumber(token) {
			num, err := strconv.ParseFloat(token, 64)
			if err != nil {
				return 0, err
			}
			output = append(output, num)
		} else if token == "(" {
			operators = append(operators, token)
		} else if token == ")" {
			for len(operators) > 0 && operators[len(operators)-1] != "(" {
				output = append(output, applyOperator(operators[len(operators)-1], &output))
				operators = operators[:len(operators)-1]
			}
			if len(operators) == 0 {
				return 0, errors.New("mismatched parentheses")
			}
			operators = operators[:len(operators)-1]
		} else if isOperator(token) {
			for len(operators) > 0 && precedence(operators[len(operators)-1]) >= precedence(token) {
				output = append(output, applyOperator(operators[len(operators)-1], &output))
				operators = operators[:len(operators)-1]
			}
			operators = append(operators, token)
		} else {
			return 0, fmt.Errorf("invalid character")
		}
	}
	// if len(output) != 1 {
	// 	return 0, errors.New("invalid expression")
	// }
	for len(operators) > 0 {
		if len(output) < 2 {
			return 0, errors.New("unmatched parentheses")
		}
		output = append(output, applyOperator(operators[len(operators)-1], &output))
		operators = operators[:len(operators)-1]
	}

	if len(output) != 1 {
		return 0, errors.New("invalid expression")
	}
	return output[0], nil
}

func tokenize(expression string) []string {
	var tokens []string
	var sb strings.Builder
	for _, char := range expression {
		if isWhitespace(char) {
			continue
		}
		if isOperator(string(char)) || char == '(' || char == ')' {
			if sb.Len() > 0 {
				tokens = append(tokens, sb.String())
				sb.Reset()
			}
			tokens = append(tokens, string(char))
		} else {
			sb.WriteRune(char)
		}
	}
	if sb.Len() > 0 {
		tokens = append(tokens, sb.String())
	}
	return tokens
}

func isOperator(s string) bool {
	return s == "+" || s == "-" || s == "*" || s == "/"
}

func isWhitespace(char rune) bool {
	return char == ' ' || char == '\t' || char == '\n'
}

func isNumber(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

func precedence(op string) int {
	switch op {
	case "+", "-":
		return 1
	case "*", "/":
		return 2
	}
	return 0
}

func applyOperator(op string, output *[]float64) float64 {
	if len(*output) < 2 {
		panic("not enough values in stack")
	}
	b := (*output)[len(*output)-1]
	*output = (*output)[:len(*output)-1]
	a := (*output)[len(*output)-1]
	*output = (*output)[:len(*output)-1]

	switch op {
	case "+":
		return a + b
	case "-":
		return a - b
	case "*":
		return a * b
	case "/":
		if b == 0 {
			panic("division by zero")
		}
		return a / b
	}
	return 0
}
