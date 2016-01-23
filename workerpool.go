package blogalert

import "sync"

// Worker defines work to be done
type Worker func(*WorkerPool)

// WorkerPool defines workerpool
type WorkerPool struct {
	wg      sync.WaitGroup
	workers chan Worker
}

// NewWorkerPool creates new worker group with n runner
func NewWorkerPool(n int) *WorkerPool {
	wp := &WorkerPool{
		workers: make(chan Worker),
	}
	for i := 0; i < n; i++ {
		go wp.runner()
	}

	return wp
}

// Do adds work to queue
func (wp *WorkerPool) Do(w Worker) {
	if w != nil {
		wp.wg.Add(1)
		// Stop blocking
		go func() {
			wp.workers <- w
		}()
	}
}

func (wp *WorkerPool) runner() {
	for worker := range wp.workers {
		worker(wp)
		wp.wg.Done()
	}
}

// Wait for all current actions to finish
func (wp *WorkerPool) Wait() {
	wp.wg.Wait()
}

// Close worker channel
func (wp *WorkerPool) Close() {
	close(wp.workers)
}
