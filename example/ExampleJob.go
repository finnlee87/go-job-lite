package main

import "github.com/finnlee87/go-job-lite/lib"

type ExampleJob struct {

}

func (example ExampleJob) Name() string {
	return "example-job"
}

func (example ExampleJob) Cron() string  {
	return "0/1 * * * * * *"
}

func (example ExampleJob) Execute() error {
	lib.Logger.Infoln("example job run.")
	return nil
}