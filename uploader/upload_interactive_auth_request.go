package uploader

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"github.com/zalando/go-keyring"
)

type authorisationRequestType struct {
	ContextIdentifier struct {
		Type       string `json:"type"`
		Identifier string `json:"identifier"`
	} `json:"contextIdentifier"`
}

func (u *Uploader) prepareAuthorizationChallengeData() (io.Reader, error) {
	// step 1 - fetch the token from keyring
	var err error
	u.token, err = keyring.Get(u.host, u.issuer)
	if u.token == "" || err != nil {
		// if it does not exist then there's no point in continuing
		return nil, fmt.Errorf("token nierozpoznany. proszę ustawić token za pomocą komendy set-token")
	}

	// step 2 - call authorisationRequest

	var authorisationRequest authorisationRequestType
	authorisationRequest.ContextIdentifier.Type = "onip"
	authorisationRequest.ContextIdentifier.Identifier = u.issuer

	var requestBuffer bytes.Buffer
	err = json.NewEncoder(&requestBuffer).Encode(authorisationRequest)

	if err != nil {
		return nil, fmt.Errorf("błąd kodowania authorisationRequest do JSON: %v", err)
	}

	return &requestBuffer, nil
}
