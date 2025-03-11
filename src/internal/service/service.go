package service

import "time"

type ReservationCreateInput struct {
	AccountId int
	ProductId int
	OrderId   int
	Amount    int
}

type OperationHistoryInput struct {
	AccountId int
	SortType  string
	Offset    int
	Limit     int
}

type OperationHistoryOutput struct {
	Amount      int       `json:"amount"`
	Operation   string    `json:"operation"`
	Time        time.Time `json:"time"`
	Product     string    `json:"product,omitempty"`
	Order       *int      `json:"order,omitempty"`
	Description string    `json:"description,omitempty"`
}
