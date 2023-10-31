package uploader

import (
	"bytes"
	"crypto/aes"
	"embed"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"ksef/common"
	"net/http"
	"text/template"
	"time"
)

type authorisationResponseType struct {
	Timestamp time.Time `json:"timestamp"`
	Challenge string    `json:"challenge"`
}

//go:embed "interactive_auth_challenge.xml"
var templateInitToken embed.FS

func (u *Uploader) initSession() (*Session, error) {
	authorizationChallengeData, err := u.prepareAuthorizationChallengeData()
	if err != nil {
		return nil, fmt.Errorf("nie udało się przygotować pakietu do AuthorizationChallenge: %v", err)
	}

	request, err := http.NewRequest("POST", u.host+"api/online/Session/AuthorisationChallenge", authorizationChallengeData)
	if err != nil {
		return nil, fmt.Errorf("błąd przygotowywania requestu AuthorisationChallenge: %v", err)
	}
	request.Header.Add("Content-Type", "application/json")

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("błąd wywołania AuthorisationChallenge: %v", err)
	}
	defer response.Body.Close()

	var authorisationResponse authorisationResponseType
	if err = json.NewDecoder(response.Body).Decode(&authorisationResponse); err != nil {
		return nil, fmt.Errorf("błąd dekodowania AuthorisationResponse: %v", err)
	}

	return u.sendEncryptedToken(authorisationResponse)
}

func (u *Uploader) sendEncryptedToken(authorisationResponse authorisationResponseType) (*Session, error) {
	var challengePlaintext = fmt.Sprintf("%s|%d", u.token, authorisationResponse.Timestamp.UnixMilli())

	encryptedBytes, err := common.EncryptMessageWithCertificate(u.certificateFile, []byte(challengePlaintext))
	if err != nil {
		return nil, fmt.Errorf("błąd szyfrowania tokenu: %v", err)
	}

	var funcMap = template.FuncMap{
		"base64": base64.StdEncoding.EncodeToString,
	}

	tmpl, err := template.New("interactive_auth_challenge.xml").Funcs(funcMap).ParseFS(templateInitToken, "interactive_auth_challenge.xml")

	if err != nil {
		return nil, fmt.Errorf("błąd inicjalizacji szablonu: %v", err)
	}
	var authChallengeDataBuffer bytes.Buffer
	type templateDataType struct {
		Issuer         string
		Challenge      string
		EncryptedToken []byte
		Cipher         struct {
			IV            []byte
			EncryptionKey []byte
		}
	}

	var templateData templateDataType = templateDataType{
		Issuer:         u.issuer,
		Challenge:      authorisationResponse.Challenge,
		EncryptedToken: encryptedBytes,
	}

	encryptedKeyBytes, err := common.EncryptMessageWithCertificate(u.certificateFile, u.cipher.Key)
	if err != nil {
		return nil, fmt.Errorf("błąd szyfrowania klucza za pomocą klucza RSA ministerstwa: %v", err)
	}
	templateData.Cipher.EncryptionKey = make([]byte, len(encryptedKeyBytes))
	templateData.Cipher.IV = make([]byte, aes.BlockSize)

	copy(templateData.Cipher.IV, u.cipher.IV)
	copy(templateData.Cipher.EncryptionKey, encryptedKeyBytes)

	if err = tmpl.Execute(
		&authChallengeDataBuffer,
		templateData,
	); err != nil {
		return nil, fmt.Errorf("błąd generowania szablonu authRequest: %v\n", err)
	}

	resp, err := http.DefaultClient.Post(u.host+"api/online/Session/InitToken", "application/octet-stream", &authChallengeDataBuffer)
	if err != nil {
		return nil, fmt.Errorf("błąd wywołanai initToken: %v\n", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("body: %s\n", string(body))

	if resp.StatusCode/100 != 2 {
		return nil, fmt.Errorf("nieoczekiway kod odpowiedzi: %d vs 2xx", resp.StatusCode)
	}

	var session Session

	if err = json.Unmarshal(body, &session); err != nil {
		return nil, fmt.Errorf("błąd dekodowania initToken: %v", err)
	}

	return &session, nil
}
