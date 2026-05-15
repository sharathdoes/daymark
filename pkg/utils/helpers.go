package utils

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
)

// GenerateOTP returns a 6-digit numeric OTP as a string.
func GenerateOTP() string {
	var n uint32
	_ = binary.Read(rand.Reader, binary.LittleEndian, &n)
	code := n % 1000000
	return fmt.Sprintf("%06d", code)
}
