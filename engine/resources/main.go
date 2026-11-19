package resources

import (
	"errors"
	"fmt"
	"regexp"
)

var regex_invalid = regexp.MustCompile("[^a-zA-Z0-9._]")

// IsNameInvalid is used to check the incoming names for validity.
//
// A name can only contain [a-zA-Z0-9._] and cant have a length of 0.
func IsNameInvalid(name string) bool {
	return regex_invalid.MatchString(name) || len(name) == 0
}

// IsNameInvalid is used to check the incoming names for validity.
//
// A name can only contain [a-zA-Z0-9._] and cant have a length of 0.
func IsNameInvalidError(name string) error {
	errs := []error{}
	if regex_invalid.MatchString(name) {
		errs = append(errs, fmt.Errorf("name %q contains invalid characters", name))
	}
	if len(name) == 0 {
		errs = append(errs, fmt.Errorf("name %q has a length of zero", name))
	}
	return errors.Join(errs...)
}
