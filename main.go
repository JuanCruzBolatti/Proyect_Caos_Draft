package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"database/sql"

	_ "github.com/mattn/go-sqlite3"
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
	db *sql.DB
)

func conexionDB() {
	var err error
	db, err = sql.Open("sqlite3", "./banphase.db")
	if err != nil {
		fmt.Println("Error al abrir la DB:", err)
		return
	}

	// Crear tablas si no existen
	sqlJuegos := `
	CREATE TABLE IF NOT EXISTS juegos (
		id TEXT PRIMARY KEY,
		titulo TEXT NOT NULL
	);
	`
	sqlMapas := `
	CREATE TABLE IF NOT EXISTS mapas (
		id TEXT PRIMARY KEY,
		juego_id TEXT NOT NULL,
		titulo TEXT NOT NULL,
		baneado INTEGER DEFAULT 0,
		FOREIGN KEY(juego_id) REFERENCES juegos(id)
	);
	`

	_, err = db.Exec(sqlJuegos)
	if err != nil {
		fmt.Println("Error al crear tabla juegos:", err)
		return
	}

	_, err = db.Exec(sqlMapas)
	if err != nil {
		fmt.Println("Error al crear tabla mapas:", err)
		return
	}

	fmt.Println("Conexión a SQLite exitosa y tablas listas")
}

func getJuegos(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Traer todos los juegos
	rows, err := db.Query("SELECT id, titulo FROM juegos")
	if err != nil {
		http.Error(w, "Error al leer juegos", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var juegosResp []Juego

	for rows.Next() {
		var j Juego
		if err := rows.Scan(&j.ID, &j.Titulo); err != nil {
			continue
		}

		// Inicializamos el slice para que nunca sea nil
		j.Mapas = []Mapa{}

		// Traer mapas de este juego
		mapasRows, err := db.Query("SELECT id, titulo FROM mapas WHERE juego_id = ?", j.ID)
		if err == nil {
			for mapasRows.Next() {
				var m Mapa
				if err := mapasRows.Scan(&m.ID, &m.Titulo); err != nil {
					continue
				}
				j.Mapas = append(j.Mapas, m)
			}
			mapasRows.Close()
		}

		juegosResp = append(juegosResp, j)
	}

	// Devolver JSON
	json.NewEncoder(w).Encode(juegosResp)
}

func addJuego(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Titulo string `json:"titulo"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil || input.Titulo == "" {
		http.Error(w, "Título requerido", http.StatusBadRequest)
		return
	}

	// Generar un ID simple
	id := fmt.Sprintf("juego-%d", time.Now().UnixNano())

	_, err := db.Exec("INSERT INTO juegos(id, titulo) VALUES (?, ?)", id, input.Titulo)
	if err != nil {
		http.Error(w, "Error al guardar juego", http.StatusInternalServerError)
		return
	}

	// Devolver el juego creado
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(Juego{ID: id, Titulo: input.Titulo, Mapas: []Mapa{}})
}

func addMapa(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		JuegoID string `json:"juegoId"`
		Titulo  string `json:"titulo"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil || payload.Titulo == "" || payload.JuegoID == "" {
		http.Error(w, "Datos incompletos", http.StatusBadRequest)
		return
	}

	// Generar ID único
	mapaID := fmt.Sprintf("mapa-%d", time.Now().UnixNano())

	_, err := db.Exec("INSERT INTO mapas(id, juego_id, titulo) VALUES (?, ?, ?)", mapaID, payload.JuegoID, payload.Titulo)
	if err != nil {
		http.Error(w, "Error al guardar mapa", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(Mapa{ID: mapaID, Titulo: payload.Titulo})
}

func main() {
	conexionDB()

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
