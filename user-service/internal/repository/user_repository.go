package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"user-service/internal/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (repo *UserRepository) SaveOrUpdateUser(user models.User) (*models.User, bool, error) {
	exists, err := repo.IsUserExist(user.ID)
	if err != nil {
		return nil, false, fmt.Errorf("failed to check if user exists: %w", err)
	}

	if exists {
		// Update existing user
		query := "UPDATE users SET name = $1, email = $2, phone = $3 WHERE id = $4 RETURNING id, name, email, phone"
		log.Printf("Executing query: %s with values: %v, %v, %v, %v", query, user.Name, user.Email, user.Phone, user.ID)
		err := repo.db.QueryRow(query, user.Name, user.Email, user.Phone, user.ID).Scan(&user.ID, &user.Name, &user.Email, &user.Phone)
		if err != nil {
			return nil, false, fmt.Errorf("failed to update user: %w", err)
		}
		log.Println("User updated successfully")
		return &user, false, nil
	} else {
		// Insert new user
		query := "INSERT INTO users (id, name, email, phone) VALUES ($1, $2, $3, $4) RETURNING id, name, email, phone"
		log.Printf("Executing query: %s with values: %v, %v, %v, %v", query, user.ID, user.Name, user.Email, user.Phone)
		err := repo.db.QueryRow(query, user.ID, user.Name, user.Email, user.Phone).Scan(&user.ID, &user.Name, &user.Email, &user.Phone)
		if err != nil {
			return nil, false, fmt.Errorf("failed to save user: %w", err)
		}
		log.Println("User saved successfully")
		return &user, true, nil
	}
}

func (repo *UserRepository) IsUserExist(id int64) (bool, error) {
	query := `SELECT COUNT(1) FROM users WHERE id = $1`
	log.Printf("Executing query: %s with value: %v", query, id)

	var count int
	err := repo.db.QueryRow(query, id).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check if user exists: %w", err)
	}

	return count > 0, nil
}

func (repo *UserRepository) GetUser(id int64) (*models.User, error) {
	query := `SELECT id, name, phone, email FROM users WHERE id = $1`
	log.Printf("Executing query: %s with value: %v", query, id)
	var user models.User
	err := repo.db.QueryRow(query, id).Scan(&user.ID, &user.Name, &user.Phone, &user.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, nil
}

func GetAndConnectDB() *sql.DB {
	dsn := "postgres://jamil:jamil123@host.docker.internal:5432/chatapp?sslmode=disable"
	log.Printf("Connecting to database: %s", dsn)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatalf("Failed to close the database connection: %v", err)
		}
	}(db)

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping the database: %v", err)
	}
	log.Println("Connected to the database successfully")

	return db
}
