package models

import (
	"context"
	"fmt"
	"loan-service/services/auth"
	"loan-service/services/email"
	"loan-service/utils/errs"
	"os"
	"time"

	"github.com/sendgrid/sendgrid-go/helpers/mail"
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
	Role           Role         `json:"role" gorm:"foreignKey:RoleID"`
	InvestedLoans  []Loan       `json:"invested_loans" gorm:"many2many:investments;foreignKey:ID;joinForeignKey:InvestorID;references:ID;joinReferences:LoanID"`
	Investments    []Investment `json:"investments" gorm:"foreignKey:InvestorID"`
	BorrowedLoans  []Loan       `json:"borrowed_loans" gorm:"foreignKey:BorrowerID"`
}

type LoginResponse struct {
	UserID uint   `json:"user_id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
}

func (User) TableName() string {
	return "users"
}

// SetNewPassword set a new hashsed password to user.
func (u *User) SetNewPassword(passwordString string) {
	bcryptPassword, _ := bcrypt.GenerateFromPassword([]byte(passwordString), bcrypt.DefaultCost)
	u.HashedPassword = bcryptPassword
}

// Notify investor by email when loan is fully invested
func (u *User) NotifyEmailLoanFunded(ctx context.Context, emailService email.EmailService, loan *Loan, totalInvested float64, attachment *os.File) error {
	subject := "Loan has been fully funded!"
	body := fmt.Sprintf(`
			<h1>A loan you financed has been fully funded!</h1>
			<p>You've invested Rp %.2f into "%s", a loan requested by %s.</p>
			<p>We're proud to let you know that the loan has been fully funded. Thank you for your contribution!.</p>
			<p>They will now receive the full amount of Rp %s once it has been disbursed by our staff.</p>
			<strong>You will earn Rp %s total interest, that's %s%% return on investment.</strong>
			<p>Attached is the loan agreement letter to sign.</p>
			<p>Thank you for trusting LoanService.io!</p>
		`,
		totalInvested, loan.Name, loan.Borrower.Name, loan.PrincipalAmount, loan.TotalInterest, loan.ROI,
	)

	err := emailService.SendMail(
		ctx,
		subject,
		body,
		mail.Email{Name: emailService.DefaultSenderName(), Address: emailService.DefaultSenderAddress()},
		mail.Email{Name: u.Name, Address: u.Email},
		&email.AttachmentOpts{
			File:        attachment,
			ContentType: email.AttachmentTypePDF,
			Filename: fmt.Sprintf("Loan_Agreement_Letter-%s-%s-%s",
				loan.Name, loan.Borrower.Name, loan.UpdatedAt.Format(time.RFC3339),
			),
		},
	)
	if err != nil {
		return errs.Wrap(err)
	}

	return nil
}

type ViewUsersOpt struct {
	RoleType auth.RoleType
	UserID   uint
}

type UserRepository interface {
	CreateUser(ctx context.Context, user *User) error
	FetchUserByID(ctx context.Context, userID uint) (*User, error)
	FetchUsers(ctx context.Context, allowedRoles []auth.RoleType, allowedLoanIDs []uint) ([]User, error)
	UpdateUser(ctx context.Context, user *User) error
	FetchRoleByRoleType(ctx context.Context, roleType auth.RoleType) (*Role, error)
	FetchUserByEmail(ctx context.Context, email string) (*User, error)
}

type UserUsecase interface {
	Login(ctx context.Context, email, password string) (LoginResponse, string, string, error)
	ViewUsers(ctx context.Context, opts ViewUsersOpt) ([]User, error)
	FetchUserByID(ctx context.Context, userID uint) (*User, error)
	RegisterUser(ctx context.Context, user *User) error
	UpdateProfile(ctx context.Context, user *User) error
}
