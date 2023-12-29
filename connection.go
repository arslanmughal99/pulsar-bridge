package main

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

var (
	_             = godotenv.Load()
	brokers       = os.Getenv("PULSAR_BROKERS")
	opTimeout     = os.Getenv("PULSAR_OPERATION_TIMEOUT_SEC")
	clientTimeout = os.Getenv("PULSAR_CLIENT_TIMEOUT_SEC")
)

// Connection instance
type Connection struct {
	//lock *sync.Cond
	sync.Mutex
	con pulsar.Client
}

// NewConnection create a new connection instance
func NewConnection() *Connection {
	conn := new(Connection)
	//conn.lock = sync.NewCond(&sync.Mutex{})
	conn.Connect()
	return conn
}

// Get retrieve pulsar client instance
func (c *Connection) Get() pulsar.Client {
	return c.con
}

// Close drops pulsar connection
func (c *Connection) Close() {
	c.con.Close()
}

// Connect connect pulsar client
func (c *Connection) Connect() {
	c.Lock()
	defer c.Unlock()
	clientOpt := pulsar.ClientOptions{
		ConnectionMaxIdleTime:   -1,
		MaxConnectionsPerBroker: 50,
		EnableTransaction:       true,
		KeepAliveInterval:       time.Second * 5,
		URL:                     fmt.Sprintf("pulsar://%s", brokers),
	}
	if operationTimeout, err := strconv.
		ParseInt(opTimeout, 10, 64); err != nil {
		clientOpt.OperationTimeout = time.Duration(operationTimeout) * time.Second
	}
	if cTimeout, err := strconv.
		ParseInt(clientTimeout, 10, 64); err != nil {
		clientOpt.ConnectionTimeout = time.Duration(cTimeout) * time.Second
	}

	client, err := pulsar.NewClient(clientOpt)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create pulsar client")
	}

	c.con = client
}
