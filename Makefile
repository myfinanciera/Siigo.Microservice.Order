all: generate install clean  build run

clean:
	@echo "  >  Cleaning Project ..."
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go clean

build:
	@echo "  >  Building binary..."
	mockery --all --recursive --keeptree
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) GOPRIVATE=dev.azure.com go build -buildvcs=false

run:
	@echo "  >  Running Siigo Bolt..."
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go install github.com/cosmtrek/air@latest
	air || echo "you must configure GOBIN and GOPATH environment variables. please visit https://www.programming-books.io/essential/go/gopath-goroot-gobin-d6da4b8481f94757bae43be1fdfa9e73"; exit 1 

generate:
	buf format -o src/api/proto
	buf mod update -v
	buf lint
	buf generate
	go install github.com/favadi/protoc-go-inject-tag@latest
	go install github.com/vektra/mockery/v2@latest
	protoc-go-inject-tag -input="./src/api/proto/order/v1/order.pb.go"

BUF_VERSION:=1.3.1

install:
	GOPRIVATE=dev.azure.com go get ./...
	go mod download

load_test_grpc_get:
	ghz --call=services.OrderService/LoadOrder -d '{"id":"ab5bfecb-ff6a-43c1-b0c2-78c7ef4da600"}' \
	--skipFirst=0 \
  	--insecure -c 40 -z 10s \
  	--proto src/api/proto/order/v1/order.proto \
  	-i third_party localhost:10000

load_test_grpc_post:
	ghz --call=services.OrderService/AddOrder -d '{"cost": 9,"activated": true,"address": "rsdhrtfdh"}' \
	--skipFirst=0 \
  	--insecure -c 40 -z 10s \
  	--proto proto/services.proto \
  	-i proto localhost:10000

load_test_http_get:
	siege  -c200 -t15S \
	--content-type "application/json" 'http://localhost:11000/api/v1/orders'

load_test_http_post:
	siege -c120 -t15S \
	--content-type "application/json" 'http://localhost:11000/api/v1/order POST {"cost": 9,"activated": true,"address": "rsdhrtfdh"}'

test.unit:
	@echo "  >  Running Tests ..."
	go install github.com/vektra/mockery/v2@latest
	mockery --recursive --keeptree --all
	go test ./... -test.v; echo "Unit Test finished."

test.it:
	go test -tags=integration ./it -v

test.bench:
	go test -test.bench=. -tags=bench ./it/benchmark

test.e2e:
	go test -test.e2e=. -tags=e2e ./it/e2e

test.all: test.unit test.it test.bench


