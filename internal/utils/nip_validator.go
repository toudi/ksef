package utils

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

var (
	errNIPInvalid    = errors.New("nieprawidłowy numer NIP")
	errNIPTooShort   = errors.New("numer NIP zbyt krótki")
	errNotADigit     = errors.New("nieprawidłowa cyfra")
	errInvalidModulo = errors.New("nieprawidłowa cyfra kontrolna")
)

type NIPValidatorType func(nip string) error

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
func NIPLengthValidator(nip string) error {
	if !nipRegexp.MatchString(normalizeNIP(nip)) {
		return errNIPInvalid
	}
	return nil
}

var digitWeights = []int{6, 5, 7, 2, 3, 4, 5, 6, 7}

func NIPValidator(nip string) error {
	nip = normalizeNIP(nip)
	var checksum int

	if len(nip) != 10 {
		return errors.Join(errNIPInvalid, errNIPTooShort)
	}

	for digitNo := range 9 {
		digit, err := strconv.Atoi(string(nip[digitNo]))
		if err != nil {
			return errors.Join(errNIPInvalid, errNotADigit)
		}
		checksum += digit * digitWeights[digitNo]
	}

	var expectedLastDigit = strconv.Itoa(checksum % 11)
	if string(nip[9]) != expectedLastDigit {
		return errors.Join(errNIPInvalid, errInvalidModulo)
	}

	return nil
}
