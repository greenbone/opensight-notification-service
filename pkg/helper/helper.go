// SPDX-FileCopyrightText: 2024 Greenbone AG <https://greenbone.net>
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package helper

func ToPtr[T any](val T) *T {
	return &val
}

// ToNullablePtr converts a value to a pointer. If the passed value is the zero value of its type, it returns nil.
func ToNullablePtr[T comparable](val T) *T {
	var zeroT T
	if val == zeroT {
		return nil
	}

	return &val
}

// SafeDereference return the value ptr points to. If ptr is nil, it returns the default value if the type instead.
func SafeDereference[T any](ptr *T) T {
	var zeroT T
	if ptr == nil {
		return zeroT
	}
	return *ptr
}
