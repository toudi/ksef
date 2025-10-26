package mnemonics

type FieldMnemonic struct {
	Name     string
	Mnemonic string
}

var (
	// item mnemonics
	Item           = FieldMnemonic{Name: "p_7", Mnemonic: "item"}
	Units          = FieldMnemonic{Name: "p_8a", Mnemonic: "units"}
	Quantity       = FieldMnemonic{Name: "p_8b", Mnemonic: "quantity"}
	UnitPriceNet   = FieldMnemonic{Name: "p_9a", Mnemonic: "unit-price-net"}
	UnitPriceGross = FieldMnemonic{Name: "p_9b", Mnemonic: "unit-price-gross"}
	VatRate        = FieldMnemonic{Name: "p_12", Mnemonic: "vat-rate"}

	ItemMnemonics = []FieldMnemonic{
		Item, Units, Quantity, UnitPriceNet, UnitPriceGross, VatRate,
	}

	// invoice mnemonics
	TotalNet   = FieldMnemonic{Name: "p_13_1", Mnemonic: "total-net"}
	TotalVat   = FieldMnemonic{Name: "p_14_1", Mnemonic: "total-vat"}
	TotalGross = FieldMnemonic{Name: "p_15", Mnemonic: "total-gross"}

	InvoiceMnemonics = []FieldMnemonic{
		TotalNet, TotalVat, TotalGross,
	}
)
