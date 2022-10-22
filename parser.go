package ksef

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"ksef/common"
	"os"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
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
	section                string
	state                  uint8
	headerMap              map[string]map[int]string
	LineHandler            HookFunc
	Comma                  string
	EncodingConversionFile string
	encodingConversion     map[byte]string
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
	line = p.convertEncoding(line)

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

func (p *Parser) convertEncoding(data string) string {
	if p.EncodingConversionFile == "" {
		return data
	}

	if p.encodingConversion == nil {
		p.prepareEncodingConversionTable()
	}

	inputBytes := []byte(data)
	outputBytes := []byte{}
	var dstByte []byte

	for _, srcByte := range inputBytes {
		dstByte = []byte{srcByte}
		if dstChar, exists := p.encodingConversion[srcByte]; exists {
			dstByte = bytes.Replace(dstByte, []byte{srcByte}, []byte(dstChar), 1)
		}
		outputBytes = append(outputBytes, dstByte...)
	}

	return string(outputBytes)

}

func (p *Parser) prepareEncodingConversionTable() {
	fileBytes, err := ioutil.ReadFile(p.EncodingConversionFile)

	if err != nil {
		log.Errorf("error opening encoding conversion file")
		return
	}

	p.encodingConversion = make(map[byte]string)

	for _, line := range strings.Split(string(fileBytes), common.LineBreak) {
		mapping := strings.Split(line, ":")

		if len(mapping) == 2 {
			srcByteHex := strings.Trim(mapping[0], " ")
			dstChar := strings.Trim(mapping[1], " \r")

			srcByte, err := strconv.ParseUint(srcByteHex, 0, 8)
			if err == nil {
				p.encodingConversion[byte(srcByte)] = dstChar
			}
		}
	}

	log.Debugf("successfuly read conversion table file: %v\n", p.encodingConversion)
}
