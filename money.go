/*
Package money is a library to deal with money and currency representation.
Inspired by ruby money library http://rubymoney.github.io/money.

Defaults

    Options{
      "currency":                 "usd",
      "with_cents":               true,
      "with_currency":            false,
      "with_symbol":              true,
      "with_symbol_space":        false,
      "with_thousands_separator": true,
    }

Usage

    New(10)                                               // "$10.00"
    New(10, Options{"currency": "eur"})                   // "€10.00"
    New(10, Options{"with_cents": false})                 // "$10"
    New(10, Options{"with_currency:" true })              // "$10.00 USD"
    New(10, Options{"with_symbol": false})                // "10.00"
    New(10, Options{"with_symbol_space":true})            // "$ 10.00"
    New(1000)                                             // "$1,000.00"
    New(1000, Options{"with_thousands_separator": false}) // "$1000.00"
*/
package money

import (
	"fmt"
	"math"
	"strings"
)

// New returns a formatted price string according to currency rules and options
func New(val float64, opts ...Options) (result string) {
	options := defaults()

	if len(opts) > 0 {
		options = override(options, opts[0])
	}

	currency := currencies[options["currency"].(string)]

	integer, fractional := splitValue(val)

	if options["with_thousands_separator"].(bool) {
		result = separateThousands(integer, currency["thousands_separator"].(string))
	} else {
		result = integer
	}

	if options["with_cents"].(bool) && currency["subunit"] != nil {
		result = fmt.Sprintf("%s%s%s", result, currency["decimal_mark"].(string), fractional)
	}

	if options["with_symbol"].(bool) {
		result = addSymbol(result, currency, options)
	}

	if options["with_currency"].(bool) {
		result = fmt.Sprintf("%s %s", result, currency["iso_code"])
	}

	return result
}

func addSymbol(result string, currency map[string]interface{}, options Options) string {
	var space string

	if options["with_symbol_space"].(bool) {
		space = " "
	}

	if currency["symbol_first"].(bool) {
		result = fmt.Sprintf("%s%s%s", currency["symbol"], space, result)
	} else {
		result = fmt.Sprintf("%s%s%s", result, space, currency["symbol"])
	}

	return result
}

func separateThousands(value, separator string) string {
	s := len(value) / 3
	m := int(math.Mod(float64(len(value)), 3))

	if m > 0 {
		s++
	}

	if s == 0 {
		return value
	}

	r := make([]string, s)

	for i := 0; i < len(r); i++ {
		if i == 0 && m > 0 {
			r[i] = value[i : i+m]
		} else {
			r[i] = value[i : i+3]
		}
	}

	return strings.Join(r, separator)
}

func splitValue(val float64) (integer, fractional string) {
	i, f := math.Modf(val)

	integer = fmt.Sprintf("%.0f", i)
	fractional = fmt.Sprintf("%.2f", f)[2:]

	return
}