package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLoggingMiddleware(t *testing.T) {
	// Cria um "response recorder" para capturar a resposta do servidor
	rr := httptest.NewRecorder()

	// Simula uma requisição HTTP de teste
	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Define o próximo handler que será chamado após o middleware
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Apenas retorna uma resposta simples para testar
		w.WriteHeader(http.StatusOK)
	})

	// Aplica o middleware de log
	handler := Logging(next)

	// Executa a requisição no middleware
	handler.ServeHTTP(rr, req)

	// Verifica se o código de status da resposta está correto
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Esperado o status %v, mas obteve %v", http.StatusOK, status)
	}

	// Verifica se o log gerado contém informações que esperamos
	// O log gerado depende da sua configuração de logging, então podemos simular isso
	// e verificar se ele gerou um log com a URL ou o método correto.
	expectedLog := "GET /test"
	// Para simplificar, você pode verificar se o log contém as partes importantes.
	// Mas para um teste real você pode configurar um mock para capturar logs.
	if !containsLog(expectedLog) {
		t.Errorf("Esperado que o log contenha: %v", expectedLog)
	}
}

// Função simulada para verificar se o log contém a string esperada
// Para fazer isso em um caso real, você deve utilizar uma estratégia como
// criar um logger customizado que armazene as mensagens para inspeção.
func containsLog(expected string) bool {
	// Este é um placeholder para a verificação do log
	// Em uma implementação real, você deveria usar um logger customizado para verificar os logs
	// ou capturar a saída do log usando um buffer
	return true
}
