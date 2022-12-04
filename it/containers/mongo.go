package containers

import (
	"context"
	"fmt"
	"github.com/imdario/mergo"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var mongoMux = new(sync.Mutex)

// ContainerOptions ...
type ContainerOptionsBase struct {
	testcontainers.ContainerRequest
	CollectLogs    bool
	StartupTimeout time.Duration
}

// ContainerOptions ...
type ContainerOptions struct {
	ContainerOptionsBase
	User     string
	Password string
}

// DBConfig ...
type DBConfig struct {
	Host     string
	Port     uint
	User     string
	Password string
}

// MergeRequest ...
func MergeRequest(c *testcontainers.ContainerRequest, override *testcontainers.ContainerRequest) {
	if err := mergo.Merge(c, override, mergo.WithOverride); err != nil {
		panic(err)
	}
}

// ConnectionURI ...
func (c DBConfig) ConnectionURI() string {
	var databaseAuth string
	if c.User != "" && c.Password != "" {
		databaseAuth = fmt.Sprintf("%s:%s@", c.User, c.Password)
	}
	databaseHost := fmt.Sprintf("%s:%d", c.Host, c.Port)
	return fmt.Sprintf("mongodb://%s%s/?connect=direct", databaseAuth, databaseHost)
}

const defaultMongoDBPort = 27017

// StartMongoContainer ...
func StartMongoContainer(ctx context.Context, options ContainerOptions) (mongoC testcontainers.Container, Config DBConfig, err error) {
	mongoMux.Lock()
	defer mongoMux.Unlock()
	mongoPort, _ := nat.NewPort("", strconv.Itoa(defaultMongoDBPort))

	env := make(map[string]string)
	if options.User != "" && options.Password != "" {
		env["MONGO_INITDB_ROOT_USERNAME"] = options.User
		env["MONGO_INITDB_ROOT_PASSWORD"] = options.Password
	}

	timeout := options.StartupTimeout
	if int64(timeout) < 1 {
		timeout = 5 * time.Minute // Default timeout
	}

	req := testcontainers.ContainerRequest{
		Image:        "mongo:4.4.3",
		Env:          env,
		ExposedPorts: []string{string(mongoPort)},
		WaitingFor:   wait.ForLog("Waiting for connections").WithStartupTimeout(timeout),
	}

	MergeRequest(&req, &options.ContainerRequest)

	mongoC, err = testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	if err != nil {
		err = fmt.Errorf("Failed to start mongo container: %v", err)
		return
	}

	host, err := mongoC.Host(ctx)
	if err != nil {
		err = fmt.Errorf("Failed to get mongo container host: %v", err)
		return
	}

	port, err := mongoC.MappedPort(ctx, mongoPort)
	if err != nil {
		err = fmt.Errorf("Failed to get exposed mongo container port: %v", err)
		return
	}

	Config = DBConfig{
		Host:     host,
		Port:     uint(port.Int()),
		User:     options.User,
		Password: options.Password,
	}

	return
}

// BuildMongoInfrastructure create a docker container with mongodb
func BuildMongoInfrastructure(t *testing.T) (*mongo.Collection, testcontainers.Container) {
	// Start mongo container
	mongoC, mongoConn, err := StartMongoContainer(context.Background(), ContainerOptions{})
	if err != nil {
		t.Fatalf("Failed to start mongoDB container: %v", err)
	}

	// Connect to the database
	mongoURI := mongoConn.ConnectionURI()
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
	if err != nil {
		t.Fatalf("Failed to create mongo client (%s): %v", mongoURI, err)
	}
	mctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	client.Connect(mctx)

	err = client.Ping(mctx, readpref.Primary())
	if err != nil {
		t.Fatalf("Could not ping database within %d seconds (%s): %v", 20, mongoURI, err)
	}
	database := client.Database("testdatabase")
	collection := database.Collection("order-collection")

	return collection, mongoC
}
