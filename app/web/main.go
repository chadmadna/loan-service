package main

import (
	authMiddleware "loan-service/app/web/middleware"
	"loan-service/config"
	"loan-service/database"
	"loan-service/models"
	loansModule "loan-service/modules/loans"
	"loan-service/services/email"
	"loan-service/utils/resp"
	"loan-service/utils/tern"
	"net/http"
	"os"

	_userHandlers "loan-service/modules/users/handlers"

	"github.com/apsystole/log"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/samber/do"
	"golang.org/x/time/rate"
	"gorm.io/gorm"
)

var injector *do.Injector

func main() {
	log.Println("Running in stage:", config.Data.Stage)

	// Load env variables
	envPath := tern.String(os.Getenv("ENV_PATH"), ".env")
	err := config.Load(envPath)
	if err != nil {
		panic(err)
	}

	// Get database instance
	db, err := database.GetDB()
	if err != nil {
		panic(err)
	}

	// Dependency injection
	injector = do.New()

	// Database
	do.Provide[*gorm.DB](injector, func(i *do.Injector) (*gorm.DB, error) {
		return db, nil
	})

	// Services
	do.Provide[email.EmailService](injector, func(i *do.Injector) (email.EmailService, error) {
		return email.NewEmailService(config.Data.EmailSendGridAPIKey), nil
	})

	// Loans module
	do.Provide[models.LoanRepository](injector, func(i *do.Injector) (models.LoanRepository, error) {
		return loansModule.NewLoansRepository(db), nil
	})

	do.Provide[models.LoanUsecase](injector, func(i *do.Injector) (models.LoanUsecase, error) {
		return loansModule.NewLoansUsecase(
			do.MustInvoke[models.LoanRepository](i),
		), nil
	})

	// HTTP server instance
	e := echo.New()
	e.HTTPErrorHandler = httpErrorHandler
	e.HideBanner = true

	e.Pre(middleware.RemoveTrailingSlash())

	// Middleware
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*://localhost:*"}, // add allowed domains
	}))
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(
		rate.Limit(50),
	)))

	// Register router groups
	mainGroup := e.Group("/", authMiddleware.JWTAuth(do.MustInvoke[models.UserRepository](injector)))
	// staffGroup := mainGroup.Group("/admin", authMiddleware.AllowOnlyRoles(
	// 	auth.RoleTypeSuperuser, auth.RoleTypeStaff,
	// ))
	// fieldValidatorGroup := mainGroup.Group("/field-validation", authMiddleware.AllowOnlyRoles(
	// 	auth.RoleTypeSuperuser, auth.RoleTypeStaff, auth.RoleTypeFieldValidator,
	// ))
	// investorGroup := mainGroup.Group("/invest", authMiddleware.AllowOnlyRoles(
	// 	auth.RoleTypeSuperuser, auth.RoleTypeInvestor,
	// ))
	// borrowGroup := mainGroup.Group("/user", authMiddleware.AllowOnlyRoles(
	// 	auth.RoleTypeSuperuser, auth.RoleTypeBorrower,
	// ))

	// Register endpoints for each group
	_userHandlers.NewCommonUsersHandler(mainGroup, do.MustInvoke[models.UserUsecase](injector))
}

func httpErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
	}

	if code != http.StatusInternalServerError {
		_ = c.JSON(code, err)
	} else {
		log.Error(err)
		_ = resp.HTTPServerError(c)
	}
}
