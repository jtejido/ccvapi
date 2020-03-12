package validation

import (
	"strconv"
	"strings"
	"testing"
	"unicode"
)

func TestLoadCardWwithoutPath(t *testing.T) {
	err := LoadCards("")

	if err == nil {
		t.Errorf("Should throw an error at an empty path")
	}
}

func TestSpacesIncluded(t *testing.T) {
	err := LoadCards("../card_types.json")

	if err != nil {
		t.Errorf("Failed to open json file")
	}

	validNumbers := []string{
		"5101 1800 0000 0007",
		"2222 4000 7000 0005",
	}

	for i, v := range validNumbers {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			if isValidInputType(v) {
				t.Errorf("Mismatch. Case %d, Space should not be allowed", i)
			}

		})
	}
}

func TestVisaNumbers(t *testing.T) {
	err := LoadCards("../card_types.json")

	if err != nil {
		t.Errorf("Failed to open json file")
	}

	validNumbers := []string{
		"4111 1111 4555 1142",
		"4988 4388 4388 4305",
		"4166 6766 6766 6746",
		"4646 4646 4646 4644",
		"4000 6200 0000 0007",
		"4000 0600 0000 0006",
		"4293 1891 0000 0008",
		"4988 0800 0000 0000",
		"4111 1111 1111 1111",
		"4444 3333 2222 1111",
		"4001 5900 0000 0001",
		"4000 1800 0000 0002",
		"4000 0200 0000 0000",
		"4000 1600 0000 0004",
		"4002 6900 0000 0008",
		"4400 0000 0000 0008",
		"4484 6000 0000 0004",
		"4607 0000 0000 0009",
		"4977 9494 9494 9497",
		"4000 6400 0000 0005",
		"4003 5500 0000 0003",
		"4000 7600 0000 0001",
		"4017 3400 0000 0003",
		"4005 5190 0000 0006",
		"4131 8400 0000 0003",
		"4035 5010 0000 0008",
		"4151 5000 0000 0008",
		"4571 0000 0000 0001",
		"4199 3500 0000 0002",
	}

	for i, v := range validNumbers {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			s := strings.Map(func(r rune) rune {
				if unicode.IsSpace(r) {
					return -1
				}
				return r
			}, v)

			if getCreditCardType(s)[0].name != "Visa" {
				t.Errorf("Mismatch. Case %d, Card is a Visa", i)
			}

		})
	}
}

func TestMasterCardNumbers(t *testing.T) {
	err := LoadCards("../card_types.json")

	if err != nil {
		t.Errorf("Failed to open json file")
	}

	validNumbers := []string{
		"5101 1800 0000 0007",
		"2222 4000 7000 0005",
		"5100 2900 2900 2909",
		"5555 3412 4444 1115",
		"5577 0000 5577 0004",
		"5136 3333 3333 3335",
		"5585 5585 5585 5583",
		"5555 4444 3333 1111",
		"2222 4107 4036 0010",
		"5555 5555 5555 4444",
		"2222 4107 0000 0002",
		"2222 4000 1000 0008",
		"2223 0000 4841 0010",
		"2222 4000 6000 0007",
		"2223 5204 4356 0010",
		"5500 0000 0000 0004",
		"2222 4000 3000 0004",
		"5100 0600 0000 0002",
		"5100 7050 0000 0002",
		"5103 2219 1119 9245",
		"5424 0000 0000 0015",
		"2222 4000 5000 0009",
		"5106 0400 0000 0008",
	}

	for i, v := range validNumbers {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			s := strings.Map(func(r rune) rune {
				if unicode.IsSpace(r) {
					return -1
				}
				return r
			}, v)

			if getCreditCardType(s)[0].name != "MasterCard" {
				t.Errorf("Mismatch. Case %d, Card is a MasterCard", i)
			}

		})
	}
}
