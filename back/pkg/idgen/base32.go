package idgen

// Crockford Base32 alphabet (32 characters, excluding I, L, O, U to avoid confusion)
const crockfordAlphabet = "0123456789ABCDEFGHJKMNPQRSTVWXYZ"

// EncodeCrockford encodes a uint64 number to Crockford Base32 string
// Returns a minimum of 6 characters (capacity: 32^6 = 1,073,741,824)
func EncodeCrockford(num uint64) string {
	if num == 0 {
		return "000000"
	}

	var result []byte
	for num > 0 {
		result = append([]byte{crockfordAlphabet[num%32]}, result...)
		num /= 32
	}

	// Pad with leading zeros to ensure minimum 6 characters
	for len(result) < 6 {
		result = append([]byte{'0'}, result...)
	}

	return string(result)
}

// DecodeCrockford decodes a Crockford Base32 string to uint64
// Returns 0 if the string contains invalid characters
func DecodeCrockford(encoded string) uint64 {
	var result uint64

	for _, char := range encoded {
		result *= 32

		// Find position in alphabet
		var value int = -1
		for i, c := range crockfordAlphabet {
			if c == char {
				value = i
				break
			}
		}

		// Invalid character
		if value == -1 {
			return 0
		}

		result += uint64(value)
	}

	return result
}
