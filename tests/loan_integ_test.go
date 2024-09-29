package integration

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"loan-service/app"
	"loan-service/models"
	loanModule "loan-service/modules/loans"
	_loanHandlers "loan-service/modules/loans/handlers"
	"loan-service/services/auth"
	_emailMock "loan-service/services/email/mocks"
	_uploadMock "loan-service/services/upload/mocks"
	"loan-service/utils/jsonutil"
	"loan-service/utils/ptr"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/samber/do"
	_assert "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type loanIntegrationTestSuite struct {
	suite.Suite
	db                        *gorm.DB
	rest                      *echo.Echo
	borrowerLoanHandler       *_loanHandlers.BorrowerLoanHandler
	fieldValidatorLoanHandler *_loanHandlers.FieldValidatorLoanHandler
	staffLoanHandler          *_loanHandlers.StaffLoanHandler
	investorLoanHandler       *_loanHandlers.InvestorLoanHandler
	models                    []interface{}
	emailSvc                  *_emailMock.EmailService
	uploadSvc                 *_uploadMock.UploadService
	injector                  *do.Injector
	imageFixture              *os.File
}

func TestIntegrationLoan(t *testing.T) {
	suite.Run(t, new(loanIntegrationTestSuite))
}

func (s *loanIntegrationTestSuite) SetupSuite() {
	var err error
	s.db, err = InitDB()
	if err != nil {
		panic(err)
	}

	s.rest = SetupEcho()

	s.emailSvc = _emailMock.NewEmailService(s.T())
	s.uploadSvc = _uploadMock.NewUploadService(s.T())

	s.injector = app.SetupInjections(s.db, s.rest, s.emailSvc, s.uploadSvc)

	s.borrowerLoanHandler = &_loanHandlers.BorrowerLoanHandler{
		Usecase:        do.MustInvoke[models.LoanUsecase](s.injector),
		UserUsecase:    do.MustInvoke[models.UserUsecase](s.injector),
		ProductUsecase: do.MustInvoke[models.ProductUsecase](s.injector),
		CommonHandler: _loanHandlers.NewCommonLoanHandler(
			do.MustInvoke[models.LoanUsecase](s.injector),
			do.MustInvoke[models.UserUsecase](s.injector),
		),
	}

	s.fieldValidatorLoanHandler = &_loanHandlers.FieldValidatorLoanHandler{
		Usecase:     do.MustInvoke[models.LoanUsecase](s.injector),
		UserUsecase: do.MustInvoke[models.UserUsecase](s.injector),
		CommonHandler: _loanHandlers.NewCommonLoanHandler(
			do.MustInvoke[models.LoanUsecase](s.injector),
			do.MustInvoke[models.UserUsecase](s.injector),
		),
	}

	s.staffLoanHandler = &_loanHandlers.StaffLoanHandler{
		Usecase:     do.MustInvoke[models.LoanUsecase](s.injector),
		UserUsecase: do.MustInvoke[models.UserUsecase](s.injector),
		CommonHandler: _loanHandlers.NewCommonLoanHandler(
			do.MustInvoke[models.LoanUsecase](s.injector),
			do.MustInvoke[models.UserUsecase](s.injector),
		),
	}

	s.investorLoanHandler = &_loanHandlers.InvestorLoanHandler{
		Usecase:     do.MustInvoke[models.LoanUsecase](s.injector),
		UserUsecase: do.MustInvoke[models.UserUsecase](s.injector),
		CommonHandler: _loanHandlers.NewCommonLoanHandler(
			do.MustInvoke[models.LoanUsecase](s.injector),
			do.MustInvoke[models.UserUsecase](s.injector),
		),
	}

	s.models = []any{
		&models.Role{},
		&models.User{},
		&models.Product{},
		&models.Loan{},
		&models.Investment{},
	}

	s.imageFixture, err = os.Open("fixtures/example-attachment.jpg")
	if err != nil {
		panic(err)
	}
}

func (s *loanIntegrationTestSuite) TestIntegration_StartLoan() {
	tests := []struct {
		name    string
		userID  uint
		reqStr  string
		want    string
		wantErr error
	}{
		{
			name:   "returns results given valid request and no existing loan",
			userID: 6,
			reqStr: `
				{
					"name": "Beli furnitur",
    			"product_id": 2
				}
			`,
			want: `
				{
					"data": {
						"id": 5,
						"name": "Beli furnitur",
						"status": "proposed",
						"principal_amount": "Rp10.000.000,00",
						"remaining_amount": "Rp10.000.000,00",
						"interest_rate": "8%",
						"total_interest": "Rp400.000,00",
						"loan_term": "6 months"
					}
				}
			`,
			wantErr: nil,
		},
		{
			name:   "throws error given valid request and an existing loan",
			userID: 7,
			reqStr: `
				{
					"name": "Beli furnitur",
    			"product_id": 2
				}
				`,
			wantErr: loanModule.ErrLoanAlreadyExists,
		},
		{
			name:   "throws error given invalid request",
			userID: 8,
			reqStr: `
				{
					"name": "",
				}
			`,
			wantErr: errors.New("invalid request parameters"),
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			assert := _assert.New(s.T())

			// Build request and its context
			bodyReader := strings.NewReader(tt.reqStr)
			req := httptest.NewRequest(http.MethodPost, "/loans", bodyReader)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ctx := s.rest.NewContext(req, rec)
			ctx.Set(auth.AuthClaimsCtxKey, auth.AuthClaims{
				UserID: uint(tt.userID),
			})

			// Do test and assert
			want, _ := jsonutil.Compact(tt.want)
			err := s.borrowerLoanHandler.StartLoan(ctx)
			got := strings.TrimSpace(rec.Body.String())

			if tt.wantErr == nil {
				assert.NoError(err)
				assert.Equal(http.StatusOK, rec.Code)
				assert.Equal(want, got)
			} else {
				assert.Contains(got, tt.wantErr.Error())
			}
		})
	}
}

func (s *loanIntegrationTestSuite) TestIntegration_MarkLoanBorrowerVisited() {
	tests := []struct {
		name    string
		userID  uint
		params  map[string]string
		wantErr error
	}{
		{
			name:   "returns results given valid request and existing proposed loan",
			userID: 7,
			params: map[string]string{
				"loan_id": "1",
			},
			wantErr: nil,
		},
		{
			name:   "throws error given valid request and already visited loan",
			userID: 9,
			params: map[string]string{
				"loan_id": "3",
			},
			wantErr: loanModule.ErrLoanAlreadyVisited,
		},
		{
			name:    "throws error given invalid request",
			userID:  8,
			params:  nil,
			wantErr: errors.New("invalid request parameters"),
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			assert := _assert.New(s.T())

			// Build request and its context
			body := new(bytes.Buffer)
			bodyWriter := multipart.NewWriter(body)

			file, err := bodyWriter.CreateFormFile("attachment", "example-attachment.jpg")
			assert.NoError(err)

			_, err = io.Copy(file, s.imageFixture)
			assert.NoError(err)
			assert.NoError(bodyWriter.Close())

			req := httptest.NewRequest(http.MethodPost, "/loans/:loan_id/visit", body)
			req.Header.Set(echo.HeaderContentType, bodyWriter.FormDataContentType())
			rec := httptest.NewRecorder()
			ctx := s.rest.NewContext(req, rec)
			ctx.Set(auth.AuthClaimsCtxKey, auth.AuthClaims{
				UserID: uint(tt.userID),
			})

			for k, v := range tt.params {
				ctx.SetParamNames(k)
				ctx.SetParamValues(v)
			}

			// Mock services
			s.uploadSvc.On("UploadFile", mock.Anything, mock.AnythingOfType("string"), "image/jpeg").
				Return("tmp/attachment-path.jpg", nil)

			// Do test and assert
			err = s.fieldValidatorLoanHandler.MarkLoanBorrowerVisited(ctx)
			got := strings.TrimSpace(rec.Body.String())

			if tt.wantErr == nil {
				assert.NoError(err)
				assert.Equal(http.StatusOK, rec.Code)
			} else {
				assert.Contains(got, tt.wantErr.Error())
			}
		})
	}
}

func (s *loanIntegrationTestSuite) TestIntegration_ApproveLoan() {
	tests := []struct {
		name    string
		userID  uint
		params  map[string]string
		wantErr error
	}{
		{
			name:   "returns results given valid request and existing visited loan",
			userID: 9,
			params: map[string]string{
				"loan_id": "2",
			},
			wantErr: nil,
		},
		{
			name:   "throws error given valid request and not visited loan",
			userID: 6,
			params: map[string]string{
				"loan_id": "1",
			},
			wantErr: models.NewNextStateError(models.LoanStatusProposed, models.LoanStatusApproved, "ApproveLoan"),
		},
		{
			name:    "throws error given invalid request",
			userID:  7,
			params:  nil,
			wantErr: errors.New("invalid request parameters"),
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			assert := _assert.New(s.T())

			// Build request and its context
			req := httptest.NewRequest(http.MethodPatch, "/loans/:loan_id/approve", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ctx := s.rest.NewContext(req, rec)
			ctx.Set(auth.AuthClaimsCtxKey, auth.AuthClaims{
				UserID: uint(tt.userID),
			})

			for k, v := range tt.params {
				ctx.SetParamNames(k)
				ctx.SetParamValues(v)
			}

			// Do test and assert
			err := s.staffLoanHandler.ApproveLoan(ctx)
			got := strings.TrimSpace(rec.Body.String())

			if tt.wantErr == nil {
				assert.NoError(err)
				assert.Equal(http.StatusOK, rec.Code)
			} else {
				assert.Contains(got, tt.wantErr.Error())
			}
		})
	}
}

func (s *loanIntegrationTestSuite) TestIntegration_InvestInLoan() {
	tests := []struct {
		name    string
		userID  uint
		params  map[string]string
		reqStr  string
		want    string
		wantErr error
	}{
		{
			name:   "returns results given valid request and amount, and approved loan",
			userID: 4,
			params: map[string]string{
				"loan_id": "3",
			},
			reqStr: `
				{
					"amount": 10000000
				}
			`,
			wantErr: nil,
		},
		{
			name:   "throws error given too large amount in request and approved loan",
			userID: 4,
			params: map[string]string{
				"loan_id": "3",
			},
			reqStr: `
				{
					"amount": 30000000
				}
			`,
			wantErr: loanModule.ErrInvestmentAmountExceedsPrincipal,
		},
		{
			name:   "throws error given valid request and an unapproved loan",
			userID: 4,
			params: map[string]string{
				"loan_id": "1",
			},
			reqStr: `
				{
					"amount": 15000000
				}
			`,
			wantErr: loanModule.ErrLoanNotInvestable,
		},
		{
			name:   "throws error given invalid request",
			userID: 8,
			reqStr: `
				{
					"name": "",
				}
			`,
			wantErr: errors.New("invalid request parameters"),
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			assert := _assert.New(s.T())

			// Build request and its context
			bodyReader := strings.NewReader(tt.reqStr)
			req := httptest.NewRequest(http.MethodPost, "/loans/:loan_id/invest", bodyReader)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ctx := s.rest.NewContext(req, rec)
			ctx.Set(auth.AuthClaimsCtxKey, auth.AuthClaims{
				UserID: uint(tt.userID),
			})

			for k, v := range tt.params {
				ctx.SetParamNames(k)
				ctx.SetParamValues(v)
			}

			// Do test and assert
			err := s.investorLoanHandler.InvestInLoan(ctx)
			got := strings.TrimSpace(rec.Body.String())

			if tt.wantErr == nil {
				assert.NoError(err)
				assert.Equal(http.StatusCreated, rec.Code)
			} else {
				assert.Contains(got, tt.wantErr.Error())
			}
		})
	}
}

func (s *loanIntegrationTestSuite) TestIntegration_DisburseLoan() {
	tests := []struct {
		name    string
		userID  uint
		params  map[string]string
		wantErr error
	}{
		{
			name:   "returns results given valid request and existing invested loan",
			userID: 2,
			params: map[string]string{
				"loan_id": "4",
			},
			wantErr: nil,
		},
		{
			name:   "throws error given valid request and not invested loan",
			userID: 2,
			params: map[string]string{
				"loan_id": "3",
			},
			wantErr: loanModule.ErrLoanNotDisbursable,
		},
		{
			name:    "throws error given invalid request",
			userID:  2,
			params:  nil,
			wantErr: errors.New("invalid request parameters"),
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			assert := _assert.New(s.T())

			// Build request and its context
			req := httptest.NewRequest(http.MethodPatch, "/loans/:loan_id/disburse", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ctx := s.rest.NewContext(req, rec)
			ctx.Set(auth.AuthClaimsCtxKey, auth.AuthClaims{
				UserID: uint(tt.userID),
			})

			for k, v := range tt.params {
				ctx.SetParamNames(k)
				ctx.SetParamValues(v)
			}

			// Do test and assert
			err := s.fieldValidatorLoanHandler.DisburseLoan(ctx)
			got := strings.TrimSpace(rec.Body.String())

			if tt.wantErr == nil {
				assert.NoError(err)
				assert.Equal(http.StatusOK, rec.Code)
			} else {
				assert.Contains(got, tt.wantErr.Error())
			}
		})
	}
}

func (s *loanIntegrationTestSuite) SeedData() {
	products := []models.Product{
		{
			Name:            "Dana Fleksibel",
			PrincipalAmount: "5000000.0",
			InterestRate:    0.1,
			Term:            models.TermLength3Month,
		},
		{
			Name:            "Dana Sejahtera",
			PrincipalAmount: "10000000.0",
			InterestRate:    0.08,
			Term:            models.TermLength6Month,
		},
		{
			Name:            "Dana Usaha",
			PrincipalAmount: "100000000.0",
			InterestRate:    0.06942,
			Term:            models.TermLength12Month,
		},
	}

	if err := s.db.Clauses(clause.OnConflict{DoNothing: true}).Create(&products).Error; err != nil {
		panic(fmt.Errorf("cannot bulk insert products: %v", err))
	}

	roles := []models.Role{
		{Name: "Superuser", RoleType: auth.RoleTypeSuperuser},
		{Name: "Staff", RoleType: auth.RoleTypeStaff},
		{Name: "Field Validator", RoleType: auth.RoleTypeFieldValidator},
		{Name: "Investor", RoleType: auth.RoleTypeInvestor},
		{Name: "Borrower", RoleType: auth.RoleTypeBorrower},
	}

	if err := s.db.Clauses(clause.OnConflict{DoNothing: true}).Create(&roles).Error; err != nil {
		panic(fmt.Errorf("cannot bulk insert roles: %v", err))
	}

	users := []models.User{
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
		{
			Name:     "Rayfe Hamid",
			Email:    "rayfe@indonesia.go.id",
			Password: "rayfe@borrower",
			IsActive: true,
			RoleID:   5,
		},
		{
			Name:     "Irsyad Nabil",
			Email:    "chadmadna@gmail.com",
			Password: "irsyad@borrower",
			IsActive: true,
			RoleID:   5,
		},
	}

	for i := range users {
		users[i].SetNewPassword(users[i].Password)
	}

	if err := s.db.Clauses(clause.OnConflict{DoNothing: true}).Create(&users).Error; err != nil {
		panic(fmt.Errorf("cannot bulk insert users: %v", err))
	}

	loans := []models.Loan{
		{
			Name:            "Crowdfunding bayar kosan sama cicilan",
			Status:          models.LoanStatusProposed,
			BorrowerID:      7,
			ProductID:       1,
			PrincipalAmount: "5000000.0",
			RemainingAmount: "5000000.0",
			InterestRate:    0.1,
			TotalInterest:   "500000.0",
			ROI:             "1",
			LoanTerm:        int(models.TermLength3Month),
		},
		{
			Name:                       "Pinjem dulu seratus",
			Status:                     models.LoanStatusProposed,
			BorrowerID:                 8,
			ProductID:                  3,
			PrincipalAmount:            "100000000.0",
			RemainingAmount:            "15000000.0",
			InterestRate:               0.06942,
			TotalInterest:              "6942000.0",
			ROI:                        "6.94",
			LoanTerm:                   int(models.TermLength12Month),
			VisitorID:                  ptr.NewUintPtr(2),
			ProofOfVisitAttachmentFile: "https://picsum.photos/seed/loanservice/900/1600",
		},
		{
			Name:                       "Beli gelar",
			Status:                     models.LoanStatusApproved,
			BorrowerID:                 9,
			ProductID:                  3,
			PrincipalAmount:            "100000000.0",
			RemainingAmount:            "15000000.0",
			InterestRate:               0.06942,
			TotalInterest:              "6942000.0",
			ROI:                        "6.94",
			LoanTerm:                   int(models.TermLength12Month),
			VisitorID:                  ptr.NewUintPtr(2),
			ProofOfVisitAttachmentFile: "https://picsum.photos/seed/loanservice/900/1600",
			ApproverID:                 ptr.NewUintPtr(1),
		},
		{
			Name:                       "Biaya rekaman album baru",
			Status:                     models.LoanStatusInvested,
			BorrowerID:                 10,
			ProductID:                  2,
			PrincipalAmount:            "10000000.0",
			RemainingAmount:            "10000000.0",
			InterestRate:               0.08,
			TotalInterest:              "800000.0",
			ROI:                        "8",
			LoanTerm:                   int(models.TermLength6Month),
			VisitorID:                  ptr.NewUintPtr(2),
			ProofOfVisitAttachmentFile: "https://picsum.photos/seed/loanservice/900/1600",
			ApproverID:                 ptr.NewUintPtr(1),
		},
	}

	if err := s.db.Clauses(clause.OnConflict{DoNothing: true}).Create(&loans).Error; err != nil {
		panic(fmt.Errorf("cannot bulk insert loans: %v", err))
	}

	investments := []models.Investment{
		{
			InvestorID: 6,
			LoanID:     3,
			Amount:     "60000000",
		},
		{
			InvestorID: 5,
			LoanID:     3,
			Amount:     "25000000",
		},
		{
			InvestorID: 5,
			LoanID:     4,
			Amount:     "10000000",
		},
	}

	if err := s.db.Clauses(clause.OnConflict{DoNothing: true}).Create(&investments).Error; err != nil {
		panic(fmt.Errorf("cannot bulk insert investments: %v", err))
	}
}

func (s *loanIntegrationTestSuite) SetupTest() {
	AutoMigrate(s.db, s.models...)
	s.SeedData()
}

func (s *loanIntegrationTestSuite) TearDownTest() {
	for _, model := range s.models {
		err := s.db.Migrator().DropTable(model)
		if err != nil {
			panic(err)
		}
	}
}

func (s *loanIntegrationTestSuite) TearDownSuite() {
	s.imageFixture.Close()
}
