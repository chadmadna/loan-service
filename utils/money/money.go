package money

import (
	"errors"
	"loan-service/utils/errs"
	"math"
	"strconv"

	gomoney "github.com/Rhymond/go-money"
)

type MoneyString string

func (m MoneyString) ToFloat64() (float64, error) {
	result, err := strconv.ParseFloat(string(m), 64)
	if err != nil {
		return 0, errs.Wrap(errors.New("cannot parse money as float"))
	}

	return result, nil
}

// Rounded to 2 decimal points
func (m MoneyString) ToMoney(currencyStr string) (*gomoney.Money, error) {
	moneyFloat, err := m.ToFloat64()
	if err != nil {
		return nil, errs.Wrap(err)
	}

	roundedMoneyFloat := math.Floor(moneyFloat*100) / 100
	return gomoney.NewFromFloat(roundedMoneyFloat, currencyStr), nil
}
