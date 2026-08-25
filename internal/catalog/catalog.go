package catalog

import (
	"errors"
	"fmt"
)

var ErrInvalidPrice = errors.New("price must be greater than zero")

type Product struct {
	Name string
	Price int64
}

func NewProduct(name string, price int64) (*Product, error) {
	if name == "" {
		return nil, errors.New("name is required")
	}
	if price <=0 {
		return nil, ErrInvalidPrice
	}
	return &Product{Name: name, Price: price}, nil
}

func (p Product) DisplayPrice() string {
	return fmt.Sprintf("$%d.%02d", p.Price/100, p.Price%100)
}