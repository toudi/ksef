package types

type VatRate uint8

var VatRates map[string]VatRate

// as extracted from FA_3 xsd schema (TStawkaPodatku)
const (
	VatRateZw      VatRate = iota + 1 // zwolnione
	VatRateNpI                        // np I  niepodlegające opodatkowaniu - dostawy towarów oraz świadczenia usług poza terytorium kraju (z wyłączeniem krajów UE)
	VatRateNpII                       // np II niepodlegające opodatkowaniu - dostawy towarów oraz świadczenia usług poza terytorium kraju (kraje UE)
	VatRateZeroKR                     // Stawka 0% w przypadku sprzedaży towarów i świadczenia usług na terytorium kraju (z wyłączeniem WDT i eksportu)
	VatRateZeroWDT                    // Stawka 0% w przypadku wewnątrzwspólnotowej dostawy towarów (WDT)
	VatRateZeroEX                     // Stawka 0% w przypadku eksportu towarów
	VatRateRC                         // odwrotne obciążenie
	VatRate3
	VatRate4
	VatRate5
	VatRate7
	VatRate8
	VatRate22
	VatRate23
)

func init() {
	// the reason for this utterly stupid map is very simple - VAT Rates come from the parsed XML invoices.
	// on the other hand, JPK has specific fields for reporting specific vat rates and the aggregated amounts
	// associated with them. And so we need to recognize which one is it so that we could later put it into
	// the dedicated fields. I've got no idea how to do it nicer - hence the following abomination.
	VatRates = map[string]VatRate{
		"zw":    VatRateZw,
		"np I":  VatRateNpI,
		"np II": VatRateNpII,
		"0 KR":  VatRateZeroKR,
		"0 WDT": VatRateZeroWDT,
		"0 EX":  VatRateZeroEX,
		"oo":    VatRateRC,
		"3":     VatRate3,
		"4":     VatRate4,
		"5":     VatRate5,
		"7":     VatRate7,
		"8":     VatRate8,
		"22":    VatRate22,
		"23":    VatRate23,
	}
}
