package main

import (
	"log"

	"backend-app/internal/domain"
	"backend-app/internal/config"
	"backend-app/pkg/db"
	"gorm.io/gorm"
)

func main() {
	cfg := config.Load()
	database, err := db.NewPostgres(cfg)
	if err != nil {
		log.Fatal(err)
	}

	// Создаём ENUM тип для пользователей
	createUserStatusEnum(database)

	// Миграции
	err = database.AutoMigrate(
		&domain.User{},
		&domain.Group{},
		&domain.Role{},
	)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Migration completed successfully!")
}

func createUserStatusEnum(db *gorm.DB) {
	db.Exec(`
	DO $$
	BEGIN
	    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'user_status') THEN
	        CREATE TYPE user_status AS ENUM ('active', 'banned', 'deleted');
	    END IF;
	END$$;
	`)
}