package main

import (
	app "doctor/internal"
	"fmt"
)

func main() {
	app.RunServer()
	fmt.Println("RunServer")
}
