package app

import (
	calculation "doctor/pkg"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func CalcHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"Wrong Method"}`, http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, `{"error":"Invalid Body"}`, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var request struct {
		Expression string `json:"expression"`
	}
	err = json.Unmarshal(body, &request)
	if err != nil || request.Expression == "" {
		http.Error(w, `{"error":"Invalid Body"}`, http.StatusBadRequest)
		return
	}

	result, err := calculation.Calc(request.Expression)
	if err != nil {
		var errorMsg string
		statusCode := http.StatusUnprocessableEntity

		switch err {
		case calculation.ErrInvalidExpression:
			errorMsg = "Error calculation"
		case calculation.ErrDivisionByZero:
			errorMsg = "Division by zero"
		case calculation.ErrMismatchedParentheses:
			errorMsg = "Mismatched parentheses"
		case calculation.ErrInvalidNumber:
			errorMsg = "Invalid number"
		case calculation.ErrUnexpectedToken:
			errorMsg = "Unexpected token"
		case calculation.ErrNotEnoughValues:
			errorMsg = "Not enough values"
		case calculation.ErrInvalidOperator:
			errorMsg = "Invalid operator"
		case calculation.ErrOperatorAtEnd:
			errorMsg = "Operator at end"
		case calculation.ErrMultipleDecimalPoints:
			errorMsg = "Multiple decimal points"
		case calculation.ErrEmptyInput:
			errorMsg = "Empty input"
		default:
			errorMsg = "Error calculation"
			statusCode = http.StatusUnprocessableEntity
		}

		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, errorMsg), statusCode)
		return
	}

	response := struct {
		Result string `json:"result"`
	}{
		Result: fmt.Sprintf("%v", result),
	}

	responseJson, err := json.Marshal(response)
	if err != nil {
		log.Printf("Error while marshaling response: %v", err)
		http.Error(w, `{"error":"Unknown error occurred"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(responseJson)
	if err != nil {
		log.Printf("Error writing response: %v", err)
	}
}

func RunServer() error {
	http.HandleFunc("/api/v1/calculate", CalcHandler)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, `{"error":"Not Found"}`, http.StatusNotFound)
	})

	return http.ListenAndServe(":8080", nil)
}
