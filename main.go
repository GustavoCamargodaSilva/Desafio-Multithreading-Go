package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type BrasilAPICep struct {
	Cep          string `json:"cep"`
	State        string `json:"state"`
	City         string `json:"city"`
	Neighborhood string `json:"neighborhood"`
	Street       string `json:"street"`
	Service      string `json:"service"`
}

type ViaCep struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Estado      string `json:"estado"`
	Regiao      string `json:"regiao"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

type Resultado struct {
	API      string
	Endereco string
}

func buscaCepAPI(baseURL string, cep string, nomeAPI string, parser func([]byte) (string, error), ch chan<- Resultado) {
	resp, err := http.Get(baseURL + cep)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	endereco, err := parser(body)
	if err != nil {
		return
	}

	ch <- Resultado{
		API:      nomeAPI,
		Endereco: endereco,
	}
}

func parseBrasilAPI(body []byte) (string, error) {
	var dados BrasilAPICep
	if err := json.Unmarshal(body, &dados); err != nil {
		return "", err
	}
	return fmt.Sprintf("CEP: %s, Rua: %s, Bairro: %s, Cidade: %s, Estado: %s",
		dados.Cep, dados.Street, dados.Neighborhood, dados.City, dados.State), nil
}

func parseViaCep(body []byte) (string, error) {
	var dados ViaCep
	if err := json.Unmarshal(body, &dados); err != nil {
		return "", err
	}
	return fmt.Sprintf("CEP: %s, Rua: %s, Bairro: %s, Cidade: %s, Estado: %s",
		dados.Cep, dados.Logradouro, dados.Bairro, dados.Localidade, dados.Uf), nil
}

func BuscaBrasilAPI(cep string, ch chan<- Resultado) {
	buscaCepAPI("https://brasilapi.com.br/api/cep/v1/", cep, "BrasilAPI", parseBrasilAPI, ch)
}

func BuscaViaCep(cep string, ch chan<- Resultado) {
	buscaCepAPI("http://viacep.com.br/ws/", cep, "ViaCEP", parseViaCep, ch)
}

func main() {
	cep := "01153000"
	if len(os.Args) > 1 {
		cep = os.Args[1]
	}

	ch := make(chan Resultado, 2)

	go BuscaBrasilAPI(cep, ch)
	go BuscaViaCep(cep, ch)

	select {
	case resultado := <-ch:
		fmt.Printf("API mais rápida: %s\n", resultado.API)
		fmt.Printf("Endereço: %s\n", resultado.Endereco)
	case <-time.After(1 * time.Second):
		fmt.Println("Erro: timeout - nenhuma API respondeu em 1 segundo")
	}
}
