package domain

import (
	"encoding/csv"
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

func ReadBetsFile(client_id string) (*csv.Reader, *os.File, error) {
	file, err := os.Open("./.data/agency-" + client_id + ".csv")
	if err != nil {
		return nil, nil, fmt.Errorf("error opening file: %v", err)
	}
	csvReader := csv.NewReader(file)
	return csvReader, file, nil
}

func ReadBet(csvReader *csv.Reader) (*Bet, error) {
	agency := os.Getenv("CLI_ID")
	agencyParsed, err := strconv.Atoi(agency)
	if err != nil {
		return nil, fmt.Errorf("error parsing agency: %v", err)
	}
	record, err := csvReader.Read()
	if err != nil {
		return nil, fmt.Errorf("error reading record: %v", err)
	}
	id, err := strconv.Atoi(record[2])
	if err != nil {
		return nil, fmt.Errorf("error parsing id: %v", err)
	}
	betNumber, err := strconv.Atoi(record[4])
	if err != nil {
		return nil, fmt.Errorf("error parsing bet number: %v", err)
	}
	return &Bet{
		Agency: agencyParsed,
		Name: record[0],
		Surname: record[1],
		ID: id,
		BirthDate: record[3],
		BetNumber: betNumber,
	}, nil
}