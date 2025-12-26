package keyring

const (
	AppPrefix         = "github.com/toudi/ksef"
	keySessionTokens  = "sessionTokens"
	KeyBackupPassword = "backup-password"
)

func SessionTokensKey(certId string) string {
	return keySessionTokens + "-" + certId
}
