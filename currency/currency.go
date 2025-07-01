package currency

import "strconv"

// ID represents currency coinmarketcap enternal id.
type ID string

// IDFromInt convert int to ID.
func IDFromInt(id int) ID {
	return ID(strconv.Itoa(id))
}

// String convert ID to string.
func (id ID) String() string {
	return string(id)
}

// Currency struct for currency representation.
type Currency struct {
	ID     ID
	Symbol string
	Slug   string
}

// FromID create currency from id.
func FromID(id int) Currency {
	return Currency{
		ID: IDFromInt(id),
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
