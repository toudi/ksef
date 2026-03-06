package flags

// module created stricly for the reason of avoiding import cycles
const (
	CfgKeyKeyringEngine             = "keyring.engine"
	CfgKeyKeyringFileLocation       = "keyring.file.path"
	CfgKeyKeyringFileBuffered       = "keyring.file.buffered"
	CfgKeyKeyringFileAskPassword    = "keyring.file.ask-password"
	CfgKeyKeyringFilePasswordFile   = "keyring.file.password-file"
	CfgKeyKeyringFilePasswordEnvVar = "keyring.file.password-env-var"

	KeyringEngineSystem = "system"
	KeyringEngineFile   = "file"
)
