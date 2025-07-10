package currency

import (
	"strconv"
)

// Currency struct for currency representation.
type Currency struct {
	ID     string
	Symbol string
	Slug   string
}

// ID create currency from id.
func ID(id int) Currency {
	return Currency{
		ID: strconv.Itoa(id),
	}
}

// Symbol create currency from symbol.
func Symbol(symbol string) Currency {
	return Currency{
		Symbol: symbol,
	}
}

// Slug create currency from slug.
func Slug(slug string) Currency {
	return Currency{
		Slug: slug,
	}
}
