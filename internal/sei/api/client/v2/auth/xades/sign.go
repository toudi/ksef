package xades

import (
	"bytes"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"io"
	"ksef/internal/certsdb"
	"os"
	"strings"
	"time"

	"github.com/beevik/etree"
	"github.com/google/uuid"
)

type TemplateVars struct {
	ID          uuid.UUID
	SigningTime time.Time
	Certificate struct {
		Hash []byte
		Raw  *x509.Certificate
	}
	SrcDoc struct {
		Hash []byte
	}
	SignedInfo struct {
		Content string
		Hash    []byte
	}
	SignedProperties struct {
		Content string
		Hash    []byte
	}
}

func SignAuthChallenge(challenge io.Reader, cert certsdb.Certificate, dest io.Writer) error {
	var templateVars = TemplateVars{
		SigningTime: time.Now().UTC(),
	}
	templateVars.Certificate.Hash = make([]byte, 32)
	templateVars.SrcDoc.Hash = make([]byte, 32)
	templateVars.SignedInfo.Hash = make([]byte, 32)
	templateVars.SignedProperties.Hash = make([]byte, 32)
	var hash [32]byte

	// read certificate into memory
	certificateBytes, err := os.ReadFile(cert.Filename())
	if err != nil {
		return err
	}
	certificateBlock, _ := pem.Decode(certificateBytes)
	hash = sha256.Sum256(certificateBlock.Bytes)
	copy(templateVars.Certificate.Hash[:], hash[:])
	certificate, err := x509.ParseCertificate(certificateBlock.Bytes)
	if err != nil {
		return err
	}
	templateVars.Certificate.Raw = certificate

	// let's deal with the original message
	signedDocument := etree.NewDocument()
	if _, err := signedDocument.ReadFrom(challenge); err != nil {
		return err
	}
	// now we need to dump it to a temporary buffer so that we could hash the content
	// this is important as the "<?xml .." cannot become part of the hash
	var buffer bytes.Buffer
	signedDocument.Root().WriteTo(&buffer, &etree.WriteSettings{})
	// and we can hash it:
	hash = sha256.Sum256(buffer.Bytes())
	copy(templateVars.SrcDoc.Hash[:], hash[:])
	if templateVars.ID, err = uuid.NewRandom(); err != nil {
		return err
	}

	// we can construct the signedProperties element
	var signedProperties string
	if signedProperties, err = renderTemplate(signedPropertiesTemplate, templateVars); err != nil {
		return err
	}
	// let's hash it:
	hash = sha256.Sum256([]byte(signedProperties))
	copy(templateVars.SignedProperties.Hash[:], hash[:])
	templateVars.SignedProperties.Content = signedProperties

	// now we can render signed info
	var signedInfo string
	if signedInfo, err = renderTemplate(signedInfoTemplate, templateVars); err != nil {
		return err
	}
	templateVars.SignedInfo.Content = signedInfo
	hash = sha256.Sum256([]byte(signedInfo))
	copy(templateVars.SignedInfo.Hash[:], hash[:])

	// finally, build the signature XML
	var signatureXML string
	if signatureXML, err = renderTemplate(signatureTemplate, templateVars); err != nil {
		return err
	}

	// now for the final part
	// 1. read the signature XML and parse it as an etree node:
	signatureDoc := etree.NewDocument()
	signatureDoc.ReadFromString(signatureXML)
	// 2. calculate the signature of the document. We will use signedInfo node as our
	// input hash
	var signature []byte
	if signature, err = cert.SignContent([]byte(signedInfo)); err != nil {
		return err
	}

	// now we need to locate the signature node and replace it's content by the actual signature
	signatureValueNode := signatureDoc.FindElement(`[@Id="signature-value-` + templateVars.ID.String() + `"]`)
	if signatureValueNode == nil {
		return errors.New("unable to find signature value node")
	}
	signatureValueNode.SetText(base64.StdEncoding.EncodeToString(signature))

	// whew! that was a lot. now for the easy part
	// 1. output signature document to a buffer (you'll find out why in a moment)
	buffer.Reset()
	// required for c14n
	signatureDoc.WriteSettings.CanonicalEndTags = true
	if _, err = signatureDoc.WriteTo(&buffer); err != nil {
		return err
	}

	// and now we can replace the closing tag of source document's node with the signature:
	content, err := signedDocument.WriteToString()
	if err != nil {
		return err
	}
	// this may seem very weird and rightfully so
	// it all has to do with c14n. Essentially, there isn't any c14n that would discard whitespaces
	// between the nodes therefore if we start with:
	//
	// <foo>
	//  <bar></bar>
	// </foo>
	//
	// and then sign it, the c14n (when rendering "source" document) will yield the following:
	//
	// <foo>
	//  <bar></bar>
	//
	// </foo>
	// since the empty space would usually be occupied by the signature node.
	// however if we manually replace it it will end up being on the same line thus not suffering
	// from this problem.
	// why there isn't a c14n that gets rid of all whitespaces between nodes is just beyond me.
	content = strings.Replace(content, "</AuthTokenRequest>", buffer.String()+"</AuthTokenRequest>", 1)
	_, err = io.WriteString(dest, content)

	return err
}
