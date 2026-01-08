package crawler

import "sync"

// WaitGroupWrapper wraps sync.WaitGroup
// to avoid misuse across components
type WaitGroupWrapper struct {
	wg sync.WaitGroup
}

func (w *WaitGroupWrapper) Add(delta int) {
	w.wg.Add(delta)
}

func (w *WaitGroupWrapper) Done() {
	w.wg.Done()
}

func (w *WaitGroupWrapper) Wait() {
	w.wg.Wait()
}
