package xades

import (
	"bytes"
	"crypto/x509/pkix"
	_ "embed"
	"encoding/base64"
	"fmt"
	"text/template"
)

var (
	//go:embed "templates/AuthTokenRequest.xml"
	challengeRequest         string
	challengeRequestTemplate *template.Template
	//go:embed "templates/signed/signature.xml"
	signatureNode     string
	signatureTemplate *template.Template
	//go:embed "templates/signed/signed-info.xml"
	signedInfoNode     string
	signedInfoTemplate *template.Template
	//go:embed "templates/signed/signed-properties.xml"
	signedPropertiesNode     string
	signedPropertiesTemplate *template.Template
)

func init() {
	var funcMap = template.FuncMap{
		"base64": func(input []byte) string {
			return base64.StdEncoding.EncodeToString(input)
		},
		"rdname": func(name pkix.Name) string {
			// s := name.ToRDNSequence().String()
			var s = ""

			// fmt.Printf("s? %v\n", s)

			// List the extra attributes that should be added.
			attributeTypeNames := map[string]string{
				"2.5.4.42": "GN",
				"2.5.4.4":  "surname",
				"2.5.4.5":  "serialNumber",
				"2.5.4.3":  "CN",
				"2.5.4.6":  "C",
			}

			for typ, typeName := range attributeTypeNames {
				for _, atv := range name.Names {
					oidString := atv.Type.String()
					if oidString == typ {
						// To keep this demo simple, I just call fmt.Sprint to get the string.
						// Maybe you want to escape some of the characters.
						// See https://github.com/golang/go/blob/1db23771afc7b9b259e926db35602ecf5047ae23/src/crypto/x509/pkix/pkix.go#L67-L86
						if s != "" {
							s += ", "
						}
						s += typeName + "=" + fmt.Sprint(atv.Value)
						break
					}
				}
			}
			return s
		},
	}
	challengeRequestTemplate = template.Must(template.New("challenge-request").Parse(challengeRequest))
	signatureTemplate = template.Must(template.New("signature").Funcs(funcMap).Parse(signatureNode))
	signedInfoTemplate = template.Must(template.New("signed-info").Funcs(funcMap).Parse(signedInfoNode))
	signedPropertiesTemplate = template.Must(template.New("signed-properties").Funcs(funcMap).Parse(signedPropertiesNode))
}

func renderTemplate(tmpl *template.Template, vars TemplateVars) (string, error) {
	var buffer bytes.Buffer
	err := tmpl.Execute(&buffer, vars)
	return buffer.String(), err
}
