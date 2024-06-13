package vo

import "errors"

// Currency 这里看作是一个枚举
type Currency struct {
	name string
}

// 值对象
type MonetaryValue struct {
	amount   int
	currency Currency
}

func (v MonetaryValue) Amount() int {
	return v.amount
}

func NewMonetaryValue(amount int, currency Currency) (MonetaryValue, error) {
	if amount < 0 {
		return MonetaryValue{}, errors.New("amount must ge 0")
	}
	return MonetaryValue{
		amount:   amount,
		currency: currency,
	}, nil
}
