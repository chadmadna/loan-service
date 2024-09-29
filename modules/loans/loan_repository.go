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
func (r *repository) FetchLoans(ctx context.Context, opts *models.FetchLoanOpts) ([]models.Loan, error) {
	var results []models.Loan
	query := r.db.WithContext(ctx).Model(&models.Loan{}).
		Preload("Borrower").
		Preload("Product")

	if opts != nil && len(opts.Status) > 0 {
		query = query.Where("status IN (?)", opts.Status)
	}

	if opts != nil && opts.UserID > 0 && opts.RoleType != "" {
		query = scopeLoanQuery(query, opts.UserID, opts.RoleType)
	}

	err := query.Find(&results).Error
	if err != nil {
		return nil, err
	}

	return results, nil
}

// FetchLoanByID implements models.LoanRepository.
func (r *repository) FetchLoanByID(ctx context.Context, loanID uint, opts *models.FetchLoanOpts) (*models.Loan, error) {
	var result *models.Loan
	query := r.db.WithContext(ctx).Model(&models.Loan{}).
		Preload("Borrower").
		Preload("Product")

	if opts != nil && len(opts.Status) > 0 {
		query = query.Where("status IN (?)", opts.Status)
	}

	if opts != nil && opts.UserID > 0 && opts.RoleType != "" {
		query = scopeLoanQuery(query, opts.UserID, opts.RoleType)
	}

	err := query.Where("loans.id = ?", loanID).First(&result).Error
	if err != nil {
		return nil, err
	}

	return result, nil
}

// CreateLoan implements models.LoanRepository.
func (r *repository) CreateLoan(ctx context.Context, loan *models.Loan) error {
	err := r.db.WithContext(ctx).Model(&models.Loan{}).Create(loan).Error
	if err != nil {
		return err
	}

	return nil
}

// InvestInLoan implements models.LoanRepository.
func (r *repository) InvestInLoan(ctx context.Context, loan *models.Loan, investor *models.User, amount float64) error {
	txErr := r.db.Debug().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Check for existing investments
		var existingInvestments []models.Investment
		err := tx.Model(&models.Investment{}).Where("loan_id = ?", loan.ID).Find(&existingInvestments).Error
		if err != nil {
			return errs.Wrap(err)
		}

		// Sum up the amounts of existing investments
		var existingInvestmentsAmount float64
		for _, investment := range existingInvestments {
			amountFloat, err := strconv.ParseFloat(investment.Amount, 64)
			if err != nil {
				return errs.Wrap(err)
			}

			existingInvestmentsAmount += amountFloat
		}

		// Check if amount invested will exceed remaining principal
		principalAmountFloat, err := strconv.ParseFloat(loan.PrincipalAmount, 64)
		if err != nil {
			return errs.Wrap(err)
		}

		if amount+existingInvestmentsAmount > principalAmountFloat {
			return errs.Wrap(ErrInvestmentAmountExceedsPrincipal)
		}

		// Create investment
		if err := tx.Model(&models.Investment{}).Create(&models.Investment{
			InvestorID: investor.ID,
			LoanID:     loan.ID,
			Amount:     fmt.Sprintf("%.2f", amount),
		}).Error; err != nil {
			return errs.Wrap(err)
		}

		// Update loan status if completely invested
		if money.NearlyEqual((amount + existingInvestmentsAmount), principalAmountFloat) {
			if err := loan.AdvanceState(models.LoanStatusInvested, "InvestInLoan"); err != nil {
				return errs.Wrap(err)
			}
		}

		// Update remaining amount
		remainingAmountFloat, err := strconv.ParseFloat(loan.RemainingAmount, 64)
		if err != nil {
			return errs.Wrap(err)
		}
		loan.RemainingAmount = fmt.Sprintf("%.2f", remainingAmountFloat-amount)

		// Save loan, update only selected columns
		if err := tx.Model(loan).Updates(map[string]any{
			"remaining_amount": loan.RemainingAmount,
			"status":           loan.Status,
		}).Error; err != nil {
			return errs.Wrap(err)
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
	err := r.db.Debug().WithContext(ctx).Model(loan).Updates(map[string]any{
		"name":                           loan.Name,
		"status":                         loan.Status,
		"remaining_amount":               loan.RemainingAmount,
		"visitor_id":                     loan.VisitorID,
		"approver_id":                    loan.ApproverID,
		"disburser_id":                   loan.DisburserID,
		"proof_of_visit_attachment_file": loan.ProofOfVisitAttachmentFile,
	}).Error
	if err != nil {
		return err
	}

	return nil
}

func NewLoanRepository(db *gorm.DB) models.LoanRepository {
	return &repository{db}
}

func scopeLoanQuery(query *gorm.DB, userID uint, roleType auth.RoleType) *gorm.DB {
	switch roleType {
	// Fetch loans that an investor has funded
	case auth.RoleTypeInvestor:
		query = query.Preload("Visitor").
			Preload("Approver").
			Preload("Investors").
			Preload("Disburser").
			Joins("LEFT JOIN investments ON investments.investor_id = ?", userID).
			Where("status != ?", models.LoanStatusProposed)
	// Fetch loans that a field validator has worked on
	case auth.RoleTypeFieldValidator:
		query = query.Preload("Visitor").
			Preload("Approver").
			Preload("Investors").
			Preload("Disburser").
			Where("visitor_id = ? OR disburser_id = ?", userID, userID).
			Or("status = ?", models.LoanStatusProposed)
	// Fetch loans that a borrower has requested
	case auth.RoleTypeBorrower:
		query = query.Preload("Visitor").
			Preload("Disburser").
			Where("borrower_id = ?", userID)
	default:
		query = query.Preload("Visitor").
			Preload("Approver").
			Preload("Investors").
			Preload("Disburser")
	}

	return query
}
