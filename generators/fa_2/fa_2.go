package fa_2

import (
	"fmt"
	"ksef/common/xml"
	"os"
	"strconv"
	"time"
)

func FA_2(invoice *xml.Node, dest string) error {
	var err error
	var numItems int
	var totalAmountNet int
	var totalAmountGross int
	var totalVatAmount int

	var totalNetAmountPerVATRate map[string]int = make(map[string]int)
	var totalVatAmountPerVATRate map[string]int = make(map[string]int)

	var net, gross, vat int
	var vatRate string
	var invoiceIsBasedOnGross bool
	var vatRateDivisor float32

	// set common data, this can be later overriden by the data from parser.
	invoice.SetData("Faktura.Fa.Adnotacje", map[string]string{
		"P_16":     "2",
		"P_17":     "2",
		"P_18":     "2",
		"P_18A":    "2",
		"P_19":     "2",
		"P_22":     "2",
		"P_23":     "2",
		"P_PMarzy": "2",
	})
	invoice.SetValue("Faktura.Naglowek.DataWytworzeniaFa", time.Now().Format(time.RFC3339))

	var totalNetAmountFieldPerNetAmountMapping = map[string][]string{
		"P_13_1": {"22", "23"},
		"P_13_2": {"8", "7"},
		"P_13_3": {"5"},
		"P_13_6": {"0"},
		"P_13_7": {"zw"},
	}
	var totalVATAmountFieldPerVATRatesMapping = map[string][]string{}

	itemsNode, err := invoice.LocateNode("Faktura.Fa.FaWiersze")
	if err != nil {
		return fmt.Errorf("cannot locaate row with items: %v", err)
	}

	for _, child := range itemsNode.Children {
		if child.Name == "FaWiersz" {
			child.SetValue("NrWierszaFa", fmt.Sprintf("%d", numItems+1))
			numItems += 1

			if vatRate, err = child.ValueOf("P_12"); err != nil {
				return fmt.Errorf("cannot parse VAT rate: %v", err)
			}
			vatRateDivisor = -1

			if vatRateNumber, err := strconv.Atoi(vatRate); err == nil {
				vatRateDivisor = 1 + float32(vatRateNumber)/100
			}

			if net, err = getInt(child.ValueOf("P_11")); err != nil {
				// no net amount. is there a gross amount?
				if gross, err = getInt(child.ValueOf("P_11A")); err != nil {
					return fmt.Errorf("each invoice item has to have either net or gross amount")
				}
				invoiceIsBasedOnGross = true
			}

			if net > 0 && gross > 0 {
				return fmt.Errorf("invoice cannot have both net and gross amount (P_11/P11_A). Only one field has to be specified")
			}

			if invoiceIsBasedOnGross {
				child.DeleteChild("P_9A")
				net = gross
				if vatRateDivisor > -1 {
					net = int(float32(gross) / vatRateDivisor)
				}
			} else {
				child.DeleteChild("P_11A")

				gross = net
				if vatRateDivisor > -1 {
					gross = int(float32(net) * vatRateDivisor)
				}
			}

			vat = gross - net
			totalVatAmount += vat
			totalAmountNet += net
			totalAmountGross += gross
			totalVatAmountPerVATRate[vatRate] += vat
			totalNetAmountPerVATRate[vatRate] += net

			fmt.Printf("net=%v; vat=%v; gross=%v\n", net, vat, gross)
		}
	}

	fmt.Printf("totalVatAmount=%v; totalGrossAmount=%v\n", totalVatAmount, totalAmountGross)

	invoice.SetValue("Faktura.Fa.P_14_1", fmt.Sprintf("%.2f", float32(totalVatAmount)/100))
	invoice.SetValue("Faktura.Fa.P_15", fmt.Sprintf("%.2f", float32(totalAmountGross)/100))

	// set agregates of net amounts per VAT rate
	for field, vatRates := range totalNetAmountFieldPerNetAmountMapping {
		net = 0
		for _, rate := range vatRates {
			net += totalNetAmountPerVATRate[rate]
		}
		if net > 0 {
			invoice.SetValue("Faktura.Fa."+field, fmt.Sprintf("%.2f", float32(net)/100))
		}
	}
	// set agregates of VAT amounts per VAT rate
	for field, vatRates := range totalVATAmountFieldPerVATRatesMapping {
		vat = 0
		for _, rate := range vatRates {
			vat += totalVatAmountPerVATRate[rate]
		}
		if vat > 0 {
			invoice.SetValue("Faktura.Fa."+field, fmt.Sprintf("%.2f", float32(vat)/100))
		}
	}

	//itemsNode.Children = append(itemsNode.Children, &xml.Node{Name: "WartoscWierszyFaktury2", Value: fmt.Sprintf("%.2f", float32(totalAmountGross/100))})

	if err = invoice.ApplyOrdering(FA_2ChildrenOrder); err != nil {
		return fmt.Errorf("unable to apply ordering: %v", err)
	}
	destFile, err := os.OpenFile(dest, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("cannot create target file: %v", err)
	}
	destFile.Truncate(0)
	defer destFile.Close()
	destFile.WriteString(`<?xml version="1.0" encoding="UTF-8"?>`)
	destFile.WriteString("\n")
	return invoice.DumpToFile(destFile, 0)
}
