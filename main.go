package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

type CepBrasilApiResponse struct {
	Cep          string `json:"cep"`
	Street       string `json:"street"`
	Neighborhood string `json:"neighborhood"`
	City         string `json:"city"`
	State        string `json:"state"`
}

type CepViaCepResponse struct {
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

func buscarDadosBrasilApiCEP(cep string, c chan<- *CepBrasilApiResponse) {
	url := fmt.Sprintf("https://brasilapi.com.br/api/cep/v2/%s", cep)
	response, err := http.Get(url)

	if err != nil {
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return
	}
	var cepResponse CepBrasilApiResponse
	err = json.NewDecoder(response.Body).Decode(&cepResponse)

	if err != nil {
		return
	}

	c <- &cepResponse
}

func buscarDadosViaCepCEP(cep string, c chan<- *CepViaCepResponse) {
	url := fmt.Sprintf("http://viacep.com.br/ws/%s/json/", cep)
	response, err := http.Get(url)

	if err != nil {
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return
	}

	var cepResponse CepViaCepResponse
	err = json.NewDecoder(response.Body).Decode(&cepResponse)
	if err != nil {
		return
	}

	c <- &cepResponse
}

func main() {
	arg := os.Args[1]

	c1 := make(chan *CepBrasilApiResponse)
	c2 := make(chan *CepViaCepResponse)

	go buscarDadosBrasilApiCEP(arg, c1)
	go buscarDadosViaCepCEP(arg, c2)

	select {
	case cepBrasil := <-c1:
		fmt.Println(cepBrasil, "Brasil API")
	case cepViaCep := <-c2:
		fmt.Println(cepViaCep, "Via Cep")
	case <-time.After(time.Second):
		fmt.Println("Timeout")
	}
}
