package certsdb

import "errors"

var ErrUnsupportedPrivateKeyType = errors.New("unsupported private key type (not EC PRIVATE KEY)")
