// SPDX-FileCopyrightText: 2023 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package mockitmetrics

import (
	"errors"
	"fmt"
	"strings"
)

type tuple struct {
	label string
	value string
}

var (
	errInvalidLabelValues = errors.New("labelValues is invalid")
)

func convert(s []string) ([]tuple, error) {
	if len(s)%2 != 0 {
		return nil, fmt.Errorf("%w - must be a multiple of 2, 'label1', 'value1', 'label2', 'value2', ...", //nolint:staticcheck
			errInvalidLabelValues)
	}

	rv := make([]tuple, 0, len(s)/2)
	for i := 0; i < len(s); i += 2 {
		if s[i] == "" {
			return nil, fmt.Errorf("%w - the label must not be empty", errInvalidLabelValues)
		}
		if s[i+1] == "" {
			return nil, fmt.Errorf("%w - the value must not be empty", errInvalidLabelValues)
		}
		rv = append(rv, tuple{
			label: s[i],
			value: s[i+1],
		})
	}

	return rv, nil
}

func validateLabels(expected *[]string, actual []tuple, exact bool) error {
	list := make([]string, 0, len(actual))
	for _, t := range actual {
		list = append(list, t.label)
	}

	// Only validate if expected is not nil.
	if expected == nil {
		return nil
	}

	wanted := *expected

	if exact && len(wanted) != len(actual) {
		return fmt.Errorf("%w - expected labels: want '%s', got '%s'",
			errInvalidLabelValues,
			strings.Join(wanted, "', '"),
			strings.Join(list, "', '"))
	}

	if !exact && len(wanted) < len(actual) {
		return fmt.Errorf("%w - too many labels: want '%s', got '%s'",
			errInvalidLabelValues,
			strings.Join(wanted, "', '"),
			strings.Join(list, "', '"))
	}

	for i := range actual {
		if wanted[i] != actual[i].label {
			return fmt.Errorf("%w - the labels do not match: want '%s', got '%s'",
				errInvalidLabelValues,
				strings.Join(wanted, "', '"),
				strings.Join(list, "', '"))
		}
	}

	return nil
}

func joinValues(t []tuple, delimiter string) string {
	if len(t) == 0 {
		return ""
	}

	rv := make([]string, 0, len(t))
	for _, v := range t {
		rv = append(rv, v.value)
	}

	return strings.Join(rv, delimiter)
}
