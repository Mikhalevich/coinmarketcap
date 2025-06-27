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

// Currency currency representation.
type Currency struct {
	ID     ID
	Symbol string
	Slug   string
}

func FromID(id int) Currency {
	return Currency{
		ID: IDFromInt(id),
	}
}

func Symbol(symbol string) Currency {
	return Currency{
		Symbol: symbol,
	}
}

func Slug(slug string) Currency {
	return Currency{
		Slug: slug,
	}
}
