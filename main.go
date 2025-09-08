package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
)

type Mapa struct {
	ID     string `json:"id"`
	Titulo string `json:"titulo"`
}

type Juego struct {
	ID     string `json:"id"`
	Titulo string `json:"titulo"`
	Mapas  []Mapa `json:"mapas"`
}

var (
	juegos   []Juego
	mutex    sync.Mutex
	dataFile = "data.json"
)

func cargarDatos() {
	file, err := os.ReadFile(dataFile)
	if err == nil {
		json.Unmarshal(file, &juegos)
	}
}

func guardarDatos() {
	data, _ := json.MarshalIndent(juegos, "", "  ")
	_ = ioutil.WriteFile(dataFile, data, 0644)
}

func generarIDJuego() string {
	return fmt.Sprintf("juego-%d", len(juegos)+1)
}

func generarIDMapa(juego *Juego) string {
	return fmt.Sprintf("mapa-%d", len(juego.Mapas)+1)
}

func getJuegos(w http.ResponseWriter, _ *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(juegos)
}

func addJuego(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Titulo string `json:"titulo"`
	}

	_ = json.NewDecoder(r.Body).Decode(&input)

	mutex.Lock()
	defer mutex.Unlock()

	j := Juego{
		ID:     generarIDJuego(),
		Titulo: input.Titulo,
		Mapas:  []Mapa{},
	}

	juegos = append(juegos, j)
	guardarDatos()
	w.WriteHeader(http.StatusCreated)
}

func addMapa(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		JuegoID string `json:"juegoId"`
		Titulo  string `json:"titulo"`
	}

	_ = json.NewDecoder(r.Body).Decode(&payload)

	mutex.Lock()
	defer mutex.Unlock()

	for i := range juegos {
		if juegos[i].ID == payload.JuegoID {
			mapa := Mapa{
				ID:     generarIDMapa(&juegos[i]),
				Titulo: payload.Titulo,
			}
			juegos[i].Mapas = append(juegos[i].Mapas, mapa)
			guardarDatos()
			w.WriteHeader(http.StatusCreated)
			return
		}
	}
	http.Error(w, "Juego no encontrado", http.StatusNotFound)
}

func main() {
	cargarDatos()

	http.Handle("/", http.FileServer(http.Dir("./public")))

	http.HandleFunc("/api/juegos", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			getJuegos(w, r)
		} else if r.Method == "POST" {
			addJuego(w, r)
		}
	})
	http.HandleFunc("/api/mapas", addMapa)

	fmt.Println("Servidor corriendo en http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
