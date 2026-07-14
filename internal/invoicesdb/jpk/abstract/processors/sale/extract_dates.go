package sale

import (
	"errors"
	"ksef/internal/invoicesdb/jpk/abstract/types"
	monthlyregistry "ksef/internal/invoicesdb/monthly-registry"

	"github.com/beevik/etree"
)

const (
	xpathIssued     = "//Faktura/Fa/P_1"
	xpathSale       = "//Faktura/Fa/P_6"
	xpathSaleRanged = "//Faktura/Fa/OkresFa/P_6_Do"
)

func ExtractDates(
	invoice *monthlyregistry.Invoice,
	doc *etree.Document,
	salesRow *types.SaleItem,
) error {
	salesRow.IssueDate = doc.FindElement(xpathIssued).Text()
	var saleDateElement *etree.Element = doc.FindElement(xpathSale)
	if saleDateRangeElement := doc.FindElement(xpathSaleRanged); saleDateRangeElement != nil {
		saleDateElement = saleDateRangeElement
	}
	if saleDateElement == nil {
		return errors.New("Nie udało się znaleźć daty sprzedaży faktury")
	}
	salesRow.SaleDate = saleDateElement.Text()
	return nil
}
