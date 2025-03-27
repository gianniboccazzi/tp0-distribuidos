package common

import (
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/communication"
	"github.com/op/go-logging"
)

const MAX_PACKET_SIZE = 1024 * 8
var log = logging.MustGetLogger("log")

// ClientConfig Configuration used by the client
type ClientConfig struct {
	ID            string
	ServerAddress string
	LoopAmount    int
	LoopPeriod    time.Duration
	BatchMaxAmount int
}

// Client Entity that encapsulates how
type Client struct {
	config ClientConfig
	conn   net.Conn
	endProgram chan bool
}

// NewClient Initializes a new client receiving the configuration
// as a parameter
func NewClient(config ClientConfig) *Client {
	client := &Client{
		config: config,
		endProgram: make(chan bool),
	}
	client.handleSignal()
	return client
}

// CreateClientSocket Initializes client socket. In case of
// failure, error is printed in stdout/stderr and exit 1
// is returned
func (c *Client) createClientSocket() error {
	conn, err := net.Dial("tcp", c.config.ServerAddress)
	if err != nil {
		log.Criticalf(
			"action: connect | result: fail | client_id: %v | error: %v",
			c.config.ID,
			err,
		)
		return err
	}
	c.conn = conn
	c.conn.SetDeadline(time.Now().Add(20 * time.Second))
	return nil
}


func (c *Client) handleSignal() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGTERM)
	go func() {
		select {
			case <-signalChan:
				c.conn.Close()
			case <-c.endProgram:
				return
		}
	}()
}

func (c *Client) SendBets() {
	err := c.createClientSocket() 
	if err != nil {
		log.Criticalf("action: conexion_socket | result: fail | client_id: %d | error: %v", c.config.ID, err)
		return
	}
	defer c.conn.Close()
	betProtocol := communication.NewBetProtocol(c.conn, c.config.BatchMaxAmount, c.config.ID, MAX_PACKET_SIZE)
	betProtocol.SendBatches()
}

func (c *Client) RequestWinners() bool {
	err := c.createClientSocket()
	if err != nil {
		log.Criticalf("action: conexion_socket | result: fail | client_id: %d | error: %v", c.config.ID, err)
		return false
	}
	defer c.conn.Close()
	betProtocol := communication.NewBetProtocol(c.conn, c.config.BatchMaxAmount, c.config.ID, MAX_PACKET_SIZE)
	return betProtocol.RequestWinners()
}


func (c *Client) Run() {
	c.SendBets()
	time.Sleep(2 * time.Second)
	var lotteryFinished bool
	for !lotteryFinished {
		lotteryFinished = c.RequestWinners()
		time.Sleep(2 * time.Second)
	}
}