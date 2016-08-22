package money

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseAndString(t *testing.T) {
	cases := []struct {
		input  string
		output string
		err    bool
	}{
		{"3500000 €", "3.5M €", false},
		{"3500000€", "3.5M €", false},
		{"3500000 EUR", "3.5M eur", false},
		{"$ 3500000", "3.5M $", false},
		{"$3500000", "3.5M $", false},
		{"USD 3500000", "3.5M usd", false},
		{"USD 35k", "35K usd", false},
		{"USD 3.5M", "3.5M usd", false},
		{"USD 3.5 M", "3.5M usd", false},
		{"3.5 M EUR", "3.5M eur", false},
		{"3500 EUR", "3.5K eur", false},
		{"500 EUR", "500 eur", false},
		{"500.567 EUR", "500.57 eur", false},
		{"€ 3500000 €", "3.5M €", true},
		{"€ 35.000.00 €", "3.5M €", true},
	}

	for _, c := range cases {
		a, err := Parse(c.input)
		if c.err {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err)
			assert.Equal(t, c.output, a.String())
		}
	}
}

func TestStringComma(t *testing.T) {
	a, err := Parse("3.5M€")
	assert.Nil(t, err)
	assert.Equal(t, "3,5M €", a.StringComma())
}

func TestStringBefore(t *testing.T) {
	a, err := Parse("3.5M$")
	assert.Nil(t, err)
	assert.Equal(t, "$3.5M", a.StringBefore())
}
