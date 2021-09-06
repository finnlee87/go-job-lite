package main

import "com.github/finnlee87/go-job-lite/lib"

type ExampleJob1 struct {

}

func (example ExampleJob1) Name() string {
	return "example-job1"
}

func (example ExampleJob1) Cron() string  {
	return "0/5 * * * * * *"
}

func (example ExampleJob1) Execute() error {
	lib.Logger.Infoln("example job 1 run.")
	return nil
}