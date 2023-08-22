package models

type AddCartRequestBody struct {
	Product int `json:"idProduct"`
	Qty     int `json:"qty"`
}

type RemoveCartRequestBody struct {
	IdCart string `json:"idCart"`
}

type Cart struct {
	IdCart        string  `json:"idCart"`
	IdProduct     int     `json:"idProduct"`
	QtyPurchase   int     `json:"qtyPurchase"`
	PricePurchase float64 `json:"pricePurchase"`
	NamaProduct   string  `json:"nameProduct"`
	QtyProduct    int     `json:"qtyProduct"`
}
