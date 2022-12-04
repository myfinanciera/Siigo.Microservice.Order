package main

import (
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	log "github.com/sirupsen/logrus"
	_ "go.uber.org/automaxprocs"
	"go.uber.org/fx"
	"siigo.com/order/src/api/boot"
	"siigo.com/order/src/api/config"
	"siigo.com/order/src/api/controller"
	fxmodule2 "siigo.com/order/src/api/fxmodule"
	"siigo.com/order/src/api/logger"
)

func init() {

}

func main() {
	// start app
	newFxApp().Run()
}

func newFxApp() *fx.App {

	return fx.New(

		// Create Struct FX Providers
		fx.Provide(
			config.NewViperConfig,
			config.NewConfiguration,
			log.New,
		),

		fx.Decorate(logger.NewLogrus),

		fx.Provide(
			boot.CreateGrpcServer,
			boot.CreateGrpcClient,
			boot.NewNetListener,
			runtime.NewServeMux,
			controller.NewController,
		),

		// Load Module in bottom order
		fxmodule2.CQRSDDDModule,
		fxmodule2.CacheModule,
		fxmodule2.BrokerModule,
		fxmodule2.InfrastructureModule,
		fxmodule2.ApplicationModule,

		// Invoke to init functions to start
		fx.Invoke(
			boot.RegisterGrpcServers,
			boot.StartGrpcServer,
			boot.RegisterGrpcHandlers,
			boot.StartHttpServer,
		),
	)
}
