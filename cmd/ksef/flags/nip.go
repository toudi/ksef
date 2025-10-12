package flags

import (
	"errors"
	"strconv"

	"github.com/spf13/pflag"
)

type NIP string

const (
	FlagNameNIP = "nip"
)

var (
	errNIPInvalid    = errors.New("nieprawidłowy numer NIP")
	errNIPTooShort   = errors.New("numer NIP zbyt krótki")
	errNotADigit     = errors.New("nieprawidłowa cyfra")
	errInvalidModulo = errors.New("nieprawidłowa cyfra kontrolna")
)

func (n *NIP) Type() string {
	return "string"
}

func (n *NIP) String() string {
	return string(*n)
}

func (n *NIP) Set(value string) error {
	if err := validateNIP(value); err != nil {
		return err
	}
	*n = NIP(value)
	return nil
}

func validateNIP(input string) error {
	var checksum int
	var digitWeights = []int{6, 5, 7, 2, 3, 4, 5, 6, 7}

	if len(input) != 10 {
		return errors.Join(errNIPInvalid, errNIPTooShort)
	}

	for digitNo := range 9 {
		digit, err := strconv.Atoi(string(input[digitNo]))
		if err != nil {
			return errors.Join(errNIPInvalid, errNotADigit)
		}
		checksum += digit * digitWeights[digitNo]
	}

	var expectedLastDigit = strconv.Itoa(checksum % 11)
	if string(input[9]) != expectedLastDigit {
		return errors.Join(errNIPInvalid, errInvalidModulo)
	}

	return nil
}

func NIPVarP(nip *NIP, flagSet *pflag.FlagSet) {
	flagSet.VarP(nip, "nip", "n", "numer NIP podmiotu")
}
