package certsdb

import "errors"

const (
	// certyfikaty wystawiane przez KSeF
	UsageAuthentication Usage = "Authentication"
	UsageOffline        Usage = "Offline"
)

var (
	errUnexpectedUsageValue = errors.New("unexpected usage value")
)

func ValidateUsage(usageString string) (Usage, error) {
	var u = Usage(usageString)
	if u != UsageAuthentication && u != UsageOffline {
		return Usage(""), errUnexpectedUsageValue
	}
	return u, nil
}
