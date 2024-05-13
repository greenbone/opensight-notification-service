// SPDX-FileCopyrightText: 2024 Greenbone AG <https://greenbone.net>
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package errs

import (
	"errors"
	"fmt"
)

// ErrItemNotFound is the error used when looking up an item by ID, e.g.
// OID for VTs, fails because the item cannot be found.
var ErrItemNotFound = errors.New("item not found")

// embed this error to mark an error as retryable
var ErrRetryable = errors.New("(retryable error)")

// ErrConflict indicates a conflict. If there are certain fields conflicting which are meaningful to the client,
// set the individual error message for a property via `Errors`, otherwise just set `Message`.
type ErrConflict struct {
	Message string
	Errors  map[string]string // maps property to specific error message
}

func (e *ErrConflict) Error() string {
	message := e.Message
	if len(e.Errors) > 0 {
		message += fmt.Sprintf(", specific errors: %v", e.Errors)
	}
	return message
}

type ErrValidation struct {
	Message string
	Errors  map[string]string // maps property to specific error message
}

func (e *ErrValidation) Error() string {
	message := e.Message
	if len(e.Errors) > 0 {
		message += fmt.Sprintf(", specific errors: %v", e.Errors)
	}
	return message
}
