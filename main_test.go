package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	loadNotas("itens.json")
	os.Exit(m.Run())
}

func TestHealthEndpoint(t *testing.T) {
	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	healthHandler(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("Esperado status 200, obtido %d", res.StatusCode)
	}

	body, _ := ioutil.ReadAll(res.Body)
	if string(body) != "OK" {
		t.Errorf("Resposta inesperada: %s", string(body))
	}
}

func TestNotaValida(t *testing.T) {
	req := httptest.NewRequest("GET", "/notas/12345/itens", nil)
	writer := httptest.NewRecorder()

	itensHandler(writer, req)

	res := writer.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("Esperado status 200, obtido %d", res.StatusCode)
	}

	var data map[string]interface{}
	json.NewDecoder(res.Body).Decode(&data)

	if data["numero_nota"] != "12345" {
		t.Errorf("Esperado numero_nota 12345, obtido %v", data["numero_nota"])
	}

	itens := data["itens"].([]interface{})
	if len(itens) != 2 {
		t.Errorf("Esperado 2 itens, obtido %d", len(itens))
	}
}

func TestNotaInvalida(t *testing.T) {
	req := httptest.NewRequest("GET", "/notas/99999/itens", nil)
	w := httptest.NewRecorder()

	itensHandler(w, req)

	res := w.Result()
	if res.StatusCode != http.StatusNotFound {
		t.Errorf("Esperado status 404, obtido %d", res.StatusCode)
	}
}
