package validation

import "strconv"

type Luhn struct {
}

func (l Luhn) IsValid(number string) bool {
	var sum int
	var alternate bool

	numberLen := len(number)

	if numberLen < 13 || numberLen > 19 {
		return false
	}

	for i := numberLen - 1; i > -1; i-- {
		mod, err := strconv.Atoi(string(number[i]))
		if err != nil {
			return false
		}

		if alternate {
			mod *= 2
			if mod > 9 {
				mod = (mod % 10) + 1
			}
		}

		alternate = !alternate

		sum += mod
	}

	return sum%10 == 0
}
