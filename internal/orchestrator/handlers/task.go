package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

func GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleGetTask(w)
	case http.MethodPost:
		handlePostResult(w, r)
	default:
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func handleGetTask(w http.ResponseWriter) {
	mutex.Lock()
	defer mutex.Unlock()

	var taskToSend *Task
	var taskID string

	// Ищем первую не-nil задачу
	for id, task := range Tasks {
		if task != nil {
			taskToSend = task
			taskID = id
			break
		}
	}

	if taskToSend == nil {
		// Если задач нет, отправляем сообщение
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "No tasks available"})
		return
	}

	// Логируем задачу
	log.Printf("Отправляем задачу агенту: ID=%s, ExpressionID=%s, Arg1=%v, Arg2=%v, Operation=%s",
		taskToSend.ID, taskToSend.ExpressionID, taskToSend.Arg1, taskToSend.Arg2, taskToSend.Operation)

	// Отправляем задачу клиенту
	response := map[string]*Task{"task": taskToSend}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Удаляем задачу из мапы
	delete(Tasks, taskID)
}

func handlePostResult(w http.ResponseWriter, r *http.Request) {
	var result Result
	if err := json.NewDecoder(r.Body).Decode(&result); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	log.Printf("Получен результат: %+v", result)

	mutex.Lock()

	expr, exists := Expressions[result.ExpressionID]
	if !exists {
		mutex.Unlock()
		log.Printf("Ошибка: Expression not found (ID: %s)", result.ExpressionID)
		http.Error(w, "Expression not found", http.StatusNotFound)
		return
	}

	expr.TaskResults[result.ID] = result.Result
	log.Printf("Результат сохранён: ExpressionID=%s, TaskID=%s, Result=%f", result.ExpressionID, result.ID, result.Result)

	if len(expr.TaskResults) == len(expr.Tasks) {
		expr.Result = &result.Result
		expr.Status = "completed"
		log.Printf("Все задачи завершены. Статус выражения изменен на 'completed' (ExpressionID: %s)", result.ExpressionID)
	}

	mutex.Unlock()

	w.WriteHeader(http.StatusOK)
}
