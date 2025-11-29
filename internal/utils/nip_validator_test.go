package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type testCase struct {
	nip string
	err error
}

func TestNIPLengthValidator(t *testing.T) {
	for _, test := range []testCase{
		{nip: "1111111111"},
		{nip: "111-111-11-11"},
		{nip: "111 111 11 11"},
		{nip: "2222222222"},
		{nip: "00000000000", err: errNIPInvalid},
		{nip: "000000000a", err: errNIPInvalid},
	} {
		err := NIPLengthValidator(test.nip)
		require.Equal(t, test.err, err)
	}
}

func TestNIPValidator(t *testing.T) {
	for _, test := range []testCase{
		{nip: "1234567901", err: errInvalidModulo},
		{nip: "1234563218"},
		{nip: "123-456-32-18"},
	} {
		err := NIPValidator(test.nip)
		if test.err != nil {
			require.ErrorIs(t, err, test.err)
		}
	}
}
