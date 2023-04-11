package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

// O client.go deverá realizar uma requisição HTTP no server.go solicitando a cotação do dólar.
// O client.go precisará receber do server.go apenas o valor atual do câmbio (campo "bid" do JSON).
// Utilizando o package "context"
// o client.go terá um timeout máximo de 300ms para receber o resultado do server.go.
// O client.go terá que salvar a cotação atual em um arquivo "cotacao.txt" no formato: Dólar: {valor}
// O endpoint necessário gerado pelo server.go para este desafio será: /cotacao e
// a porta a ser utilizada pelo servidor HTTP será a 8080.

const filename = "cotacao.txt"
const url = "http://localhost:8080/cotacao"

type Exchange struct {
	ID         string
	Code       string `json:"code"`
	Codein     string `json:"codein"`
	Name       string `json:"name"`
	Bid        string `json:"bid"`
	Timestamp  string `json:"timestamp"`
	CreateDate string `json:"create_date"`
}

func init() {
	os.Remove(filename)
}

func main() {

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*3000)
	defer cancel()

	var exchange Exchange
	err := PerformGET(ctx, &exchange)
	if err != nil {
		log.Fatal("[Client] error getting exchange by Server \n", err.Error())
		panic(err)
	}
	RecordOnFile(&exchange)
}

func PerformGET(ctx context.Context, target interface{}) error {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		log.Fatal("[Client] ", err.Error())
		return err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal("[Client] ", err.Error())
		return err
	}
	fmt.Println("[Client] StatusCode: ", res.StatusCode)
	defer res.Body.Close()
	return json.NewDecoder(res.Body).Decode(target)
}

func RecordOnFile(ex *Exchange) {
	file, err := os.Create(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[Client] Create file error: %v \n", err)
	}
	defer file.Close()
	_, err = file.WriteString(fmt.Sprintf("Dólar: %s", ex.Bid))
	if err != nil {
		fmt.Fprintf(os.Stderr, "[Client] Write file error: %v \n", err)
	}
	fmt.Println("[Client] Success file Created")
}
