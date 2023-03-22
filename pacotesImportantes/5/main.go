package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

const base = "https://viacep.com.br/ws/"
const filename = "mail.txt"

type Address struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

func SearchAddress(s string) Address {
	url := base + s + "/json/"

	req, err := http.Get(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Request error: %v \n", err)
	}
	defer req.Body.Close()

	res, err := io.ReadAll((req.Body))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Response error: %v \n", err)
	}
	var output Address
	err = json.Unmarshal(res, &output)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Parse error: %v \n", err)
	}
	return output
}

func RecordOnFile(a *Address) {
	file, err := os.Create(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Create file error: %v \n", err)
	}
	defer file.Close()
	_, err = file.WriteString(fmt.Sprintf("CEP: %s, Localidade: %s, UF: %s", a.Cep, a.Localidade, a.Uf))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Write file error: %v \n", err)
	}
	fmt.Println("Success file Created")
}

func main() {
	var myAddress Address
	for _, cep := range os.Args[1:] {
		myAddress = SearchAddress(cep)
		fmt.Println(myAddress)
		RecordOnFile(&myAddress)
	}
}
