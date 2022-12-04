//go:build e2e
// +build e2e

package e2e

import (
	"bytes"
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"net/http/httptest"
	"siigo.com/order/src/api/controller"
	orderv1 "siigo.com/order/src/api/proto/order/v1"
	"testing"
)

type HealthTestSuite struct {
	suite.Suite
	server *grpc.Server
	mux    *runtime.ServeMux
}

// run health suite test
func TestE2EHealthSuite(t *testing.T) {
	suite.Run(t, new(HealthTestSuite))
}

// test health endpoint
func (suite *HealthTestSuite) TestHealth() {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", fmt.Sprintf("http://localhost:%v/api/health", 11000), bytes.NewReader(nil))

	suite.mux.ServeHTTP(w, r)

	assert.Equal(suite.T(), w.Code, http.StatusOK)
}

// setup grpc and http server
func (suite *HealthTestSuite) SetupTest() {

	ctx := context.Background()

	mux := runtime.NewServeMux()
	server := grpc.NewServer()

	addr := "[::]:10001"
	lis, err := net.Listen("tcp", addr)

	conn, err := grpc.DialContext(ctx, addr, grpc.WithInsecure())
	if err != nil {
		suite.Error(err)
		suite.FailNowf("error DialContext %s", err.Error())
	}

	controller := &controller.Controller{}

	orderv1.RegisterHealthServiceServer(server, controller)

	err = orderv1.RegisterHealthServiceHandler(ctx, mux, conn)
	if err != nil {
		suite.Error(err)
		suite.FailNow(err.Error())
	}

	go func() {
		err := server.Serve(lis)
		if err != nil {
			suite.Error(err)
		}
	}()

	suite.server = server
	suite.mux = mux

}

// shutdown servers
func (suite *HealthTestSuite) TearDownTest() {
	suite.server.GracefulStop()
}
