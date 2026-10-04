package certsdb

type CertificateUsage string

const (
	// publicly available, issued by the ministry of finance
	CertiticateUsageKsefTokenEncryption    CertificateUsage = "KsefTokenEncryption"
	CertificateUsageSymmetricKeyEncryption CertificateUsage = "SymmetricKeyEncryption"
	// issued by MCU (certificates and permissions module), per user
	CertificateUsageAuthentication CertificateUsage = "Authentication"
	CertificateUsageOffline        CertificateUsage = "Offline"
)
