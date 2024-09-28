package loans

import (
	"context"
	"database/sql"
	"fmt"
	"loan-service/models"
	"loan-service/services/auth"
	"loan-service/utils/errs"
	"loan-service/utils/money"
	"strconv"

	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

// FetchLoans implements models.LoanRepository.
// FetchLoanOpts is for the system to fetch loans associated a certain user according to their role type
func (r *repository) FetchLoans(ctx context.Context, opts models.FetchLoanOpts) ([]models.Loan, error) {
	var results []models.Loan
	query := r.db.WithContext(ctx).Model(&models.Loan{})

	if opts.WithPreloads {
		query = query.Preload("Borrower").
			Preload("Investors").
			Preload("Product").
			Preload("Visitor").
			Preload("Approver").
			Preload("Disburser")
	}

	if len(opts.Status) > 0 {
		query = query.Where("status IN (?)", opts.Status)
	}

	if opts.UserID > 0 && opts.RoleType != "" {
		switch opts.RoleType {
		// Fetch loans that an investor has funded
		case auth.RoleTypeInvestor:
			query = query.Joins("JOIN users ON users.id = ?", opts.UserID).
				Joins("JOIN roles ON roles.id = users.role_id").
				Where("roles.role_type = ?", opts.RoleType)
		// Fetch loans that a field validator has worked on
		case auth.RoleTypeFieldValidator:
			query = query.Where("visitor_id = ? OR disburser_id = ?", opts.UserID, opts.UserID)
		// Fetch loans that a borrower has requested
		case auth.RoleTypeBorrower:
			query = query.Where("borrower_id = ?", opts.UserID)
		default:
		}
	}

	err := query.Find(&results).Error
	if err != nil {
		return nil, err
	}

	return results, nil
}

// FetchLoanByID implements models.LoanRepository.
func (r *repository) FetchLoanByID(ctx context.Context, loanID uint) (*models.Loan, error) {
	var result *models.Loan
	err := r.db.WithContext(ctx).Model(&models.Loan{}).Where("id = ?", loanID).First(&result).Error
	if err != nil {
		return nil, err
	}

	return result, nil
}

// CreateLoan implements models.LoanRepository.
func (r *repository) CreateLoan(ctx context.Context, loan *models.Loan) error {
	err := r.db.WithContext(ctx).Model(&models.Loan{}).Save(loan).Error
	if err != nil {
		return err
	}

	return nil
}

// InvestInLoan implements models.LoanRepository.
func (r *repository) InvestInLoan(ctx context.Context, loan *models.Loan, investor *models.User, amount float64) error {
	txErr := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Check for existing investments
		var existingInvestments []models.Investment
		err := tx.Model(&models.Investment{}).Where("borrower_id = ?", loan.BorrowerID).Find(&existingInvestments).Error
		if err != nil {
			return err
		}

		// Sum up the amounts of existing investments
		var existingInvestmentsAmount float64
		for _, investment := range existingInvestments {
			amountFloat, err := strconv.ParseFloat(investment.Amount, 64)
			if err != nil {
				return err
			}

			existingInvestmentsAmount += amountFloat
		}

		// Check if amount invested will exceed remaining principal
		principalAmountFloat, err := strconv.ParseFloat(loan.PrincipalAmount, 64)
		if err != nil {
			return err
		}

		if amount+existingInvestmentsAmount > principalAmountFloat {
			return ErrInvestmentAmountExceedsPrincipal
		}

		// Create investment
		if err := tx.Model(&models.Investment{}).Create(&models.Investment{
			InvestorID: investor.ID,
			LoanID:     loan.ID,
			Amount:     fmt.Sprintf("%.2f", amount),
		}).Error; err != nil {
			return err
		}

		// Update loan status if completely invested
		if money.NearlyEqual((amount + existingInvestmentsAmount), principalAmountFloat) {
			if err := loan.AdvanceState(models.LoanStatusInvested, "InvestInLoan"); err != nil {
				return err
			}

			if err := tx.Model(&models.Loan{}).Save(loan).Error; err != nil {
				return err
			}
		}

		return nil
	}, &sql.TxOptions{Isolation: sql.LevelSerializable})

	if txErr != nil {
		return errs.Wrap(txErr)
	}

	return nil
}

// GetTotalInvestedAmount implements models.LoanRepository.
func (r *repository) GetTotalInvestedAmount(ctx context.Context, investorID *uint) (float64, error) {
	var results []models.Investment
	err := r.db.WithContext(ctx).Model(&models.Investment{}).Where("investor_id = ?", investorID).Find(&results).Error
	if err != nil {
		return 0, err
	}

	var totalAmount float64
	for _, investment := range results {
		amountFloat, err := strconv.ParseFloat(investment.Amount, 64)
		if err != nil {
			return 0, err
		}

		totalAmount += amountFloat
	}

	return totalAmount, nil
}

// UpdateLoan implements models.LoanRepository.
func (r *repository) UpdateLoan(ctx context.Context, loan *models.Loan) error {
	err := r.db.WithContext(ctx).Model(&models.Loan{}).Save(loan).Error
	if err != nil {
		return err
	}

	return nil
}

func NewLoanRepository(db *gorm.DB) models.LoanRepository {
	return &repository{db}
}
