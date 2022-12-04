package finder

import (
	"context"
	"dev.azure.com/SiigoDevOps/Siigo/_git/go-cqrs.git/cqrs/uuid"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"siigo.com/order/src/domain/order"
)

func (c OrderFinder) GetAll() chan order.OrdersResponse {

	result := make(chan order.OrdersResponse)
	ordersResponse := order.OrdersResponse{}

	go func() {
		defer close(result)

		cursor, err := c.collection.Find(context.Background(), bson.D{{}})
		if err != nil {
			c.logger.Error("[Mongo] Error find orders:", err)
			ordersResponse.Error = err
		} else {
			defer cursor.Close(context.Background())
			var orders []*order.Order

			// iterate for each order
			for cursor.Next(context.Background()) {
				var ct *order.Order
				err := cursor.Decode(&ct)
				if err != nil {
					c.logger.Error("[Mongo] Error decoding orders:", err)
					ordersResponse.Error = err
					result <- ordersResponse
				} else {
					orders = append(orders, ct)
				}
			}

			ordersResponse.Orders = orders
		}

		result <- ordersResponse

	}()

	return result
}

func (c OrderFinder) Get(id uuid.UUID) chan *order.OrderResponse {

	result := make(chan *order.OrderResponse)
	response := &order.OrderResponse{}

	go func() {
		defer close(result)

		orderResponse := &order.Order{}

		if findError := c.collection.
			FindOne(context.Background(), bson.D{{"_id", id.String()}}).
			Decode(orderResponse); findError != nil {

			c.logger.Error(findError)

			if errors.Is(findError, mongo.ErrNoDocuments) {
				response.Error = errors.New("order not found")
			} else {
				response.Error = findError
			}
		}

		response.Order = orderResponse
		result <- response
	}()

	return result

}
