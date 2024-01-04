package utils

import "testing"

type testCase struct {
	nip            string
	expectedResult bool
}

func TestNIPLengthValidator(t *testing.T) {
	for _, test := range []testCase{
		{nip: "1111111111", expectedResult: true},
		{nip: "111-111-11-11", expectedResult: true},
		{nip: "111 111 11 11", expectedResult: true},
		{nip: "2222222222", expectedResult: true},
		{nip: "00000000000", expectedResult: false},
		{nip: "000000000a", expectedResult: false},
	} {
		if NIPLengthValidator(test.nip) != test.expectedResult {
			t.Fatalf("unexpected result when calling NIPLengthValidator(%s): %v != %v", test.nip, !test.expectedResult, test.expectedResult)
		}
	}
}

func TestNIPValidator(t *testing.T) {
	for _, test := range []testCase{
		{nip: "1234567901", expectedResult: false},
		{nip: "1234563218", expectedResult: true},
		{nip: "123-456-32-18", expectedResult: true},
	} {
		if NIPValidator(test.nip) != test.expectedResult {
			t.Fatalf("unexpected result when calling NIPValidator(%s): %v != %v", test.nip, !test.expectedResult, test.expectedResult)
		}
	}

}
