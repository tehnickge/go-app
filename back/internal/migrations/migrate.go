package migrations

import (
	"log"

	"backend-app/internal/domain"

	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) {
	// включаем расширение uuid-ossp
	enableUUIDExtension(db)

	// создаём ENUM для статусов пользователей
	createUserStatusEnum(db)

	// Автоматическая миграция всех таблиц
	if err := db.AutoMigrate(
		&domain.User{},
		&domain.Group{},
		&domain.Role{},
		&domain.UsersGroups{},
		&domain.UsersRoles{},
	); err != nil {
		log.Fatal("Migration failed:", err)
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

func enableUUIDExtension(db *gorm.DB) {
	db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`)
}