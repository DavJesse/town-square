package utils

import (
	"testing"
	"time"
)

/*
* Tests ConvertToEAT function.
* 1. Valid UTC time: We provide a valid UTC time in RFC3339 format (2025-02-16T12:00:00Z) and check that it converts correctly to East African Time (EAT).
* 2. Invalid UTC time format: We provide an invalid format ("2025-02-16 12:00:00") and ensure that the function returns an error.
* 3. Empty UTC time: We test an empty string and expect an error.
 */
func TestConvertToEAT(t *testing.T) {
	table := []struct {
		name           string
		utcTime        string
		expectedResult time.Time
		expectError    bool
	}{
		{
			name:           "Valid UTC time",
			utcTime:        "2025-02-16T12:00:00Z",                                                         // UTC time
			expectedResult: time.Date(2025, time.February, 16, 15, 0, 0, 0, time.FixedZone("EAT", 3*3600)), // EAT time (UTC+3)
			// includes a time.FixedZone to ensure that we are correctly comparing the EAT time, which is UTC +3 hours.
			expectError: false,
		},
		{
			name:           "Invalid UTC time format",
			utcTime:        "2025-02-16 12:00:00",
			expectedResult: time.Time{},
			expectError:    true,
		},
		{
			name:           "Empty UTC time",
			utcTime:        "",
			expectedResult: time.Time{},
			expectError:    true,
		},
	}

	for _, entry := range table {
		t.Run(entry.name, func(t *testing.T) {
			result, err := ConvertToEAT(entry.utcTime)

			// check if error is expected
			if (err != nil) != entry.expectError {
				t.Errorf("expected error: %v, got: %v", entry.expectError, err)
			}

			// check if the result matches expected
			if !entry.expectedResult.IsZero() && !result.Equal(entry.expectedResult) {
				t.Errorf("expected %v, got %v", entry.expectedResult, result)
			}
		})
	}
}
