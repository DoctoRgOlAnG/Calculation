package calculation

import (
	"testing"
)

func TestCalc(t *testing.T) {
	testCasesSuccess := []struct {
		name           string
		expression     string
		expectedResult float64
	}{
		{
			name:           "simple",
			expression:     "1+1",
			expectedResult: 2,
		},
		{
			name:           "priority",
			expression:     "(2+2)*2",
			expectedResult: 8,
		},
		{
			name:           "priority",
			expression:     "2+2*2",
			expectedResult: 6,
		},
		{
			name:           "division",
			expression:     "1/2",
			expectedResult: 0.5,
		},
	}

	for _, testCase := range testCasesSuccess {
		t.Run(testCase.name, func(t *testing.T) {
			val, err := Calc(testCase.expression)
			if err != nil {
				t.Fatalf("successful case %s returns error: %v", testCase.expression, err)
			}
			if val != testCase.expectedResult {
				t.Fatalf("%f should be equal %f", val, testCase.expectedResult)
			}
		})
	}

	testCasesFail := []struct {
		name        string
		expression  string
		expectedErr string
	}{
		{
			name:        "simple",
			expression:  "1+1*",
			expectedErr: "unmatched parentheses",
		},
		{
			name:        "unmatched parentheses",
			expression:  "(1+2",
			expectedErr: "unmatched parentheses",
		},
		{
			name:        "invalid character",
			expression:  "1+2a",
			expectedErr: "invalid character",
		},
		{
			name:        "empty expression",
			expression:  "",
			expectedErr: "empty expression",
		},
	}

	for _, testCase := range testCasesFail {
		t.Run(testCase.name, func(t *testing.T) {
			val, err := Calc(testCase.expression)
			if err == nil {
				t.Fatalf("expected error for case %s, got result %f", testCase.expression, val)
			}
			if err.Error() != testCase.expectedErr {
				t.Fatalf("expected error '%v', got '%v'", testCase.expectedErr, err)
			}
		})
	}
}
