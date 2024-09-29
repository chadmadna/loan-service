package money

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Rhymond/go-money"
)

func DisplayMoney(moneyStr string) string {
	moneyFloat, _ := strconv.ParseFloat(moneyStr, 64)
	return money.NewFromFloat(moneyFloat, money.IDR).Display()
}

func DisplayAsPercentage(fl float64) string {
	floatStr := fmt.Sprintf("%f", fl*100)
	for {
		trim := strings.TrimSuffix(floatStr, "0")
		if floatStr == trim {
			break
		}

		floatStr = trim
	}

	floatStr = strings.TrimSuffix(floatStr, ".")

	return fmt.Sprintf("%s%%", floatStr)
}
