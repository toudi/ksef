package parser

import (
	"fmt"
	"ksef/internal/invoice"
	"ksef/internal/money"
	"strconv"
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

type HookFunc func(*invoice.Invoice, string, map[string]string, func() error) error

const section string = "sekcja"

type Parser struct {
	section          string
	state            uint8
	headerMap        map[string]map[int]string
	LineHandler      HookFunc
	InvoiceReadyFunc func(invoice *invoice.Invoice) error
	invoice          *invoice.Invoice
}

func (p *Parser) ProcessLine(fields []string) error {
	var err error
	if p.invoice == nil {
		p.invoice = &invoice.Invoice{}
		p.invoice.Clear()
	}
	if p.headerMap == nil {
		p.headerMap = make(map[string]map[int]string)
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

			data[header] = strings.TrimSpace(value)
		}

		if !emptyLine {
			// check if we have to parse .decimal-places cells
			for header, value := range data {
				decimalPlacesTest := strings.Split(header, ".decimal-places")
				if len(decimalPlacesTest) == 1 {
					// nope, this is just a regular header
					continue
				}
				// we have a match. we now have to re-parse the value according
				// to decimal places.
				tmpStruct := map[string]interface{}{
					"value":          data[decimalPlacesTest[0]],
					"decimal-places": value,
				}
				if valueAsInt, err := strconv.Atoi(data[decimalPlacesTest[0]]); err == nil {
					tmpStruct["value"] = valueAsInt
				}
				parsedMonetaryValue, success, err := money.ParseMonetaryValue(tmpStruct)
				if !success || err != nil {
					return fmt.Errorf("error during ParseMonetaryValue: %v", err)
				}
				data[decimalPlacesTest[0]] = parsedMonetaryValue
				delete(data, header)
			}
			// check if we need to create new invoice
			if err = p.LineHandler(p.invoice, p.section, data, p.InvoiceReady); err != nil {
				return err
			}
		}
	}

	return nil
}

func (p *Parser) InvoiceReady() error {
	if p.InvoiceReadyFunc != nil {
		return p.InvoiceReadyFunc(p.invoice)
	}

	return nil
}

func (p *Parser) ClearInvoiceData() error {
	if len(p.invoice.Items) > 0 {
		if err := p.InvoiceReady(); err != nil {
			return err
		}
	}
	p.invoice.Clear()
	return nil
}
