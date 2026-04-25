package main

import (
	"errors"
	"strings"
	"testing"
)

func TestValidateName(t *testing.T) {
	cases := []struct {
		label   string
		in      string
		want    string
		wantErr error
	}{
		{"simple", "alpha", "alpha", nil},
		{"trims surrounding whitespace", "  alpha  ", "alpha", nil},
		{"unicode passes", "Círculo", "Círculo", nil},
		{"exact max length", strings.Repeat("a", 50), strings.Repeat("a", 50), nil},
		{"max length by rune count", strings.Repeat("é", 50), strings.Repeat("é", 50), nil},
		{"empty rejected", "", "", errNameRequired},
		{"whitespace only rejected", "   ", "", errNameRequired},
		{"too long by one rejected", strings.Repeat("a", 51), "", errNameTooLong},
		{"newline in middle rejected", "foo\nbar", "", errNameInvalidChar},
		{"tab in middle rejected", "foo\tbar", "", errNameInvalidChar},
		{"null byte rejected", "foo\x00bar", "", errNameInvalidChar},
	}

	for _, tc := range cases {
		t.Run(tc.label, func(t *testing.T) {
			got, err := validateName(tc.in)
			if !errors.Is(err, tc.wantErr) {
				t.Errorf("err=%v, want %v", err, tc.wantErr)
			}
			if got != tc.want {
				t.Errorf("got=%q, want %q", got, tc.want)
			}
		})
	}
}
