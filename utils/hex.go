package utils

import (
	"encoding/hex"
)

func IsValidHex(s string) bool {
	if len(s) >= 2 && s[:2] == "0x" {
		s = s[2:]
	}

	_, err := hex.DecodeString(s)
	return err == nil
}
