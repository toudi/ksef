package vat

import (
	"bytes"
	"encoding/xml"
	"ksef/internal/money"
	"math"
	"strconv"

	"github.com/beevik/etree"
)

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

func (c *Calculator) GetVat(itemNode *etree.Element) (*VatInfo, error) {
	var err error
	doc := etree.NewDocument()
	doc.SetRoot(itemNode)
	c.buffer.Reset()
	if _, err = doc.WriteTo(&c.buffer); err != nil {
		return nil, err
	}
	invoiceItem := &InvoiceItem{}
	if err = xml.Unmarshal(c.buffer.Bytes(), invoiceItem); err != nil {
		return nil, err
	}

	vatMultiplier := 100
	var vatRateNumber int
	vatRateNumber, err = strconv.Atoi(invoiceItem.TaxRate)

	var netAmount money.MonetaryValue
	var grossAmount money.MonetaryValue

	if err != nil {
		// it just means that the vat is not a number - could be "zw", "np" and so on.
		amountStr := invoiceItem.NetAmount
		if amountStr == nil {
			amountStr = invoiceItem.GrossAmount
		}
		if err = netAmount.LoadFromString(*amountStr); err != nil {
			return nil, err
		}

		return &VatInfo{
			NetAmount: netAmount,
			VatAmount: money.MonetaryValue{},
			VatRate:   invoiceItem.TaxRate,
		}, nil
	}

	vatMultiplier += vatRateNumber
	if invoiceItem.NetAmount != nil {
		// we calculate stuff from net to gross
		if err = netAmount.LoadFromString(*invoiceItem.NetAmount); err != nil {
			return nil, err
		}
		netAmount = netAmount.ToDecimalPlaces(2)
		grossAmount = money.MonetaryValue{
			DecimalPlaces: netAmount.DecimalPlaces,
			Amount:        int(math.Round(float64(netAmount.Amount*vatMultiplier) / 100)),
		}
	} else {
		// we calculate stuff from gross to net
		grossAmount.LoadFromString(*invoiceItem.GrossAmount)
		netAmount = money.MonetaryValue{
			DecimalPlaces: grossAmount.DecimalPlaces,
			Amount:        int(math.Round(float64(grossAmount.Amount/vatMultiplier) * 100)),
		}.ToDecimalPlaces(2)
	}

	vatAmount := money.MonetaryValue{
		DecimalPlaces: netAmount.DecimalPlaces,
		Amount:        grossAmount.Amount - netAmount.Amount,
	}

	return &VatInfo{
		NetAmount: netAmount,
		VatAmount: vatAmount,
		VatRate:   invoiceItem.TaxRate,
	}, nil
}
