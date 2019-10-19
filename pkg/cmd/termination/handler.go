package termination

import "sync"

func NewTerminationHandler(attempts int) *Handler {
	return &Handler{terminationAttemptsLeft: attempts}
}

type Handler struct {
	terminationAttemptsLeft int
	mx                      sync.RWMutex
}

func (h *Handler) SignalTermination() {
	h.mx.Lock()
	defer h.mx.Unlock()

	if h.terminationAttemptsLeft <= 0 {
		return
	}

	h.terminationAttemptsLeft--
}

func (h *Handler) ShouldTerminate() bool {
	h.mx.RLock()
	defer h.mx.RUnlock()

	return h.terminationAttemptsLeft <= 0
}
