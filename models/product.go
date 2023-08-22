package models

type Product struct {
	IdProduct       int     `json:"idProduct"`
	NamaProduct     string  `json:"namaProduct"`
	QtyProduct      int     `json:"qtyProduct"`
	PriceProduct    float64 `json:"priceProduct"`
	CategoryProduct string  `json:"categoryProduct"`
}
