package services

import (
	"encoding/json"

	"dev.azure.com/SiigoDevOps/Siigo/_git/go-cqrs.git/cqrs"
	"dev.azure.com/SiigoDevOps/Siigo/_git/go-cqrs.git/cqrs/uuid"
	"github.com/go-redis/redis"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"siigo.com/order/src/domain/order"
)

// OrderService is a services specialized for persistence events of Orders and invoke infrastructure implementation.
// While it is not required to construct a services specialized for a
// specific aggregate type, it is better to do so. There can be quite a lot of
// services configuration that is specific to a type and it is cleaner if that
// code is contained in a specialized repository as shown here.
// Also because the CommonDomainRepository Load method returns an interface{}, a
// type assertion is required. Here the type assertion is contained in this specialized
// eventStore and a *Order is returned from the eventStore.
type OrderService struct {
	events         *cqrs.DomainEventsDispatcher
	redis          *redis.Client
	orderRepo   order.IOrderRepository
	orderFinder order.IOrderFinder
}

type IOrderService interface {
	LoadAll() ([]*order.Order, error)
	Get(id uuid.UUID) (*order.Order, error)
	Save(aggregate cqrs.AggregateRoot, expectedVersion *int) error
	Delete(id string) error
	Update(aggregate cqrs.AggregateRoot, expectedVersion *int) error
}

// NewOrderService constructs a new NewOrderService.
func NewOrderService(eventStore *cqrs.DomainEventsDispatcher, redisCliente *redis.Client, orderRepo order.IOrderRepository, orderFinder order.IOrderFinder) IOrderService {

	ret := &OrderService{
		events:         eventStore,
		redis:          redisCliente,
		orderRepo:   orderRepo,
		orderFinder: orderFinder,
	}

	// An aggregate factory creates an aggregate instance given the name of an aggregate.
	aggregateFactory := cqrs.NewDelegateAggregateFactory()
	err := aggregateFactory.RegisterDelegate(&order.Order{}, func(id string) cqrs.AggregateRoot { return order.NewOrder(id) })
	if err != nil {
		panic(err)
	}
	if ret.events.EventStore != nil {
		ret.events.EventStore.SetAggregateFactory(aggregateFactory)
	}

	// A stream name delegate constructs a stream name.
	// A common way to construct a stream name is to use a bounded context and
	// an aggregate id.
	// The interface for a stream name delegate takes a two strings. One may be
	// the aggregate type and the other the aggregate id. In this case the first
	// argument and the second argument are concatenated with a hyphen.
	streamNameDelegate := cqrs.NewDelegateStreamNamer()
	err = streamNameDelegate.RegisterDelegate(func(t string, id string) string {
		return t + "-" + id
	}, &order.Order{})
	if err != nil {
		panic(err)
	}
	if ret.events.EventStore != nil {
		ret.events.EventStore.SetStreamNameDelegate(streamNameDelegate)
	}

	// An event factory creates an instance of an event given the name of an event
	// as a string.
	eventFactory := cqrs.NewDelegateEventFactory()
	err = eventFactory.RegisterDelegate(&order.OrderCreatedDomainEvent{}, func() interface{} { return &order.OrderCreatedDomainEvent{} })
	if err != nil {
		panic(err)
	}

	err = eventFactory.RegisterDelegate(&order.OrderUpdatedDomainEvent{}, func() interface{} { return &order.OrderUpdatedDomainEvent{} })
	if err != nil {
		panic(err)
	}

	if ret.events.EventStore != nil {
		ret.events.EventStore.SetEventFactory(eventFactory)
	}

	return ret
}

// LoadAll Load
func (r *OrderService) LoadAll() ([]*order.Order, error) {

	orderResponse := <-r.orderFinder.GetAll()

	if orderResponse.Error != nil {
		return nil, orderResponse.Error
	}

	return orderResponse.Orders, nil
}

// Get Load Returns an *Aggregate.
func (r *OrderService) Get(id uuid.UUID) (*order.Order, error) {

	//Evaluate redis properties
	if r.redis != nil {
		uid := id.String()
		val, err := r.redis.Get(uid).Result()
		if err == nil {
			subscribers := order.OrderResponse{}
			errorReflect := json.Unmarshal([]byte(val), &subscribers)
			if errorReflect == nil && subscribers.Order != nil {
				return subscribers.Order, nil
			}
		}

	}

	// load document by bson id
	orderResponse := <-r.orderFinder.Get(id)

	if orderResponse.Error != nil {
		return nil, orderResponse.Error
	}

	//Put item redis
	json, err := json.Marshal(orderResponse)
	if err == nil {
		if r.redis != nil {
			go r.redis.Set(orderResponse.Order.Id, json, 0).Err()
		}
	}

	return orderResponse.Order, nil
}

//logic domain
// Save persists an aggregate.
func (r *OrderService) Save(aggregate cqrs.AggregateRoot, expectedVersion *int) error {

	ctr := aggregate.(*order.Order)

	// validate cost of order
	/*if err := ctr.IsWithinBudget(); err != nil {
		return status.Error(codes.InvalidArgument, err.Error())
	}*/

	// use infrastructure for persist info agregate
	errno := <-r.orderRepo.Save(ctr)
	if errno != nil {
		return status.Error(codes.FailedPrecondition, errno.Error())
	}

	// add domain event
	ctr.Apply(order.NewOrderCreatedEvent(aggregate), true)

	// dispatch domain events
	if r.events.EventStore != nil {
		go r.events.SaveAndPublish(ctr, nil)
	}

	return nil
}

// Delete element.
func (r *OrderService) Delete(Id string) error {

	id, invalidIdError := uuid.FromString(Id)
	if invalidIdError != nil {
		return invalidIdError
	}

	orderDocument, err := r.Get(id)
	if err != nil {
		return err
	}

	// use infrastructure
	errno := <-r.orderRepo.Delete(orderDocument)
	if errno != nil {
		return errno
	}

	//Delete reference redis
	if r.redis != nil {
		go r.redis.Del(Id)
	}

	return nil
}

// Update element.
func (r *OrderService) Update(aggregate cqrs.AggregateRoot, expectedVersion *int) error {

	ctr := aggregate.(*order.Order)

	// validate cost of order
	/*if err := ctr.IsWithinBudget(); err != nil {
		return status.Error(codes.InvalidArgument, err.Error())
	}*/

	//Validate item exist
	id, invalidIdError := uuid.FromString(ctr.Id)
	if invalidIdError != nil {
		return invalidIdError
	}
	_, err := r.Get(id)
	if err != nil {
		return status.Error(codes.NotFound, "order not found")
	}

	// use infrastructure
	errno := <-r.orderRepo.Update(ctr)
	if errno != nil {
		return errno
	}

	// add domain events
	ctr.Apply(order.NewOrderUpdateEvent(aggregate), false)

	// dispatch domain events
	if r.events.EventStore != nil {
		go r.events.SaveAndPublish(ctr, nil)
	}

	//Delete reference redis
	if r.redis != nil {
		go r.redis.Del(ctr.Id)
	}
	return nil
}
