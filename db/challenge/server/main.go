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
	Server()

}

func Server() {
	fmt.Println("Run Server at port 8080")
	http.HandleFunc("/cotacao", handler)
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	fmt.Fprint(os.Stdout, "processing request\n")
	defer log.Println("[Server] End Request")

	exchange, err := client(ctx)
	if err != nil {
		fmt.Println("[Server] connection error ", err.Error())
		failRequest(w)
	}
	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(exchange)
	if err != nil {
		log.Fatalf("[Server] Error happened in JSON marshal. Err: %s", err)
		failRequest(w)
	}
	// fmt.Println(string(jsonResp))
	err = insertExchange(ctx, exchange)
	if err != nil {
		log.Fatalf("[Server] Error to insert data. Err: %s", err)
		failRequest(w)
	}

	w.Write(jsonResp)

	return
}

func failRequest(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("500 - Something bad happened!"))
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

	createTable(db)
	return nil
}

func createTable(db *sql.DB) error {
	DB := db
	log.Println("[SQL] Creating exchange table...")
	statement, err := DB.Prepare(`CREATE TABLE exchange (
		id integer PRIMARY KEY,
		hash text,
		code text,
		codein text,
		name text,
		bid text NOT NULL,
		timestamp text,
		create_date text
	);`) // Prepare SQL Statement
	if err != nil {
		log.Fatal(err.Error())
		return err
	}
	statement.Exec() // Execute SQL Statements
	log.Println("[SQL] exchange table created")

	return nil
}

func insertExchange(ctx context.Context, ex *Exchange) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	log.Println("[SQL] statement: insert")
	stmt, err := DB.PrepareContext(ctx, "insert into exchange(hash, code, codein, name, bid, timestamp, create_date) values(?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatalf("[SQL] Error to statement insert. Err: %s", err)
		return err
	}
	log.Println("[SQL] exec query: insert")
	_, err = stmt.Exec(ex.ID, ex.Code, ex.Codein, ex.Name, ex.Bid, ex.Timestamp, ex.CreateDate)
	if err != nil {
		log.Fatalf("[SQL] Error to exec insert. Err: %s", err)
		return err
	}
	defer stmt.Close()
	defer DB.Close()
	return nil
}

func client(ctx context.Context) (*Exchange, error) {
	var last LastOccurence
	// var exRate ExchangeRate
	err := PerformGET(ctx, &last)
	if err != nil {
		log.Fatal("[Server] error getting exchange: ", err.Error())
		return nil, err
	}

	exchange := NewExchange(last.ExchangeRate)
	fmt.Println("[Server] Rx Hash: ", exchange.ID)
	return exchange, nil
}

func PerformGET(ctx context.Context, target interface{}) error {
	start := time.Now()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		log.Fatal("[Server] #1 ", err.Error())
		return err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("[Server] time expanse ", start)
		log.Fatal("[Server] #2 ", err.Error())
		return err
	}
	defer res.Body.Close()
	fmt.Println("[Server] time expanse ", start)
	return json.NewDecoder(res.Body).Decode(target)
}
