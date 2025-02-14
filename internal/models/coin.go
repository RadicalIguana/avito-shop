package models

type TransferRequest struct {
	ToUser string `json:"toUser"`
	Amount int    `json:"amount"`
}

type Transfer struct {
	ID        string `json:"id"`
	FromUser  string `json:"fromUser"`
	ToUser    string `json:"toUser"`
	Amount    int    `json:"amount"`
	CreatedAt string `json:"createdAt"`
}


// TODO: Просмотреть правильность использования UserCoin вместо User
type UserCoin struct {
	ID    string `json:"id"`
	Coins int `json:"coins"`
}

// TODO: Может добавить ErrorResponse struct?