package order

import (
	"dev.azure.com/SiigoDevOps/Siigo/_git/go-cqrs.git/cqrs"
	"dev.azure.com/SiigoDevOps/Siigo/_git/go-cqrs.git/cqrs/uuid"
)

// Events are just plain structs

// OrderCreatedDomainEvent CreatedEvent OrderCreatedEvent OrderCreated event
type OrderCreatedDomainEvent struct {
	Order *Order
}

// OrderUpdatedDomainEvent UpdatedEvent CreatedEvent OrderCreatedEvent OrderCreated event
type OrderUpdatedDomainEvent struct {
	Order *Order
}

// NewOrderCreatedEvent Creates a new order created domain event struct
func NewOrderCreatedEvent(aggregate cqrs.AggregateRoot) *cqrs.EventDescriptor {
	eventid := uuid.NewV4().String()
	return cqrs.NewEventMessagePutKindEvent(
		aggregate.AggregateID(),
		eventid,
		&OrderCreatedDomainEvent{Order: aggregate.(*Order)},
		&OrderCreatedDomainEvent{},
		cqrs.Int(aggregate.CurrentVersion()),
	)
}

// NewOrderUpdateEvent  Creates a new order updated domain event struct
func NewOrderUpdateEvent(aggregate cqrs.AggregateRoot) *cqrs.EventDescriptor {
	eventid := uuid.NewV4().String()
	return cqrs.NewEventMessagePutKindEvent(
		aggregate.AggregateID(),
		eventid,
		&OrderUpdatedDomainEvent{Order: aggregate.(*Order)},
		&OrderUpdatedDomainEvent{},
		cqrs.Int(aggregate.CurrentVersion()),
	)
}
