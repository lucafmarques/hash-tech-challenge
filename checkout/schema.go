package checkout

type CheckoutResponse struct {
	TotalAmount             int               `json:"total_amount"`
	TotalAmountWithDiscount int               `json:"total_amount_with_discount"`
	TotalDiscount           int               `json:"total_discount"`
	Products                []ProductResponse `json:"products"`
}

type ProductResponse struct {
	ID          int  `json:"id"`
	Quantity    int  `json:"quantity"`
	UnitAmount  int  `json:"unit_amount"`
	TotalAmount int  `json:"total_amount"`
	Discount    int  `json:"discount"`
	Gift        bool `json:"is_gift"`
}

type CheckoutRequest struct {
	Products []ProductRequest `json:"products"`
}

type ProductRequest struct {
	ID       int `json:"id"`
	Quantity int `json:"quantity"`
}

type DiscountResponse struct {
	Percentage float32 `json:"percentage"`
}

type HTTPError struct {
	Message string `json:"message"`
}
