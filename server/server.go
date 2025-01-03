package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Cotacao struct {
	Bid string `json:"bid"`
}

func main() {
	// Configura o banco de dados SQLite
	db, err := sql.Open("sqlite3", "./cotacoes.db")
	if err != nil {
		log.Fatalf("Erro ao abrir banco de dados: %v", err)
	}
	defer db.Close()

	// Cria tabela se não existir
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS cotacoes (id INTEGER PRIMARY KEY AUTOINCREMENT, bid TEXT, timestamp DATETIME DEFAULT CURRENT_TIMESTAMP)")
	if err != nil {
		log.Fatalf("Erro ao criar tabela: %v", err)
	}

	http.HandleFunc("/cotacao", func(w http.ResponseWriter, r *http.Request) {
		// Contexto para consultar a API
		ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
		defer cancel()

		cotacao, err := buscarCotacao(ctx)
		if err != nil {
			http.Error(w, fmt.Sprintf("Erro ao buscar cotação: %v", err), http.StatusInternalServerError)
			log.Println("Erro ao buscar cotação:", err)
			return
		}

		// Contexto para salvar no banco de dados
		dbCtx, dbCancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
		defer dbCancel()

		if err := salvarCotacao(dbCtx, db, cotacao); err != nil {
			http.Error(w, fmt.Sprintf("Erro ao salvar cotação: %v", err), http.StatusInternalServerError)
			log.Println("Erro ao salvar cotação:", err)
			return
		}

		// Responde ao cliente com o valor da cotação
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(cotacao)
	})

	log.Println("Servidor iniciado na porta 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func buscarCotacao(ctx context.Context) (*Cotacao, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data map[string]map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	return &Cotacao{Bid: data["USDBRL"]["bid"]}, nil
}

func salvarCotacao(ctx context.Context, db *sql.DB, cotacao *Cotacao) error {
	query := "INSERT INTO cotacoes (bid) VALUES (?)"

	ch := make(chan error, 1)
	go func() {
		_, err := db.Exec(query, cotacao.Bid)
		ch <- err
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-ch:
		return err
	}
}