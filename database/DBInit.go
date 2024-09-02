package database

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"users-api/models"
)

var DB *gorm.DB

func DBInit() {
	envVars := LoadEnvVariables()

	url := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=America/New_York",
		envVars.DatabaseHost, envVars.DatabaseUsername, envVars.DatabasePassword,
		envVars.DatabaseName, envVars.DatabasePort)
	fmt.Println("This is the current URL that has been gotten: " + url)
	db, err := gorm.Open(postgres.Open(url), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}
	DB = db
	err = DB.AutoMigrate(&models.User{}, &models.Group{}, &models.Role{})
	if err != nil {
		log.Fatal(err)
		return
	}
	db.Exec("CREATE INDEX IF NOT EXISTS roles_uid_idx ON roles USING BTREE (uid);")
	db.Exec("CREATE INDEX IF NOT EXISTS users_uid_idx ON users USING BTREE (uid);")
	db.Exec("CREATE INDEX IF NOT EXISTS users_role_id_idx ON users USING BTREE (role_id);")
	db.Exec("CREATE INDEX IF NOT EXISTS groups_uid_idx ON groups USING BTREE (uid);")
}
