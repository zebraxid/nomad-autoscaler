package modes

import (
	"strings"
	"testing"

	"github.com/hashicorp/go-multierror"
	"github.com/stretchr/testify/assert"
)

func TestChecker_ValidateStruct(t *testing.T) {
	testCases := []struct {
		name            string
		enabled         []string
		input           interface{}
		expectedInvalid []string
	}{
		{
			name:            "nil",
			enabled:         []string{},
			input:           nil,
			expectedInvalid: []string{},
		},
		{
			name:            "empty",
			enabled:         []string{},
			input:           &TestStruct{},
			expectedInvalid: []string{},
		},
		{
			name:            "not pointer",
			enabled:         []string{},
			input:           TestStruct{},
			expectedInvalid: []string{},
		},
		{
			name:            "not struct",
			enabled:         []string{},
			input:           "hi",
			expectedInvalid: []string{},
		},
		{
			name:    "no mode enabled",
			enabled: []string{},
			input:   NewTestStructFull(),
			expectedInvalid: []string{
				"top_level_ent",
				"top_level_expert",
				"top_level_ent_expert",
				"top_level_pro",
				"nested_none -> nested_field_ent",
				"nested_none -> deep_nested -> deep_nested_pro",
				"nested_pro -> nested_field_none",
				"nested_pro -> nested_field_ent",
				"nested_pro -> deep_nested -> deep_nested_pro",
				"nested_pro_expert -> nested_field_none",
				"nested_pro_expert -> nested_field_ent",
				"nested_pro_expert -> deep_nested -> deep_nested_pro",
				"nested_multiple -> deep_nested -> deep_nested_pro",
				"nested_multiple -> deep_nested -> deep_nested_pro",
			},
		},
		{
			name:    "ent enabled",
			enabled: []string{"ent"},
			input:   NewTestStructFull(),
			expectedInvalid: []string{
				"top_level_pro",
				"top_level_expert",
				"nested_none -> deep_nested -> deep_nested_pro",
				"nested_pro -> nested_field_none",
				"nested_pro -> deep_nested -> deep_nested_pro",
				"nested_pro_expert -> nested_field_none",
				"nested_pro_expert -> deep_nested -> deep_nested_pro",
				"nested_multiple -> deep_nested -> deep_nested_pro",
				"nested_multiple -> deep_nested -> deep_nested_pro",
			},
		},
		{
			name:    "pro and expert enabled",
			enabled: []string{"pro", "expert"},
			input:   NewTestStructFull(),
			expectedInvalid: []string{
				"top_level_ent",
				"nested_none -> nested_field_ent",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c := NewChecker(TestModes, tc.enabled)
			err := c.ValidateStruct(tc.input)

			if len(tc.expectedInvalid) == 0 {
				assert.NoError(t, err)
				return
			}

			assert.Error(t, err)

			if mErr, ok := err.(*multierror.Error); ok {
				// Check if the expected errors are present in error list.
			OUTER:
				for _, invalid := range tc.expectedInvalid {
					for _, e := range mErr.Errors {
						if strings.Contains(e.Error(), invalid) {
							continue OUTER
						}
					}
					t.Errorf("expected error for %q", invalid)
				}

				// Check that _only_ expected errors are present in error list.
			OUTER_2:
				for _, e := range mErr.Errors {
					for _, invalid := range tc.expectedInvalid {
						if strings.Contains(e.Error(), invalid) {
							continue OUTER_2
						}
					}
					t.Errorf("unexpected error: %v", e)
				}
			}

		})
	}
}
