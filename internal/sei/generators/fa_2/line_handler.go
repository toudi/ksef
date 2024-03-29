package fa_2

import (
	"fmt"
	"ksef/internal/invoice"
	"strconv"
	"strings"
)

func (fg *FA2Generator) LineHandler(
	inv *invoice.Invoice,
	section string,
	data map[string]string,
	invoiceReady func() error,
) error {
	var err error

	if fg.isCommonData(section) {
		for key, value := range data {
			fg.commonData[section+"."+key] = value
		}
		return nil
	} else if fg.isItemSection(section) {
		item := &invoice.InvoiceItem{Attributes: make(map[string]string)}

		for field, value := range data {
			field_lowercase := strings.ToLower(field)
			switch field_lowercase {
			case "p_7", "item":
				item.Description = value
			case "p_8a", "units":
				item.Unit = value
			case "p_8b", "quantity":
				if err = item.Quantity.LoadFromString(value); err != nil {
					return fmt.Errorf("cannot parse item quantity: %v", err)
				}
			case "p_9a", "unit-price-net":
				if err = item.UnitPrice.LoadFromString(value); err != nil {
					return fmt.Errorf("cannot parse item net price: %v", err)
				}
			case "p_9b", "unit-price-gross":
				if err = item.UnitPrice.LoadFromString(value); err != nil {
					return fmt.Errorf("cannot parse item gross price: %v", err)
				}
				item.UnitPrice.IsGross = true
			case "p_12", "vat-rate":
				item.UnitPrice.Vat.Description = value
				if vatRate, err := strconv.ParseInt(value, 10, 32); err == nil {
					item.UnitPrice.Vat.Rate = int(vatRate)
				}
			case "vat-rate.except", "p_12.except":
				if value == "1" {
					item.UnitPrice.Vat.Except = true
				}
			default:
				item.Attributes[field] = value
			}
		}
		if err = inv.AddItem(item); err != nil {
			return err
		}
	} else {
		for key, value := range data {
			inv.Attributes[section+"."+key] = value
		}
	}
	return nil
}
