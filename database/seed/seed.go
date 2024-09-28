package main

import (
	"loan-service/config"
	database "loan-service/database"
	"loan-service/models"
	"loan-service/services/auth"

	"github.com/apsystole/log"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func main() {
	log.Println("Seeding database..")

	// Load env variables
	err := config.LoadFromEnv()
	if err != nil {
		panic(err)
	}

	db, err := database.GetDB()
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(
		&models.Role{},
		&models.User{},
		// &models.Product{},
		&models.Loan{},
	)

	roles := []models.Role{
		{Model: gorm.Model{ID: 1}, Name: "Superuser", RoleType: auth.RoleTypeSuperuser},
		{Model: gorm.Model{ID: 2}, Name: "Staff", RoleType: auth.RoleTypeStaff},
		{Model: gorm.Model{ID: 3}, Name: "Field Validator", RoleType: auth.RoleTypeFieldValidator},
		{Model: gorm.Model{ID: 4}, Name: "Investor", RoleType: auth.RoleTypeInvestor},
		{Model: gorm.Model{ID: 5}, Name: "Borrower", RoleType: auth.RoleTypeBorrower},
	}

	if err := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&roles).Error; err != nil {
		log.Fatalf("cannot bulk insert roles: %v", err)
	}

	users := []models.User{
		{
			Name:     "Angela Merkel",
			Email:    "admin@loanservice.io",
			Password: "@superuser",
			IsActive: true,
			RoleID:   1,
		},
		{
			Name:     "Emmanuel Macron",
			Email:    "staff@loanservice.io",
			Password: "@staff",
			IsActive: true,
			RoleID:   2,
		},
		{
			Name:     "Silvio Berlusconi",
			Email:    "field.validator@loanservice.io",
			Password: "@field.validator",
			IsActive: true,
			RoleID:   3,
		},
		{
			Name:     "Larry Fink",
			Email:    "larryfink@blackrock.com",
			Password: "larry@investor",
			IsActive: true,
			RoleID:   4,
		},
		{
			Name:     "Luke Sarsfield",
			Email:    "lukesarsfield@goldmansachs.com",
			Password: "luke@investor",
			IsActive: true,
			RoleID:   4,
		},
		{
			Name:     "Zulhas Hasan",
			Email:    "zulhashasan@indonesia.go.id",
			Password: "zulhas@borrower",
			IsActive: true,
			RoleID:   5,
		},
	}

	for i := range users {
		users[i].SetNewPassword(users[i].Password)
	}

	if err := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&users).Error; err != nil {
		log.Fatalf("cannot bulk insert users: %v", err)
	}
}
