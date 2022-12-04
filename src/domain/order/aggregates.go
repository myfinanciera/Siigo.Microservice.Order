package order

import (
	"dev.azure.com/SiigoDevOps/Siigo/_git/go-cqrs.git/cqrs"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Order struct {
	Id                 string `bson:"_id"`
	cqrs.AggregateBase `bson:"-"`
	Idcustomer         string   `json:"idcustomer" validate:"required"`
	Number             int32  `json:"number"`
	Occurred           *timestamppb.Timestamp
	State              bool `json:"state"`
    Product          []string `json:"product"`
	Email	           string   `json:"email" validate:"required"`
}

type OrderResponse struct {
	Error    error
	Order *Order
}

type OrdersResponse struct {
	Error     error
	Orders []*Order
}

// NewOrder constructs a new order aggregate.
//
// Importantly it embeds a new AggregateBase.
func NewOrder(id string) *Order {
	aggregate := &Order{
		Id:            id,
		AggregateBase: *cqrs.NewAggregateBase(id),
		Occurred:    timestamppb.Now(),
	}

	return aggregate
}

func NewOrderWithUuid() *Order {
	return NewOrder(cqrs.NewUUID())
}


// Apply handles the logic of events on the aggregate.
func (a *Order) Apply(message cqrs.EventMessage, isNew bool) {
	a.TrackChange(message)
}
