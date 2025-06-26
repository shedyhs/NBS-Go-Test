package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

type Item struct {
	Codigo        string  `json:"codigo"`
	Descricao     string  `json:"descricao"`
	Quantidade    int     `json:"quantidade"`
	ValorUnitario float64 `json:"valor_unitario"`
}

type Nota struct {
	NumeroNota string `json:"numero_nota"`
	Itens      []Item `json:"itens"`
}

var cache map[string][]Item

func loadNotas(path string) error {
	file, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var notas []Nota
	if err := json.Unmarshal(file, &notas); err != nil {
		return err
	}

	cache = make(map[string][]Item)
	for _, nota := range notas {
		var itensValidos []Item
		for _, item := range nota.Itens {
			if item.Quantidade > 0 && item.ValorUnitario > 0 {
				itensValidos = append(itensValidos, item)
			}
		}
		cache[nota.NumeroNota] = itensValidos
	}
	return nil
}

func itensHandler(writer http.ResponseWriter, req *http.Request) {
	partes := strings.Split(req.URL.Path, "/")
	if len(partes) != 4 {
		http.NotFound(writer, req)
		return
	}
	numero := partes[2]
	itens, found := cache[numero]
	if !found {
		http.NotFound(writer, req)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(map[string]interface{}{
		"numero_nota": numero,
		"itens":       itens,
	})
}

func healthHandler(writer http.ResponseWriter, req *http.Request) {
	writer.Write([]byte("OK"))
}

func main() {
	if err := loadNotas("itens.json"); err != nil {
		log.Fatalf("Erro ao carregar itens.json: %v", err)
	}

	http.HandleFunc("/notas/", itensHandler)
	http.HandleFunc("/health", healthHandler)

	fmt.Println("Servidor iniciado em http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
