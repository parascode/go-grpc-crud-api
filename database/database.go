package database

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Company struct {
	ID			string 		`gorm:"primarykey"`
	Name		string 
	Address		string 
	Location	string
}

func DatabaseConnection() *gorm.DB{
	// host := "localhost"
	// port := "5432"
	// dbName := "postgres"
	// dbUser := "postgres"
	// password := "Paras@435"
	var DB *gorm.DB
	var err error


	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	dbUser := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
	host,
	port,
	dbUser,
	dbName,
	password,	
	)
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	DB.AutoMigrate(Company{})
	if err != nil {
		log.Fatal("Error connecting to the database!", err)
	}
	fmt.Println("Database connection succesfull!!")
	fmt.Println("database setup from database package")

	return DB
}