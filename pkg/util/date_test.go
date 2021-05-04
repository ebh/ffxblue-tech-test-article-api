package util_test

import (
	"awesomeProject/pkg/util"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestConvertDateParam_Valid(t *testing.T) {
	tcs := []struct {
		input    string
		expected time.Time
	}{
		{
			input:    "20200101",
			expected: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			input:    "20201230",
			expected: time.Date(2020, 12, 30, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, tc := range tcs {
		t.Run(tc.input, func(t *testing.T) {
			assert.NotPanics(t, func() {
				actual, err := util.ConvertCondensedDateString(tc.input)

				assert.NoError(t, err)
				assert.Equal(t, tc.expected, actual)
			})
		})
	}

}

func TestConvertDateParam_Invalid(t *testing.T) {
	tcs := []string{
		"0200101",
		"2020010",
		"20201301",
		"20200132",
		"",
	}

	for _, tc := range tcs {
		t.Run(tc, func(t *testing.T) {
			assert.NotPanics(t, func() {
				_, err := util.ConvertCondensedDateString(tc)

				assert.Error(t, err)
			})
		})
	}
}
