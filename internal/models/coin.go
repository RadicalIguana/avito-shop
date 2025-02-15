package models

type TransferRequest struct {
	ToUser int `json:"toUser"`
	Amount int `json:"amount"`
}

type Transfer struct {
	ID        int 	 `json:"id"`
	FromUser  int	 `json:"fromUser"`
	ToUser    int 	 `json:"toUser"`
	Amount    int    `json:"amount"`
	CreatedAt string `json:"createdAt"`
}


// TODO: Просмотреть правильность использования UserCoin вместо User
type UserCoin struct {
	ID    int `json:"id"`
	Coins int `json:"coins"`
}

// TODO: Может добавить ErrorResponse struct?