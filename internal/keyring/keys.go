package keyring

const (
	appPrefix        = "github.com/toudi/ksef"
	keySessionTokens = "sessionTokens"
)

func SessionTokensKey(certId string) string {
	return keySessionTokens + "-" + certId
}
