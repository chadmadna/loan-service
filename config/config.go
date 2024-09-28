package config

import (
	"loan-service/utils/errs"

	"github.com/ilyakaznacheev/cleanenv"
)

type contextKey string

const ContextKeyDBTransaction contextKey = "ContextKeyDBTransaction"

type Config struct {
	AppPort string `env:"APP_PORT" env-default:"8080"`
	Stage   string `env:"STAGE" env-default:"development"`

	DBHost     string `env:"DB_HOST" env-default:"postgres"`
	DBPort     string `env:"DB_PORT" env-default:"5432"`
	DBUser     string `env:"DB_USER" env-default:"rw_loanservice"`
	DBPassword string `env:"DB_PASSWORD" env-default:"LoanService@2024!"`
	DBName     string `env:"DB_NAME" env-default:"loanservice_db"`

	AppSecret string `env:"APP_SECRET" env-required:"true"`

	EmailSendGridAPIKey  string `env:"EMAIL_SENDGRID_API_KEY" env-required:"true"`
	DefaultSenderAddress string `env:"DEFAULT_SENDER_ADDRESS" env-required:"true"`
	DefaultSenderName    string `env:"DEFAULT_SENDER_NAME" env-required:"true"`
}

var (
	Data Config
)

func LoadFromFile(path string) error {
	err := cleanenv.ReadConfig(path, &Data)
	if err != nil {
		return errs.Wrap(err)
	}

	return nil
}

func LoadFromEnv() error {
	err := cleanenv.ReadEnv(&Data)
	if err != nil {
		return errs.Wrap(err)
	}

	return nil
}
