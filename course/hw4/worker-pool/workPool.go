package workpool

import (
	"fmt"
	"os"
	"sync"
	"time"

	pq "github.com/jupp0r/go-priority-queue"
)

// Task comment
type Task interface {
	Execute(workItem *WorkerItem) ([]*WorkerItem, error)
}

// WorkerItem comment
type WorkerItem struct {
	Value    string // The value of the item; arbitrary.
	Priority int    // The priority of the item in the queue.
}

// WorkerPool comment
type WorkerPool struct {
	Deadline    time.Time
	workQueue   pq.PriorityQueue
	mux         sync.Mutex
	task        Task
	maxRoutines int
}

// NewWorkerPool comment
func NewWorkerPool(task Task, maxRoutines int, deadline time.Time) *WorkerPool {
	wq := pq.New()
	ret := &WorkerPool{
		workQueue:   wq,
		task:        task,
		maxRoutines: maxRoutines,
		Deadline:    deadline,
	}

	return ret
}

// AddWork comment
func (p *WorkerPool) AddWork(item *WorkerItem) {
	p.workQueue.Insert(item, float64(item.Priority))
}

// HasWork comment
func (p *WorkerPool) HasWork() bool {
	return p.workQueue.Len() > 0 && time.Now().Before(p.Deadline)
}

// ClearWorkQueue comment
func (p *WorkerPool) ClearWorkQueue() {
	p.workQueue = pq.New()
}

// Work comment
func (p *WorkerPool) Work() {
	if !p.HasWork() {
		return
	}

	workerCount := 0
	if p.workQueue.Len() < p.maxRoutines {
		workerCount = p.workQueue.Len()
	} else {
		workerCount = p.maxRoutines
	}

	workCh := make(chan *WorkerItem, workerCount)
	moreWork := make([]*WorkerItem, 0)
	var wg sync.WaitGroup
	wg.Add(workerCount)

	for i := 0; i < workerCount; i++ {
		go func() {
			defer wg.Done()

			for { // Loop forever
				workItem, ok := <-workCh
				if !ok {
					return
				}

				moreWorkArr, err := p.task.Execute(workItem)
				if err != nil {
					fmt.Fprintf(os.Stderr, "%s\n", err.Error())
				}

				if moreWorkArr != nil {
					for _, m := range moreWorkArr {
						p.mux.Lock()
						moreWork = append(moreWork, m)
						p.mux.Unlock()
					}
				}
			}
		}()
	}

	for p.HasWork() {
		item, err := p.workQueue.Pop()
		if err != nil {
			panic(err)
		}

		workItem := item.(*WorkerItem)
		workCh <- workItem
	}

	close(workCh)
	wg.Wait()

	for _, item := range moreWork {
		p.AddWork(item)
	}
}
