package sanitation

import (
	"strconv"
	"testing"
)

func TestValidNumbers(t *testing.T) {
	validNumbers := []string{
		`371144371144376`,
		`341134113411347`,
		`370000000000002`,
		`378282246310005`,
		`6011016011016011`,
		`6559906559906557`,
		`6011000000000012`,
		`6011111111111117`,
		`5111005111051128`,
		`5112345112345114`,
		`5424000000000015`,
		`5105105105105100`,
		`4112344112344113`,
		`4007000000027`,
		`4111111111111111`,
		`4110144110144115`,
		`4114360123456785`,
		`4061724061724061`,
		`5115915115915118`,
		`5116601234567894`,
		`36111111111111`,
		`36110361103612`,
		`36438936438936`,
		`30569309025904`,
	}

	luhn := Luhn{}
	for i, v := range validNumbers {
		t.Run(strconv.Itoa(i), func(t *testing.T) {

			if !luhn.IsValid(v) {
				t.Errorf("Mismatch. Case %d, want: true, got: false (%v)", i, v)
			}

		})
	}
}

func TestInvalidNumbers(t *testing.T) {
	invalidNumbers := []string{
		`4111111111111`,
		`5105105105105106`,
		`111`,
	}

	luhn := Luhn{}
	for i, v := range invalidNumbers {
		t.Run(strconv.Itoa(i), func(t *testing.T) {

			if luhn.IsValid(v) {
				t.Errorf("Mismatch. Case %d, want: false, got: true (%v)", i, v)
			}

		})
	}

}
