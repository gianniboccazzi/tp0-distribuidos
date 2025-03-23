package domain

import (
	"fmt"
	"os"
	"strconv"
)

type Bet struct {
	Agency int
	Name string
	Surname string
	ID int
	BirthDate string
	BetNumber int
}

func LoadBet() (*Bet, error) {
	agency := os.Getenv("CLI_ID")
	name := os.Getenv("NOMBRE")
	surname := os.Getenv("APELLIDO")
	id := os.Getenv("DOCUMENTO")
	birthDate := os.Getenv("NACIMIENTO")
	betNumber := os.Getenv("NUMERO")
	agencyParsed, err := strconv.Atoi(agency)
	if err != nil {
		return nil, fmt.Errorf("error parsing agency: %v", err)
	}
	idParsed, err := strconv.Atoi(id)
	if err != nil {
		return nil, fmt.Errorf("error parsing id: %v", err)
	}
	betNumberParsed, err := strconv.Atoi(betNumber)
	if err != nil {
		return nil, fmt.Errorf("error parsing bet number: %v", err)
	}
	return &Bet{
		Agency: agencyParsed,
		Name: name,
		Surname: surname,
		ID: idParsed,
		BirthDate: birthDate,
		BetNumber: betNumberParsed,
	}, nil
}