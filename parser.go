package ksef

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"strings"
)

const (
	stateParseHeaders = iota
	stateParseData
)

const (
	lineWithHeaders = iota
	lineWithData    = iota
)

var arrayNodes = map[string]bool{"Faktura.Fa.FaWiersze.FaWiersz": true}

type HookFunc func(string, map[string]string) error

const section string = "sekcja"

type Parser struct {
	section     string
	state       uint8
	headerMap   map[string]map[int]string
	LineHandler HookFunc
	Comma       string
}

func (p *Parser) Parse(input *os.File) error {
	scanner := bufio.NewScanner(input)
	p.headerMap = make(map[string]map[int]string)
	for scanner.Scan() {
		p.parseLine(scanner.Text())
	}
	return nil
}

func (p *Parser) parseLine(line string) error {
	var err error

	csvReader := csv.NewReader(strings.NewReader(line))
	csvReader.Comma = rune(p.Comma[0])

	fields, err := csvReader.Read()
	if err != nil {
		return fmt.Errorf("cannot read line from csv: %v", err)
	}

	if len(fields) > 0 && strings.ToLower(fields[0]) == "-- ignore --" {
		return nil
	}

	if len(fields) > 1 && strings.ToLower(fields[0]) == section {
		p.state = stateParseHeaders
		p.section = fields[1]
	} else if p.state == stateParseHeaders {
		p.state = stateParseData
		// if _, exists := p.headerMap[p.section]; !exists {
		// na wszelki wypadek parsujemy nagłówki ponownie ponieważ w kolejnych fakturach
		// kolejność lub ilość pól może być zamieniona
		p.headerMap[p.section] = make(map[int]string)
		for idx, header := range fields {
			if header == "" {
				continue
			}
			p.headerMap[p.section][idx] = header
		}
		// }
	} else if p.state == stateParseData {
		var header string
		var emptyLine bool = true

		data := make(map[string]string)

		for idx, value := range fields {
			header = p.headerMap[p.section][idx]
			if header == "" || value == "" {
				continue
			}

			emptyLine = false

			data[header] = value
		}

		if !emptyLine {
			p.LineHandler(p.section, data)
		}
	}

	return nil
}
