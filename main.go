// package main

// import (
// 	app "doctor/internal"
// 	"fmt"
// )

// func main() {
// 	app.RunServer()
// 	fmt.Println("RunServer")
// }

package main

import (
	"doctor/internal/agent"
	"doctor/internal/orchestrator"
)

func main() {
	orchestrator.Start()
	agent.StartWorker()
}
