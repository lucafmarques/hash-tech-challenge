package main

import "gitlab.com/lucafmarques/hash-test/repository"

func CalculateDiscount(productValue int, discountPerc float32) int {
	return int(float32(productValue) * discountPerc)
}

func NewProductResponse(productData *repository.Product, productRequest ProductRequest) ProductResponse {
	return ProductResponse{
		Discount:    0,
		Quantity:    productRequest.Quantity,
		ID:          productData.Id,
		UnitAmount:  productData.Amount,
		TotalAmount: productData.Amount * productRequest.Quantity,
		Gift:        productData.Gift,
	}
}
