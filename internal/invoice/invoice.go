package invoice

import (
	"ksef/internal/money"
	"strconv"
	"time"
)

type VAT struct {
	Rate        int
	Description string
	Except      bool
}

type Price struct {
	money.MonetaryValue
	IsGross bool
	Vat     VAT
}

type Amount struct {
	// these ints indicate amounts in groszy's (1/100'ths of 1PLN)
	Net   int
	Gross int
	VAT   int
}

func (a *Amount) Add(other Amount) {
	a.Gross += other.Gross
	a.Net += other.Net
	a.VAT += other.VAT
}

type KSeFFlags struct {
	Offline bool
}

func (kf *KSeFFlags) Load(data map[string]string) {
	if isOfflineStr, exists := data["offline"]; exists {
		if isOffline, err := strconv.ParseBool(isOfflineStr); err == nil {
			kf.Offline = isOffline
		}
	}
}

type Invoice struct {
	IssuerNIP        string
	RecipientName    string
	GenerationTime   time.Time
	Number           string
	Issued           time.Time
	Items            []*InvoiceItem
	TotalPerVATRate  map[string]Amount
	Total            Amount
	Attributes       map[string]string
	Meta             map[string]string
	BasedOnNetPrices bool
	KSeFFlags        *KSeFFlags
}

func (i *Invoice) Clear() {
	i.Items = make([]*InvoiceItem, 0)
	i.TotalPerVATRate = make(map[string]Amount)
	i.Attributes = make(map[string]string)
	i.Meta = make(map[string]string)
	i.Total = Amount{}
	i.KSeFFlags = &KSeFFlags{}
}
