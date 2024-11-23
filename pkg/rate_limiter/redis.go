package rate_limiter

import (
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"golang.org/x/net/context"
)

type redisLimiter struct {
	client *redis.Client
	limit  int
	window time.Duration
}

// NewRedisLimiter cria um novo limitador baseado em Redis.
func NewRedisLimiter(client *redis.Client, limit int, windowSeconds int) Limiter {
	return &redisLimiter{
		client: client,
		limit:  limit,
		window: time.Duration(windowSeconds) * time.Second,
	}
}

// Allow verifica se o cliente pode fazer a requisição.
func (l *redisLimiter) Allow(clientID string) bool {
	ctx := context.Background()
	key := fmt.Sprintf("rate_limiter:%s", clientID)
	now := time.Now().Unix()

	// Inicia um pipeline Redis para eficiência
	pipe := l.client.TxPipeline()

	// Adiciona o timestamp atual
	pipe.ZAdd(ctx, key, redis.Z{Score: float64(now), Member: now})

	// Remove timestamps fora da janela
	pipe.ZRemRangeByScore(ctx, key, "0", fmt.Sprintf("%d", now-int64(l.window.Seconds())))

	// Conta o número de requisições na janela
	countCmd := pipe.ZCard(ctx, key)

	// Define a expiração da chave
	pipe.Expire(ctx, key, l.window)

	// Executa o pipeline
	_, err := pipe.Exec(ctx)
	if err != nil {
		log.Printf("Erro ao executar pipeline Redis: %v", err)
		return false
	}

	log.Printf("countCmd: %v: %v", countCmd.Val(), countCmd.Val() <= int64(l.limit))
	// Verifica o número de requisições
	return countCmd.Val() <= int64(l.limit)
}

// SetRate permite configurar o limite dinamicamente.
func (l *redisLimiter) SetRate(limit int, windowSeconds int) {
	l.limit = limit
	l.window = time.Duration(windowSeconds) * time.Second
}
