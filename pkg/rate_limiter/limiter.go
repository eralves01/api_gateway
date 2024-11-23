package rate_limiter

// Limiter define a interface para limitar requisições.
type Limiter interface {
	Allow(clientID string) bool           // Verifica se o cliente pode fazer uma requisição
	SetRate(limit int, windowSeconds int) // Configura os limites
}
