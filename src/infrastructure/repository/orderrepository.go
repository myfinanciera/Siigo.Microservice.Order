// Copyright 2020 Siigo. All rights reserved.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"siigo.com/order/src/domain/order"
)

func (c OrderRepository) Save(order *order.Order) chan error {

	result := make(chan error)
	//ejecuta asincronamente
	go func() {
		defer close(result)

		_, insertErr := c.collection.InsertOne(context.TODO(), order)

		if insertErr != nil {
			c.logger.Error("InsertOne Order ERROR:", insertErr)
		}
		//en canal result graba el error
		result <- insertErr
	}()

	return result

}

func (c OrderRepository) Delete(order *order.Order) chan error {

	result := make(chan error)

	go func() {
		defer close(result)
		_, deleteErr := c.collection.DeleteOne(context.TODO(), order)

		if deleteErr != nil {
			c.logger.Error("DeleteOne Order ERROR:", deleteErr)
		}

		result <- deleteErr
	}()

	return result

}

func (c OrderRepository) Update(order *order.Order) chan error {

	result := make(chan error)

	go func() {
		defer close(result)

		_, updateErr := c.collection.UpdateOne(context.TODO(),
			bson.M{"_id": order.Id},
			bson.D{{"$set", order}})

		if updateErr != nil {
			c.logger.Error("UpdateOne Order ERROR:", updateErr)
		}

		result <- updateErr
	}()

	return result

}
