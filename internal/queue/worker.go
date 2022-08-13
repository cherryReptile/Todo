package queue

import (
	"github.com/cherryReptile/Todo/internal/interfaces"
)

type JobWorker struct {
	ch chan interfaces.JobInterface
}

func Run() JobWorker {
	ch := make(chan interfaces.JobInterface)

	go worker(ch)

	return JobWorker{ch: ch}
}

func (w *JobWorker) Add(job interfaces.JobInterface) {
	w.ch <- job
}

func worker(ch <-chan interfaces.JobInterface) {
	for job := range ch {
		go job.Run()
	}
}
