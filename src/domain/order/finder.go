package order

import "dev.azure.com/SiigoDevOps/Siigo/_git/go-cqrs.git/cqrs/uuid"

type IOrderFinder interface {
	GetAll() chan OrdersResponse
	Get(id uuid.UUID) chan *OrderResponse
}
