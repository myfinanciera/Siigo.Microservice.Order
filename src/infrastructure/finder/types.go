package finder

import (
	"context"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"siigo.com/order/src/domain/order"
)

// OrderFinder OrderRepository provides methods for processing mongo order queries
type OrderFinder struct {
	collection *mongo.Collection
	context    context.Context
	logger     *logrus.Logger
}

// NewOrderFinder a new OrderFinder
func NewOrderFinder(context context.Context, collection *mongo.Collection, logger *logrus.Logger) order.IOrderFinder {
	return OrderFinder{context: context, collection: collection, logger: logger}
}
