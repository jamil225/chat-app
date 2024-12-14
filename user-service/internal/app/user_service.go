package user_service

//
//import (
//	"chat-app/user-service/internal/repository"
//	"database/sql"
//	"github.com/gin-gonic/gin"
//	"log"
//	"net/http"
//	"strconv"
//	"time"
//)
//
//// User structure
//type User struct {
//	ID     int64  `json:"id" binding:"required"`
//	Name   string `json:"name" binding:"required"`
//	Mobile string `json:"mobile" binding:"required"`
//	Email  string `json:"email" binding:"required,email"`
//}
//
//// In-memory store
//var users = make(map[int64]User)
//
//// Save User
//func saveUser(c *gin.Context) {
//	defer LogExecutionTime(time.Now(), "saveUser")
//	var user User
//
//	if err := c.ShouldBindJSON(&user); err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		log.Printf("Error while creating user  : %v\n", err.Error())
//		return
//	}
//	userRepo := repository.NewUserRepository(db)
//	if _, err := userRepo.SaveOrUpdateUser(user); err != nil {
//		log.Printf("Error saving user: %v", err)
//		http.Error(w, "Failed to save user", http.StatusInternalServerError)
//		return
//	}
//	if _, exists := users[user.ID]; exists {
//		c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
//		log.Printf("User already exists with id : %v\n", user.ID)
//		return
//	}
//
//	users[user.ID] = user
//	log.Printf("User created with id : %v\n", user.ID)
//	c.JSON(http.StatusCreated, gin.H{"message": "User saved successfully", "user": user})
//}
//
//// Get User by ID
//func getUser(c *gin.Context) {
//	defer LogExecutionTime(time.Now(), "getUser")
//	idParam := c.Param("id")
//	id, err := strconv.ParseInt(idParam, 10, 64)
//	if err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
//		log.Printf("Invalid user ID : %v\n", idParam)
//		return
//	}
//
//	user, exists := users[id]
//	if !exists {
//		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
//		log.Printf("User not found with id : %v\n", id)
//		return
//	}
//	log.Println("User found with id : ", id)
//	c.JSON(http.StatusOK, user)
//}
//
//// LogExecutionTime logs the time taken by a function to execute
//func LogExecutionTime(start time.Time, functionName string) {
//	elapsed := time.Since(start)
//	log.Printf("Execution time for %s: %s\n", functionName, elapsed)
//}
//
//func getRepo() *repository.UserRepository {
//	dsn := "postgres://jamil:jamil123@host.docker.internal:5432/chatapp?sslmode=disable"
//	log.Printf("Connecting to database: %s", dsn)
//
//	db, err := sql.Open("postgres", dsn)
//	if err != nil {
//		log.Fatalf("Failed to connect to the database: %v", err)
//	}
//	defer db.Close()
//
//	if err := db.Ping(); err != nil {
//		log.Fatalf("Failed to ping the database: %v", err)
//	}
//	log.Println("Connected to the database successfully")
//
//	userRepo := repository.NewUserRepository(db)
//	return userRepo
//}
