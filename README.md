# Projeto Desafio Client-Server-API

Este projeto é composto por dois sistemas desenvolvidos em Go:
- **server.go**: Responsável por consultar a cotação do dólar em uma API externa e registrar as informações em um banco de dados SQLite.
- **client.go**: Realiza uma requisição HTTP ao servidor para obter o valor atual da cotação e salva a informação em um arquivo de texto.

## Estrutura de Diretórios
```
project/
├── client/
│   ├── client.go
│   └── cotacao.txt
├── server/
│   ├── server.go
│   └── db/
│       └── cotacoes.db
├── go.mod
└── README.md
```

## Requisitos
- Go 1.20 ou superior
- Biblioteca SQLite para Go ([github.com/mattn/go-sqlite3](https://github.com/mattn/go-sqlite3))
- Conexão com a internet para acessar a API de cotação ([AwesomeAPI](https://economia.awesomeapi.com.br/json/last/USD-BRL))

## Instalação

1. **Clone o repositório**:
   ```bash
   git clone <url-do-repositorio>
   cd project
   ```

2. **Instale a dependência do SQLite**:
   ```bash
   go get github.com/mattn/go-sqlite3
   ```

3. **Inicialize o banco de dados** (opcional):
   - O servidor criará automaticamente o arquivo `cotacoes.db` dentro do diretório `server/db` na primeira execução.

## Execução

### Iniciar o Servidor
1. Navegue até o diretório `server/`:
   ```bash
   cd server
   ```
2. Execute o servidor:
   ```bash
   go run server.go
   ```
   - O servidor será iniciado na porta `8080`.
   - Endpoint disponível: `http://localhost:8080/cotacao`

### Executar o Cliente
1. Navegue até o diretório `client/`:
   ```bash
   cd client
   ```
2. Execute o cliente:
   ```bash
   go run client.go
   ```
   - O cliente salvará o valor da cotação no arquivo `cotacao.txt` no formato:
     ```
     Dólar: {valor}
     ```

## Testes
1. **Verifique o funcionamento do servidor**:
   - Acesse o endpoint diretamente no navegador ou via `curl`:
     ```bash
     curl http://localhost:8080/cotacao
     ```
   - Exemplo de resposta:
     ```json
     {
       "bid": "6.2437"
     }
     ```

2. **Verifique o arquivo `cotacao.txt`**:
   - Após executar o cliente, confirme que o arquivo foi criado e contém o valor da cotação:
     ```bash
     cat cotacao.txt
     ```

## Erros Comuns
- **Erro ao baixar a dependência do SQLite**:
  - Certifique-se de que o Go está configurado corretamente em seu sistema e tente novamente o comando:
    ```bash
    go get github.com/mattn/go-sqlite3
    ```

- **Erro de "context deadline exceeded"**:
  - A API externa pode estar demorando mais do que o tempo configurado. Verifique sua conexão com a internet ou ajuste o timeout no `server.go`.
  
## Licença
Este projeto está sob a licença MIT. Veja o arquivo LICENSE para mais detalhes.