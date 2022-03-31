package git

import (
	"bytes"
	"fmt"
	"time"
)

// FallbackTimeValue is the value returned by `SafeTimeParse` in case it
// encounters a parse error. It's the maximum time value possible in golang.
// See https://gitlab.com/gitlab-org/gitaly/issues/556#note_40289573
var FallbackTimeValue = time.Unix(1<<63-62135596801, 999999999)

func validateRevision(revision []byte, allowEmpty bool) error {
	if !allowEmpty && len(revision) == 0 {
		return fmt.Errorf("empty revision")
	}
	if bytes.HasPrefix(revision, []byte("-")) {
		return fmt.Errorf("revision can't start with '-'")
	}
	if bytes.Contains(revision, []byte(" ")) {
		return fmt.Errorf("revision can't contain whitespace")
	}
	if bytes.Contains(revision, []byte("\x00")) {
		return fmt.Errorf("revision can't contain NUL")
	}
	if bytes.Contains(revision, []byte(":")) {
		return fmt.Errorf("revision can't contain ':'")
	}
	return nil
}

// ValidateRevisionAllowEmpty checks if a revision looks valid, but allows
// empty strings
func ValidateRevisionAllowEmpty(revision []byte) error {
	return validateRevision(revision, true)
}

// ValidateRevision checks if a revision looks valid
func ValidateRevision(revision []byte) error {
	return validateRevision(revision, false)
}
