package functionaltests

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/fbveronez/go-test-case-api/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var TestDB *gorm.DB

func ConnectTestDB() *gorm.DB {
	if TestDB != nil {
		return TestDB
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_USER", "test"),
		getEnv("DB_PASSWORD", "test"),
		getEnv("DB_NAME", "testdb"),
		getEnv("DB_PORT", "5433"),
	)

	var db *gorm.DB
	var err error
	for i := 0; i < 10; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		})
		if err == nil {
			break
		}
		log.Println("Waiting for test DB...", err)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		log.Fatal("Failed to connect test DB:", err)
	}

	// AutoMigrate garante que as tabelas existam
	if err := db.AutoMigrate(&model.Account{}, &model.Transaction{}); err != nil {
		log.Fatal("Failed to migrate test DB:", err)
	}

	TestDB = db
	return TestDB
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}
