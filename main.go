// to open postgres cli use 'psql'

package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"time"
)

type Project struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Location string `json:"location"`
	Data     string `json:"date"`
	Imageurl string `json:"imageurl"`
}

var db *sql.DB

func initDB() {

	err1 := godotenv.Load() // This will look for the .env file in the current directory
	if err1 != nil {
		log.Fatal("Error loading .env file")
	}

	dbpath := fmt.Sprintf("user=%s dbname=%s host=%s port=%s sslmode=%s",
		os.Getenv("DBUSER"),
		os.Getenv("DBNAME"),
		os.Getenv("DBHOST"),
		os.Getenv("DBPORT"),
		os.Getenv("DBSSL"))

	fmt.Println("connected to database")

	connStr := dbpath //"user=axellmartinez dbname=mydb  host=localhost port=5432 sslmode=disable"
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
}

func postProject(c *gin.Context) {
	var newProject Project

	// Bind the received JSON to newPerson
	if err := c.ShouldBindJSON(&newProject); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Insert newPerson into the database
	query := `INSERT INTO projects (name, location, imageurl) VALUES ($1, $2, $3) RETURNING id`
	var id int
	err := db.QueryRow(query, newProject.Name, newProject.Location, newProject.Imageurl).Scan(&id)

	if err != nil {
		log.Printf("Error while inserting new person: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add new person"})
		return
	}

	// Return the new person as JSON
	c.JSON(http.StatusCreated, newProject)
}

func getPeople(c *gin.Context) {
	rows, err := db.Query("SELECT id, name, location, imageurl FROM projects")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query people"})
		return
	}
	defer rows.Close()

	var projects []Project
	for rows.Next() {
		var p Project
		if err := rows.Scan(&p.ID, &p.Name, &p.Location, &p.Imageurl); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan person"})
			return
		}
		projects = append(projects, p)
	}
	fmt.Println(projects)

	c.IndentedJSON(http.StatusOK, projects)
}

func deleteProject(c *gin.Context) {
	// Get the id from the URL parameter
	id := c.Query("id") // This method is used for query parameters
	fmt.Println(id)

	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing ID in the URL parameter"})
		return
	}

	// Execute the delete query
	query := `DELETE FROM projects WHERE id = $1`
	result, err := db.Exec(query, id)
	if err != nil {
		log.Printf("Error while deleting person: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete person"})
		return
	}

	// Check how many rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error getting rows affected: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking deletion result"})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No project found with the provided ID"})
		return
	}

	// Return success message
	c.JSON(http.StatusOK, gin.H{"message": "Person deleted successfully"})
}

func main() {

	initDB()
	defer db.Close()

	router := gin.Default()

	// middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.GET("/projects", getPeople)
	router.POST("/projects", postProject)
	router.DELETE("/projects", deleteProject)
	err := router.Run("localhost:8080")

	if err != nil {
		return
	}
}
