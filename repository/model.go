package repository

type Product struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Amount      int    `json:"amount"`
	Gift        bool   `json:"is_gift"`
}
