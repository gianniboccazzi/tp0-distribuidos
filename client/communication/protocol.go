package communication

import (
	"fmt"
	"net"
	"strings"
	"github.com/op/go-logging"

	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/domain"
)

var log = logging.MustGetLogger("log")
const (
	EOF_DELIMITER = "|||"
)

type BetProtocol struct {
	Conn net.Conn
	MaxAmountOfBets int
	ClientID string
	MaxPacketSize int

}

func NewBetProtocol(conn net.Conn, maxAmountOfBets int, clientID string, maxPacketSize int) *BetProtocol {
	return &BetProtocol{
		Conn: conn,
		MaxAmountOfBets: maxAmountOfBets,
		ClientID: clientID,
		MaxPacketSize: maxPacketSize,
	}
}
func (b *BetProtocol) SendBatches() {
	csvReader, file, err := domain.ReadBetsFile(b.ClientID)
	if err != nil {
		log.Criticalf("action: open_file | result: fail | client_id: %v | error: %v", b.ClientID, err)
		return
	}
	defer file.Close()

	isEOF := false
	var packet, pendingBetPacket []byte
	hasPendingBet := false

	for !isEOF {
		betsQuantity := 0
		batchMaxAmount := b.MaxAmountOfBets
		packet = []byte{} 

		// if there is a pending bet from the previous batch, include it
		if hasPendingBet {
			packet = pendingBetPacket
			pendingBetPacket = nil
			hasPendingBet = false
			betsQuantity++
			batchMaxAmount-- 
		}

		// read bets until the batch is full
		for i := 0; i < batchMaxAmount; i++ {
			bet, err := domain.ReadBet(csvReader)
			if err != nil {
				if strings.Contains(err.Error(), "EOF") {
					isEOF = true
					break
				}
				log.Criticalf("action: read_bet | result: fail | client_id: %s | error: %v", b.ClientID, err)
				return
			}

			betString := PrepareBetToBatchMessage(*bet)
			betPacket := []byte(betString)

			if len(packet)+len(betPacket) > b.MaxPacketSize {
				// if the packet exceeds the maximum size, save the bet for the next batch
				pendingBetPacket = betPacket
				hasPendingBet = true
				break
			}

			packet = append(packet, betPacket...)
			betsQuantity++
		}

		if isEOF {
			packet = append(packet, []byte(EOF_DELIMITER)...)
		}

		batchMessage := PrepareBatchMessage(packet)

		if err := b.SendBatch(batchMessage); err != nil {
			log.Criticalf("action: send_batch | result: fail | client_id: %v | error: %v", b.ClientID, err)
			return
		}

		if err := b.ReceiveAck(); err != nil {
			log.Criticalf("action: receive_ack | result: fail | client_id: %v | error: %v", b.ClientID, err)
			return
		}

		log.Infof("action: apuesta_enviada | result: success | client_id: %s | cantidad: %d", b.ClientID, betsQuantity)
	}
}




func(b *BetProtocol) SendBatch(batchMessage []byte) error {
	bytesSent := 0
	for bytesSent < len(batchMessage) {
		n, err := b.Conn.Write(batchMessage[bytesSent:])
		if err != nil {
			return err
		}
		bytesSent += n
	}
	return nil
}


func(b *BetProtocol) ReceiveAck() error {
	// wait for response
	resLength := 3
	bytesReceived := 0
	bytesAck := make([]byte, resLength)
	for bytesReceived < resLength {
		n, err := b.Conn.Read(bytesAck[bytesReceived:])
		if err != nil {
			return err
		}
		bytesReceived += n
	}
	if string(bytesAck) != "ACK" {
		return fmt.Errorf("ACK not received")
	}
	return nil
}

func PrepareBetToBatchMessage(bet domain.Bet) string {
	/// Protocol: 
	///MESSAGE_LENGTH|AGENCY|NAME|SURNAME|ID|BIRTHDATE|BET_NUMBER

	betMessage := fmt.Sprintf("%d|%s|%s|%d|%s|%d||",
		bet.Agency,
		bet.Name,
		bet.Surname,
		bet.ID,
		bet.BirthDate,
		bet.BetNumber,
	)

	return betMessage
}

func PrepareBatchMessage(message []byte) []byte {
	/// Protocol: 
	///MESSAGE_LENGTH|BATCH
	message = message[:len(message)-2]

	header := fmt.Sprintf("%d|", len(message))

	return append([]byte(header), message...)
}
