package types

type invoiceSubjectQueryResponse struct {
	IssuedBy struct {
		TIN string `json:"identifier"`
	} `json:"issuedByIdentifier"`
	Issuer struct {
		FullName string `json:"fullName"`
	} `json:"issuedByName"`
}

type InvoiceSubject struct {
	invoiceSubjectQueryResponse `       yaml:"-"`
	TIN                         string `yaml:"NIP"`
	FullName                    string `yaml:"fullName,omitempty"`
}
