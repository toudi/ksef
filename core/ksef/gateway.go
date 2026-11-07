package ksef

import "ksef/internal/utils"

type KSeFGateway string

const (
	KSeFGatewayInvalid KSeFGateway = ""
	KSeFGatewayTest    KSeFGateway = "test"
	KSeFGatewayDemo    KSeFGateway = "demo"
	KSeFGatewayProd    KSeFGateway = "prod"
)

type KSeFGatewayConfig struct {
	Gateway       KSeFGateway
	APIBaseURL    string // bramka API
	QRCodeBaseURL string // url do weryfikacji kodów QR
	NIPValidator  utils.NIPValidatorType
}

var (
	testEnvironment = KSeFGatewayConfig{
		Gateway:       KSeFGatewayTest,
		APIBaseURL:    "https://api-test.ksef.mf.gov.pl",
		QRCodeBaseURL: "https://qr-test.ksef.mf.gov.pl",
		NIPValidator:  utils.NIPLengthValidator,
	}
	demoEnvironment = KSeFGatewayConfig{
		Gateway:       KSeFGatewayDemo,
		APIBaseURL:    "https://api-demo.ksef.mf.gov.pl",
		QRCodeBaseURL: "https://qr-demo.ksef.mf.gov.pl",
		NIPValidator:  utils.NIPValidator,
	}
	prodEnvironment = KSeFGatewayConfig{
		Gateway:       KSeFGatewayProd,
		APIBaseURL:    "https://api.ksef.mf.gov.pl",
		QRCodeBaseURL: "https://qr.ksef.mf.gov.pl",
		NIPValidator:  utils.NIPValidator,
	}

	gateways = map[KSeFGateway]KSeFGatewayConfig{
		KSeFGatewayTest: testEnvironment,
		KSeFGatewayDemo: demoEnvironment,
		KSeFGatewayProd: prodEnvironment,
	}
)

// ParseGateway reports whether s is a valid KSeFGateway identifier.
func ParseGateway(gatewayId string) (KSeFGateway, error) {
	gateway := KSeFGateway(gatewayId)
	if _, ok := gateways[gateway]; !ok {
		return KSeFGatewayInvalid, ErrUnknownGateway
	}
	return gateway, nil
}

func GatewayConfig(gatewayId string) (KSeFGatewayConfig, error) {
	gateway, err := ParseGateway(gatewayId)
	if err != nil {
		return KSeFGatewayConfig{}, err
	}

	return gateways[gateway], nil
}
