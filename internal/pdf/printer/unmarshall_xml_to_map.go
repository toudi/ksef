package printer

import (
	"encoding/xml"
	"fmt"
	"io"
	"ksef/internal/invoice"
	"ksef/internal/sei/generators/fa/mnemonics"
	"strconv"
	"strings"
)

type Amounts struct {
	Total     *invoice.Amount
	ByVatRate map[string]*invoice.Amount
}

type Invoice struct {
	Xml           map[string]any
	Amounts       Amounts
	ArrayElements map[string]bool
}

const (
	pathItem = "Faktura.Fa.FaWiersz"
	pathFa   = "Faktura.Fa"
)

// original code:
// https://habr.com/en/articles/847854/
// I just modified it slightly:
// * so that it can support array elements
// * so that it can convert invoice properties to mnemonics
// * so that it can calculate totals and totals per vat rate as mnemonics
func (x *Invoice) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {
	x.Xml = map[string]any{"_": start.Name.Local}
	x.Amounts = Amounts{
		Total:     &invoice.Amount{},
		ByVatRate: make(map[string]*invoice.Amount),
	}
	path := []map[string]any{x.Xml}
	dottedPath := []string{start.Name.Local}

	for {
		token, err := decoder.Token()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		switch elem := token.(type) {
		case xml.StartElement:
			newMap := map[string]any{"_": elem.Name.Local}
			dottedPath = append(dottedPath, elem.Name.Local)
			fullName := strings.Join(dottedPath, ".")
			if x.ArrayElements[fullName] {
				var elemSlice []any
				if sliceI, ok := path[len(path)-1][elem.Name.Local].([]any); ok {
					elemSlice = sliceI
					elemSlice = append(elemSlice, newMap)
				} else {
					elemSlice = []any{newMap}
					path[len(path)-1][elem.Name.Local] = elemSlice
				}
			} else {
				path[len(path)-1][elem.Name.Local] = newMap
			}
			path = append(path, newMap)
		case xml.EndElement:
			dottedPath = dottedPath[:len(path)-1]
			path = path[:len(path)-1]
		case xml.CharData:
			val := strings.TrimSpace(string(elem))
			if val == "" {
				break
			}
			curName := path[len(path)-1]["_"].(string)
			var nodeValue = typeConvert(val)
			path[len(path)-2][curName] = nodeValue
			fullName := strings.Join(dottedPath[:len(dottedPath)-1], ".")
			if fullName == pathItem {
				for _, mnemonic := range mnemonics.ItemMnemonics {
					if strings.ToLower(curName) == mnemonic.Name {
						path[len(path)-2][mnemonic.Mnemonic] = nodeValue
						break
					}
				}
			}
			if fullName == pathFa {
				for _, mnemonic := range mnemonics.InvoiceMnemonics {
					if strings.ToLower(curName) == mnemonic.Name {
						path[len(path)-2][mnemonic.Mnemonic] = nodeValue
						break
					}
				}
			}
		}
	}

	items := x.Xml["Fa"].(map[string]any)["FaWiersz"].([]any)
	for _, item := range items {
		var ii invoice.InvoiceItem

		itemMap := item.(map[string]any)
		vatRate := itemMap["vat-rate"]
		quantity := itemMap["quantity"].(float64)
		unitPriceNet, calculateFromNet := itemMap["unit-price-net"]
		unitPriceGross := itemMap["unit-price-gross"]

		ii.Quantity.LoadFromString(strconv.FormatFloat(quantity, 'f', -1, 64))
		if vatFloat, ok := vatRate.(float64); ok {
			vatInt := int(vatFloat)
			ii.UnitPrice.Vat.Rate = vatInt
			ii.UnitPrice.Vat.Description = fmt.Sprintf("%d %%", vatInt)
		} else {
			ii.UnitPrice.Vat.Description = vatRate.(string)
		}

		if calculateFromNet {
			ii.UnitPrice.LoadFromString(strconv.FormatFloat(unitPriceNet.(float64), 'f', -1, 64))
		} else {
			ii.UnitPrice.LoadFromString(strconv.FormatFloat(unitPriceGross.(float64), 'f', -1, 64))
			ii.UnitPrice.IsGross = true
		}

		amt := ii.Amount()

		var itemM = item.(map[string]any)
		itemM["vat-desc"] = ii.UnitPrice.Vat.Description
		itemM["amount-net"] = amt.Net
		itemM["amount-gross"] = amt.Gross
		itemM["amount-vat"] = amt.VAT

		x.Amounts.Total.Add(amt)
		if _, exists := x.Amounts.ByVatRate[ii.UnitPrice.Vat.Description]; !exists {
			x.Amounts.ByVatRate[ii.UnitPrice.Vat.Description] = &invoice.Amount{}
		}
		x.Amounts.ByVatRate[ii.UnitPrice.Vat.Description].Add(amt)
	}
	return nil
}

func typeConvert(s string) any {
	f, err := strconv.ParseFloat(s, 64)
	if err == nil {
		return f
	}
	return s
}
