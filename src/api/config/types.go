package config

import (
	"time"
)

type Configuration struct {
	AppName        string                 `yaml:"appName"`
	Kafka          map[string]interface{} `yaml:"kafka"`
	Mongo          MongoConfiguration     `yaml:"mongo"`
	Log            LogConfiguration       `yaml:"log"`
	Redis          RedisConfiguration     `yaml:"redis"`
	EventStore     EventStore             `yaml:"eventStore"`
	BusinessLogger BusinessLogger         `yaml:"businessLogger"`
	Grpc           Grpc                   `yaml:"grpc"`
	HttpServer     HttpServer             `yaml:"httpServer"`
}

type LogConfiguration struct {
	Level        string `yaml:"level" `
	ReportCaller bool   `yaml:"reportCaller" `
}

type MongoConfiguration struct {
	ConnectionString string `yaml:"connectionString" `
	Database         string `yaml:"database" `
	Collection       string `yaml:"collection"`
}

type RedisConfiguration struct {
	Addr     string        `yaml:"addr" validate:"required" `
	Password string        `yaml:"password" `
	Db       int           `yaml:"db" `
	Timeout  time.Duration `yaml:"timeOut"`
}

type EventStore struct {
	Url string `yaml:"url"`
}

type BusinessLogger struct {
	Enabled       bool              `yaml:"enabled"`
	LokiUrl       string            `yaml:"lokiUrl"`
	DefaultLabels map[string]string `yaml:"defaultLabels"`
}

type Grpc struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`

	ServerMinTime time.Duration `yaml:"serverMinTime"` // if a client pings more than once every 5 minutes (default), terminate the connection
	ServerTime    time.Duration `yaml:"serverTime"`    // ping the client if it is idle for 2 hours (default) to ensure the connection is still active
	ServerTimeout time.Duration `yaml:"serverTimeout"` // wait 20 second (default) for the ping ack before assuming the connection is dead
	ConnTime      time.Duration `yaml:"connTime"`      // send pings every 10 seconds if there is no activity
	ConnTimeout   time.Duration `yaml:"connTimeout"`   // wait 20 second for ping ack before considering the connection dead
}

type HttpServer struct {
	Port int `yaml:"port"`
}
