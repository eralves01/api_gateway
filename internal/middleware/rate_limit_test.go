package middleware

import (
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/eralves01/api_gateway/pkg/rate_limiter"
	"github.com/redis/go-redis/v9"
)

func TestRateLimitMiddleware(t *testing.T) {
	// Configura um cliente Redis mockado para o teste
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // O Redis mockado pode ser usado no ambiente de teste
	})

	// Testa a conexão com o Redis
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		log.Fatalf("Erro ao conectar ao Redis: %v", err)
	}

	log.Println("Conectado ao Redis com sucesso")

	// Define a quantidade de requisições permitidas e a janela de tempo
	limit := 3
	windowSeconds := 60

	// Instancia o rate limiter (redis ou in-memory)
	limiter := rate_limiter.NewRedisLimiter(client, limit, windowSeconds)

	log.Println(limiter)

	// Cria o "response recorder" para capturar a resposta
	rr := httptest.NewRecorder()

	// Simula uma requisição HTTP
	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Define o próximo handler que será chamado após o middleware
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Apenas retorna uma resposta simples para testar
		w.WriteHeader(http.StatusOK)
	})

	// Aplica o middleware de rate limiting
	handler := RateLimitMiddleware(limiter)(next)

	// Testa a primeira requisição (deve passar)
	handler.ServeHTTP(rr, req)

	// Verifica se o código de status da resposta é 200 OK
	if rr.Code != http.StatusOK {
		t.Errorf("Esperado o status %v, mas obteve %v", http.StatusOK, rr.Code)
	}

	time.Sleep(time.Duration(1) * time.Second)

	// Simula uma segunda requisição (deve passar)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("Esperado o status %v, mas obteve %v", http.StatusOK, rr.Code)
	}

	time.Sleep(time.Duration(1) * time.Second)

	// Simula uma terceira requisição (deve passar)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("Esperado o status %v, mas obteve %v", http.StatusOK, rr.Code)
	}

	time.Sleep(time.Duration(1) * time.Second)

	// Simula uma quarta requisição (deve retornar 429 Too Many Requests)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTooManyRequests {
		t.Errorf("Esperado o status %v, mas obteve %v", http.StatusTooManyRequests, rr.Code)
	}
}

func TestRateLimitMiddleware_ResetAfterWindow(t *testing.T) {
	// Configura um cliente Redis mockado para o teste
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // O Redis mockado pode ser usado no ambiente de teste
	})

	// Define a quantidade de requisições permitidas e a janela de tempo
	limit := 3
	windowSeconds := 60

	// Instancia o rate limiter (redis ou in-memory)
	limiter := rate_limiter.NewRedisLimiter(client, limit, windowSeconds)

	// Cria o "response recorder" para capturar a resposta
	rr := httptest.NewRecorder()

	// Simula uma requisição HTTP
	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Define o próximo handler que será chamado após o middleware
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Apenas retorna uma resposta simples para testar
		w.WriteHeader(http.StatusOK)
	})

	// Aplica o middleware de rate limiting
	handler := RateLimitMiddleware(limiter)(next)

	// Simula as 3 primeiras requisições (devem passar)
	for i := 0; i < 3; i++ {
		rr = httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
		if rr.Code != http.StatusOK {
			t.Errorf("Esperado o status %v, mas obteve %v", http.StatusOK, rr.Code)
		}
	}

	// Espera o tempo da janela para resetar as contagens
	time.Sleep(time.Duration(windowSeconds) * time.Second)

	// Simula uma nova requisição após o reset da janela (deve passar)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("Esperado o status %v, mas obteve %v", http.StatusOK, rr.Code)
	}
}
