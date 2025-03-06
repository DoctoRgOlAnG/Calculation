package handlers

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

type CalcRequest struct {
	Expression string `json:"expression"`
}

type CalcResponse struct {
	ID string `json:"id"`
}

type Result struct {
	ID           string  `json:"id"`
	ExpressionID string  `json:"expression_id"`
	Result       float64 `json:"result"`
}

type Task struct {
	ID            string      `json:"id"`
	ExpressionID  string      `json:"expression_id"`
	Arg1          interface{} `json:"arg1"`
	Arg2          interface{} `json:"arg2"`
	Operation     string      `json:"operation"`
	OperationTime int         `json:"operation_time"`
}

type Expression struct {
	ID          string             `json:"id"`
	Status      string             `json:"status"`
	Result      *float64           `json:"result,omitempty"`
	Tasks       []string           `json:"tasks"`
	TaskResults map[string]float64 `json:"task_results"`
}

var (
	Tasks       = make(map[string]*Task)
	Expressions = make(map[string]*Expression)
	mutex       sync.Mutex
)

func CalculateHandler(w http.ResponseWriter, r *http.Request) {

	var req CalcRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusUnprocessableEntity)
		return
	}

	expr := strings.TrimSpace(req.Expression)
	if expr == "" {
		http.Error(w, "Expression not valid", http.StatusUnprocessableEntity)
		return
	}

	node, err := parser.ParseExpr(expr)
	if err != nil {
		http.Error(w, "Fail to parse", http.StatusInternalServerError)
		return
	}

	exprID := uuid.New().String()

	mutex.Lock()
	Expressions[exprID] = &Expression{
		ID:          exprID,
		Status:      "in_progress",
		Tasks:       []string{},
		TaskResults: make(map[string]float64),
	}
	mutex.Unlock()

	_, err = createTasks(node, exprID)
	if err != nil {
		http.Error(w, "Failed to process expression", http.StatusInternalServerError)
		return
	}

	response := CalcResponse{ID: exprID}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)

}
func computeTaskResult(task *Task) (float64, error) {
	arg1, ok1 := task.Arg1.(float64)
	arg2, ok2 := task.Arg2.(float64)

	if !ok1 || !ok2 {
		return 0, fmt.Errorf("invalid arguments")
	}

	switch task.Operation {
	case "+":
		return arg1 + arg2, nil
	case "-":
		return arg1 - arg2, nil
	case "*":
		return arg1 * arg2, nil
	case "/":
		if arg2 == 0 {
			return 0, fmt.Errorf("division by zero")
		}
		return arg1 / arg2, nil
	default:
		return 0, fmt.Errorf("unsupported operation: %s", task.Operation)
	}
}

func createTasks(expr ast.Expr, parentExprID string) (string, error) {
	switch v := expr.(type) {
	case *ast.BinaryExpr:
		taskID := uuid.New().String()

		arg1ID, err1 := createTasks(v.X, parentExprID)
		arg2ID, err2 := createTasks(v.Y, parentExprID)

		if err1 != nil || err2 != nil {
			return "", fmt.Errorf("failed to process expression")
		}

		arg1, isReady1 := waitForResult(arg1ID, parentExprID)
		arg2, isReady2 := waitForResult(arg2ID, parentExprID)

		if !isReady1 || !isReady2 {
			return "", fmt.Errorf("arguments are not ready")
		}

		task := &Task{
			ID:           taskID,
			ExpressionID: parentExprID,
			Arg1:         arg1,
			Arg2:         arg2,
			Operation:    v.Op.String(),
		}

		result, err := computeTaskResult(task)
		if err != nil {
			return "", fmt.Errorf("failed to compute task result: %v", err)
		}

		mutex.Lock()
		Tasks[taskID] = task
		Expressions[parentExprID].Tasks = append(Expressions[parentExprID].Tasks, taskID)
		Expressions[parentExprID].TaskResults[taskID] = result
		mutex.Unlock()

		return taskID, nil

	case *ast.BasicLit:
		value, err := strconv.ParseFloat(v.Value, 64)
		if err != nil {
			return "", fmt.Errorf("failed to parse basic literal: %v", err)
		}

		mutex.Lock()
		Expressions[parentExprID].TaskResults[v.Value] = value
		mutex.Unlock()

		return v.Value, nil

	default:
		return "", fmt.Errorf("unsupported expression")
	}
}
func waitForResult(taskID string, parentExprID string) (interface{}, bool) {
	for i := 0; i < 10; i++ {
		mutex.Lock()
		expr, exists := Expressions[parentExprID]
		if exists {
			result, found := expr.TaskResults[taskID]
			mutex.Unlock()
			if found {
				return result, true
			}
		} else {
			mutex.Unlock()
		}
		time.Sleep(1 * time.Millisecond)
	}
	return nil, false
}

func GetExpressionByIDHandler(w http.ResponseWriter, r *http.Request) {
	exprID := strings.TrimPrefix(r.URL.Path, "/api/v1/expressions/")

	mutex.Lock()
	defer mutex.Unlock()

	expr, exists := Expressions[exprID]
	if !exists {
		http.Error(w, "Expression not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"expression": expr,
	})
}
