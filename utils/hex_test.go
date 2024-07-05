package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValidHex_HappyPath(t *testing.T) {
	valid := IsValidHex("0x388c818ca8b9251b393131c08a736a67ccb19297")

	assert.True(t, valid)
}

func TestIsValidHex_HappyPath_Without_0x(t *testing.T) {
	valid := IsValidHex("388c818ca8b9251b393131c08a736a67ccb19297")

	assert.True(t, valid)
}

func TestIsValidHex_InvalidHex(t *testing.T) {
	invalid := IsValidHex("potato")

	assert.False(t, invalid)
}
