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
	ctx, cancel := context.WithTimeout(context.Background(), 3000*time.Millisecond)
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

    // Lê e decodifica a resposta do servidor
    var cotacao Cotacao
    if err := json.NewDecoder(resp.Body).Decode(&cotacao); err != nil {
        body, _ := ioutil.ReadAll(resp.Body)
        return nil, fmt.Errorf("Erro ao decodificar resposta: %w. JSON recebido: %s", err, string(body))
    }

    if cotacao.Bid == "" {
        return nil, fmt.Errorf("Campo 'bid' vazio na resposta do servidor. JSON recebido: %+v", cotacao)
    }

    return &cotacao, nil
}

func salvarCotacaoEmArquivo(valor string) error {
    conteudo := fmt.Sprintf("Dólar: %s\n", valor)
    caminho := "client/cotacao.txt"
    return ioutil.WriteFile(caminho, []byte(conteudo), 0644)
}