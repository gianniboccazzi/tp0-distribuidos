package communication

import (
	"bytes"
	"fmt"
	"net"
	"strconv"
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
	err := b.SendStartBatch()
	if err != nil {
		log.Criticalf("action: send_start_batch | result: fail | client_id: %v | error: %v", b.ClientID, err)
		return
	}
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

		if err := b.SendMessage(batchMessage); err != nil {
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

func (b *BetProtocol) SendStartBatch() error {
	payload := fmt.Sprintf("%s|BETS", b.ClientID)
	header := fmt.Sprintf("%d|", len(payload))
	message := append([]byte(header), []byte(payload)...)
	return b.SendMessage(message)
}




func(b *BetProtocol) SendMessage(message []byte) error {
	bytesSent := 0
	for bytesSent < len(message) {
		n, err := b.Conn.Write(message[bytesSent:])
		if err != nil {
			return err
		}
		bytesSent += n
	}
	return nil
}


func(b *BetProtocol) ReceiveAck() error {
	buffer, err := b.receiveUntilDelimiter()
	if err != nil {
		return err
	}
	delimiterIndex := bytes.Index(buffer, []byte("|"))
	if delimiterIndex == -1 {
		return fmt.Errorf("error finding delimiter")
	}
	bytesToRead := buffer[:delimiterIndex]

	remainingData := buffer[delimiterIndex+1:]


	messageLength, err := strconv.Atoi(string(bytesToRead))
	if err != nil {
		return fmt.Errorf("error parsing message length: %w", err)
	}
	bytesReceived := len(remainingData)
	bufferToRead := make([]byte, messageLength)
	remainingDataToRead, err := b.ReceiveMessage(messageLength - bytesReceived, bufferToRead, bytesReceived)
	if err != nil {
		return err
	}
	remainingData = append(remainingData, remainingDataToRead...)
	remainingDataString := string(remainingData)
	remainingDataString = strings.TrimRight(remainingDataString, "\x00")
	if strings.TrimSpace(remainingDataString) != "ACK" {
		return fmt.Errorf("error receiving ack: %s", remainingDataString)
	}
	return nil
}

func (b *BetProtocol) ReceiveMessage(resLength int, buffer []byte, offset int) ([]byte, error) {
	for offset < resLength {
		n, err := b.Conn.Read(buffer[offset:])
		if err != nil {
			return nil, err
		}
		offset += n
	}
	return buffer, nil
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

func (b *BetProtocol) RequestWinners() {
	err := b.SendRequestWinners()
	if err != nil {
		log.Criticalf("action: send_request_winners | result: fail | client_id: %v | error: %v", b.ClientID, err)
		return
	}
	err = b.ReceiveWinners()
	if err != nil {
		log.Criticalf("action: receive_winners | result: fail | client_id: %v | error: %v", b.ClientID, err)
		return
	}
}

func (b *BetProtocol) SendRequestWinners() error {
	payload := fmt.Sprintf("%s|WINNERS", b.ClientID)
	header := fmt.Sprintf("%d|", len(payload))
	message := append([]byte(header), []byte(payload)...)
	return b.SendMessage(message)
}

func (b *BetProtocol) ReceiveWinners() error {
	buffer, err := b.receiveUntilDelimiter()
	if err != nil {
		return err
	}
	delimiterIndex := bytes.Index(buffer, []byte("|"))
	if delimiterIndex == -1 {
		return fmt.Errorf("error finding delimiter")
	}
	bytesToRead := buffer[:delimiterIndex]

	remainingData := buffer[delimiterIndex+1:]

	messageLength, err := strconv.Atoi(string(bytesToRead))
	if err != nil {
		return fmt.Errorf("error parsing message length: %w", err)
	}
	bytesReceived := len(remainingData)
	bufferToRead := make([]byte, messageLength)
	remainingDataToRead, err := b.ReceiveMessage(messageLength - bytesReceived, bufferToRead, bytesReceived)
	if err != nil {
		return err
	}
	remainingData = append(remainingData, remainingDataToRead...)
	remainingDataString := string(remainingData)
	remainingDataString = strings.TrimRight(remainingDataString, "\x00")
	if strings.TrimSpace(remainingDataString) == "ERR" {
		log.Infof("action: consulta_ganadores | result: fail | client_id: %s | error: el torneo no fue realizado aun", b.ClientID)
		return nil
	}
	if strings.TrimSpace(remainingDataString) == "NONE" {
		log.Infof("action: consulta_ganadores | result: success | cant_ganadores: 0")
		return nil
	}
	winners := strings.Split(remainingDataString, "|")
	log.Infof("action: consulta_ganadores | result: success | cant_ganadores: %d", len(winners))
	return nil	
}




func (b *BetProtocol) receiveUntilDelimiter() ([]byte, error) {
	var buffer bytes.Buffer
	chunkSize := 2

	for {
		chunk := make([]byte, chunkSize)
		n, err := b.Conn.Read(chunk)
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				return nil, fmt.Errorf("timeout waiting for delimiter")
			}
			return nil, fmt.Errorf("error reading from client: %w", err)
		}
		if n == 0 {
			return nil, fmt.Errorf("client disconnected before sending message")
		}

		buffer.Write(chunk[:n]) 

		if strings.Contains(buffer.String(), "|") {
			break
		}
	}

	return buffer.Bytes(), nil
}
