package purchase

import (
	"errors"
	"ksef/internal/invoicesdb/jpk/abstract/processors/vat"
	"ksef/internal/invoicesdb/jpk/abstract/types"
	"ksef/internal/invoicesdb/jpk/manager"
	jpkTypes "ksef/internal/invoicesdb/jpk/types"
	monthlyregistry "ksef/internal/invoicesdb/monthly-registry"
	"ksef/internal/invoicesdb/shared"
	"ksef/internal/logging"

	"github.com/beevik/etree"
)

const (
	xpathItem = "//Faktura/Fa/FaWiersz"
)

var errUnableToCalculateVAT = errors.New("unable to calculate VAT for invoice item")

func AggregateItems(
	manager *manager.JPKManager,
	invoice *monthlyregistry.Invoice,
	doc *etree.Document,
	purchase *types.PurchaseItem,
) error {
	items := doc.FindElements(xpathItem)

	vatCalculator := &vat.Calculator{}

	logger := logging.JPKLogger.With("faktura zakupowa", invoice.RefNo)

	for _, item := range items {
		itemHash := itemHashFromXML(item)
		// check if item should be excluded from the report
		if manager.ItemHasRule(invoice, itemHash, func(rule shared.JPKItemRule) bool { return rule.Exclude }) {
			logger.Info("przedmiot faktury ma flagę wyłączenia z raportu", "hash", itemHash)
			continue
		}

		vatInfo, err := vatCalculator.GetVat(item)
		if err != nil {
			return errors.Join(errUnableToCalculateVAT, err)
		}

		// apply 50% VAT if applicable
		if manager.ItemHasRule(invoice, itemHash, func(rule shared.JPKItemRule) bool { return rule.Vat50Percent }) {
			logger.Info("przedmiot faktury ma flagę zastosowania 50% VAT", "hash", itemHash)
			vatInfo.VatAmount = vatInfo.VatAmount.HalveAndRoundUp()
		}

		purchaseAttributes := jpkTypes.PurchaseAttributes{}
		vatRate, exists := jpkTypes.VatRates[vatInfo.VatRate]
		if !exists {
			return errors.New("unknown VAT rate")
		}

		purchaseAttributes.VatRate = vatRate
		// decide if this is a fixed asset or not
		if manager.ItemHasRule(invoice, itemHash, func(rule shared.JPKItemRule) bool { return rule.FixedAsset }) {
			logger.Info("przedmiot faktury ma flagę środków trwałych", "hash", itemHash)
			purchaseAttributes.FixedAssets = true
		}

		purchase.AddAmount(purchaseAttributes, *vatInfo)
	}

	return nil
}

const (
	pathName  = "P_7"
	pathIndex = "Indeks"
	pathGTIN  = "GTIN"
	pathPKWiU = "PKWiU"
)

func itemHashFromXML(node *etree.Element) shared.ItemHash {
	var hash shared.ItemHash

	if nodeName := node.FindElement(pathName); nodeName != nil {
		hash.Name = nodeName.Text()
	}
	if nodeIndex := node.FindElement(pathName); nodeIndex != nil {
		hash.Index = nodeIndex.Text()
	}
	if nodeGTIN := node.FindElement(pathGTIN); nodeGTIN != nil {
		hash.GTIN = nodeGTIN.Text()
	}
	if nodePKWiU := node.FindElement(pathPKWiU); nodePKWiU != nil {
		hash.PKWiU = nodePKWiU.Text()
	}

	return hash
}
