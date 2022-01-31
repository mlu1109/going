package filesys

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseFileName_whenValid_thenReturnExpectedPartsAndNoError(t *testing.T) {
	tests := []struct {
		input               string
		expectedVersion     string
		expectedDescription string
	}{
		{"V2__this_is_version_2.sql", "2", "this_is_version_2"},
		{"V2_1__this__is__version_2_1.sql", "2_1", "this__is__version_2_1"},
		{"V11.222.3333__can't_be___arsed_to_write_version.sql", "11.222.3333", "can't_be___arsed_to_write_version"},
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
		actualVersion, actualDescription, actualError := parseFileName(test.input)
		assert.Len(t, actualVersion, 0)
		assert.Len(t, actualDescription, 0)
		assert.Equal(t, test.expectedError, actualError)
	}
}

func TestLoadMigrationsFromDir_whenValid_thenReturnExpectedMigrations(t *testing.T) {

}
