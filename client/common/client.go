package common

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/communication"
	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/domain"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("log")

// ClientConfig Configuration used by the client
type ClientConfig struct {
	ID            string
	ServerAddress string
	LoopAmount    int
	LoopPeriod    time.Duration
}

// Client Entity that encapsulates how
type Client struct {
	config ClientConfig
	conn   net.Conn
	active bool
}

// NewClient Initializes a new client receiving the configuration
// as a parameter
func NewClient(config ClientConfig) *Client {
	client := &Client{
		config: config,
	}
	client.active = true
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
	c.conn.SetDeadline(time.Now().Add(5 * time.Second))
	return nil
}

// StartClientLoop Send messages to the client until some time threshold is met
func (c *Client) StartClientLoop() {
	// There is an autoincremental msgID to identify every message sent
	// Messages if the message amount threshold has not been surpassed
	for msgID := 1; msgID <= c.config.LoopAmount && c.active; msgID++ {
		// Create the connection the server in every loop iteration. Send an
		c.createClientSocket()

		// TODO: Modify the send to avoid short-write
		fmt.Fprintf(
			c.conn,
			"[CLIENT %v] Message N°%v\n",
			c.config.ID,
			msgID,
		)
		msg, err := bufio.NewReader(c.conn).ReadString('\n')
		c.conn.Close()

		if err != nil {
			log.Errorf("action: receive_message | result: fail | client_id: %v | error: %v",
				c.config.ID,
				err,
			)
			return
		}

		log.Infof("action: receive_message | result: success | client_id: %v | msg: %v",
			c.config.ID,
			msg,
		)

		// Wait a time between sending one message and the next one
		time.Sleep(c.config.LoopPeriod)

	}
	log.Infof("action: loop_finished | result: success | client_id: %v", c.config.ID)
}

func (c *Client) handleSignal() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGTERM)
	go func() {
		<-signalChan
		c.active = false
	}()
}

/// SendBet sends a bet to the server and waits for the ACK
func (c *Client) SendBet(bet *domain.Bet) {
	err := c.createClientSocket() 
	if err != nil {
		log.Criticalf("action: conexion_socket | result: fail | dni: %d | error: %v", bet.ID, err)
		return
	}
	message := communication.PrepareBetMessage(*bet)
	bytesMessage := []byte(message)
	bytesSent := 0
	// avoid short-write
	for bytesSent < len(bytesMessage) {
		n, err := c.conn.Write(bytesMessage[bytesSent:])
		if err != nil {
			log.Criticalf("action: envio_apuesta | result: fail | dni: %d | error: %v", bet.ID, err)
			c.conn.Close()
			return
		}
		bytesSent += n
	}
	
	// wait for ACK
	ackLength := 4
	bytesReceived := 0
	bytesAck := make([]byte, ackLength)
	for bytesReceived < ackLength {
		n, err := c.conn.Read(bytesAck[bytesReceived:])
		if err != nil {
			log.Errorf("action: recibo_ack | result: fail | dni: %d | error: %v", bet.ID, err)
			c.conn.Close()
			return
		}
		bytesReceived += n
	}
	
	// Check ACK content
	if strings.TrimSpace(string(bytesAck)) != "ACK" {
		log.Errorf("action: recibo_ack | result: fail | dni: %d | error: ACK not received", bet.ID)
		return
	}

	log.Infof("action: apuesta_enviada | result: success | dni: %d | numero: %d", bet.ID, bet.BetNumber)
	c.conn.Close()
}