package interceptor

import (
	"dev.azure.com/SiigoDevOps/Siigo/_git/go-cqrs.git/cqrs"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/golang/protobuf/proto"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var validate = validator.New()

func NewInMemoryDispatcherWithInterceptors() cqrs.Dispatcher {
	return cqrs.NewInMemoryDispatcher(DispatcherValidationInterceptor)
}

// DispatcherValidationInterceptor validate commands and queries
func DispatcherValidationInterceptor(message cqrs.RequestMessage) error {

	// struct validations by tags https://github.com/go-playground/validator
	err := validate.Struct(message.Request())

	// validation ok
	if err == nil {
		return nil
	}

	// map errors to grpc-errors
	var ve validator.ValidationErrors

	if errors.As(err, &ve) {
		out := make([]proto.Message, len(ve))
		for i, fe := range ve {
			out[i] = &errdetails.ErrorInfo{Domain: "Api", Reason: fmt.Sprintf("field: %s. %s", fe.Field(), msgForTag(fe))}
		}

		stat, e := status.
			New(codes.InvalidArgument, "Invalid request").
			WithDetails(out...)

		if e != nil {
			return status.Errorf(codes.Internal, "unexpected error adding details: %s", e)
		}

		return stat.Err()
	}

	return nil
}

func msgForTag(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email"
	}
	return fe.Error() // default error
}
