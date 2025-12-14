package types

type Invoice struct {
	Filename string
	Checksum string
	Offline  bool
}

type InvoiceFormCode struct {
	SystemCode    string `xml:"kodSystemowy,attr" json:"systemCode"`
	SchemaVersion string `xml:"wersjaSchemy,attr" json:"schemaVersion"`
	Value         string `xml:",chardata" json:"value"`
}

type UploadPayload map[InvoiceFormCode][]Invoice
