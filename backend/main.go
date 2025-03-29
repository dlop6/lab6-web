package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	ginSwagger "github.com/swaggo/gin-swagger"
	swaggerFiles "github.com/swaggo/files"
)

// Match representa un partido de fútbol
// @Summary Estructura de un partido
// @Description Información detallada sobre un partido de La Liga
type Match struct {
	ID         string `json:"id"`
	Team1      string `json:"team1" binding:"required"`
	Team2      string `json:"team2" binding:"required"`
	Score1     int    `json:"score1"`
	Score2     int    `json:"score2"`
	Date       string `json:"date" binding:"required"`
	YellowCards int    `json:"yellowCards"`
	RedCards    int    `json:"redCards"`
	ExtraTime   int    `json:"extraTime"`
}

var db *sql.DB

// @title La Liga Tracker API
// @version 1.0
// @description API para gestión de partidos de fútbol
// @host localhost:8080
// @BasePath /api
func main() {
	// Conexión a PostgreSQL con reintentos
	connStr := "user=postgres dbname=laliga password=postgres host=db sslmode=disable"
	var err error

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

	// Configura el router Gin con CORS
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Endpoints
	r.GET("/api/matches", getMatches)
	r.GET("/api/matches/:id", getMatchByID)
	r.POST("/api/matches", createMatch)
	r.PUT("/api/matches/:id", updateMatch)
	r.DELETE("/api/matches/:id", deleteMatch)

	// Nuevos endpoints PATCH
	r.PATCH("/api/matches/:id/goals", updateGoals)
	r.PATCH("/api/matches/:id/yellowcards", updateYellowCards)
	r.PATCH("/api/matches/:id/redcards", updateRedCards)
	r.PATCH("/api/matches/:id/extratime", updateExtraTime)

	// Documentación Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Sirve el frontend estático
	r.StaticFile("/", "./frontend/LaLigaTracker.html")

	r.Run(":8080")
}

// @Summary Obtener todos los partidos
// @Description Devuelve la lista completa de partidos
// @Produce json
// @Success 200 {array} Match
// @Router /api/matches [get]
func getMatches(c *gin.Context) {
	rows, err := db.Query("SELECT id, team1, team2, score1, score2, date, yellow_cards, red_cards, extra_time FROM matches")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var matches []Match
	for rows.Next() {
		var m Match
		if err := rows.Scan(&m.ID, &m.Team1, &m.Team2, &m.Score1, &m.Score2, &m.Date, &m.YellowCards, &m.RedCards, &m.ExtraTime); err != nil {
			log.Printf("Error escaneando fila: %v", err)
			continue
		}
		matches = append(matches, m)
	}
	c.JSON(http.StatusOK, matches)
}

// @Summary Actualizar goles
// @Description Incrementa los goles de un equipo
// @Param id path string true "ID del partido"
// @Param team body string true "Equipo (team1 o team2)"
// @Param goals body int true "Goles a añadir"
// @Produce json
// @Success 200 {object} Match
// @Router /api/matches/{id}/goals [patch]
func updateGoals(c *gin.Context) {
	id := c.Param("id")
	var update struct {
		Team  string `json:"team" binding:"required,oneof=team1 team2"`
		Goals int    `json:"goals" binding:"required,min=1"`
	}

	if err := c.ShouldBindJSON(&update); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var field string
	if update.Team == "team1" {
		field = "score1"
	} else {
		field = "score2"
	}

	_, err := db.Exec(fmt.Sprintf("UPDATE matches SET %s = %s + $1 WHERE id = $2", field, field), update.Goals, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Devuelve el partido actualizado
	var match Match
	err = db.QueryRow("SELECT id, team1, team2, score1, score2, date FROM matches WHERE id = $1", id).Scan(
		&match.ID, &match.Team1, &match.Team2, &match.Score1, &match.Score2, &match.Date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, match)
}

// @Summary Registrar tarjeta amarilla
// @Description Incrementa el contador de tarjetas amarillas
// @Param id path string true "ID del partido"
// @Produce json
// @Success 200 {object} Match
// @Router /api/matches/{id}/yellowcards [patch]
func updateYellowCards(c *gin.Context) {
	id := c.Param("id")
	_, err := db.Exec("UPDATE matches SET yellow_cards = yellow_cards + 1 WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Tarjeta amarilla registrada"})
}

// @Summary Registrar tarjeta roja
// @Description Incrementa el contador de tarjetas rojas
// @Param id path string true "ID del partido"
// @Produce json
// @Success 200 {object} Match
// @Router /api/matches/{id}/redcards [patch]
func updateRedCards(c *gin.Context) {
	id := c.Param("id")
	_, err := db.Exec("UPDATE matches SET red_cards = red_cards + 1 WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Tarjeta roja registrada"})
}

// @Summary Actualizar tiempo extra
// @Description Establece minutos de tiempo extra
// @Param id path string true "ID del partido"
// @Param minutes body int true "Minutos de tiempo extra"
// @Produce json
// @Success 200 {object} Match
// @Router /api/matches/{id}/extratime [patch]
func updateExtraTime(c *gin.Context) {
	id := c.Param("id")
	var update struct {
		Minutes int `json:"minutes" binding:"required,min=1,max=15"`
	}

	if err := c.ShouldBindJSON(&update); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := db.Exec("UPDATE matches SET extra_time = $1 WHERE id = $2", update.Minutes, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Tiempo extra actualizado"})
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