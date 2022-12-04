package order

type IOrderRepository interface {
	Save(order *Order) chan error
	Delete(order *Order) chan error
	Update(order *Order) chan error
}
