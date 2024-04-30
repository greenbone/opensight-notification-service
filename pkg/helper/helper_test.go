// SPDX-FileCopyrightText: 2024 Greenbone AG <https://greenbone.net>
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package helper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSafeDereference(t *testing.T) {
	type someStruct struct {
		foo, bar string
	}

	zeroValue := someStruct{}
	nonZeroValue := someStruct{foo: "aa", bar: "bb"}

	tests := []struct {
		name  string
		input *someStruct
		want  someStruct
	}{
		{
			name:  "yields zero value with nil pointer",
			input: nil,
			want:  zeroValue,
		},
		{
			name:  "yields object pointed on for non-nil pointer",
			input: &nonZeroValue,
			want:  nonZeroValue,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SafeDereference(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}
