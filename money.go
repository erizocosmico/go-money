package money

import (
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

var (
	moneyRegex       = regexp.MustCompile(`^([^0-9 ]*) *([0-9][0-9,\.]*) *(k|m{1,2})? *([^0-9]*)$`)
	errInvalid       = errors.New("invalid money format given")
	errInvalidNumber = errors.New("invalid number format given")
)

// Parse returns an Amount of money with quantity and currency (if identified) and an error.
func Parse(text string) (*Amount, error) {
	return parse(text, ".")
}

// ParseComma parses the amount of money in the text with "," instead of "." as decimal separator.
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

var currencies = map[string][]string{
	"Lek":  {"all", "lek"},
	"؋":    {"؋", "afn"},
	"$":    {"$", "usd", "svc", "ars", "aud", "bsd", "bbd", "bmd", "bnd", "cad", "kyd", "clp", "cop", "xcd", "svc", "fjd", "gyd", "hkd", "lrd", "mxn", "nad", "nzd", "sgd", "sbd", "srd", "tvd"},
	"ƒ":    {"ƒ", "awg", "ang"},
	"₼":    {"₼", "azn"},
	"p.":   {"p.", "byr"},
	"BZ$":  {"bz$", "bzd"},
	"$b":   {"$b", "bob"},
	"KM":   {"km", "bam"},
	"P":    {"p", "bwp"},
	"лв":   {"лв", "bgn", "kzt", "uzs"},
	"R$":   {"r$", "brl"},
	"៛":    {"៛", "khr"},
	"¥":    {"¥", "dny", "jpy"},
	"₡":    {"₡", "crc"},
	"kn":   {"kn", "hrk"},
	"₩":    {"₩", "kpw", "krw"},
	"₱":    {"₱", "cup", "php"},
	"Kč":   {"Kč", "czk"},
	"kr":   {"kr", "dkk", "eek", "isk", "nok", "sek"},
	"RD$":  {"rd$", "dop"},
	"£":    {"£", "gbp", "egp", "fkp", "gip", "ggp", "imp", "jep", "lbp", "shp", "syp"},
	"GEL":  {"gel"},
	"€":    {"€", "eur"},
	"¢":    {"¢", "ghc"},
	"Q":    {"q", "gtq"},
	"L":    {"l", "hnl"},
	"Ft":   {"ft", "huf"},
	"₹":    {"₹", "inr"},
	"Rp":   {"rp", "idr"},
	"﷼":    {"﷼", "irr", "omr", "qar", "sar", "yer"},
	"₪":    {"₪", "ils"},
	"J$":   {"j$", "jmd"},
	"₭":    {"₭", "lak"},
	"Ls":   {"ls", "lvl"},
	"Bs":   {"bs", "vef"},
	"Lt":   {"lt", "ltl"},
	"ден":  {"ден", "mkd"},
	"RM":   {"rm", "myr"},
	"Rs":   {"rs", "mur", "npr", "pkr", "scr", "lkr"},
	"₮":    {"₮", "mnt"},
	"MT":   {"mt", "mzn"},
	"C$":   {"c$", "nio"},
	"₦":    {"₦", "ngn"},
	"B/.":  {"b/.", "pab"},
	"S/.":  {"s/.", "pen"},
	"Gs":   {"gs", "pyg"},
	"zł":   {"zł", "pln"},
	"lei":  {"lei", "ron"},
	"Дин.": {"Дин.", "rsd"},
	"S":    {"s", "sos", "zar"},
	"CHF":  {"chf"},
	"NT$":  {"nt$", "twd"},
	"Z$":   {"z$", "zwd"},
	"TT$":  {"tt$", "ttd"},
	"฿":    {"฿", "tbh"},
	"₺":    {"₺", "tlr"},
	"₴":    {"₴", "uah"},
	"$U":   {"$u", "uyu"},
	"₫":    {"₫", "vnd"},
}

var currenciesByAlias = make(map[string]string)

func init() {
	for c, aliases := range currencies {
		for _, a := range aliases {
			currenciesByAlias[a] = c
		}
	}
}

func identifyCurrency(curr string) string {
	c, ok := currenciesByAlias[curr]
	if !ok {
		return ""
	}
	return c
}
