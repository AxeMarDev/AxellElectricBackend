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

type Message struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Location string `json:"location"`
	Message  string `json:"message"`
	Read     bool   `json:"read"`
}

type Employee struct {
	ID        string `json:"id"`
	CompanyId string `json:"company_id"`
	Firstname string `json:"first_name"`
	Lastname  string `json:"last_name"`
	Email     string `json:"email"`
	Ismaster  bool   `json:"is_master"`
	Username  string `json:"username"`
}

type Company struct {
	ID          string `json:"id"`
	Companyname string `json:"company_name"`
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

func getProjects(c *gin.Context) {
	rows, err := db.Query("SELECT id, name, location, imageurl FROM projects ORDER BY id ASC")
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

func getMessages(c *gin.Context) {
	rows, err := db.Query("SELECT id, email, name, location, message, read FROM messages ORDER BY id ASC")

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query messages"})
		return
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var p Message
		if err := rows.Scan(&p.ID, &p.Email, &p.Name, &p.Location, &p.Message, &p.Read); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan message"})
			return
		}
		messages = append(messages, p)
	}
	fmt.Println(messages)

	c.IndentedJSON(http.StatusOK, messages)
}

func addMessage(c *gin.Context) {
	var newMessage Message

	// Bind the received JSON to newPerson
	if err := c.ShouldBindJSON(&newMessage); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Insert newPerson into the database
	query := `INSERT INTO messages (email, name, location, message, read) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	var id int
	err := db.QueryRow(query, newMessage.Email, newMessage.Name, newMessage.Location, newMessage.Message, newMessage.Read).Scan(&id)

	if err != nil {
		log.Printf("Error while inserting new message: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add new message"})
		return
	}

	// Return the new person as JSON
	c.JSON(http.StatusCreated, newMessage)
}

func updateProject(c *gin.Context) {

	var updatedProject Project

	// Bind the received JSON to newPerson
	if err := c.ShouldBindJSON(&updatedProject); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error when binding json"})
		return
	}

	// Get the id from the URL parameter
	id := c.Query("id") // This method is used for query parameters
	fmt.Println(id)

	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing ID in the URL parameter"})
		return
	}

	query := `UPDATE projects SET name = $2, location = $3, data = $4, imageurl = $5 WHERE id = $1`
	result, err := db.Exec(query, updatedProject.ID, updatedProject.Name, updatedProject.Location, updatedProject.Data, updatedProject.Imageurl)

	if err != nil {
		log.Printf("Error while deleting person: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update project"})
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

	c.JSON(http.StatusOK, gin.H{"message": "Project updated successfully"})
}

func updateMessage(c *gin.Context) {
	var updatedMessage Message

	// Bind the received JSON to newPerson
	if err := c.ShouldBindJSON(&updatedMessage); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error when binding json"})
		return
	}

	query := `UPDATE messages SET email = $2, name = $3, location = $4, message = $5, read = $6 WHERE id = $1`
	result, err := db.Exec(query, updatedMessage.ID, updatedMessage.Email, updatedMessage.Name, updatedMessage.Location, updatedMessage.Message, updatedMessage.Read)

	if err != nil {
		log.Printf("Error while deleting person: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update message"})
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

	c.JSON(http.StatusOK, gin.H{"message": "Message updated successfully"})

}

func getCompanies(c *gin.Context) {
	rows, err := db.Query("SELECT id, company_name read FROM companies ORDER BY id ASC")

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query companies"})
		return
	}
	defer rows.Close()

	var companies []Company
	for rows.Next() {
		var p Company
		if err := rows.Scan(&p.ID, &p.Companyname); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan message"})
			return
		}
		companies = append(companies, p)
	}
	fmt.Println(companies)

	c.IndentedJSON(http.StatusOK, companies)
}

type Mastercompanyemployee struct {
	ID          string `json:"id"`
	Companyname string `json:"company_name"`
	Email       string `json:"email"`
	Username    string `json:"username"`
}

func addCompany(c *gin.Context) {
	var newCompany Mastercompanyemployee

	// Bind the received JSON to newPerson
	if err := c.ShouldBindJSON(&newCompany); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Insert newPerson into the database
	query := `INSERT INTO companies (company_name) VALUES ($1) RETURNING id`
	var id int
	err := db.QueryRow(query, newCompany.Companyname).Scan(&id)

	if err != nil {
		log.Printf("Error while inserting new person: %v", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add new company"})
		return
	}

	var idNew int
	queryMasterEmployee := `INSERT INTO employees (company_id, first_name, last_name, email, is_master, username) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	errEmployee := db.QueryRow(queryMasterEmployee, id, "MASTER", "ADMIN", newCompany.Email, true, newCompany.Username).Scan(&idNew)

	if errEmployee != nil {
		log.Printf("Error while inserting new person: %v", errEmployee.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add master Employee"})
		return
	}

	// Return the new person as JSON
	c.JSON(http.StatusCreated, gin.H{"id": id, "company_name": newCompany.Companyname})
}

func getEmployees(c *gin.Context) {
	rows, err := db.Query("SELECT id, company_id, first_name, last_name, email, is_master read FROM employees ORDER BY id ASC")

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query employees"})
		return
	}
	defer rows.Close()

	var employees []Employee
	for rows.Next() {
		var p Employee
		if err := rows.Scan(&p.ID, &p.CompanyId, &p.Firstname, &p.Lastname, &p.Email, &p.Ismaster); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan employees"})
			return
		}
		employees = append(employees, p)
	}
	fmt.Println(employees)

	c.IndentedJSON(http.StatusOK, employees)
}

func main() {

	initDB()
	defer db.Close()

	router := gin.Default()

	// middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// company
	router.GET("/companies", getCompanies)
	router.POST("/companies", addCompany)

	// employee
	router.GET("/employees", getEmployees)

	// routes for projects
	router.GET("/projects", getProjects)
	router.POST("/projects", postProject)
	router.DELETE("/projects", deleteProject)
	router.PATCH("/projects", updateProject)
	// routes for messages
	router.GET("/messages", getMessages)
	router.POST("/messages", addMessage)
	router.PATCH("/messages", updateMessage)

	err := router.Run("localhost:8080")

	if err != nil {
		return
	}
}
