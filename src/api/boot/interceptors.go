package boot

import (
	"context"
	"dev.azure.com/SiigoDevOps/Siigo/_git/go-cqrs.git/cqrs/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// RequestIDInterceptor is a interceptor of access control list.
func RequestIDInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		requestID := requestIDFromContext(ctx)

		header := metadata.Pairs(XRequestIDKey, requestID)
		grpc.SendHeader(ctx, header)

		ctx = context.WithValue(ctx, XRequestIDKey, requestID)
		return handler(ctx, req)
	}
}

//
func requestIDFromContext(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return unknownRequestID
	}

	header, ok := md[XRequestIDKey]
	if !ok || len(header) == 0 {
		//  generate request id if not exist
		return uuid.NewUUID()
	}

	requestID := header[0]
	if requestID == "" {
		return unknownRequestID
	}

	return requestID
}
