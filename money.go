package money

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Amount represents a quantity of money with a currency.
type Amount struct {
	// Quantity is the quantity of the amount.
	Quantity float64
	// Currency is the currency of the amount.
	Currency string
}

// NewAmount creates a new amount with the given quantity and currency.
func NewAmount(q float64, currency string) *Amount {
	return &Amount{q, currency}
}

// String returns the money as string in the format "[QUANTITY](.DECIMAL)(M|k)( CURRENCY)".
func (a Amount) String() string {
	var n float64
	var suffix string
	if a.Quantity >= 1000000 {
		n = a.Quantity / 1000000
		suffix = "M"
	} else if a.Quantity >= 1000 {
		n = a.Quantity / 1000
		suffix = "K"
	} else {
		n = a.Quantity
	}
	num := strings.TrimRight(strings.TrimRight(fmt.Sprintf("%.2f", n), "0"), ".")
	if a.Currency != "" {
		return fmt.Sprintf("%s%s %s", num, suffix, a.Currency)
	}
	return fmt.Sprintf("%s%s", num, suffix)
}

// StringComma returns the money as string in the format "[QUANTITY](,DECIMAL)(M|k)( CURRENCY)".
func (a Amount) StringComma() string {
	return strings.Replace(a.String(), ".", ",", -1)
}

// StringBefore returns the money as string in the format "[CURRENCY][QUANTITY](.DECIMAL)(M|k)".
func (a Amount) StringBefore() string {
	parts := strings.Split(a.String(), " ")
	if len(parts) > 1 {
		return parts[1] + parts[0]
	}
	return parts[0]
}

var zero = []byte("0")

func trailingZeros(n int) string {
	var buf bytes.Buffer
	for i := 0; i < n; i++ {
		buf.Write(zero)
	}
	return buf.String()
}

var (
	moneyRegex       = regexp.MustCompile(`^([^0-9 ]*) *([0-9][0-9,\.]*) *(k|m|mm)? *([^0-9]*)$`)
	errInvalid       = errors.New("invalid money format given")
	errInvalidNumber = errors.New("invalid number format given")
)

func Parse(text string) (*Amount, error) {
	return parse(text, ".")
}

func ParseComma(text string) (*Amount, error) {
	return parse(text, ",")
}

func parse(text string, decimalSep string) (*Amount, error) {
	text = strings.TrimSpace(strings.ToLower(strings.Replace(text, sepInverse(decimalSep), "", -1)))
	if !moneyRegex.MatchString(text) {
		return nil, errInvalid
	}

	parts := moneyRegex.FindStringSubmatch(text)
	if parts[1] != "" && parts[4] != "" {
		return nil, errInvalid
	}

	if strings.Count(parts[2], decimalSep) > 1 {
		return nil, errInvalidNumber
	}

	var quantity = parts[2]
	if decimalSep == "," {
		quantity = strings.Replace(quantity, ",", ".", -1)
	}

	q, err := strconv.ParseFloat(quantity, 64)
	if err != nil {
		return nil, err
	}

	switch parts[3] {
	case "m", "mm":
		q *= 1000000.0
	case "k":
		q *= 1000.0
	}

	var currency string
	switch true {
	case parts[1] != "":
		currency = identifyCurrency(parts[1])
	case parts[4] != "":
		currency = identifyCurrency(parts[4])
	}

	return &Amount{
		Quantity: q,
		Currency: currency,
	}, nil
}

func sepInverse(sep string) string {
	if sep == "," {
		return "."
	}
	return ","
}

func identifyCurrency(curr string) string {
	// TODO: Correctly identify currencies
	return curr
}
