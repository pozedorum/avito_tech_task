package service

type ReservationCreateInput struct {
	AccountId int
	ProductId int
	OrderId   int
	Amount    int
}
