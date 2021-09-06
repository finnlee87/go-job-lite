package main

import (
	"github.com/finnlee87/go-job-lite/engine"
)

func main() {
	jobEngine := engine.Default()
	jobEngine.Register(ExampleJob{})
	jobEngine.Register(ExampleJob1{})
	jobEngine.Run()

	ch := make(chan string)
	<-ch
}
