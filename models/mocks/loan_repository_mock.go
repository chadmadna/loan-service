// Code generated by mockery v2.32.4. DO NOT EDIT.

package mocks

import (
	context "context"
	models "loan-service/models"

	mock "github.com/stretchr/testify/mock"
)

// LoanRepository is an autogenerated mock type for the LoanRepository type
type LoanRepository struct {
	mock.Mock
}

// CreateLoan provides a mock function with given fields: ctx, loan
func (_m *LoanRepository) CreateLoan(ctx context.Context, loan *models.Loan) error {
	ret := _m.Called(ctx, loan)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *models.Loan) error); ok {
		r0 = rf(ctx, loan)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FetchLoanByID provides a mock function with given fields: ctx, loanID, opts
func (_m *LoanRepository) FetchLoanByID(ctx context.Context, loanID uint, opts *models.FetchLoanOpts) (*models.Loan, error) {
	ret := _m.Called(ctx, loanID, opts)

	var r0 *models.Loan
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint, *models.FetchLoanOpts) (*models.Loan, error)); ok {
		return rf(ctx, loanID, opts)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint, *models.FetchLoanOpts) *models.Loan); ok {
		r0 = rf(ctx, loanID, opts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Loan)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint, *models.FetchLoanOpts) error); ok {
		r1 = rf(ctx, loanID, opts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FetchLoans provides a mock function with given fields: ctx, opts
func (_m *LoanRepository) FetchLoans(ctx context.Context, opts *models.FetchLoanOpts) ([]models.Loan, error) {
	ret := _m.Called(ctx, opts)

	var r0 []models.Loan
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *models.FetchLoanOpts) ([]models.Loan, error)); ok {
		return rf(ctx, opts)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *models.FetchLoanOpts) []models.Loan); ok {
		r0 = rf(ctx, opts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.Loan)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *models.FetchLoanOpts) error); ok {
		r1 = rf(ctx, opts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetTotalInvestedAmount provides a mock function with given fields: ctx, investorID
func (_m *LoanRepository) GetTotalInvestedAmount(ctx context.Context, investorID *uint) (float64, error) {
	ret := _m.Called(ctx, investorID)

	var r0 float64
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *uint) (float64, error)); ok {
		return rf(ctx, investorID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *uint) float64); ok {
		r0 = rf(ctx, investorID)
	} else {
		r0 = ret.Get(0).(float64)
	}

	if rf, ok := ret.Get(1).(func(context.Context, *uint) error); ok {
		r1 = rf(ctx, investorID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// InvestInLoan provides a mock function with given fields: ctx, loan, investor, amount
func (_m *LoanRepository) InvestInLoan(ctx context.Context, loan *models.Loan, investor *models.User, amount float64) error {
	ret := _m.Called(ctx, loan, investor, amount)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *models.Loan, *models.User, float64) error); ok {
		r0 = rf(ctx, loan, investor, amount)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateLoan provides a mock function with given fields: ctx, loan
func (_m *LoanRepository) UpdateLoan(ctx context.Context, loan *models.Loan) error {
	ret := _m.Called(ctx, loan)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *models.Loan) error); ok {
		r0 = rf(ctx, loan)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewLoanRepository creates a new instance of LoanRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewLoanRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *LoanRepository {
	mock := &LoanRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
