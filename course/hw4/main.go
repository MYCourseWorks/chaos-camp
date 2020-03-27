package main

import (
	"time"

	"github.com/homework/hw4/settings"
	workpool "github.com/homework/hw4/worker-pool"
	"github.com/homework/hw4/worker-pool/tasks"
)

func main() {
	st := settings.GetSettings()
	deadline := time.Now().Add(time.Duration(st.MaxIndexingTimeInMs) * time.Millisecond)

	p := workpool.NewWorkerPool(tasks.NewScrapperTask(5*time.Second), st.MaxRoutines, deadline)
	for _, u := range st.ListOfURLs {
		p.AddWork(&workpool.WorkerItem{Value: u, Priority: 1})
	}

	defer p.ClearWorkQueue()
	for p.HasWork() {
		p.Work()
	}
}
