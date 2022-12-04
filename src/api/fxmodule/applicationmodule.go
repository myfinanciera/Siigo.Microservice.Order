// Copyright 2020 Siigo. All rights reserved.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package fxmodule

import (
	"dev.azure.com/SiigoDevOps/Siigo/_git/go-cqrs.git/cqrs"
	"go.uber.org/fx"
	"siigo.com/order/src/application/command"
	"siigo.com/order/src/application/query"
)

// ApplicationModule Application Module Commands And Queries
//modulo fx con dependencias
var ApplicationModule = fx.Options(
	fx.Provide(
		command.NewOrderCommandHandler,
		query.NewOrderQueryHandler,
	),
	fx.Invoke(
		RegisterHandlers,
	),
)

// RegisterHandlers Register Commands and Queries
func RegisterHandlers(dispatcher cqrs.Dispatcher,
	orderCommandHandler *command.OrderCommandHandler,
	orderQueryHandler *query.OrderQueryHandler,
) {

	// Configure Commands
	HandleErrorRegister(
		dispatcher.RegisterHandler(orderCommandHandler, &command.CreateOrderCommand{}),
		dispatcher.RegisterHandler(orderCommandHandler, &command.DeleteOrderCommand{}),
		dispatcher.RegisterHandler(orderCommandHandler, &command.UpdateOrderCommand{}),
	)

	// Configure Queries
	HandleErrorRegister(
		dispatcher.RegisterHandler(orderQueryHandler, &query.LoadOrderQuery{}),
		dispatcher.RegisterHandler(orderQueryHandler, &query.LoadAllOrderQuery{}),
	)

}

func HandleErrorRegister(errors ...error) {
	for _, err := range errors {
		if err != nil {
			panic(err)
		}
	}
}
