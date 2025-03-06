package main

import (
	"doctor/internal/agent"
	"doctor/internal/orchestrator"
)

func main() {
	orchestrator.Start()
	agent.StartWorker()
}
