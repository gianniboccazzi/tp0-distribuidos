package communication

import (
	"fmt"

	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/domain"
)

const (
	BET uint8 = 1
)

func PrepareBetMessage(bet domain.Bet) string {
	/// Protocol: 
	///MESSAGE_LENGTH|AGENCY|NAME|SURNAME|ID|BIRTHDATE|BET_NUMBER

	payload := fmt.Sprintf("%d|%s|%s|%d|%s|%d",
		bet.Agency,
		bet.Name,
		bet.Surname,
		bet.ID,
		bet.BirthDate,
		bet.BetNumber,
	)

	header := fmt.Sprintf("%d|", len(payload))

	return header + payload


	
}