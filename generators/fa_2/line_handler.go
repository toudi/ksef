package fa_2

import (
	"fmt"
	"ksef/common"
	"strconv"
	"strings"
)

func (fg *FA2Generator) LineHandler(invoice *common.Invoice, section string, data map[string]string, invoiceReady func() error) error {
	var err error

	if fg.isCommonData(section) {
		if fg.commonData == nil {
			fg.commonData = make(map[string]string)
		}
		for key, value := range data {
			fg.commonData[section+"."+key] = value
		}
		return nil
	}
	if fg.createNewInvoice(section) {
		if len(invoice.Items) > 0 {
			if err = invoiceReady(); err != nil {
				return err
			}
		}
		invoice.Clear()
		fmt.Printf("new invoice: %v\n", data)
		for key, value := range fg.commonData {
			invoice.Attributes[key] = value
		}
		for key, value := range data {
			invoice.Attributes[fmt.Sprintf("Faktura.Fa.%s", key)] = value
		}
	}
	if fg.isItemSection(section) {
		item := &common.InvoiceItem{Attributes: make(map[string]string)}
		for field, value := range data {
			field_lowercase := strings.ToLower(field)
			switch field_lowercase {
			case "p_7", "description":
				item.Description = value
			case "p_8a", "unit":
				item.Unit = value
			case "p_8b", "quantity":
				if item.Quantity, err = strconv.ParseFloat(value, 64); err != nil {
					return fmt.Errorf("cannot parse item quantity: %v", err)
				}
			case "p_9a", "unit-price-net":
				var unitPrice int64
				if unitPrice, err = strconv.ParseInt(value, 10, 64); err != nil {
					return fmt.Errorf("cannot parse item net price: %v", err)
				}
				item.UnitPrice.Value = int(unitPrice)
			case "p_12", "vat-description":
				item.UnitPrice.Vat.Description = value
				if vatRate, err := strconv.ParseInt(value, 10, 32); err == nil {
					item.UnitPrice.Vat.Rate = int(vatRate)
				}

			default:
				item.Attributes[field] = value
			}
		}
		if err = invoice.AddItem(item); err != nil {
			return err
		}
	}
	return nil
}
