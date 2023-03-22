package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
)

const base = "https://viacep.com.br/ws/"

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

func main() {
	http.HandleFunc("/", SearchAddress)
	http.ListenAndServe(":8080", nil)
}

func SearchAddress(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	param := DoFilterParam(r.URL.Query().Get("cep"))
	if param == "" || len(param) < 8 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	address, error := SearchAddressByVendor(param)
	if error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(address)
	// w.Write([]byte("Buscouuu"))
}

func DoFilterParam(param string) string {
	pex := regexp.MustCompile(`[^a-zA-Z0-9 ]+`).ReplaceAllString(param, "")
	p := strings.ReplaceAll(pex, "-", "")
	fmt.Println("Param: ", p)
	return p
}

func SearchAddressByVendor(zipCode string) (*Address, error) {
	url := base + zipCode + "/json/"

	req, err := http.Get(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[Client] Request error: %v \n", err)
		return nil, err
	}
	defer req.Body.Close()

	res, err := io.ReadAll((req.Body))
	if err != nil {
		fmt.Fprintf(os.Stderr, "[Client] Response error: %v \n", err)
		return nil, err
	}
	var output Address
	err = json.Unmarshal(res, &output)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[Client] Parse error: %v \n", err)
		return nil, err
	}
	return &output, nil
}
