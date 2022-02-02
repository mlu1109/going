package filesys

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseFileName_whenValid_thenReturnExpectedPartsAndNoError(t *testing.T) {
	tests := []struct {
		input               string
		expectedVersion     uint
		expectedDescription string
	}{
		{"V2__this_is_version_2.sql", 2, "this_is_version_2"},
		{"V3__this__is__version_3.sql", 3, "this__is__version_3"},
	}
	for _, test := range tests {
		actualVersion, actualDescription, actualError := parseFileName(test.input)
		assert.Equal(t, test.expectedVersion, actualVersion)
		assert.Equal(t, test.expectedDescription, actualDescription)
		assert.Nil(t, actualError)
	}
}

func TestParseFileName_whenInvalid_thenReturnBlanksAndError(t *testing.T) {
	tests := []struct {
		input         string
		expectedError error
	}{
		{"V__this_is_version_2.sql", fmt.Errorf("invalid filename")},
		{"V_2_1__this__is__version_2_1.sql", fmt.Errorf("invalid filename")},
		{"VA__can't_be___arsed_to_write_version.sql", fmt.Errorf("invalid filename")},
	}
	for _, test := range tests {
		_, actualDescription, actualError := parseFileName(test.input)
		assert.Len(t, actualDescription, 0)
		assert.Equal(t, test.expectedError, actualError)
	}
}
