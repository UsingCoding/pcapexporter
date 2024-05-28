package analyzer

import "sync"

func NewWP(size uint) *WP {
	c := make(chan Job, size)
	wg := &sync.WaitGroup{}

	for range size {
		worker(c, wg)
	}

	return &WP{
		c:  c,
		wg: wg,
	}
}

// WP - worker pool
type WP struct {
	c  chan Job
	wg *sync.WaitGroup
}

func (w *WP) Post(j Job) {
	w.wg.Add(1)
	w.c <- j
}

// Wait until all workers complete
func (w *WP) Wait() {
	w.wg.Wait()
}

// Close job chan
func (w *WP) Close() {
	close(w.c)
}

type Job func()

func worker(c <-chan Job, wg *sync.WaitGroup) {
	go func() {
		for job := range c {
			job()
			wg.Done()
		}
	}()
}
