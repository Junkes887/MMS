package database

import (
	"fmt"
	"os"

	"github.com/Junkes887/MMS/internal/database/entity"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DBConnection struct {
	DB *gorm.DB
}

func NewConfig() *DBConnection {
	fmt.Println("Conectando ao banco de dados...")

	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	database := os.Getenv("DB_NAME")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s TimeZone=America/Sao_Paulo",
		host, user, password, database, port,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	fmt.Println("Connected to the database")

	err = db.AutoMigrate(
		&entity.MMSEntity{},
	)
	if err != nil {
		panic(fmt.Sprintf("Falha ao migrar tabelas: %v", err))
	}

	fmt.Println("Migração concluída com sucesso!")

	return &DBConnection{
		DB: db,
	}
}
