package keyring

const (
	AppPrefix         = "github.com/toudi/ksef"
	keySessionTokens  = "sessionTokens"
	KeyBackupPassword = "backup-password"
	keyPrivateKeyAES  = "privateKeyAES"
	// primary key encryption cipher - base64 encoded since zalando lib uses strings internally rather than bytes
)

func SessionTokensKey(certId string) string {
	return keySessionTokens + "-" + certId
}

func PrivateKeyEncryptionKey(keyId string) string {
	return keyPrivateKeyAES + "-" + keyId
}
