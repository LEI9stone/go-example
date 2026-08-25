package catalog

import (
	"errors"
	"testing"
)

func TestNewProduct(t *testing.T) {
	product, err := NewProduct("Test Product", 1099)
	if err != nil{
		t.Fatalf("expected no error, got %v", err)
	}
	if (product.DisplayPrice() != "$10.99") {
		t.Fatalf("unexpected price to be $10.99, got %s", product.DisplayPrice())
	}
}

func TestNewProducRejectsInvalidPrice(t *testing.T) {
	_, err := NewProduct("Test Product", 0)
	if !errors.Is(err, ErrInvalidPrice) {
		t.Fatalf("expected ErrInvalidPrice, got  %v", err)
	}
}