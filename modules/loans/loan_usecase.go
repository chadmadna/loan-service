package loans

import (
	"context"
	"errors"
	"fmt"
	"io"
	"loan-service/models"
	"loan-service/services/auth"
	"loan-service/services/email"
	"loan-service/services/upload"
	"loan-service/utils/errs"
	"os"
	"strconv"
	"time"

	"github.com/gabriel-vasile/mimetype"
	"github.com/subosito/gozaru"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
)

type usecase struct {
	repo          models.LoanRepository
	userUsecase   models.UserUsecase
	emailService  email.EmailService
	uploadService upload.UploadService
}

// FetchLoanByID implements models.LoanUsecase.
func (u *usecase) FetchLoanByID(ctx context.Context, loanID uint, opts *models.FetchLoanOpts) (*models.Loan, error) {
	loan, err := u.repo.FetchLoanByID(ctx, loanID, opts)
	if err != nil {
		return nil, errs.Wrap(err)
	}

	return loan, nil
}

// FetchLoans implements models.LoanUsecase.
func (u *usecase) FetchLoans(ctx context.Context, opts *models.FetchLoanOpts) ([]models.Loan, error) {
	if opts != nil && !(opts.UserID > 0 && opts.RoleType != "") {
		return nil, errs.Wrap(ErrInvalidParams)
	}

	loans, err := u.repo.FetchLoans(ctx, opts)
	if err != nil {
		return nil, errs.Wrap(err)
	}

	return loans, nil
}

// FetchLoansByUserID implements models.LoanUsecase.
func (u *usecase) FetchLoansByUserID(ctx context.Context, userID uint) ([]models.Loan, error) {
	user, err := u.userUsecase.FetchUserByID(ctx, userID, &models.FetchUserByIDOpts{IncludeBorrowedLoans: true})
	if err != nil {
		return nil, errs.Wrap(err)
	}

	if user == nil {
		return nil, nil
	}

	var loans []models.Loan
	if user.Role.RoleType == auth.RoleTypeBorrower {
		loans = append(loans, user.BorrowedLoans...)
	}

	if user.Role.RoleType == auth.RoleTypeInvestor {
		loans = append(loans, user.InvestedLoans...)
	}

	if user.Role.RoleType == auth.RoleTypeFieldValidator {
		loans, err = u.repo.FetchLoans(ctx, nil)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.Wrap(err)
		}
	}

	return loans, nil
}

// StartLoan implements models.LoanUsecase.
func (u *usecase) StartLoan(ctx context.Context, name string, product *models.Product, borrower *models.User) (*models.Loan, error) {
	if product == nil || borrower == nil {
		return nil, errs.Wrap(ErrInvalidParams)
	}

	existingLoans, err := u.repo.FetchLoans(ctx, &models.FetchLoanOpts{
		UserID:   borrower.ID,
		RoleType: borrower.Role.RoleType,
	})
	if err != nil {
		return nil, errs.Wrap(err)
	}

	if len(existingLoans) > 0 {
		return nil, errs.Wrap(ErrLoanAlreadyExists)
	}

	loan := models.NewLoan(name, product, borrower)
	err = u.repo.CreateLoan(ctx, loan)
	if err != nil {
		return nil, errs.Wrap(err)
	}

	return loan, nil
}

// MarkLoanBorrowerVisited implements models.LoanUsecase.
func (u *usecase) MarkLoanBorrowerVisited(ctx context.Context, loan *models.Loan, visitor *models.User, attachment io.Reader) error {
	if loan == nil || visitor == nil || attachment == nil {
		return errs.Wrap(ErrInvalidParams)
	}

	if loan.VisitorID != nil {
		return errs.Wrap(ErrLoanAlreadyVisited)
	}

	mimeType, err := mimetype.DetectReader(attachment)
	if err != nil {
		return errs.Wrap(err)
	}

	var extension string
	switch mimeType.String() {
	case "image/png":
		extension = "png"
	case "image/jpeg":
		extension = "jpg"
	default:
		// wrong mimetype
		return errs.Wrap(ErrInvalidParams)
	}

	safePath := gozaru.Sanitize(fmt.Sprintf("ProofOfVisit_%s.%s", time.Now().Format(time.RFC3339), extension))

	attachmentPath, err := u.uploadService.UploadFile(
		attachment,
		safePath,
		mimeType.String(),
	)
	if err != nil {
		return errs.Wrap(err)
	}

	loan.Visitor = visitor
	loan.ProofOfVisitAttachmentFile = attachmentPath

	err = u.repo.UpdateLoan(ctx, loan)
	if err != nil {
		return errs.Wrap(err)
	}

	return nil
}

// ApproveLoan implements models.LoanUsecase.
func (u *usecase) ApproveLoan(ctx context.Context, loan *models.Loan, approver *models.User) error {
	loan.Approver = approver
	err := loan.AdvanceState(models.LoanStatusApproved, "ApproveLoan")
	if err != nil {
		return errs.Wrap(err)
	}

	err = u.repo.UpdateLoan(ctx, loan)
	if err != nil {
		return errs.Wrap(err)
	}

	return nil
}

// InvestInLoan implements models.LoanUsecase.
func (u *usecase) InvestInLoan(ctx context.Context, loan *models.Loan, investor *models.User, amount float64) error {
	if loan.Status != models.LoanStatusApproved {
		return ErrLoanNotInvestable
	}

	err := u.repo.InvestInLoan(ctx, loan, investor, amount)
	if err != nil {
		return errs.Wrap(err)
	}

	remainingAmountFloat, _ := strconv.ParseFloat(loan.RemainingAmount, 64)
	if remainingAmountFloat > 0 {
		return nil
	}

	// observer and fan-out pattern
	eg, egCtx := errgroup.WithContext(ctx)
	for _, investor := range loan.Investors {
		func(loan *models.Loan, investor *models.User) {
			eg.Go(func() error {
				// TODO: Generate an actual loan agreement PDF letter for attachment
				file, err := os.Open("public/loan-agreement-letter.pdf")
				if err != nil {
					return errs.Wrap(err)
				}

				err = investor.NotifyEmailLoanFunded(egCtx, u.emailService, loan, file)
				if err != nil {
					return errs.Wrap(err)
				}

				return nil
			})
		}(loan, &investor)
	}

	return eg.Wait()
}

// DisburseLoan implements models.LoanUsecase.
func (u *usecase) DisburseLoan(ctx context.Context, loan *models.Loan, disburser *models.User) error {
	if loan.Status != models.LoanStatusInvested {
		return ErrLoanNotDisbursable
	}

	loan.Disburser = disburser
	err := loan.AdvanceState(models.LoanStatusDisbursed, "DisburseLoan")
	if err != nil {
		return errs.Wrap(err)
	}

	err = u.repo.UpdateLoan(ctx, loan)
	if err != nil {
		return errs.Wrap(err)
	}

	return nil
}

func NewLoanUsecase(
	repo models.LoanRepository,
	userUC models.UserUsecase,
	emailService email.EmailService,
	uploadService upload.UploadService,
) models.LoanUsecase {
	return &usecase{repo, userUC, emailService, uploadService}
}
