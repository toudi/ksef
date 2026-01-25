package purchase

import (
	"errors"
	"ksef/internal/invoicesdb/jpk/interfaces"
	"ksef/internal/invoicesdb/jpk/processors/vat"
	"ksef/internal/invoicesdb/jpk/types"
	monthlyregistry "ksef/internal/invoicesdb/monthly-registry"
	"ksef/internal/invoicesdb/shared"
	"ksef/internal/logging"
	"ksef/internal/money"

	"github.com/beevik/etree"
)

const (
	xpathItem          = "//Faktura/Fa/FaWiersz"
	vatRateFixedAssets = ":fixed-assets:" // ficticious vat rate so that we can separately aggregate amounts for fixed assets
	vatRateOther       = ":other:"
)

var errUnableToCalculateVAT = errors.New("unable to calculate VAT for invoice item")

func ProcessVatInfo(
	dest *types.Invoice,
	invoiceXML *etree.Document,
	registryInvoice *monthlyregistry.Invoice,
	manager interfaces.JPKManager,
) error {
	items := invoiceXML.FindElements(xpathItem)

	vatCalculator := &vat.Calculator{}

	logger := logging.JPKLogger.With("faktura zakupowa", registryInvoice.RefNo)

	for _, item := range items {
		itemHash := itemHashFromXML(item)
		// check if item should be excluded from the report
		if manager.ItemHasRule(registryInvoice, itemHash, func(rule shared.JPKItemRule) bool { return rule.Exclude }) {
			logger.Info("przedmiot faktury ma flagę wyłączenia z raportu", "hash", itemHash)
			continue
		}

		vatInfo, err := vatCalculator.GetVat(item)
		if err != nil {
			return errors.Join(errUnableToCalculateVAT, err)
		}

		// apply 50% VAT if applicable
		if manager.ItemHasRule(registryInvoice, itemHash, func(rule shared.JPKItemRule) bool { return rule.Vat50Percent }) {
			logger.Info("przedmiot faktury ma flagę zastosowania 50% VAT", "hash", itemHash)
			vatInfo.VatAmount = vatInfo.VatAmount.HalveAndRoundUp()
		}

		if dest.VAT == nil {
			dest.VAT = &types.VATInfo{
				Base: money.MonetaryValue{},
				Vat:  money.MonetaryValue{},
			}
		}

		dest.VAT.Base = dest.VAT.Base.Add(vatInfo.NetAmount)
		dest.VAT.Vat = dest.VAT.Vat.Add(vatInfo.VatAmount)

		// decide if this is a fixed asset or not
		vatRateId := vatRateOther
		if manager.ItemHasRule(registryInvoice, itemHash, func(rule shared.JPKItemRule) bool { return rule.FixedAsset }) {
			logger.Info("przedmiot faktury ma flagę środków trwałych", "hash", itemHash)
			vatRateId = vatRateFixedAssets
		}
		if _, exists := dest.VATByRate[vatRateId]; !exists {
			dest.VATByRate[vatRateId] = &types.VATInfo{}
		}
		dest.VATByRate[vatRateId].Base = dest.VATByRate[vatRateId].Base.Add(vatInfo.NetAmount)
		dest.VATByRate[vatRateId].Vat = dest.VATByRate[vatRateId].Vat.Add(vatInfo.VatAmount)
	}

	if vatInfoFixedAssets, exists := dest.VATByRate[vatRateFixedAssets]; exists {
		dest.Attributes["K_40"] = vatInfoFixedAssets.Base.Format(2) // netto dla środków trwałych
		dest.Attributes["K_41"] = vatInfoFixedAssets.Vat.Format(2)  // vat dla środków trwałych
	}
	if vatInfoOther, exists := dest.VATByRate[vatRateOther]; exists {
		dest.Attributes["K_42"] = vatInfoOther.Base.Format(2) // netto - pozostałe towary i usługi
		dest.Attributes["K_43"] = vatInfoOther.Vat.Format(2)  // vat - pozostałe towary i usługi
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
