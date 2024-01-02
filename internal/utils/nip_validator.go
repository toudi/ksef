package utils

import (
	"regexp"
	"strings"
)

type NIPValidatorType func(nip string) bool

// check whether NIP matches the predefined pattern. Note that
// the NIP must already be normalized at this point
var nipRegexp = regexp.MustCompile("^[0-9]{10}$")

// replacer used to remove separators
var nipNormalizer = strings.NewReplacer("-", "", " ", "")

func normalizeNIP(nip string) string {
	return nipNormalizer.Replace(nip)
}

// NIPLengthValidator should only be used on the test environment
// where one can use completely fake NIP numbers, thus the only
// thing that we actually can validate is the length of the NIP itself
func NIPLengthValidator(nip string) bool {
	return nipRegexp.MatchString(normalizeNIP(nip))
}

var nipValidationWeights = [...]int{6, 5, 7, 2, 3, 4, 5, 6, 7}

// ASCII starts the numbers at offset 48 which happens to be
// the result of int('0') therefore we can cast the digit into
// an int by calling int(char) - offset
const zeroAsAsciiInt = int('0')

func NIPValidator(nip string) bool {
	nip = normalizeNIP(nip)

	if !NIPLengthValidator(nip) {
		return false
	}

	// https://stackoverflow.com/questions/37765687/golang-how-to-convert-string-to-int
	var checksum int = 0

	for weightIndex, weight := range nipValidationWeights {
		checksum += (int(nip[weightIndex]) - zeroAsAsciiInt) * weight
	}

	return checksum%11 == int(nip[9])-zeroAsAsciiInt
}
