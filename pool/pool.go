package worker

import (
	"context"
	"log"
	"sync"
	"time"
)

type pool struct {
	jobs       chan func()
	maxWorker  int
	sigctx     context.Context
	oncomplete func()

	wg sync.WaitGroup
	mu sync.Mutex

	closed bool
}

// var poolonce sync.Once
// var workpool *pool

func newPool(sigctx context.Context, maxWorker, maxJob int) *pool {
	workpool := &pool{
		jobs:      make(chan func(), maxJob),
		sigctx:    sigctx,
		maxWorker: maxWorker,
	}
	ctx, canc := context.WithCancel(context.Background())
	for i := 0; i < int(maxWorker); i++ {
		newWorker(ctx, workpool.jobs, i+1)
	}
	go func() {
		<-workpool.sigctx.Done()
		workpool.mu.Lock()
		workpool.closed = true
		workpool.mu.Unlock()

		workpool.wg.Wait()
		canc()
		close(workpool.jobs)
		if workpool.oncomplete != nil {
			workpool.oncomplete()
		}
	}()
	return workpool
}

func GetPoolWithSpec(sigctx context.Context, oncomplete func(), workerSize, jobSize int) *pool {
	if sigctx == nil {
		sigctx = context.Background()
	}
	p := newPool(sigctx, workerSize, jobSize)
	p.oncomplete = oncomplete
	return p
}

func (p *pool) AddJob(job func(), timeout time.Duration) bool {
	if timeout == 0 {
		timeout = time.Hour
	}
	defer func() {
		r := recover()
		if _, ok := r.(error); ok {
			log.Println("Recovered")
		}
	}()
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.closed {
		return false
	}
	p.wg.Add(1)
	timer := time.NewTimer(timeout)
	select {
	case p.jobs <- func() { job(); p.wg.Done() }:
		return true
	case <-p.sigctx.Done():
	case <-timer.C:
	}
	defer p.wg.Done()
	return false
}

type worker struct{}

func newWorker(ctx context.Context, jobs <-chan func(), wID int) *worker {
	w := worker{}
	go func() {
		log.Printf("[Pool] Worker %d has been started\n", wID)
		for {
			select {
			case job := <-jobs:
				job()
			case <-ctx.Done():
				return
			}
		}
	}()
	return &w
}
