package validation

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"strings"
)

type PatternType int

const (
	Int PatternType = iota
	Range
)

const (
	SUCCS int = iota // Success
	UKNWN            // General failure, unknown issuer, failed match and length
	INVDN            // Failed verification
)

const (
	UnknownNumber = "Unknown Card Number."                                     // Unknown Issuer
	InvalidNumber = "Issuer is found but number failed checksum verification." // Failed Luhn verification
)

var (
	cardTypes CardTypes
)

type CardTypes []CardConfig

type CardConfig struct {
	Name     string        `json:"name"`
	Patterns []CardPattern `json:"patterns"`
	Lengths  []int         `json:"lengths"`
	Error    string
}

type RangeValue struct {
	Min, Max int
}

type IntValue struct {
	Val int
}

/**
 * CardPattern Structure
 * Multiple Issuers uses different patterns, some are ranges between min...max (Mastercard), sometimes just a single number (e.g. Visa)
 * Some Issuers have even have multiple patterns/ranges (e.g., Elo cards), so we need to be able to list a bunch of these.
 */
type CardPattern struct {
	Value interface{}
	T     PatternType
}

// check if this pattern matches the cardNumber
func (cp *CardPattern) matches(cardNumber string) bool {
	if cp.T == Range {
		min_s := strconv.Itoa(cp.Value.(RangeValue).Min)
		maxLen := len(min_s) // the minimum value will be the base length
		if maxLen <= len(cardNumber) {
			substr := cardNumber[:maxLen] // get the substring from the number with the given length
			if i, err := strconv.Atoi(substr); err == nil {
				var min, max int
				var err error
				min, err = strconv.Atoi(min_s[:len(substr)])
				if err != nil {
					return false
				}

				max_s := strconv.Itoa(cp.Value.(RangeValue).Max)
				max, err = strconv.Atoi(max_s[:len(substr)])
				if err != nil {
					return false
				}

				return i >= min && i <= max
			}
		}
	} else {
		patt := strconv.Itoa(cp.Value.(IntValue).Val)
		if len(patt) <= len(cardNumber) {
			return patt == cardNumber[:len(patt)]
		}
	}

	return false

}

// The CardPattern accepts values of homogenous type (int or []int), we have to unmarshal this properly.
func (cp *CardPattern) UnmarshalJSON(b []byte) error {
	if b[0] == '[' {
		cp.T = Range
		substr := strings.Split(string(b[1:len(b)-1]), ",")
		var err error
		min, err := strconv.Atoi(strings.TrimSpace(substr[0]))
		if err != nil {
			return err
		}
		max, err := strconv.Atoi(strings.TrimSpace(substr[1]))
		if err != nil {
			return err
		}
		cp.Value = RangeValue{int(min), int(max)}
	} else {
		cp.T = Int
		i, err := strconv.Atoi(strings.TrimSpace(string(b)))
		if err != nil {
			return err
		}
		cp.Value = IntValue{i}
	}

	return nil
}

type results []*result

// This would sort the result by the highest PatternMatch, thus a number with 401178 would match an Elo more than a Visa.
func (a results) Len() int           { return len(a) }
func (a results) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a results) Less(i, j int) bool { return a[i].patternMatch > a[j].patternMatch }

// internal result placeholder
type result struct {
	name         string
	patternMatch int
	lengthMatch  int
}

type TopResult struct {
	Valid        bool
	Name         string
	Error        Error
	PatternMatch int
	LengthMatch  int
}

type Error struct {
	ErrorNo int
	Message string
}

// Load card types given a path to the json file
func LoadCards(path string) error {
	// Open our jsonFile
	if path == "" {
		return fmt.Errorf("Unable to open file from an empty path")
	}

	jsonFile, err := os.Open(path)
	// if we os.Open returns an error then handle it
	if err != nil {
		return err
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	if err := json.Unmarshal(byteValue, &cardTypes); err != nil {
		return err
	}

	return nil
}

func getCreditCardType(cardNumber string) results {
	// We need to produce a list of cards that matched, put a score of how much characters have matched, then sort them.
	results := make(results, 0)

	if !isValidInputType(cardNumber) {
		return nil
	}

	// check cardNumber per card type
	for _, cc := range cardTypes {
		clen := len(cardNumber)
		var lengthMatch bool

		// We can't deal with lengths as range of min..max because some cards have discrete values like Discover/Diner's
		for _, cl := range cc.Lengths {
			if cl == clen {
				lengthMatch = true
			}
		}

		// If we don't get any match, don't bother to go on.
		if !lengthMatch {
			continue
		}

		// Issuers/Banks have multiple patterns (especially Elo, which could start with a 4 or 5, like Visa or MC)
		// So we'll check the card number against these patterns (provided they've passed the length matching).
		// If a single pattern matches then stop, the card is of this type.
		for i := 0; i < len(cc.Patterns); i++ {
			pattern := cc.Patterns[i]

			if !pattern.matches(cardNumber) {
				continue
			}

			var patternLength int
			if pattern.T == Range {
				patternLength = len(strconv.Itoa(pattern.Value.(RangeValue).Min))
			} else {
				patternLength = len(strconv.Itoa(pattern.Value.(IntValue).Val))
			}

			r := new(result)
			r.lengthMatch = clen

			if clen >= patternLength {
				r.patternMatch = patternLength
			}

			r.name = cc.Name

			results = append(results, r)
			break
		}
	}

	return results
}

// only allow digits
func isValidInputType(cardNumber string) bool {
	isNotDigit := func(c rune) bool { return c < '0' || c > '9' }
	return strings.IndexFunc(cardNumber, isNotDigit) == -1
}

func Validate(creditCardNumber string) *TopResult {
	results := getCreditCardType(creditCardNumber)
	var e Error
	if len(results) == 0 {
		e = Error{UKNWN, UnknownNumber}
		return &TopResult{Name: "Unknown", Error: e}
	}

	var luhn Luhn
	sort.Sort(results)

	if !luhn.IsValid(creditCardNumber) {
		e = Error{INVDN, fmt.Sprintf(InvalidNumber, results[0].name)}
		return &TopResult{Name: results[0].name, PatternMatch: results[0].patternMatch, LengthMatch: results[0].lengthMatch, Error: e}
	}

	e = Error{SUCCS, "Success"}
	return &TopResult{Valid: true, Name: results[0].name, PatternMatch: results[0].patternMatch, LengthMatch: results[0].lengthMatch, Error: e}

}
