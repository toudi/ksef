package vat

import (
	"bytes"
	"encoding/xml"
	"errors"
	"ksef/internal/money"
	"math"
	"strconv"

	"github.com/beevik/etree"
)

var ErrAmountNotDefined = errors.New("amount not defined")

type InvoiceItem struct {
	TaxRate     string  `xml:"P_12"`
	NetAmount   *string `xml:"P_11"`
	GrossAmount *string `xml:"P_11A"`
}

type Calculator struct {
	buffer bytes.Buffer
}

type VatInfo struct {
	NetAmount money.MonetaryValue
	VatAmount money.MonetaryValue
	VatRate   string
}

func calculateVat(item *InvoiceItem) (*VatInfo, error) {
	var err error
	var netAmount money.MonetaryValue
	var grossAmount money.MonetaryValue
	var vatRateNumber int

	vatMultiplier := 100
	vatRateNumber, err = strconv.Atoi(item.TaxRate)
	if err != nil {
		// so basically this means that vat rate is not a number (could be NP / zw / etc)
		// which means that net = gross
		amountStr := item.NetAmount
		if amountStr == nil {
			amountStr = item.GrossAmount
		}
		if amountStr == nil {
			return nil, ErrAmountNotDefined
		}
		if err = netAmount.LoadFromString(*amountStr); err != nil {
			return nil, err
		}
		netAmount = netAmount.ToDecimalPlaces(2)

		return &VatInfo{
			NetAmount: netAmount,
			// VatAmount: money.MonetaryValue{DecimalPlaces: 2},
			VatRate: item.TaxRate,
		}, nil
	}

	vatMultiplier += vatRateNumber
	if item.NetAmount == nil && item.GrossAmount == nil {
		return nil, ErrAmountNotDefined
	}

	if item.NetAmount != nil {
		if err = netAmount.LoadFromString(*item.NetAmount); err != nil {
			return nil, err
		}
		netAmount = netAmount.ToDecimalPlaces(2)
		grossAmount = money.MonetaryValue{
			DecimalPlaces: 2,
			Amount:        int(math.Round(float64(netAmount.Amount*vatMultiplier) / 100)),
		}
	} else if item.GrossAmount != nil {
		if err = grossAmount.LoadFromString(*item.GrossAmount); err != nil {
			return nil, err
		}
		grossAmount = grossAmount.ToDecimalPlaces(2)
		netAmount = money.MonetaryValue{
			DecimalPlaces: 2,
			Amount:        int(math.Round(float64(grossAmount.Amount) / float64(vatMultiplier) * 100)),
		}
	}

	if netAmount.DecimalPlaces == 0 && grossAmount.DecimalPlaces == 0 {
		netAmount.DecimalPlaces = 2
		grossAmount.DecimalPlaces = 2
	}

	vatAmount := money.MonetaryValue{
		DecimalPlaces: 2,
		Amount:        grossAmount.Amount - netAmount.Amount,
	}

	return &VatInfo{
		NetAmount: netAmount,
		VatAmount: vatAmount,
		VatRate:   item.TaxRate,
	}, nil
}

// so basically GetVat is a function that operates on XML nodes, whereas
// the inner one (calculateVat) operates on pure go structs for easier testing.
func (c *Calculator) GetVat(itemNode *etree.Element) (*VatInfo, error) {
	doc := etree.NewDocument()
	doc.SetRoot(itemNode)
	c.buffer.Reset()
	if _, err := doc.WriteTo(&c.buffer); err != nil {
		return nil, err
	}
	invoiceItem := &InvoiceItem{}
	if err := xml.Unmarshal(c.buffer.Bytes(), invoiceItem); err != nil {
		return nil, err
	}

	return calculateVat(invoiceItem)
}
