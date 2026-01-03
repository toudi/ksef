package jpk

type FieldInfo struct {
	KSeF string
	JPK  string
}

type Extractor struct {
	Label    string
	Base     FieldInfo
	VAT      *FieldInfo
	Required bool
}

var Extractor_ZW = Extractor{
	Label: "zwolnione",
	Base: FieldInfo{
		KSeF: "P_13_7",
		JPK:  "P_10",
	},
}

var Extractor_Foreign = Extractor{
	Label: "poza terytorium",
	Base: FieldInfo{
		KSeF: "P_13_8",
		JPK:  "P_11",
	},
	Required: true,
}

var Extractor_NP = Extractor{
	Label: "np",
	Base: FieldInfo{
		KSeF: "P_13_9",
		JPK:  "P_12",
	},
}

var Extractor_0 = Extractor{
	Label: "0%",
	Base: FieldInfo{
		KSeF: "P_13_6_1",
		JPK:  "P_13",
	},
	Required: true,
}

var Extractor_5 = Extractor{
	Label: "5%",
	Base: FieldInfo{
		KSeF: "P_13_3",
		JPK:  "P_15",
	},
	VAT: &FieldInfo{
		KSeF: "P_14_3",
		JPK:  "P_16",
	},
	Required: true,
}

var Extractor_7_8 = Extractor{
	Label: "7% / 8%",
	Base: FieldInfo{
		KSeF: "P_13_2",
		JPK:  "P_17",
	},
	VAT: &FieldInfo{
		KSeF: "P_14_2",
		JPK:  "P_18",
	},
	Required: true,
}

var Extractor_22_23 = Extractor{
	Label: "22% / 23%",
	Base: FieldInfo{
		KSeF: "P_13_1",
		JPK:  "P_19",
	},
	VAT: &FieldInfo{
		KSeF: "P_14_1",
		JPK:  "P_20",
	},
	Required: true,
}

var allExtractors = []Extractor{
	Extractor_ZW,
	Extractor_Foreign,
	Extractor_NP,
	Extractor_0,
	Extractor_5,
	Extractor_7_8,
	Extractor_22_23,
}
