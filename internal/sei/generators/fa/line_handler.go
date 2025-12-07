package fa

import (
	"fmt"
	"ksef/internal/invoice"
	"ksef/internal/sei/generators/fa/mnemonics"
	"strconv"
	"strings"
	"time"
)

const (
	invoiceHeaderSection = "faktura.fa"
	invoiceMetaSection   = "meta"
	issuerDataSection    = "faktura.podmiot1.daneidentyfikacyjne"
	recipientDataSection = "faktura.podmiot2.daneidentyfikacyjne"
	ksefSection          = "ksef"
)

func (fg *FAGenerator) LineHandler(
	inv *invoice.Invoice,
	section string,
	data map[string]string,
	invoiceReady func() error,
) error {
	var err error

	if strings.ToLower(section) == invoiceMetaSection {
		inv.Meta = data
		return nil
	}

	if strings.ToLower(section) == ksefSection {
		inv.KSeFFlags.Load(data)
		return nil
	}

	if strings.HasPrefix(strings.ToLower(section), recipientDataSection) {
		inv.RecipientName = data["Nazwa"]
	}

	if strings.HasPrefix(strings.ToLower(section), issuerDataSection) {
		inv.Issuer.NIP = data["NIP"]
		inv.Issuer.Name = data["Nazwa"]
	}

	if strings.ToLower(section) == invoiceHeaderSection {
		inv.Number = data["P_2"]
		inv.Issued, err = time.Parse("2006-01-02", data["P_1"])
		if err != nil {
			return err
		}
	}
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
			case mnemonics.Item.Name, mnemonics.Item.Mnemonic:
				item.Description = value
			case mnemonics.Units.Name, mnemonics.Units.Mnemonic:
				item.Unit = value
			case mnemonics.Quantity.Name, mnemonics.Quantity.Mnemonic:
				if err = item.Quantity.LoadFromString(value); err != nil {
					return fmt.Errorf("cannot parse item quantity: %v", err)
				}
			case mnemonics.UnitPriceNet.Name, mnemonics.UnitPriceNet.Mnemonic:
				if err = item.UnitPrice.LoadFromString(value); err != nil {
					return fmt.Errorf("cannot parse item net price: %v", err)
				}
			case mnemonics.UnitPriceGross.Name, mnemonics.UnitPriceGross.Mnemonic:
				if err = item.UnitPrice.LoadFromString(value); err != nil {
					return fmt.Errorf("cannot parse item gross price: %v", err)
				}
				item.UnitPrice.IsGross = true
			case mnemonics.VatRate.Name, mnemonics.VatRate.Mnemonic:
				if value == "np" {
					value = "np II"
				}
				item.UnitPrice.Vat.Description = value
				if vatRate, err := strconv.ParseInt(value, 10, 32); err == nil {
					item.UnitPrice.Vat.Rate = int(vatRate)
				}
			case "vat-rate.except", "p_12.except":
				// old mnemonic for FA_2
				if value == "1" {
					item.UnitPrice.Vat.Description = "np I"
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
