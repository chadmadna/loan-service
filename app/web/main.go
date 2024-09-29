package main

import (
	"fmt"
	"loan-service/app"
	authMiddleware "loan-service/app/web/middleware"
	"loan-service/config"
	"loan-service/database"
	"loan-service/models"
	"loan-service/services/auth"
	"loan-service/services/email"
	"loan-service/services/upload"
	"loan-service/utils/resp"
	"loan-service/utils/tern"
	"net/http"
	"os"

	_loanHandlers "loan-service/modules/loans/handlers"
	_productHandlers "loan-service/modules/products/handlers"
	_userHandlers "loan-service/modules/users/handlers"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_log "github.com/labstack/gommon/log"
	"github.com/samber/do"
	"golang.org/x/time/rate"
)

func main() {
	fmt.Println("Running in stage:", config.Data.Stage)

	// Load env variables
	envPath := tern.String(os.Getenv("ENV_PATH"), ".env")
	err := config.LoadFromFile(envPath)
	if err != nil {
		panic(err)
	}

	// Get database instance
	db, err := database.GetDB()
	if err != nil {
		panic(err)
	}

	// HTTP server instance
	e := echo.New()
	e.HTTPErrorHandler = httpErrorHandler
	e.HideBanner = true
	e.Logger.SetLevel(_log.DEBUG)
	e.Validator = &CustomValidator{validator: validator.New()}
	e.Use(middleware.Static("/tmp"))

	e.Pre(middleware.RemoveTrailingSlash())

	// Middleware
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*://localhost:*"}, // add allowed domains
	}))
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(
		rate.Limit(50),
	)))

	// Services
	emailSvc := email.NewEmailService(config.Data.EmailSendGridAPIKey, config.Data.DefaultSenderAddress, config.Data.DefaultSenderName)
	uploadSvc := upload.NewUploadService()

	// Dependency injection
	injector := app.SetupInjections(db, e, emailSvc, uploadSvc)

	// Register router groups
	mg := e.Group("/app", authMiddleware.JWTAuth(do.MustInvoke[models.UserRepository](injector)))
	staffGroup := mg.Group("/admin", authMiddleware.AllowOnlyRoles(
		auth.RoleTypeSuperuser, auth.RoleTypeStaff,
	))
	fieldValidatorGroup := mg.Group("/field-validation", authMiddleware.AllowOnlyRoles(
		auth.RoleTypeSuperuser, auth.RoleTypeStaff, auth.RoleTypeFieldValidator,
	))
	investorGroup := mg.Group("/invest", authMiddleware.AllowOnlyRoles(
		auth.RoleTypeSuperuser, auth.RoleTypeInvestor,
	))
	borrowGroup := mg.Group("/user", authMiddleware.AllowOnlyRoles(
		auth.RoleTypeSuperuser, auth.RoleTypeBorrower,
	))

	// Healthcheck
	e.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})

	// Register endpoints for each group
	_userHandlers.NewCommonUserHandler(
		e,
		do.MustInvoke[models.UserUsecase](injector),
		authMiddleware.JWTAuth(do.MustInvoke[models.UserRepository](injector)),
	)

	_productHandlers.NewProductHandler(
		borrowGroup,
		do.MustInvoke[models.ProductUsecase](injector),
	)

	do.Provide[*_loanHandlers.CommonLoanHandler](injector, func(i *do.Injector) (*_loanHandlers.CommonLoanHandler, error) {
		return _loanHandlers.NewCommonLoanHandler(
			do.MustInvoke[models.LoanUsecase](injector),
			do.MustInvoke[models.UserUsecase](injector),
		), nil
	})

	_loanHandlers.NewStaffHandler(
		staffGroup,
		do.MustInvoke[models.LoanUsecase](injector),
		do.MustInvoke[models.UserUsecase](injector),
		do.MustInvoke[*_loanHandlers.CommonLoanHandler](injector),
	)

	_loanHandlers.NewFieldValidatorHandler(
		fieldValidatorGroup,
		do.MustInvoke[models.LoanUsecase](injector),
		do.MustInvoke[models.UserUsecase](injector),
		do.MustInvoke[*_loanHandlers.CommonLoanHandler](injector),
	)

	_loanHandlers.NewInvestorHandler(
		investorGroup, // lol
		do.MustInvoke[models.LoanUsecase](injector),
		do.MustInvoke[models.UserUsecase](injector),
		do.MustInvoke[*_loanHandlers.CommonLoanHandler](injector),
	)

	_loanHandlers.NewBorrowerHandler(
		borrowGroup,
		do.MustInvoke[models.LoanUsecase](injector),
		do.MustInvoke[models.UserUsecase](injector),
		do.MustInvoke[models.ProductUsecase](injector),
		do.MustInvoke[*_loanHandlers.CommonLoanHandler](injector),
	)

	// Start server
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", config.Data.AppPort)))
}

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func httpErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
	}

	if code != http.StatusInternalServerError {
		_ = c.JSON(code, err)
	} else {
		c.Logger().Error(err)
		_ = resp.HTTPServerError(c)
	}
}
