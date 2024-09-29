package main

import (
	"fmt"
	"loan-service/config"
	database "loan-service/database"
	"loan-service/models"
	"loan-service/services/auth"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func main() {
	fmt.Println("Seeding database..")

	// Load env variables
	err := config.LoadFromEnv()
	if err != nil {
		panic(err)
	}

	db, err := database.GetDB()
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(
		&models.Role{},
		&models.User{},
		&models.Product{},
		&models.Loan{},
		&models.Investment{},
	)
	if err != nil {
		panic(err)
	}

	products := []models.Product{
		{
			Model:           gorm.Model{ID: 1},
			Name:            "Dana Fleksibel",
			PrincipalAmount: "20000000.0",
			InterestRate:    0.1,
			Term:            models.TermLength3Month,
		},
		{
			Model:           gorm.Model{ID: 2},
			Name:            "Dana Sejahtera",
			PrincipalAmount: "10000000.0",
			InterestRate:    0.08,
			Term:            models.TermLength6Month,
		},
		{
			Model:           gorm.Model{ID: 3},
			Name:            "Dana Usaha",
			PrincipalAmount: "100000000.0",
			InterestRate:    0.06942,
			Term:            models.TermLength12Month,
		},
	}

	if err := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&products).Error; err != nil {
		panic(fmt.Errorf("cannot bulk insert products: %v", err))
	}

	roles := []models.Role{
		{Model: gorm.Model{ID: 1}, Name: "Superuser", RoleType: auth.RoleTypeSuperuser},
		{Model: gorm.Model{ID: 2}, Name: "Staff", RoleType: auth.RoleTypeStaff},
		{Model: gorm.Model{ID: 3}, Name: "Field Validator", RoleType: auth.RoleTypeFieldValidator},
		{Model: gorm.Model{ID: 4}, Name: "Investor", RoleType: auth.RoleTypeInvestor},
		{Model: gorm.Model{ID: 5}, Name: "Borrower", RoleType: auth.RoleTypeBorrower},
	}

	if err := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&roles).Error; err != nil {
		panic(fmt.Errorf("cannot bulk insert roles: %v", err))
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
			Name:     "George Soros",
			Email:    "georgesoros@sorosfund.com",
			Password: "george@investor",
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
		{
			Name:     "Nuhut Bingsar",
			Email:    "nuhutbingsar@indonesia.go.id",
			Password: "nuhut@borrower",
			IsActive: true,
			RoleID:   5,
		},
		{
			Name:     "Gibro Rakbro",
			Email:    "fufufafa@indonesia.go.id",
			Password: "fufufafa@borrower",
			IsActive: true,
			RoleID:   5,
		},
	}

	for i := range users {
		users[i].SetNewPassword(users[i].Password)
	}

	if err := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&users).Error; err != nil {
		panic(fmt.Errorf("cannot bulk insert users: %v", err))
	}

	// NOTE: Uncomment if you'd like to test with pre-existing loans and investments
	// loans := []models.Loan{
	// 	{
	// 		Model:                      gorm.Model{ID: 1},
	// 		Name:                       "Foreign investment for ada deh",
	// 		Status:                     models.LoanStatusApproved,
	// 		BorrowerID:                 7,
	// 		ProductID:                  3,
	// 		PrincipalAmount:            "100000000.0",
	// 		RemainingAmount:            "15000000.0",
	// 		InterestRate:               0.06942,
	// 		TotalInterest:              "6942000.0",
	// 		ROI:                        "6.94",
	// 		LoanTerm:                   int(models.TermLength12Month),
	// 		VisitorID:                  ptr.NewUintPtr(3),
	// 		ApproverID:                 ptr.NewUintPtr(2),
	// 		ProofOfVisitAttachmentFile: "https://picsum.photos/seed/loanservice/900/1600",
	// 	},
	// 	{
	// 		Model:           gorm.Model{ID: 2},
	// 		Name:            "Gofundme for my companies",
	// 		Status:          models.LoanStatusProposed,
	// 		BorrowerID:      8,
	// 		ProductID:       3,
	// 		PrincipalAmount: "100000000.0",
	// 		RemainingAmount: "100000000.0",
	// 		InterestRate:    0.06942,
	// 		TotalInterest:   "6942000.0",
	// 		ROI:             "6.94",
	// 		LoanTerm:        int(models.TermLength12Month),
	// 	},
	// }

	// if err := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&loans).Error; err != nil {
	// 	panic(fmt.Errorf("cannot bulk insert loans: %v", err))
	// }

	// investments := []models.Investment{
	// 	{
	// 		InvestorID: 6,
	// 		LoanID:     1,
	// 		Amount:     "60000000",
	// 	},
	// 	{
	// 		InvestorID: 5,
	// 		LoanID:     1,
	// 		Amount:     "25000000",
	// 	},
	// }

	// if err := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&investments).Error; err != nil {
	// 	panic(fmt.Errorf("cannot bulk insert investments: %v", err))
	// }
}
