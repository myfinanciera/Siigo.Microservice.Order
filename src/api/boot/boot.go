// Copyright 2020 Siigo. All rights reserved.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package boot

import (
	"context"

	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"siigo.com/order/src/api/config"
	"siigo.com/order/src/api/controller"
	orderv1 "siigo.com/order/src/api/proto/order/v1"

	"fmt"

	//interceptor "dev.azure.com/SiigoDevOps/Siigo/_git/Siigo.Golang.Security.git/Interceptor"

	"github.com/common-nighthawk/go-figure"
	middleware "github.com/grpc-ecosystem/go-grpc-middleware"

	"io/ioutil"
	"net"
	"os"

	//grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
	grpctrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/google.golang.org/grpc"
	"siigo.com/order/internal/gateway"
)

const (
	// XRequestIDKey is a key for getting request id.
	XRequestIDKey    = "x-request-id"
	unknownRequestID = "<unknown>"
)

// CreateGrpcServer Create server.Server Instance
func CreateGrpcServer(cfg *config.Configuration) *grpc.Server {

	grpcServer := grpc.NewServer(
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             cfg.Grpc.ServerMinTime, // If a client pings more than once every 5 minutes, terminate the connection
			PermitWithoutStream: true,                   // Allow pings even when there are no active streams
		}),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			Time:    cfg.Grpc.ServerTime,    // Ping the client if it is idle for 2 hours to ensure the connection is still active
			Timeout: cfg.Grpc.ServerTimeout, // Wait 20 second for the ping ack before assuming the connection is dead
		}),
		grpc.UnaryInterceptor(
			middleware.ChainUnaryServer(
				// recovery to control panics
				grpc_recovery.UnaryServerInterceptor(),

				// Generate request id
				RequestIDInterceptor(),

				//grpc_auth.UnaryServerInterceptor(interceptor.TokenValidation()),
			),
		),
	)

	return grpcServer
}

func CreateGrpcClient(cfg *config.Configuration) (conn *grpc.ClientConn) {

	// Create the client interceptor using the grpc trace package.
	si := grpctrace.StreamClientInterceptor(grpctrace.WithServiceName("ms-order-client"))
	ui := grpctrace.UnaryClientInterceptor(grpctrace.WithServiceName("ms-order-client"))

	// Create a client connection to the gRPC Server we just started.
	// This is where the gRPC-Gateway proxies the requests.
	conn, _ = grpc.DialContext(
		context.Background(), fmt.Sprintf("dns:///%s:%d", cfg.Grpc.Host, cfg.Grpc.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
		grpc.WithStreamInterceptor(si),
		grpc.WithUnaryInterceptor(ui),
	)

	return conn
}

// RegisterGrpcHandlers Register GRPC Handlers
func RegisterGrpcHandlers(conn *grpc.ClientConn, serverMux *runtime.ServeMux) {

	register := gateway.RegisterHandlers(context.Background(), serverMux, conn)

	register(
		//orderv1.RegisterHealthServiceHandler,
		orderv1.RegisterOrderServiceHandler,
	)
}

// RegisterGrpcServers Register Protobuf Services
func RegisterGrpcServers(server *grpc.Server, controller *controller.Controller) {
	// Register All GrpcServers
	//orderv1.RegisterHealthServiceServer(server, controller)
	orderv1.RegisterOrderServiceServer(server, controller)
}

// StartGrpcServer Start GrpcServers with control errors and grpc logging
func StartGrpcServer(lis net.Listener, grpcServer *grpc.Server) {

	log.Info("Starting GRPC Server ...")

	// Add Reflection Server
	reflection.Register(grpcServer)

	loggerV2 := grpclog.NewLoggerV2(os.Stdout, ioutil.Discard, ioutil.Discard)
	// Serve gRPC Server
	loggerV2.Info("Serving gRPC...")
	// Adds gRPC internal logs. This is quite verbose, so adjust as desired!
	grpclog.SetLoggerV2(loggerV2)

	go func() {
		loggerV2.Fatal(grpcServer.Serve(lis))
	}()

}

func NewNetListener(configuration *config.Configuration) net.Listener {
	addr := fmt.Sprintf("%s:%d", configuration.Grpc.Host, configuration.Grpc.Port)
	lis, err := net.Listen("tcp", addr)

	if err != nil {
		log.Fatalln("Failed to listen:", err)
	}

	return lis
}

// StartHttpServer Start Http Gateway Server
func StartHttpServer(l *log.Logger, configuration *config.Configuration, serverMux *runtime.ServeMux) {

	l.Info("Starting HTTP Server ...")

	goFigure := figure.NewFigure("Siigo Orders.", "", true)
	goFigure.Print()
	gateway.RunServerMux(configuration, serverMux)
}
