package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"strconv"
	"time"
	"user-service/internal/models"
	"user-service/internal/repository"
)

func saveUser(c *gin.Context) {
	defer LogExecutionTime(time.Now(), "saveUser")
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		log.Printf("Error while creating user: %v\n", err.Error())
		return
	}
	db := getAndConnectDB(c)
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatalf("Failed to close the database connection: %v", err)
		}
	}(db)

	userRepo := repository.NewUserRepository(db)

	updatedUser, created, err := userRepo.SaveOrUpdateUser(user)
	if err != nil {
		log.Printf("Error saving user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if created {
		log.Printf("User created with id: %v\n", updatedUser.ID)
		c.JSON(http.StatusCreated, gin.H{"message": "User created successfully", "user": updatedUser})
	} else {
		log.Printf("User updated with id: %v\n", updatedUser.ID)
		c.JSON(http.StatusOK, gin.H{"message": "User updated successfully", "user": updatedUser})
	}
}

func getUser(ginContext *gin.Context) {
	defer LogExecutionTime(time.Now(), "getUser")
	idParam := ginContext.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		log.Printf("Invalid user ID : %v\n", idParam)
		return
	}
	db := getAndConnectDB(ginContext)
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatalf("Failed to close the database connection: %v", err)
		}
	}(db)
	userRepo := repository.NewUserRepository(db)
	var user *models.User
	user, err = userRepo.GetUser(id)
	if err != nil {
		log.Printf("Error getting user: %v", err)
		ginContext.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting user", "details": err.Error()})
		return
	}

	if user == nil {
		log.Printf("User not found with id: %v", id)
		ginContext.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	log.Printf("User found with id: %v, user: %+v", id, user)
	ginContext.JSON(http.StatusOK, gin.H{"user": user})

	//err = db.Close()
	//if err != nil {
	//	log.Fatalf("Failed to close the database connection: %v", err)
	//}

}

// LogExecutionTime logs the time taken by a function to execute
func LogExecutionTime(start time.Time, functionName string) {
	elapsed := time.Since(start)
	log.Printf("Execution time for %s: %s\n", functionName, elapsed)
}

func getAndConnectDB(ginContext *gin.Context) *sql.DB {
	postgress := "postgres://jamil:jamil123@host.docker.internal:5432/chatapp?sslmode=disable"
	log.Println("Connecting to database")

	db, err := sql.Open("postgres", postgress)
	if err != nil {
		ginContext.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to the database"})
		log.Printf("Failed to connect to the database: %v\n", err)
		return nil
	}

	if err := db.Ping(); err != nil {
		ginContext.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to the database"})
		log.Printf("Failed to ping the database: %v\n", err)
		return nil
	}
	log.Println("Connected to the database successfully")

	return db
}

func main() {
	r := gin.Default()

	r.POST("/user", saveUser)   // Save user
	r.GET("/user/:id", getUser) // Get user by ID
	log.Fatal(r.Run(":8080"))   // Start server on port 8080
}
