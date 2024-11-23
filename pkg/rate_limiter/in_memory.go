package rate_limiter

import (
	"sync"
	"time"
)

// inMemoryLimiter mantém o controle de requisições em memória.
type inMemoryLimiter struct {
	limit       int
	window      time.Duration
	requests    map[string][]time.Time
	requestLock sync.Mutex
}

// NewInMemoryLimiter cria um novo limitador em memória.
func NewInMemoryLimiter(limit int, windowSeconds int) Limiter {
	return &inMemoryLimiter{
		limit:    limit,
		window:   time.Duration(windowSeconds) * time.Second,
		requests: make(map[string][]time.Time),
	}
}

// Allow verifica se o cliente pode fazer a requisição.
func (l *inMemoryLimiter) Allow(clientID string) bool {
	l.requestLock.Lock()
	defer l.requestLock.Unlock()

	now := time.Now()
	windowStart := now.Add(-l.window)

	// Filtra requisições fora da janela de tempo
	validRequests := []time.Time{}
	for _, t := range l.requests[clientID] {
		if t.After(windowStart) {
			validRequests = append(validRequests, t)
		}
	}
	l.requests[clientID] = validRequests

	// Verifica se o cliente está dentro do limite
	if len(validRequests) < l.limit {
		l.requests[clientID] = append(l.requests[clientID], now)
		return true
	}

	return false
}

// SetRate permite configurar o limite dinamicamente.
func (l *inMemoryLimiter) SetRate(limit int, windowSeconds int) {
	l.limit = limit
	l.window = time.Duration(windowSeconds) * time.Second
}
