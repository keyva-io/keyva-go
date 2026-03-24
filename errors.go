// Keyva error types.
//
// Auto-generated from keyva protocol spec. Do not edit.

package keyva

import (
	"fmt"
	"strings"
)

// KeyvaError represents an error returned by the Keyva server.
type KeyvaError struct {
	// Machine-readable error code (e.g. "NOTFOUND", "DENIED").
	Code string

	// Human-readable error message.
	Message string
}

func (e *KeyvaError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func parseError(payload string) *KeyvaError {
	code, message, _ := strings.Cut(payload, " ")
	return &KeyvaError{Code: code, Message: message}
}

// Error code constants.
const (
	// ErrBadarg — Missing or malformed command argument
	ErrBadarg = "BADARG"
	// ErrChainLimit — Refresh token chain limit exceeded
	ErrChainLimit = "CHAIN_LIMIT"
	// ErrCrypto — Cryptographic operation failed
	ErrCrypto = "CRYPTO"
	// ErrDenied — Authentication required or insufficient permissions
	ErrDenied = "DENIED"
	// ErrDisabled — Keyspace is disabled
	ErrDisabled = "DISABLED"
	// ErrExpired — Credential has expired
	ErrExpired = "EXPIRED"
	// ErrInternal — Unexpected internal error
	ErrInternal = "INTERNAL"
	// ErrLocked — Account temporarily locked due to too many failed attempts
	ErrLocked = "LOCKED"
	// ErrNotfound — Credential, keyspace, or resource does not exist
	ErrNotfound = "NOTFOUND"
	// ErrNotready — Server is not ready (still starting up)
	ErrNotready = "NOTREADY"
	// ErrReuseDetected — Refresh token reuse detected — family revoked
	ErrReuseDetected = "REUSE_DETECTED"
	// ErrStateError — Credential is in wrong state for this operation
	ErrStateError = "STATE_ERROR"
	// ErrStorage — Storage engine error
	ErrStorage = "STORAGE"
	// ErrValidationError — Metadata or claims failed schema validation
	ErrValidationError = "VALIDATION_ERROR"
	// ErrWrongtype — Operation not supported for this keyspace type
	ErrWrongtype = "WRONGTYPE"
)

// IsCode reports whether the error has the given code.
func IsCode(err error, code string) bool {
	if ke, ok := err.(*KeyvaError); ok {
		return ke.Code == code
	}
	return false
}
