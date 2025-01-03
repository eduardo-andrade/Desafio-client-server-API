package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func main() {
	// Contexto para fazer a requisição
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	cotacao, err := buscarCotacao(ctx)
	if err != nil {
		log.Fatalf("%v", err)
	}

	// Salva a cotação em um arquivo
	if err := salvarCotacaoEmArquivo(cotacao.Bid); err != nil {
		log.Fatalf("Erro ao salvar cotação no arquivo: %v", err)
	}

	log.Println("Cotação salva com sucesso.")
}

type Cotacao struct {
	Bid string `json:"bid"`
}

type Response struct {
	USDBRL Cotacao `json:"USDBRL"`
}

func buscarCotacao(ctx context.Context) (*Cotacao, error) {
    req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
    if err != nil {
        return nil, fmt.Errorf("Erro ao criar requisição: %w", err)
    }

    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return nil, fmt.Errorf("Erro ao fazer requisição: %w", err)
    }
    defer resp.Body.Close()

    // Verifica o status HTTP
    if resp.StatusCode != http.StatusOK {
        body, _ := ioutil.ReadAll(resp.Body)
        return nil, fmt.Errorf("Erro do servidor: %s", string(body))
    }

    var response Response
    if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
        return nil, fmt.Errorf("Erro ao decodificar resposta: %w", err)
    }

    return &response.USDBRL, nil
}

func salvarCotacaoEmArquivo(valor string) error {
	conteudo := fmt.Sprintf("Dólar: %s\n", valor)
	return ioutil.WriteFile("client/cotacao.txt", []byte(conteudo), 0644)
}