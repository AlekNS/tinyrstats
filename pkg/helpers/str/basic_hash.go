package str

// BasicStrHash simple hashes a string.
func BasicStrHash(str string) int {
	return BasicBytesHash([]byte(str))
}

// BasicBytesHash simple hashes bytes of string.
// Basic implementation of FNV 32 hash.
func BasicBytesHash(bytes []byte) int {
	const prime = 0x01000193
	var val = 0x811c9dc5

	for _, n := range bytes {
		val *= prime
		val ^= int(n)
	}

	if val < 0 {
		val = -val
	}

	return val
}
