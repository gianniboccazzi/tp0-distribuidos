package domain

import (
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
		return nil, err
	}
	idParsed, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}
	betNumberParsed, err := strconv.Atoi(betNumber)
	if err != nil {
		return nil, err
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