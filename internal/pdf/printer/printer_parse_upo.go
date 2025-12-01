package printer

import (
	"encoding/base64"
	"encoding/xml"
	"errors"

	"github.com/beevik/etree"
)

var (
	errContextIdNotFound     = errors.New("nie znaleziono drzewa IdKontekstu")
	errCannotFindContextType = errors.New("nie znaleziono typu kontekstu")
)

func ParseUPO(contentBase64 string) (*UPO, error) {
	upoXMLBytes, err := base64.StdEncoding.DecodeString(contentBase64)
	if err != nil {
		return nil, err
	}

	var upo UPO
	if err = xml.Unmarshal(upoXMLBytes, &upo); err != nil {
		return nil, err
	}

	// ok so this whole section of the code is here only because
	// I cannot think of some easier way to describe this xml fragment
	// in terms of structs:
	// <IdKontekstu>
	//   <Nip>1111111111</Nip>
	// </IdKontekstu>
	//
	// there can be other nodes here besides Nip as defined here:
	// https://github.com/CIRFMF/ksef-docs/blob/main/faktury/upo/schemy/upo-v4-2.xsd
	// desperate times call for desperate measures and all that.
	var upoDoc = etree.NewDocument()
	if err = upoDoc.ReadFromBytes(upoXMLBytes); err != nil {
		return nil, err
	}

	contextId := upoDoc.FindElement("//IdKontekstu")
	if contextId == nil {
		return nil, errContextIdNotFound
	}

	var contextTypes = map[string]string{
		"Nip":                   "NIP",
		"IdWewnetrzny":          "Identyfikator wewnętrzny",
		"IdZlozonyVatUE":        "Kontekst złożony: NIP + numer VAT UE",
		"IdDostawcyUslugPeppol": "Identyfikator dostawcy usług Peppol",
	}

	for nodeName, contextType := range contextTypes {
		contextNode := contextId.FindElement("//" + nodeName)
		if contextNode == nil {
			continue
		}
		upo.Auth.Context.Type = contextType
		upo.Auth.Context.Value = contextNode.Text()
	}

	if upo.Auth.Context.Type == "" {
		return nil, errCannotFindContextType
	}

	return &upo, nil
}
