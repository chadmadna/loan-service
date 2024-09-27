package models

import (
	"loan-service/services/auth"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name           string `json:"name"`
	Email          string `json:"email"`
	Password       string `json:"password" gorm:"-"`               // do not store in db
	HashedPassword []byte `json:"-" gorm:"column:hashed_password"` // do not return to api
	IsActive       bool   `json:"is_active"`
	RoleID         uint
	Role           Role   `json:"role" gorm:"foreignKey:RoleID"`
	InvestedLoans  []Loan `json:"invested_loans" gorm:"many2many:loans_investors;foreignKey:ID;joinForeignKey:InvestorID;references:ID;joinReferences:LoanID"`
	BorrowedLoans  []Loan `json:"borrowed_loans" gorm:"foreignKey:BorrowerID"`
}

type LoginResponse struct {
	UserID       uint   `json:"user_id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (User) TableName() string {
	return "users"
}

// SetNewPassword set a new hashsed password to user.
func (u *User) SetNewPassword(passwordString string) {
	bcryptPassword, _ := bcrypt.GenerateFromPassword([]byte(passwordString), bcrypt.DefaultCost)
	u.HashedPassword = bcryptPassword
}

type ViewUsersOpt struct {
	RoleType auth.RoleType
	UserID   uint
}

type UserRepository interface {
	CreateUser(user *User) error
	FetchUser(userID uint) (*User, error)
	FetchUsers(allowedRoles []auth.RoleType, allowedLoanIDs []uint) ([]User, error)
	UpdateUser(user *User) error
	FetchRoleByRoleType(roleType auth.RoleType) (*Role, error)
}

type UserUsecase interface {
	Login(email, password string) (LoginResponse, error)
	Logout(email string) error
	ViewUsers(opts ViewUsersOpt) ([]User, error)
	RegisterUser(user *User) error
	UpdateProfile(user *User) error
}
