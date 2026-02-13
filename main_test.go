package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestBuscaBrasilAPI(t *testing.T) {
	servidor := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resposta := BrasilAPICep{
			Cep:          "01153000",
			State:        "SP",
			City:         "São Paulo",
			Neighborhood: "Barra Funda",
			Street:       "Rua Vitorino Carmilo",
			Service:      "open-cep",
		}
		json.NewEncoder(w).Encode(resposta)
	}))
	defer servidor.Close()

	ch := make(chan Resultado, 1)
	go buscaCepAPI(servidor.URL+"/", "01153000", "BrasilAPI", parseBrasilAPI, ch)

	select {
	case resultado := <-ch:
		if resultado.API != "BrasilAPI" {
			t.Errorf("esperava API BrasilAPI, recebeu %s", resultado.API)
		}
		if resultado.Endereco == "" {
			t.Error("endereço não deveria estar vazio")
		}
	case <-time.After(2 * time.Second):
		t.Error("timeout esperando resposta da BrasilAPI")
	}
}

func TestBuscaViaCep(t *testing.T) {
	servidor := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resposta := ViaCep{
			Cep:        "01153-000",
			Logradouro: "Rua Vitorino Carmilo",
			Bairro:     "Barra Funda",
			Localidade: "São Paulo",
			Uf:         "SP",
		}
		json.NewEncoder(w).Encode(resposta)
	}))
	defer servidor.Close()

	ch := make(chan Resultado, 1)
	go buscaCepAPI(servidor.URL+"/", "01153000", "ViaCEP", parseViaCep, ch)

	select {
	case resultado := <-ch:
		if resultado.API != "ViaCEP" {
			t.Errorf("esperava API ViaCEP, recebeu %s", resultado.API)
		}
		if resultado.Endereco == "" {
			t.Error("endereço não deveria estar vazio")
		}
	case <-time.After(2 * time.Second):
		t.Error("timeout esperando resposta da ViaCEP")
	}
}

func TestAPIMaisRapidaVence(t *testing.T) {
	servidorLento := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(500 * time.Millisecond)
		resposta := BrasilAPICep{Cep: "01153000", City: "São Paulo", State: "SP"}
		json.NewEncoder(w).Encode(resposta)
	}))
	defer servidorLento.Close()

	servidorRapido := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resposta := ViaCep{Cep: "01153-000", Localidade: "São Paulo", Uf: "SP"}
		json.NewEncoder(w).Encode(resposta)
	}))
	defer servidorRapido.Close()

	ch := make(chan Resultado, 2)

	go buscaCepAPI(servidorLento.URL+"/", "01153000", "BrasilAPI", parseBrasilAPI, ch)
	go buscaCepAPI(servidorRapido.URL+"/", "01153000", "ViaCEP", parseViaCep, ch)

	select {
	case resultado := <-ch:
		if resultado.API != "ViaCEP" {
			t.Errorf("esperava ViaCEP (mais rápida), recebeu %s", resultado.API)
		}
	case <-time.After(2 * time.Second):
		t.Error("timeout esperando resposta")
	}
}

func TestTimeout(t *testing.T) {
	servidorLento := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
		resposta := BrasilAPICep{Cep: "01153000"}
		json.NewEncoder(w).Encode(resposta)
	}))
	defer servidorLento.Close()

	ch := make(chan Resultado, 1)
	go buscaCepAPI(servidorLento.URL+"/", "01153000", "BrasilAPI", parseBrasilAPI, ch)

	select {
	case <-ch:
		t.Error("não deveria receber resposta antes do timeout")
	case <-time.After(1 * time.Second):
		// Timeout esperado - teste passou
	}
}

func TestParseBrasilAPI(t *testing.T) {
	input := `{"cep":"01153000","state":"SP","city":"São Paulo","neighborhood":"Barra Funda","street":"Rua Vitorino Carmilo","service":"open-cep"}`
	resultado, err := parseBrasilAPI([]byte(input))
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	esperado := "CEP: 01153000, Rua: Rua Vitorino Carmilo, Bairro: Barra Funda, Cidade: São Paulo, Estado: SP"
	if resultado != esperado {
		t.Errorf("esperava %q, recebeu %q", esperado, resultado)
	}
}

func TestParseViaCep(t *testing.T) {
	input := `{"cep":"01153-000","logradouro":"Rua Vitorino Carmilo","bairro":"Barra Funda","localidade":"São Paulo","uf":"SP"}`
	resultado, err := parseViaCep([]byte(input))
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	esperado := "CEP: 01153-000, Rua: Rua Vitorino Carmilo, Bairro: Barra Funda, Cidade: São Paulo, Estado: SP"
	if resultado != esperado {
		t.Errorf("esperava %q, recebeu %q", esperado, resultado)
	}
}

func TestParseJSONInvalido(t *testing.T) {
	_, err := parseBrasilAPI([]byte("json invalido"))
	if err == nil {
		t.Error("esperava erro para JSON inválido")
	}

	_, err = parseViaCep([]byte("json invalido"))
	if err == nil {
		t.Error("esperava erro para JSON inválido")
	}
}
