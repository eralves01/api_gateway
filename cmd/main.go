package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/eralves01/api_gateway/internal/router"
	"github.com/eralves01/api_gateway/pkg/rate_limiter"
	"github.com/redis/go-redis/v9"
)

func main() {
	// Configura o cliente Redis
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // Endereço do Redis
	})

	// Testa a conexão com o Redis
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		log.Fatalf("Erro ao conectar ao Redis: %v", err)
	}

	log.Println("Conectado ao Redis com sucesso")

	// Configurações do Rate Limiter
	requestLimit := 3 // Número de requisições permitidas
	windowSeconds := 60

	// Instancia o Rate Limiter com Redis
	limiter := rate_limiter.NewRedisLimiter(client, requestLimit, windowSeconds)

	// Configura o roteador
	r := router.SetupRouter(limiter)

	// Porta do servidor (pode ser configurada via variável de ambiente)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Porta padrão
	}

	// Inicia o servidor HTTP
	log.Printf("Servidor rodando na porta %s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("Erro ao iniciar o servidor: %v", err)
	}
}
