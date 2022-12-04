//go:build e2e
// +build e2e

package e2e

import (
	"bytes"
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"siigo.com/order/it/containers"
	"siigo.com/order/src/api/boot"
	config2 "siigo.com/order/src/api/config"
	"siigo.com/order/src/api/controller"
	fxmodule2 "siigo.com/order/src/api/fxmodule"
	"sync"
	"testing"
)

var once sync.Once

func init() {
	once.Do(func() {
		// move to root path
		os.Chdir("../../")
	})
}

type ApiTestSuite struct {
	suite.Suite
	grpcServer     *grpc.Server
	mux            *runtime.ServeMux
	app            *fx.App
	mongoContainer testcontainers.Container
	httpListener   net.Listener
	port           int
}

// run api suite test
func TestE2EApiSuite(t *testing.T) {
	suite.Run(t, new(ApiTestSuite))
}

// test create order endpoint with bad request
func (suite *ApiTestSuite) TestCreateOrderWithEmptyRequest() {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", fmt.Sprintf("http://localhost:%v/api/v1/order", suite.port), bytes.NewReader(nil))

	suite.mux.ServeHTTP(w, r)

	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
}

// test health endpoint
func (suite *ApiTestSuite) TestHealth() {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", fmt.Sprintf("http://localhost:%v/api/health", suite.port), bytes.NewReader(nil))

	suite.mux.ServeHTTP(w, r)

	assert.Equal(suite.T(), w.Code, http.StatusOK)
}

// setup grpc and http server
func (suite *ApiTestSuite) SetupTest() {

	// Start mongo container
	mongoC, mongoConn, err := containers.StartMongoContainer(context.Background(), containers.ContainerOptions{})
	if err != nil {
		suite.Errorf(err, "Failed to start mongoDB container: %v", err)
	}

	// build config
	conf := config2.NewConfiguration()

	// set mongo connection string
	mongoURI := mongoConn.ConnectionURI()
	conf.Mongo.ConnectionString = mongoURI
	conf.Mongo.Database = "e2e-database"
	conf.Mongo.Collection = "e2e-collection"

	// build a server mux
	serveMux := runtime.NewServeMux()
	grpcServer := boot.CreateGrpcServer(conf)

	// config to connect to a random available port
	conf.Grpc.Port = 0
	conf.Grpc.Host = ""
	listener := boot.NewNetListener(conf)

	randomGrpcPort := listener.Addr().(*net.TCPAddr).Port
	randomHttpPort := randomGrpcPort + 1
	conf.Grpc = config2.Grpc{
		Port: randomGrpcPort,
	}
	conf.HttpServer.Port = randomHttpPort

	app := fx.New(
		// Create Struct FX Providers
		fx.Provide(
			func() *config2.Configuration { return conf },
			logrus.New,
		),

		fx.Provide(
			func() net.Listener { return listener },
			func() *grpc.Server { return grpcServer },
			boot.CreateGrpcClient,
			func() *runtime.ServeMux { return serveMux },
			controller.NewController,
		),

		// Get Module in bottom order
		fxmodule2.CQRSDDDModule,
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

	if err := app.Start(context.Background()); err != nil {
		suite.Error(err)
	}

	suite.app = app
	suite.mongoContainer = mongoC
	suite.mux = serveMux
	suite.grpcServer = grpcServer
	suite.httpListener = listener
	suite.port = randomHttpPort

}

// shutdown servers and containers
func (suite *ApiTestSuite) TearDownTest() {
	suite.mongoContainer.Terminate(context.Background())
	//suite.httpListener.Close()
	//suite.grpcServer.GracefulStop()
}
