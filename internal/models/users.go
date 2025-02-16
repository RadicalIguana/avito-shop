package models

type User struct {
	Id 		 int 	`json:"id"`
	Name 	 string `json:"name"`
	Password string `json:"password"`
	Coins 	 int 	`json:"coins"`
}

type UserResponse struct {
	Coins 		int 			`json:"coins"`
	Inventory   []InventoryItem `json:"inventory"`
	CoinHistory CoinHistory 	`json:"coinHistory"`
}

type InventoryItem struct {
	Name 	 string `json:"name"`
	Quantiry int    `json:"quantity"`
}

type CoinHistory struct {
	Received []Transfer `json:"received"`
	Sent	 []Transfer `json:"sent"` 
}