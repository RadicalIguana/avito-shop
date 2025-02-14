package models

// TODO: Merch переименовать в Item? Или в Purchase?
type Merch struct {
	Name  string `json:"name"`
	Price int    `json:"price"`
}

// TODO: Может переименовать?
type BuyRequest struct {
	ItemName string `json:"item"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}