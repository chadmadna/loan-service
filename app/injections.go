package app

import (
	"loan-service/models"
	loansModule "loan-service/modules/loans"
	productsModule "loan-service/modules/products"
	usersModule "loan-service/modules/users"
	"loan-service/services/email"
	"loan-service/services/upload"

	"github.com/labstack/echo/v4"
	"github.com/samber/do"
	"gorm.io/gorm"
)

var injector *do.Injector

func SetupInjections(
	db *gorm.DB,
	rest *echo.Echo,
	emailSvc email.EmailService,
	uploadSvc upload.UploadService,
) *do.Injector {
	injector = do.New()

	// Database
	if db != nil {
		do.Provide[*gorm.DB](injector, func(i *do.Injector) (*gorm.DB, error) {
			return db, nil
		})
	}

	// Services
	if emailSvc != nil {
		do.Provide[email.EmailService](injector, func(i *do.Injector) (email.EmailService, error) {
			return emailSvc, nil
		})
	}

	if uploadSvc != nil {
		do.Provide[upload.UploadService](injector, func(i *do.Injector) (upload.UploadService, error) {
			return uploadSvc, nil
		})
	}

	// Modules
	// Products module
	do.Provide[models.ProductRepository](injector, func(i *do.Injector) (models.ProductRepository, error) {
		return productsModule.NewProductRepository(db), nil
	})

	do.Provide[models.ProductUsecase](injector, func(i *do.Injector) (models.ProductUsecase, error) {
		return productsModule.NewProductUsecase(
			do.MustInvoke[models.ProductRepository](injector),
		), nil
	})

	// Loans module
	do.Provide[models.LoanRepository](injector, func(i *do.Injector) (models.LoanRepository, error) {
		return loansModule.NewLoanRepository(db), nil
	})

	do.Provide[models.LoanUsecase](injector, func(i *do.Injector) (models.LoanUsecase, error) {
		return loansModule.NewLoanUsecase(
			do.MustInvoke[models.LoanRepository](i),
			do.MustInvoke[models.UserUsecase](i),
			do.MustInvoke[email.EmailService](injector),
			do.MustInvoke[upload.UploadService](injector),
		), nil
	})

	// Users module
	do.Provide[models.UserRepository](injector, func(i *do.Injector) (models.UserRepository, error) {
		return usersModule.NewUserRepository(db), nil
	})

	do.Provide[models.UserUsecase](injector, func(i *do.Injector) (models.UserUsecase, error) {
		return usersModule.NewUserUsecase(
			do.MustInvoke[models.UserRepository](i),
		), nil
	})

	return injector
}
