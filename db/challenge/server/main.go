package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

// O server.go deverá consumir a API contendo o câmbio de Dólar e Real
// no endereço: https://economia.awesomeapi.com.br/json/last/USD-BRL
// em seguida deverá retornar no formato JSON o resultado para o cliente.
// Utilizando o package "context"
// Usando o package "context", o server.go deverá registrar no banco de dados SQLite cada cotação recebida
// sendo que o timeout máximo para chamar a API de cotação do dólar deverá ser de 200ms
// e o timeout máximo para conseguir persistir os dados no banco deverá ser de 10ms.
// O endpoint necessário gerado pelo server.go para este desafio será: /cotacao e
// a porta a ser utilizada pelo servidor HTTP será a 8080.

var DB *sql.DB

const db_file = "sqlite-database.db"
const url = "https://economia.awesomeapi.com.br/json/last/USD-BRL"

type LastOccurence struct {
	ExchangeRate ExchangeRate `json:"USDBRL"`
}

type ExchangeRate struct {
	Code       string `json:"code"`
	Codein     string `json:"codein"`
	Name       string `json:"name"`
	High       string `json:"high"`
	Low        string `json:"low"`
	VarBid     string `json:"varBid"`
	PctChange  string `json:"pctChange"`
	Bid        string `json:"bid"`
	Ask        string `json:"ask"`
	Timestamp  string `json:"timestamp"`
	CreateDate string `json:"create_date"`
}

type Exchange struct {
	ID         string
	Code       string `json:"code"`
	Codein     string `json:"codein"`
	Name       string `json:"name"`
	Bid        string `json:"bid"`
	Timestamp  string `json:"timestamp"`
	CreateDate string `json:"create_date"`
}

func NewExchange(exr ExchangeRate) *Exchange {
	return &Exchange{
		ID:         uuid.New().String(),
		Code:       exr.Code,
		Codein:     exr.Codein,
		Name:       exr.Name,
		Bid:        exr.Bid,
		Timestamp:  exr.Timestamp,
		CreateDate: exr.CreateDate,
	}
}

func main() {
	// DB
	RemoveDB()
	SetDBFile()
	err := OpenDB()
	if err != nil {
		panic(err)
	}
	// Server
	server()

}

func server() {
	go func() {
		http.HandleFunc("/cotacao", handler)
		http.ListenAndServe(":8080", nil)
	}()
}

func handler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	fmt.Fprint(os.Stdout, "processing request\n")
	defer log.Println("[Server] End Request")

	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*200)
	defer cancel()

	exchange, err := client(ctx)
	if err != nil {
		cancel()
	}
	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(exchange)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	w.Write(jsonResp)

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
		defer cancel()
		err = insertExchange(ctx, exchange)
		if err != nil {
			cancel()
		}
	}()

	return
}

func RemoveDB() {
	os.Remove(db_file)
}

func SetDBFile() error {
	file, err := os.Create(db_file)
	if err != nil {
		log.Fatal(err.Error())
		return err
	}
	file.Close()
	log.Println(db_file, " created")

	return nil
}

func OpenDB() error {

	path := "./" + db_file
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return err
	}
	DB = db
	defer db.Close()
	createTable(db)
	return nil
}

func createTable(db *sql.DB) {
	createExchangeTableSQL := `CREATE TABLE exchange (
		"id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
		"hash" TEXT,
		"code" TEXT,
		"codein" TEXT,
		"name" TEXT,
		"bid" TEXT,
		"timestamp" TEXT,
		"createdate" TEXT,
	  );` // SQL Statement for Create Table

	log.Println("Create exchange table...")
	statement, err := db.Prepare(createExchangeTableSQL) // Prepare SQL Statement
	if err != nil {
		log.Fatal(err.Error())
	}
	statement.Exec() // Execute SQL Statements
	log.Println("exchange table created")
	DB = db
}

func insertExchange(ctx context.Context, ex *Exchange) error {
	stmt, err := DB.PrepareContext(ctx, "insert into exchange(hash, code, codein, name, bid, timestamp, createdate) values(?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(ex.ID, ex.Code, ex.Codein, ex.Name, ex.Bid, ex.Timestamp, ex.CreateDate)
	if err != nil {
		panic(err)
	}
	return nil
}

func client(ctx context.Context) (*Exchange, error) {
	var last LastOccurence
	// var exRate ExchangeRate
	err := PerformGET(ctx, &last)
	if err != nil {
		log.Fatal("[Server] error getting exchange - %s \n", err.Error())
		return nil, err
	}

	exchange := NewExchange(last.ExchangeRate)

	return exchange, nil
}

func PerformGET(ctx context.Context, target interface{}) error {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		log.Fatal("[Server] %s \n", err.Error())
		return err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal("[Server] %s \n", err.Error())
		return err
	}
	defer res.Body.Close()
	return json.NewDecoder(res.Body).Decode(target)
}
