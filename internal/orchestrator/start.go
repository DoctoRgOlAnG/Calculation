package orchestrator

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	agent_app "doctor/internal/agent"
	orch "doctor/internal/orchestrator/handlers"
)

func Start() {
	r := mux.NewRouter()

	r.HandleFunc("/api/v1/calculate", orch.CalculateHandler).Methods("POST")
	r.HandleFunc("/api/v1/expressions", orch.ExpressionsHandler).Methods("GET")
	r.HandleFunc("/api/v1/expressions/{id}", orch.ExpressionHandler).Methods("GET")
	r.HandleFunc("/internal/task", orch.GetTaskHandler).Methods("GET", "POST")

	log.Println("Server started on :8080")

	go agent_app.StartWorker()

	log.Fatal(http.ListenAndServe(":8080", r))
}
