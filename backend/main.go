package main

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type Match struct {
	ID     string `json:"id"`
	Team1  string `json:"team1"`
	Team2  string `json:"team2"`
	Score1 int    `json:"score1"`
	Score2 int    `json:"score2"`
	Date   string `json:"date"`
}

var db *sql.DB

func main() {
	// Conexión a PostgreSQL con reintentos
	connStr := "user=postgres dbname=laliga password=postgres host=db sslmode=disable"
	var err error

	// Intenta conectar varias veces
	for i := 0; i < 5; i++ {
		db, err = sql.Open("postgres", connStr)
		if err == nil {
			err = db.Ping()
			if err == nil {
				break
			}
		}
		log.Printf("Intento de conexión %d fallido: %v", i+1, err)
		time.Sleep(5 * time.Second)
	}

	if err != nil {
		log.Fatal("No se pudo conectar a PostgreSQL:", err)
	}
	defer db.Close()

	log.Println("Conexión a PostgreSQL establecida correctamente")

	// Configura el router Gin
	r := gin.Default()

	// Endpoints
	r.GET("/api/matches", getMatches)
	r.GET("/api/matches/:id", getMatchByID)
	r.POST("/api/matches", createMatch)
	r.PUT("/api/matches/:id", updateMatch)
	r.DELETE("/api/matches/:id", deleteMatch)

	// Sirve el frontend estático
	r.StaticFile("/", "./frontend/LaLigaTracker.html")

	r.Run(":8080")
}

// Obtener todos los partidos
func getMatches(c *gin.Context) {
	if db == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection not established"})
		return
	}

	rows, err := db.Query("SELECT id, team1, team2, score1, score2, date FROM matches")
	if err != nil {
		log.Printf("Error en consulta: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var matches []Match
	for rows.Next() {
		var m Match
		if err := rows.Scan(&m.ID, &m.Team1, &m.Team2, &m.Score1, &m.Score2, &m.Date); err != nil {
			log.Printf("Error escaneando fila: %v", err)
			continue
		}
		matches = append(matches, m)
	}
	c.JSON(http.StatusOK, matches)
}

// Obtener un partido por ID
func getMatchByID(c *gin.Context) {
	id := c.Param("id")
	var match Match

	err := db.QueryRow("SELECT id, team1, team2, score1, score2, date FROM matches WHERE id = $1", id).Scan(
		&match.ID, &match.Team1, &match.Team2, &match.Score1, &match.Score2, &match.Date)
	
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Partido no encontrado"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	
	c.JSON(http.StatusOK, match)
}

// Crear un nuevo partido
func createMatch(c *gin.Context) {
	var match Match
	if err := c.ShouldBindJSON(&match); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := db.QueryRow(
		`INSERT INTO matches (team1, team2, score1, score2, date) 
		VALUES ($1, $2, $3, $4, $5) 
		RETURNING id`,
		match.Team1, match.Team2, match.Score1, match.Score2, match.Date,
	).Scan(&match.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, match)
}

// Actualizar un partido existente
func updateMatch(c *gin.Context) {
	id := c.Param("id")
	var match Match

	if err := c.ShouldBindJSON(&match); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convertir el ID del string a int para asegurar que es válido
	_, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	result, err := db.Exec(
		`UPDATE matches 
		SET team1 = $1, team2 = $2, score1 = $3, score2 = $4, date = $5 
		WHERE id = $6`,
		match.Team1, match.Team2, match.Score1, match.Score2, match.Date, id,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Partido no encontrado"})
		return
	}

	match.ID = id
	c.JSON(http.StatusOK, match)
}

// Eliminar un partido
func deleteMatch(c *gin.Context) {
	id := c.Param("id")

	result, err := db.Exec("DELETE FROM matches WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Partido no encontrado"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Partido eliminado correctamente"})
}