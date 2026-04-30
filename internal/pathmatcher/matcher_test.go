package pathmatcher

import (
	"reflect"
	"testing"
)

func TestMatcher_ExtractData(t *testing.T) {
	matcher := NewMatcher()

	// Define a structure for test cases
	type testCase struct {
		name           string
		formatTemplate string
		rawPath        string
		expectedData   map[string]string
		expectError    bool
	}

	// --- Test Cases ---
	tests := []testCase{
		{
			name:           "Standard Match Success",
			formatTemplate: "{Artist}/{Album Title} ({Release Year})/{track:00} - {Track Title}",
			rawPath:        "The Beatles/Abbey Road (1969)/01 - Come Together",
			expectedData: map[string]string{
				"Artist":       "The Beatles",
				"Album Title":  "Abbey Road",
				"Release Year": "(1969)", // Note: The regex capture will include the parentheses unless handled specifically
				"track:00":     "01",
				"Track Title":  "Come Together",
			},
			expectError: false,
		},
		{
			name:           "Complex Match Success",
			formatTemplate: "{Artist}/{Album Title} ({Release Year})/{Medium Format} {medium:00}/{track:00} - {Track Title}",
			rawPath:        "Queen/A Night at the Opera (1975)/LP 01/01 - Love of My Life",
			expectedData: map[string]string{
				"Artist":        "Queen",
				"Album Title":   "A Night at the Opera",
				"Release Year":  "(1975)",
				"Medium Format": "LP",
				"medium:00":     "01",
				"track:00":      "01",
				"Track Title":   "Love of My Life",
			},
			expectError: false,
		},
		{
			name:           "Failure_FormatMismatch",
			formatTemplate: "{Artist}/{Album Title} ({Release Year})/{track:00} - {Track Title}",
			// The path is missing the necessary / delimiters expected by the regex structure
			rawPath:      "The Beatles Abbey Road (1969) 01 - Come Together",
			expectedData: nil,
			expectError:  true,
		},
		{
			name:           "EdgeCase_EmptyFormat",
			formatTemplate: "",
			rawPath:        "some/path/here",
			expectedData:   nil,
			expectError:    true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			actualData, err := matcher.ExtractData(tc.formatTemplate, tc.rawPath)

			// Assert Error Expectation
			if tc.expectError {
				if err == nil {
					t.Errorf("Expected an error but got nil. Data: %v", actualData)
				}
				return // Stop further checks if an error was expected
			}

			// Assert No Error
			if err != nil {
				t.Fatalf("Did not expect an error, but got: %v", err)
			}

			// Assert Data Content
			if !reflect.DeepEqual(actualData, tc.expectedData) {
				t.Errorf("Extracted data mismatch.\nExpected: %#v\nActual:   %#v", tc.expectedData, actualData)
			}
		})
	}
}
